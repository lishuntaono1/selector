package internal

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
	node_type    NodeType
}

func (a *Axis) Type() AstType {
	return AstAxis
}

func (a *Axis) ReturnType() ResultType {
	return NodeSetType
}
