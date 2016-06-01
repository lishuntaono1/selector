package html

import (
	"fmt"

	"github.com/zhengchun/selector/internal"
	"github.com/zhengchun/selector/xpath"

	"golang.org/x/net/html"
)

type NodeNavigator struct {
	doc      *html.Node
	currnode *html.Node
	attindex int
}

func (n *NodeNavigator) LocalName() string {
	if n.attindex != -1 {
		return n.currnode.Attr[n.attindex].Key
	} else {
		return n.currnode.Data
	}
}

func (n *NodeNavigator) Value() string {
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

func (n *NodeNavigator) Prefix() string {
	return ""
}

func (n *NodeNavigator) Select(xpath string) xpath.NodeIterator {
	builder := &internal.QueryBuilder{}
	var q = builder.Build(xpath)
	iter := &NodeIterator{cur: n, query: q}
	//q.Evaluate(iter)
	return iter
}

func (n *NodeNavigator) Clone() xpath.Navigator {
	nav := *n
	return &nav
}

func (n *NodeNavigator) MoveTo(other xpath.Navigator) bool {
	nav, ok := other.(*NodeNavigator)
	if !ok {
		return false
	}
	if nav.doc == n.doc {
		n.currnode = nav.currnode
		n.attindex = nav.attindex
		return true
	}
	return false
}

func (n *NodeNavigator) MoveToRoot() {
	n.currnode = n.doc
}

func (n *NodeNavigator) MoveToParent() bool {
	if n.currnode.Parent == nil {
		return false
	}
	n.currnode = n.currnode.Parent
	return true
}

func (n *NodeNavigator) MoveToFirst() bool {
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

func (n *NodeNavigator) MoveToNext() bool {
	if cur := n.currnode.NextSibling; cur == nil {
		return false
	} else {
		n.currnode = cur
	}
	return true
}

func (n *NodeNavigator) MoveToFirstAttribute() bool {
	if len(n.currnode.Attr) == 0 {
		return false
	}
	n.attindex = 0
	return true
}

func (n *NodeNavigator) MoveToNextAttribute() bool {
	if n.attindex >= len(n.currnode.Attr)-1 {
		return false
	}
	n.attindex++
	return true
}

func (n *NodeNavigator) MoveToFirstChild() bool {
	if cur := n.currnode.FirstChild; cur == nil {
		return false
	} else {
		n.currnode = cur
	}
	return true
}

func (n *NodeNavigator) NodeType() xpath.NodeType {
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
