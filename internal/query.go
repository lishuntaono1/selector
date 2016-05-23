package internal

import "golang.org/x/net/html"

type Query interface {
	Select(*html.Node) bool
	Evaluate(*html.Node) interface{}
}

type funcQuery struct {
	selector  func(*html.Node) bool
	evaluator func(*html.Node) interface{}
}

func (q *funcQuery) Select(n *html.Node) bool {
	if q.selector != nil {
		return q.selector(n)
	}
	return false
}

func (q *funcQuery) Evaluate(n *html.Node) interface{} {
	if q.evaluator != nil {
		return q.evaluator(n)
	}
	return nil
}

type AncestorQuery struct {
	self bool
}
