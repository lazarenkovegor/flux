package compiler_test

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/influxdata/flux"
	"github.com/influxdata/flux/compiler"
	"github.com/influxdata/flux/internal/arrowutil"
	"github.com/influxdata/flux/runtime"
	"github.com/influxdata/flux/semantic"
	"github.com/influxdata/flux/values"
)

func TestVectorizedFns(t *testing.T) {
	testCases := []struct {
		name   string
		fn     string
		inType semantic.MonoType
		input  values.Object
		want   values.Value
	}{
		{
			name: "field access",
			fn:   `(r) => ({a: r.a, b: r.b})`,
			inType: semantic.NewObjectType([]semantic.PropertyType{
				{Key: []byte("r"), Value: semantic.NewObjectType([]semantic.PropertyType{
					{Key: []byte("a"), Value: semantic.NewVectorType(semantic.BasicInt)},
					{Key: []byte("b"), Value: semantic.NewVectorType(semantic.BasicInt)},
				})},
			}),
			input: values.NewObjectWithValues(map[string]values.Value{
				"r": values.NewObjectWithValues(map[string]values.Value{
					"a": values.NewInt(int64(1)),
					"b": values.NewInt(int64(1)),
				}),
			}),
			want: values.NewObjectWithValues(map[string]values.Value{
				"a": arrowutil.NewVectorFromSlice([]values.Value{
					values.NewInt(int64(1)),
				}, flux.TInt),
				"b": arrowutil.NewVectorFromSlice([]values.Value{
					values.NewInt(int64(1)),
				}, flux.TInt),
			}),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			pkg, err := runtime.AnalyzeSource(tc.fn)
			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			stmt := pkg.Files[0].Body[0].(*semantic.ExpressionStatement)
			fn := stmt.Expression.(*semantic.FunctionExpression)

			if fn.Vectorized == nil {
				t.Fatal("Expected to find vectorized node, but found none")
			}

			f, err := compiler.Compile(nil, fn, tc.inType)
			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			got, err := f.Eval(context.TODO(), tc.input)
			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			if !cmp.Equal(tc.want, got, CmpOptions...) {
				t.Errorf("unexpected value -want/+got\n%s", cmp.Diff(tc.want, got, CmpOptions...))
			}
		})
	}
}
