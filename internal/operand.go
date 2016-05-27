package internal

// Operand ast
type Operand struct {
	valtype ResultType
	val     interface{}
}

func (o *Operand) Type() AstNodeType {
	return ConstantOperandAst
}

func (o *Operand) ReturnType() ResultType {
	return o.valtype
}
