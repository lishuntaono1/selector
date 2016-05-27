package internal

import "github.com/zhengchun/selector/xpath"

type Query interface {
	Advance() xpath.Navigator
	Evaluate(xpath.NodeIterator) Query
}
