package xpath

type Navigator interface {
	LocalName() string
	Value() string
	Prefix() string
	NodeType() NodeType

	MoveToRoot()
	MoveToParent() bool
	MoveToNext() bool
	MoveToFirstAttribute() bool
	MoveToNextAttribute() bool
	MoveToFirstChild() bool
}
