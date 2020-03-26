// Code generated by the FlatBuffers compiler. DO NOT EDIT.

package fbast

import (
	flatbuffers "github.com/google/flatbuffers/go"
)

type LogicalExpression struct {
	_tab flatbuffers.Table
}

func GetRootAsLogicalExpression(buf []byte, offset flatbuffers.UOffsetT) *LogicalExpression {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &LogicalExpression{}
	x.Init(buf, n+offset)
	return x
}

func (rcv *LogicalExpression) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *LogicalExpression) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *LogicalExpression) BaseNode(obj *BaseNode) *BaseNode {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		x := rcv._tab.Indirect(o + rcv._tab.Pos)
		if obj == nil {
			obj = new(BaseNode)
		}
		obj.Init(rcv._tab.Bytes, x)
		return obj
	}
	return nil
}

func (rcv *LogicalExpression) Operator() LogicalOperator {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.GetInt8(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *LogicalExpression) MutateOperator(n LogicalOperator) bool {
	return rcv._tab.MutateInt8Slot(6, n)
}

func (rcv *LogicalExpression) LeftType() byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		return rcv._tab.GetByte(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *LogicalExpression) MutateLeftType(n byte) bool {
	return rcv._tab.MutateByteSlot(8, n)
}

func (rcv *LogicalExpression) Left(obj *flatbuffers.Table) bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(10))
	if o != 0 {
		rcv._tab.Union(obj, o)
		return true
	}
	return false
}

func (rcv *LogicalExpression) RightType() byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(12))
	if o != 0 {
		return rcv._tab.GetByte(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *LogicalExpression) MutateRightType(n byte) bool {
	return rcv._tab.MutateByteSlot(12, n)
}

func (rcv *LogicalExpression) Right(obj *flatbuffers.Table) bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(14))
	if o != 0 {
		rcv._tab.Union(obj, o)
		return true
	}
	return false
}

func LogicalExpressionStart(builder *flatbuffers.Builder) {
	builder.StartObject(6)
}
func LogicalExpressionAddBaseNode(builder *flatbuffers.Builder, baseNode flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(0, flatbuffers.UOffsetT(baseNode), 0)
}
func LogicalExpressionAddOperator(builder *flatbuffers.Builder, operator int8) {
	builder.PrependInt8Slot(1, operator, 0)
}
func LogicalExpressionAddLeftType(builder *flatbuffers.Builder, leftType byte) {
	builder.PrependByteSlot(2, leftType, 0)
}
func LogicalExpressionAddLeft(builder *flatbuffers.Builder, left flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(3, flatbuffers.UOffsetT(left), 0)
}
func LogicalExpressionAddRightType(builder *flatbuffers.Builder, rightType byte) {
	builder.PrependByteSlot(4, rightType, 0)
}
func LogicalExpressionAddRight(builder *flatbuffers.Builder, right flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(5, flatbuffers.UOffsetT(right), 0)
}
func LogicalExpressionEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}