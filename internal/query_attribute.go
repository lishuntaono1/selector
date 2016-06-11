package internal

import "github.com/zhengchun/selector/xpath"

type AttributeQuery struct {
	BaseAxisQuery

	onAttr bool
}

func (a *AttributeQuery) Advance() xpath.Navigator {
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

func (a *AttributeQuery) Evaluate(context NodeIterator) interface{} {
	a.Reset()
	a.qyInput.Evaluate(context)
	return a
}

func (a *AttributeQuery) Reset() {
	a.onAttr = false
	a.BaseAxisQuery.Reset()
}
