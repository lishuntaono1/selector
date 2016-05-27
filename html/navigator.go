package html

import (
	"fmt"

	"github.com/zhengchun/selector/xpath"
	"golang.org/x/net/html"
)

type nodeNavigator struct {
	doc      *html.Node
	currnode *html.Node
	attindex int
}

func (n *nodeNavigator) LocalName() string {
	return n.currnode.Data
}

func (n *nodeNavigator) Value() string {
	switch n.currnode.Type {
	case html.CommentNode:
		return n.currnode.Data
	case html.DoctypeNode, html.DocumentNode:
		return ""
	case html.ElementNode:
		if n.attindex != -1 {
			return n.currnode.Attr[n.attindex].Val
		}
		return InnerText(n.currnode)
	case html.TextNode:
		return n.currnode.Data
	default:
		panic(fmt.Sprintf("unknowed HtmlNodeType: %v", n.currnode.Type))
	}
}

func (n *nodeNavigator) Prefix() string {
	return ""
}

func (n *nodeNavigator) MoveToRoot() {
	n.currnode = n.doc
}

func (n *nodeNavigator) MoveToFirst() bool {
	if n.currnode.Parent == nil {
		return false
	}
	if cur := n.currnode.Parent.NextSibling; cur == nil {
		return false
	} else {
		n.currnode = cur
	}
	return true
}

func (n *nodeNavigator) MoveToNext() bool {
	if cur := n.currnode.NextSibling; cur == nil {
		return false
	} else {
		n.currnode = cur
	}
	return true
}

func (n *nodeNavigator) MoveToFirstAttribute() bool {
	if len(n.currnode.Attr) == 0 {
		return false
	}
	n.attindex = 0
	return true
}

func (n *nodeNavigator) MoveToNextAttribute() bool {
	if n.attindex >= len(n.currnode.Attr)-1 {
		return false
	}
	n.attindex++
	return true
}

func (n *nodeNavigator) MoveToFirstChild() bool {
	if cur := n.currnode.FirstChild; cur == nil {
		return false
	} else {
		n.currnode = cur
	}
	return true
}

func (n *nodeNavigator) NodeType() xpath.NodeType {
	switch n.currnode.Type {
	case html.CommentNode:
		return xpath.CommentNode
	case html.DoctypeNode:
		return xpath.AllNode
	case html.TextNode:
		return xpath.TextNode
	case html.DocumentNode:
		return xpath.RootNode
	case html.ElementNode:
		if n.attindex != -1 {
			return xpath.AttributeNode
		}
		return xpath.ElementNode
	default:
		panic(fmt.Sprintf("unknowed HtmlNodeType: %v", n.currnode.Type))
	}
}
