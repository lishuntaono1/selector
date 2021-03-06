package selector

import (
	"bytes"
	"fmt"
	"io"

	"github.com/zhengchun/selector/internal"
	"github.com/zhengchun/selector/xpath"

	"golang.org/x/net/html"
)

type htmlNodeNavigator struct {
	doc      *html.Node
	currnode *html.Node
	attindex int
}

func getDepth(nav xpath.Navigator) int {
	depth := 0
	for nav.MoveToParent() {
		depth++
	}
	return depth
}

func compareSiblings(n1, n2 xpath.Navigator) xpath.XmlNodeOrder {
	var cmp = 0
	switch n1.NodeType() {
	case xpath.NamespaceNode:
		break
	case xpath.AttributeNode:
		cmp += 1
	default:
		cmp += 2
	}
	switch n2.NodeType() {
	case xpath.NamespaceNode:
		if cmp == 0 {
			//not supported NamespaceNode.

			/*for n1.MoveToNextNamespace() {
				if n1.IsSamePosition(n2) {
					return XmlNodeOrderBefore
				}
			}*/
		}
	case xpath.AttributeNode:
		cmp -= 1
		if cmp == 0 {
			for n1.MoveToNextAttribute() {
				if n1.IsSamePosition(n2) {
					return xpath.XmlNodeOrderBefore
				}
			}
		}
	default:
		cmp -= 2
		if cmp == 0 {
			for n1.MoveToNext() {
				if n1.IsSamePosition(n2) {
					return xpath.XmlNodeOrderBefore
				}
			}
		}
		break
	}
	if cmp < 0 {
		return xpath.XmlNodeOrderBefore
	} else {
		return xpath.XmlNodeOrderAfter
	}
}

func (n *htmlNodeNavigator) BaseURI() string {
	return ""
}

func (n *htmlNodeNavigator) LocalName() string {
	if n.attindex != -1 && len(n.currnode.Attr) > 0 {
		return n.currnode.Attr[n.attindex].Key
	} else {
		return n.currnode.Data
	}
}

func (n *htmlNodeNavigator) Value() string {
	switch n.currnode.Type {
	case html.CommentNode:
		return n.currnode.Data
	case html.DoctypeNode, html.DocumentNode:
		return ""
	case html.ElementNode:
		if n.attindex != -1 {
			return n.currnode.Attr[n.attindex].Val
		}
		return HtmlNodeInnerText(n.currnode)
	case html.TextNode:
		return n.currnode.Data
	default:
		panic(fmt.Sprintf("unknowed HtmlNodeType: %v", n.currnode.Type))
	}
}

func (n *htmlNodeNavigator) Prefix() string {
	return ""
}

func (n *htmlNodeNavigator) Select(xpath string) xpath.NodeIterator {
	return internal.Express(n.Clone(), xpath)
}

func (n *htmlNodeNavigator) Clone() xpath.Navigator {
	nav := *n
	return &nav
}

func (n *htmlNodeNavigator) ComparePosition(other xpath.Navigator) xpath.XmlNodeOrder {
	if other == nil {
		return xpath.XmlNodeOrderUnknown
	}
	if n.IsSamePosition(other) {
		return xpath.XmlNodeOrderSame
	}
	n1 := n.Clone()
	n2 := other.Clone()
	depth1 := getDepth(n1.Clone())
	depth2 := getDepth(n2.Clone())

	if depth1 > depth2 {
		for depth1 > depth2 {
			n1.MoveToParent()
			depth1--
		}
		if n1.IsSamePosition(n2) {
			return xpath.XmlNodeOrderAfter
		}
	}

	if depth2 > depth1 {
		for depth2 > depth1 {
			n2.MoveToParent()
			depth2--
		}
		if n1.IsSamePosition(n2) {
			return xpath.XmlNodeOrderBefore
		}
	}
	parent1 := n1.Clone()
	parent2 := n2.Clone()

	for {
		if !parent1.MoveToParent() || !parent2.MoveToParent() {
			return xpath.XmlNodeOrderUnknown
		}
		if parent1.IsSamePosition(parent2) {
			return compareSiblings(n1, n2)
		}

		n1.MoveToParent()
		n2.MoveToParent()
	}
}

func (n *htmlNodeNavigator) MoveTo(other xpath.Navigator) bool {
	nav, ok := other.(*htmlNodeNavigator)
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

func (n *htmlNodeNavigator) MoveToRoot() {
	n.currnode = n.doc
}

func (n *htmlNodeNavigator) MoveToParent() bool {
	if n.currnode.Parent == nil {
		return false
	}
	n.currnode = n.currnode.Parent
	return true
}

func (n *htmlNodeNavigator) MoveToFirst() bool {
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

func (n *htmlNodeNavigator) MoveToNext() bool {
	if cur := n.currnode.NextSibling; cur == nil {
		return false
	} else {
		n.currnode = cur
	}
	return true
}

func (n *htmlNodeNavigator) MoveToFirstAttribute() bool {
	if len(n.currnode.Attr) == 0 {
		return false
	}
	n.attindex = 0
	return true
}

func (n *htmlNodeNavigator) MoveToNextAttribute() bool {
	if n.attindex >= len(n.currnode.Attr)-1 {
		return false
	}
	n.attindex++
	return true
}

func (n *htmlNodeNavigator) MoveToFirstChild() bool {
	if cur := n.currnode.FirstChild; cur == nil {
		return false
	} else {
		n.currnode = cur
	}
	return true
}

func (n *htmlNodeNavigator) NodeType() xpath.NodeType {
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

func (n *htmlNodeNavigator) IsSamePosition(other xpath.Navigator) bool {
	nav, ok := other.(*htmlNodeNavigator)
	if !ok {
		return false
	}
	return n.currnode == nav.currnode
}

func HtmlNodeInnerText(n *html.Node) string {
	var b bytes.Buffer
	var output func(*html.Node)
	output = func(node *html.Node) {
		if node.Type == html.TextNode {
			b.WriteString(node.Data)
		}
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			output(child)
		}
	}
	output(n)
	return b.String()
}

func SelectHtmlNodes(n *html.Node, xpath string) []*html.Node {
	var nav = &htmlNodeNavigator{doc: n, currnode: n, attindex: -1}
	iter := nav.Select(xpath)
	nodes := make([]*html.Node, 0)

	for iter.MoveNext() {
		cur := iter.Current().(*htmlNodeNavigator)
		nodes = append(nodes, cur.currnode)
	}
	return nodes
}

func SelectSingleHtmlNode(n *html.Node, xpath string) *html.Node {
	var nav = &htmlNodeNavigator{doc: n, currnode: n, attindex: -1}
	iter := nav.Select(xpath)
	for iter.MoveNext() {
		cur := iter.Current().(*htmlNodeNavigator)
		return cur.currnode
	}
	return nil
}

func ParseHtml(r io.Reader) (*html.Node, error) {
	doc, err := html.Parse(r)
	if err != nil {
		return nil, err
	}
	return doc, nil
}
