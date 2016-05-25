package internal

type Variable struct {
	name, prefix string
}

func (v *Variable) Type() AstNodeType {
	return VariableAst
}

func (v *Variable) ReturnType() ResultType {
	return AnyType
}
