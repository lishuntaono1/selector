package internal

import "fmt"

func CompileXPathQuery(query string) Query {
	root := ParseXPathExpression(query)
	builder := &QueryBuilder{}
	return builder.ProcessNode(root)
}

type QueryBuilder struct {
	depth int
}

func (b *QueryBuilder) ProcessNode(root AstNode) Query {
	if b.depth = b.depth + 1; b.depth > 1024 {
		panic("The xpath query is too complex.")
	}
	var result Query
	switch root.Type() {
	case AstAxis:
		result = b.ProcessAxis(root.(*Axis))
	case AstOperator:
		result = b.ProcessOperator(root.(*Operator))
	}
	b.depth--
	return result
}

func (b *QueryBuilder) ProcessAxis(root *Axis) Query {
	var qyInput Query
	if root.input != nil {
		qyInput = b.ProcessNode(root.input)
	}
	var result Query

	switch root.axis_type {
	case AxisAncestor:
		result = &AncestorQuery{qyInput: qyInput, name: root.name, prefix: root.prefix, typeTest: root.node_type}
	case AxisAncestorOrSelf:
		result = &AncestorQuery{qyInput: qyInput, name: root.name, prefix: root.prefix, typeTest: root.node_type, matchSelf: true}
	case AxisAttribute:
		result = &AttributeQuery{qyInput: qyInput, name: root.name, prefix: root.prefix, typeTest: root.node_type}
	case AxisChild:
		result = &ChildQuery{qyInput: qyInput, name: root.name, prefix: root.prefix, typeTest: root.node_type}
	case AxisDescendant:
		result = &DescendantQuery{qyInput: qyInput, name: root.name, prefix: root.prefix, typeTest: root.node_type}
	case AxisDescendantOrSelf:
		result = &DescendantQuery{qyInput: qyInput, name: root.name, prefix: root.prefix, typeTest: root.node_type, matchSelf: true}
	default:
		panic(fmt.Sprintf("Axis type[%d] not implemented.", root.axis_type))
	}
	return result
}

func (b *QueryBuilder) ProcessOperator(root *Operator) Query {
	var op1 = b.ProcessNode(root.opnd1)
	var op2 = b.ProcessNode(root.opnd2)

	switch root.op {
	case OpPLUS, OpMINUS, OpMUL, OpMOD, OpDIV:
		return &NumericExpr{root.op, op1, op2}
	case OpLT, OpGT, OpLE, OpGE, OpEQ, OpNE:
		// LogicalExpr
	case OpOR, OpAND:
	// BooleanExpr
	case OpUNION:
		// UnionExpr
	}
}
