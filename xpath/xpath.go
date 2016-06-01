package xpath

type NodeType uint

const (
	RootNode NodeType = iota
	ElementNode
	AttributeNode
	NamespaceNode
	TextNode
	SignificantWhitespaceNode
	WhitespaceNode
	ProcessingInstructionNode
	CommentNode
	AllNode
)

type Navigator interface {
	LocalName() string
	Value() string
	Prefix() string
	NodeType() NodeType
	Clone() Navigator

	MoveTo(Navigator) bool
	MoveToRoot()
	MoveToParent() bool
	MoveToNext() bool
	MoveToFirstAttribute() bool
	MoveToNextAttribute() bool
	MoveToFirstChild() bool
}

type NodeIterator interface {
	Current() Navigator
	MoveNext() bool
}
