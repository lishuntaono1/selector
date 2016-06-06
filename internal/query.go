package internal

import "github.com/zhengchun/selector/xpath"

// XPath query interface.
//
type Query interface {
	NodeIterator

	Advance() xpath.Navigator
	Evaluate(NodeIterator) interface{}
}
