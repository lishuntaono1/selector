package internal

import "golang.org/x/net/html"

type Query interface {
	Select(*html.Node) bool
	Evaluate(*html.Node) interface{}
}

// Ancestor Query
type AncestorQuery struct {
	qyInput      Query
	typeTest     NodeType
	matchSelf    bool
	name, prefix string
}

func (q *AncestorQuery) Select(n *html.Node) bool {
	if q.matchSelf && !q.matchNode(n) {
		return false
	}
	for parent := n.Parent; parent != nil; parent = parent.Parent {
		if q.matchNode(parent) {
			return true
		}
	}
	return false
}

func (q *AncestorQuery) Evaluate(n *html.Node) interface{} {
	return nil
}

func (q *AncestorQuery) matchNode(n *html.Node) bool {
	if n.Type == html.ElementNode && (q.typeTest == NodeAll || q.typeTest == NodeElement) {
		if len(q.name) > 0 {
			if q.name != n.Data {
				return false
			}
		}
		if q.qyInput != nil {
			return q.qyInput.Select(n)
		}
		return true
	}
	return false
}

// Attribute Query
type AttributeQuery struct {
	qyInput      Query
	typeTest     NodeType
	name, prefix string
}

func (q *AttributeQuery) Select(n *html.Node) bool {
	if n.Type != html.ElementNode {
		return false
	}
	for _, attr := range n.Attr {
		if attr.Key == q.name {
			if q.qyInput != nil {
				return q.qyInput.Select(n)
			}
			return true
		}
	}
	return false
}

func (q *AttributeQuery) Evaluate(n *html.Node) interface{} {
	for _, attr := range n.Attr {
		if attr.Key == q.name {
			return attr.Val
		}
	}
	return nil
}

// Child Query
type ChildQuery struct {
	qyInput      Query
	typeTest     NodeType
	name, prefix string
}

func (q *ChildQuery) Select(n *html.Node) bool {
	if n.Type != html.ElementNode || n.Data != q.name {
		return false
	}
	if q.qyInput != nil {
		return q.qyInput.Select(n)
	}
	return false
}

func (q *ChildQuery) Evaluate(n *html.Node) interface{} {
	return nil
}

// Descendant Query
type DescendantQuery struct {
	qyInput      Query
	name, prefix string
	typeTest     NodeType
	matchSelf    bool
}

func (q *DescendantQuery) Select(n *html.Node) bool {
	if q.matchSelf && !q.matchNode(n) {
		return false
	}
	for child := n.FirstChild; child != nil; child = child.NextSibling {
		if q.matchNode(child) {
			return true
		}
	}
	return false
}

func (q *DescendantQuery) Evaluate(n *html.Node) interface{} {
	return nil
}

func (q *DescendantQuery) matchNode(n *html.Node) bool {
	if n.Type == html.ElementNode {
		if q.qyInput != nil {
			return q.qyInput.Select(n)
		}
		return true
	}
	return false
}

// NumericExpr
type NumericExpr struct {
	op           OpType
	opnd1, opnd2 Query
}
