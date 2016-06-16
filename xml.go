package selector

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"strings"

	"github.com/zhengchun/selector/internal"
	"github.com/zhengchun/selector/xpath"
)

type XmlNodeType uint

const (
	XmlDocumentNode XmlNodeType = iota
	XmlDeclarationNode
	XmlElementNode
	XmlTextNode
	XmlCommentNode
	XmlProcInstNode
)

type XmlNode struct {
	Parent, FirstChild, LastChild, PrevSibling, NextSibling *XmlNode

	Type      XmlNodeType
	Data      string
	Namespace string
	Attr      []xml.Attr
	level     int
}

type xmlNodeNavigator struct {
	doc      *XmlNode
	currnode *XmlNode
	attindex int
}

func (n *xmlNodeNavigator) Select(xpath string) xpath.NodeIterator {
	return internal.Express(n.Clone(), xpath)
}

func (n *xmlNodeNavigator) BaseURI() string {
	return ""
}

func (n *xmlNodeNavigator) LocalName() string {
	if n.attindex != -1 && len(n.currnode.Attr) > 0 {
		return n.currnode.Attr[n.attindex].Name.Local
	} else {
		return n.currnode.Data
	}
}

func (n *xmlNodeNavigator) Value() string {
	switch n.currnode.Type {
	case XmlCommentNode:
		return n.currnode.Data
	case XmlDocumentNode:
		return ""
	case XmlElementNode:
		if n.attindex != -1 {
			return n.currnode.Attr[n.attindex].Value
		}
		return XmlNodeInnerText(n.currnode)
	case XmlTextNode:
		return n.currnode.Data
	default:
		panic(fmt.Sprintf("unknowed XmlNodeType: %v", n.currnode.Type))
	}
}

func (n *xmlNodeNavigator) Prefix() string {
	return ""
}

func (n *xmlNodeNavigator) NodeType() xpath.NodeType {
	switch n.currnode.Type {
	case XmlCommentNode:
		return xpath.CommentNode
	case XmlProcInstNode:
		return xpath.AllNode
	case XmlTextNode:
		return xpath.TextNode
	case XmlDeclarationNode, XmlDocumentNode:
		return xpath.RootNode
	case XmlElementNode:
		if n.attindex != -1 {
			return xpath.AttributeNode
		}
		return xpath.ElementNode
	default:
		panic(fmt.Sprintf("unknowed XmlNodeType: %v", n.currnode.Type))
	}
}

func (n *xmlNodeNavigator) Clone() xpath.Navigator {
	nav := *n
	return &nav
}

func (n *xmlNodeNavigator) IsSamePosition(other xpath.Navigator) bool {
	nav, ok := other.(*xmlNodeNavigator)
	if !ok {
		return false
	}
	return n.currnode == nav.currnode
}

func (n *xmlNodeNavigator) ComparePosition(other xpath.Navigator) xpath.XmlNodeOrder {
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

func (n *xmlNodeNavigator) MoveTo(other xpath.Navigator) bool {
	nav, ok := other.(*xmlNodeNavigator)
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

func (n *xmlNodeNavigator) MoveToRoot() {
	n.currnode = n.doc
}

func (n *xmlNodeNavigator) MoveToParent() bool {
	if n.currnode.Parent == nil {
		return false
	}
	n.currnode = n.currnode.Parent
	return true
}

func (n *xmlNodeNavigator) MoveToFirst() bool {
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

func (n *xmlNodeNavigator) MoveToNext() bool {
	if cur := n.currnode.NextSibling; cur == nil {
		return false
	} else {
		n.currnode = cur
	}
	return true
}

func (n *xmlNodeNavigator) MoveToFirstAttribute() bool {
	if len(n.currnode.Attr) == 0 {
		return false
	}
	n.attindex = 0
	return true
}

func (n *xmlNodeNavigator) MoveToNextAttribute() bool {
	if n.attindex >= len(n.currnode.Attr)-1 {
		return false
	}
	n.attindex++
	return true
}

func (n *xmlNodeNavigator) MoveToFirstChild() bool {
	if cur := n.currnode.FirstChild; cur == nil {
		return false
	} else {
		n.currnode = cur
	}
	return true
}

func SelectXmlNodes(n *XmlNode, xpath string) []*XmlNode {
	var nav = &xmlNodeNavigator{doc: n, currnode: n, attindex: -1}
	iter := nav.Select(xpath)
	nodes := make([]*XmlNode, 0)

	for iter.MoveNext() {
		cur := iter.Current().(*xmlNodeNavigator)
		nodes = append(nodes, cur.currnode)
	}
	return nodes
}

func SelectSingleXmlNode(n *XmlNode, xpath string) *XmlNode {
	var nav = &xmlNodeNavigator{doc: n, currnode: n, attindex: -1}
	iter := nav.Select(xpath)
	for iter.MoveNext() {
		cur := iter.Current().(*xmlNodeNavigator)
		return cur.currnode
	}
	return nil
}

func XmlNodeInnerText(n *XmlNode) string {
	var b bytes.Buffer
	var output func(*XmlNode)
	output = func(node *XmlNode) {
		if node.Type == XmlTextNode {
			b.WriteString(node.Data)
		}
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			output(child)
		}
	}
	output(n)
	return b.String()
}

func ParseXml(r io.Reader) (*XmlNode, error) {
	decoder := xml.NewDecoder(r)
	var nodes []*XmlNode
	nodes = append(nodes, &XmlNode{Type: XmlDocumentNode})
	var level int = 1
	for {
		t, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		switch t := t.(type) {
		case xml.EndElement:
			level--
		case xml.StartElement:
			node := &XmlNode{
				level:     level,
				Data:      t.Name.Local,
				Namespace: t.Name.Space,
				Attr:      t.Attr,
				Type:      XmlElementNode,
			}
			nodes = append(nodes, node)
			level++
		case xml.CharData:
			node := &XmlNode{
				level: level,
				Data:  string(t),
				Type:  XmlTextNode,
			}
			nodes = append(nodes, node)
		case xml.Comment:
			node := &XmlNode{
				level: level,
				Data:  string(t),
				Type:  XmlCommentNode,
			}
			nodes = append(nodes, node)
		case xml.ProcInst:
			node := &XmlNode{
				level: level,
				Data:  t.Target,
				Type:  XmlProcInstNode,
			}
			if node.Data == "xml" {
				node.Type = XmlDeclarationNode
			}
			for _, s := range strings.Split(string(t.Inst), " ") {
				s = strings.TrimSpace(s)
				if len(s) == 0 {
					continue
				}
				if i := strings.Index(s, "="); i > 0 {
					node.Attr = append(node.Attr, xml.Attr{
						Name:  xml.Name{Local: s[:i]},
						Value: strings.Trim(s[i+1:], "\""),
					})
				}
			}
			nodes = append(nodes, node)
		}
	}
	var prev *XmlNode
	var parent *XmlNode

	for i := 0; i < len(nodes); i++ {
		node := nodes[i]
		//fmt.Println(fmt.Sprintf("%d: %d - %s", i, node.level, node.Data))
		if prev == nil {
			parent = node
			prev = node
			continue
		}
		if node.level == prev.level { // sibling
			prev.NextSibling = node
			node.PrevSibling = prev
			node.Parent = parent
			if parent.level+1 == node.level {
				parent.LastChild = node
			}
		} else if node.level > prev.level { // move to children
			parent = prev
			parent.FirstChild = node
			node.Parent = parent
		} else if node.level < prev.level { // move to parent
			if node.level == parent.level {
				parent = parent.Parent
				node.PrevSibling = prev.Parent
				prev.Parent.NextSibling = node
			}
			parent.LastChild = node
			node.Parent = parent
		}
		prev = node
	}
	return nodes[0], nil
}
