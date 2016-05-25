package internal

import "github.com/zhengchun/selector/xpath"

type Query interface {
	Reset()
	MoveNext() bool
	Advance() xpath.Navigator
	Evaluate(xpath.NodeIterator) interface{}
	MatchNode(xpath.Navigator)
}
