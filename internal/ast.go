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

type NodeType uint

const (
	NodeRoot NodeType = iota
	NodeElement
	NodeAttribute
	NodeNamespace
	NodeText
	NodeSignificantWhitespace
	NodeWhitespace
	NodeProcessingInstruction
	NodeComment
	NodeAll
)

// Filter ast
type Filter struct {
	input, condition AstNode
}

func (f *Filter) Type() AstType {
	return AstFilter
}

func (f *Filter) ReturnType() ResultType {
	return NodeSetType
}

// Root ast
type Root struct{}

func (r *Root) Type() AstType {
	return AstRoot
}

func (r *Root) ReturnType() ResultType {
	return NodeSetType
}

// Group ast
type Group struct {
	group AstNode
}

func (g *Group) Type() AstType {
	return AstGroup
}

func (g *Group) ReturnType() ResultType {
	return NodeSetType
}

// Operand ast
type Operand struct {
	valType ResultType
	val     interface{}
}

func (o *Operand) Type() AstType {
	return AstConstantOperand
}

func (o *Operand) ReturnType() ResultType {
	return o.valType
}

// Variable ast
type Variable struct {
	name, prefix string
}

func (v *Variable) Type() AstType {
	return AstVariable
}

func (v *Variable) ReturnType() ResultType {
	return AnyType
}
