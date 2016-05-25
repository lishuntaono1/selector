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
