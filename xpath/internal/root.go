package internal

type Root struct{}

func (r *Root) Type() AstNodeType {
	return RootAst
}

func (r *Root) ReturnType() ResultType {
	return NodeSetType
}
