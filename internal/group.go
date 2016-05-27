package internal

type Group struct {
	group AstNode
}

func (g *Group) Type() AstNodeType {
	return GroupAst
}

func (g *Group) ReturnType() ResultType {
	return NodeSetType
}
