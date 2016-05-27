package internal

import "github.com/zhengchun/selector/xpath"

type AxisType uint

const (
	AxisAncestor AxisType = iota
	AxisAncestorOrSelf
	AxisAttribute
	AxisChild
	AxisDescendant
	AxisDescendantOrSelf
	AxisFollowing
	AxisFollowingSibling
	AxisNamespace
	AxisParent
	AxisPreceding
	AxisPrecedingSibling
	AxisSelf
)

type Axis struct {
	input        AstNode
	prefix, name string
	abbr         bool
	axis_type    AxisType
	node_type    xpath.NodeType
}

func (a *Axis) Type() AstNodeType {
	return AxisAst
}

func (a *Axis) ReturnType() ResultType {
	return NodeSetType
}
