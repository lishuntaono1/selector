package internal

import "github.com/zhengchun/selector/xpath"

type attributeQuery struct {
	qyInput  Query
	position int
	matches  func(xpath.Navigator) bool

	onAttr   bool
	currNode xpath.Navigator
}

func (a *attributeQuery) Advance() xpath.Navigator {

	for {
		if !a.onAttr {
			nav := a.qyInput.Advance()
			if nav == nil {
				return nil
			}
			a.position = 0
			a.currNode = nav.Clone()
			a.onAttr = a.currNode.MoveToFirstAttribute()
		} else {
			a.onAttr = a.currNode.MoveToNextAttribute()
		}
		if a.onAttr {
			if a.matches(a.currNode) {
				a.position++
				return a.currNode
			}
		}
	}
}

func (a *attributeQuery) Evaluate(iter NodeIterator) interface{} {
	a.qyInput.Evaluate(iter)
	return a
}

func (a *attributeQuery) Current() xpath.Navigator {
	return a.currNode
}

func (a *attributeQuery) MoveNext() bool {
	return a.Advance() != nil
}

func (a *attributeQuery) Reset() {
	a.currNode = nil
	a.position = 0
	a.qyInput.Reset()
}

func (a *attributeQuery) Count() int {
	clone := *a
	clone.Reset()
	var count int
	for clone.MoveNext() {
		count++
	}
	return count
}

func (a *attributeQuery) CurrentPosition() int {
	return a.position
}
