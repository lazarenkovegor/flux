package universe

import (
	"time"

	"github.com/influxdata/flux"
	"github.com/influxdata/flux/codes"
	"github.com/influxdata/flux/execute"
	"github.com/influxdata/flux/internal/errors"
	"github.com/influxdata/flux/internal/feature"
	"github.com/influxdata/flux/interpreter"
	"github.com/influxdata/flux/plan"
	"github.com/influxdata/flux/runtime"
	"github.com/influxdata/flux/semantic"
	"github.com/influxdata/flux/values"
)

const ShiftKind = "timeShift"

type ShiftOpSpec struct {
	Shift   flux.Duration `json:"duration"`
	Columns []string      `json:"columns"`
}

func init() {
	shiftSignature := runtime.MustLookupBuiltinType("universe", "timeShift")

	runtime.RegisterPackageValue("universe", ShiftKind, flux.MustValue(flux.FunctionValue(ShiftKind, createShiftOpSpec, shiftSignature)))
	flux.RegisterOpSpec(ShiftKind, newShiftOp)
	plan.RegisterProcedureSpec(ShiftKind, newShiftProcedure, ShiftKind)
	execute.RegisterTransformation(ShiftKind, createShiftTransformation)
}

func createShiftOpSpec(args flux.Arguments, a *flux.Administration) (flux.OperationSpec, error) {
	if err := a.AddParentFromArgs(args); err != nil {
		return nil, err
	}

	spec := new(ShiftOpSpec)

	if shift, err := args.GetRequiredDuration("duration"); err != nil {
		return nil, err
	} else {
		spec.Shift = shift
	}

	if cols, ok, err := args.GetArray("columns", semantic.String); err != nil {
		return nil, err
	} else if ok {
		columns, err := interpreter.ToStringArray(cols)
		if err != nil {
			return nil, err
		}
		spec.Columns = columns
	} else {
		spec.Columns = []string{
			execute.DefaultTimeColLabel,
			execute.DefaultStopColLabel,
			execute.DefaultStartColLabel,
		}
	}
	return spec, nil
}

func newShiftOp() flux.OperationSpec {
	return new(ShiftOpSpec)
}

func (s *ShiftOpSpec) Kind() flux.OperationKind {
	return ShiftKind
}

type ShiftProcedureSpec struct {
	plan.DefaultCost
	Shift   flux.Duration
	Columns []string
	Now     time.Time
}

// TimeBounds implements plan.BoundsAwareProcedureSpec
func (s *ShiftProcedureSpec) TimeBounds(predecessorBounds *plan.Bounds) *plan.Bounds {
	if predecessorBounds != nil {
		return predecessorBounds.Shift(values.Duration(s.Shift))
	}
	return nil
}

func newShiftProcedure(qs flux.OperationSpec, pa plan.Administration) (plan.ProcedureSpec, error) {
	spec, ok := qs.(*ShiftOpSpec)
	if !ok {
		return nil, errors.Newf(codes.Internal, "invalid spec type %T", qs)
	}

	return &ShiftProcedureSpec{
		Shift:   spec.Shift,
		Columns: spec.Columns,
		Now:     pa.Now(),
	}, nil
}

func (s *ShiftProcedureSpec) Kind() plan.ProcedureKind {
	return ShiftKind
}

func (s *ShiftProcedureSpec) Copy() plan.ProcedureSpec {
	ns := new(ShiftProcedureSpec)
	*ns = *s

	if s.Columns != nil {
		ns.Columns = make([]string, len(s.Columns))
		copy(ns.Columns, s.Columns)
	}
	return ns
}

// TriggerSpec implements plan.TriggerAwareProcedureSpec
func (s *ShiftProcedureSpec) TriggerSpec() plan.TriggerSpec {
	return plan.NarrowTransformationTriggerSpec{}
}

func createShiftTransformation(id execute.DatasetID, mode execute.AccumulationMode, spec plan.ProcedureSpec, a execute.Administration) (execute.Transformation, execute.Dataset, error) {
	s, ok := spec.(*ShiftProcedureSpec)
	if !ok {
		return nil, nil, errors.Newf(codes.Internal, "invalid spec type %T", spec)
	}

	if feature.OptimizeShiftTransformation().Enabled(a.Context()) {
		return newShiftTransformation2(id, s, a.Allocator())
	}

	cache := execute.NewTableBuilderCache(a.Allocator())
	d := execute.NewDataset(id, mode, cache)
	t := NewShiftTransformation(d, cache, s)
	return t, d, nil
}

type shiftTransformation struct {
	execute.ExecutionNode
	d       execute.Dataset
	cache   execute.TableBuilderCache
	shift   execute.Duration
	columns []string
}

func NewShiftTransformation(d execute.Dataset, cache execute.TableBuilderCache, spec *ShiftProcedureSpec) *shiftTransformation {
	return &shiftTransformation{
		d:       d,
		cache:   cache,
		shift:   execute.Duration(spec.Shift),
		columns: spec.Columns,
	}
}

func (t *shiftTransformation) RetractTable(id execute.DatasetID, key flux.GroupKey) error {
	return t.d.RetractTable(key)
}

func (t *shiftTransformation) Process(id execute.DatasetID, tbl flux.Table) error {
	key := tbl.Key()
	// Update key
	cols := make([]flux.ColMeta, len(key.Cols()))
	vs := make([]values.Value, len(key.Cols()))
	for j, c := range key.Cols() {
		if execute.ContainsStr(t.columns, c.Label) {
			if c.Type != flux.TTime {
				return errors.Newf(codes.FailedPrecondition, "column %q is not of type time", c.Label)
			}
			cols[j] = c
			vs[j] = values.NewTime(key.ValueTime(j).Add(t.shift))
		} else {
			cols[j] = c
			vs[j] = key.Value(j)
		}
	}
	key = execute.NewGroupKey(cols, vs)

	builder, created := t.cache.TableBuilder(key)
	if !created {
		return errors.Newf(codes.FailedPrecondition, "shift found duplicate table with key: %v", tbl.Key())
	}
	if err := execute.AddTableCols(tbl, builder); err != nil {
		return err
	}

	return tbl.Do(func(cr flux.ColReader) error {
		for j, c := range cr.Cols() {
			if execute.ContainsStr(t.columns, c.Label) {
				l := cr.Len()
				for i := 0; i < l; i++ {
					vs := cr.Times(j)
					if vs.IsNull(i) {
						if err := builder.AppendNil(j); err != nil {
							return err
						}
						continue
					}

					ts := execute.Time(vs.Value(i))
					if err := builder.AppendTime(j, ts.Add(t.shift)); err != nil {
						return err
					}
				}
			} else {
				if err := execute.AppendCol(j, j, cr, builder); err != nil {
					return err
				}
			}
		}
		return nil
	})
}

func (t *shiftTransformation) UpdateWatermark(id execute.DatasetID, mark execute.Time) error {
	return t.d.UpdateWatermark(mark)
}

func (t *shiftTransformation) UpdateProcessingTime(id execute.DatasetID, pt execute.Time) error {
	return t.d.UpdateProcessingTime(pt)
}

func (t *shiftTransformation) Finish(id execute.DatasetID, err error) {
	t.d.Finish(err)
}

func (t *shiftTransformation) SetParents(ids []execute.DatasetID) {}
