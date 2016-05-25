package internal

type Filter struct {
	input, condition AstNode
}

func (f *Filter) Type() AstNodeType {
	return FilterAst
}

func (f *Filter) ReturnType() ResultType {
	return NodeSetType
}
