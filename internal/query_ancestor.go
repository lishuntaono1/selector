package internal

import "github.com/zhengchun/selector/xpath"

type ancestorQuery struct {
	ParentQuery
	self bool
}

func (a *ancestorQuery) Evaluate(context NodeIterator) interface{} {
	a.Reset()
	a.qyInput.Evaluate(context)

	var ancestor, input xpath.Navigator
	for {
		input = a.qyInput.Advance()
		if input == nil {
			break
		}
		if a.self {
			if a.matches(input) {
				b, ok := insertNode(a.buff, input)
				if !ok {
					continue
				}
				a.buff = b
			}
		}
		if ancestor == nil || !ancestor.MoveTo(input) {
			ancestor = input.Clone()
		}
		for ancestor.MoveToParent() {
			if a.matches(ancestor) {
				b, ok := insertNode(a.buff, input)
				if !ok {
					break
				}
				a.buff = b
			}
		}
	}
	return a
}
