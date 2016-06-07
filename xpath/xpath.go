package xpath

type NodeType uint

type XmlNodeOrder uint

const (
	XmlNodeOrderBefore XmlNodeOrder = iota
	XmlNodeOrderAfter
	XmlNodeOrderSame
	XmlNodeOrderUnknown
)

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
	BaseURI() string
	LocalName() string
	Value() string
	Prefix() string
	NodeType() NodeType
	Clone() Navigator
	IsSamePosition(Navigator) bool
	ComparePosition(Navigator) XmlNodeOrder

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
