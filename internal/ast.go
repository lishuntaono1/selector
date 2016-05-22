package internal

type AstNode interface {
	Type() AstType
	ReturnType() ResultType
}

type AstType uint

const (
	AstAxis AstType = iota
	AstOperator
	AstFilter
	AstConstantOperand
	AstFunction
	AstGroup
	AstRoot
	AstVariable
)

type ResultType uint

const (
	NumberType ResultType = iota
	StringType
	BooleanType
	NodeSetType
	AnyType
)
