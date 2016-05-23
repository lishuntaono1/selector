package internal

import "golang.org/x/net/html"

func Build(query string) Query {
	root := ParseXPathExpression(query)
	return processNode(root, &funcQuery{
		selector:  func(n *html.Node) bool { return true },
		evaluator: func(n *html.Node) interface{} { return nil },
	})
}

func processNode(root AstNode, a Query) Query {
	var q Query
	switch root.Type() {
	case AstAxis:
		q = processAxis(root.(*Axis), a)
	}
	return q
}

func processAxis(root *Axis, a Query) Query {
	var q Query
	switch root.axis_type {
	case AxisAncestor:
		{
			q = &funcQuery{selector: func(n *html.Node) bool {
				for parent := n.Parent; parent != nil; parent = parent.Parent {
					if a.Select(parent) {
						return true
					}
				}
				return false
			}}
		}
	case AxisAncestorOrSelf:
		{
			q = &funcQuery{selector: func(n *html.Node) bool {
				if !a.Select(n) {
					return false
				}
				for parent := n.Parent; parent != nil; parent = parent.Parent {
					if a.Select(parent) {
						return true
					}
				}
				return false
			}}
		}
	case AxisAttribute:
		{
			q = &funcQuery{selector: func(n *html.Node) bool {
				if n.Type != html.ElementNode {
					return false
				}
				for _, attr := range n.Attr {
					if attr.Key == root.name && attr.Namespace == root.prefix {
						return true
					}
				}
				return false
			},evaluator:func(n *html.Node) interface{}}
		}
	}
	if root.input != nil {
		q = processNode(root.input, q)
	}
	return q
}
