package xpath

type Navigator interface {
	NodeType() NodeType
	LocalName() string

	MoveToRoot()
	Clone() Navigator
}
