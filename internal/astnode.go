package internal

type AstNode interface {
	Type() AstNodeType
	ReturnType() ResultType
}

type AstNodeType uint

const (
	AxisAst AstNodeType = iota
	OperatorAst
	FilterAst
	ConstantOperandAst
	FunctionAst
	GroupAst
	RootAst
	VariableAst
)

type ResultType uint

const (
	NumberType ResultType = iota
	StringType
	BooleanType
	NodeSetType
	AnyType
)
