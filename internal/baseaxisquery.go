package internal

import "github.com/zhengchun/selector/xpath"

// BaseAxisQuery is the base query for the XPath axes.
//
type BaseAxisQuery struct {
	qyInput  Query
	currNode xpath.Navigator
	position int
	matches  func(xpath.Navigator) bool
}

func (a *BaseAxisQuery) Reset() {
	a.position = 0
	a.currNode = nil
	a.qyInput.Reset()
}

func (a *BaseAxisQuery) Evaluate(context NodeIterator) interface{} {
	a.qyInput.Evaluate(context)
	return a
}

func (a *BaseAxisQuery) Current() xpath.Navigator {
	return a.currNode
}

func (a *BaseAxisQuery) CurrentPosition() int {
	return a.position
}
