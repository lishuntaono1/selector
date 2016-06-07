package internal

import "github.com/zhengchun/selector/xpath"

// ancestorQuery is the axis query with cache.
type ancestorQuery struct {
	parentQuery
	matchSelf bool
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
		if a.matchSelf {
			if a.matches(input) {
				b, ok := insertNode(a.outputBuffer, input)
				if !ok {
					continue
				}
				a.outputBuffer = b
			}
		}
		if ancestor == nil || !ancestor.MoveTo(input) {
			ancestor = input.Clone()
		}
		for ancestor.MoveToParent() {
			if a.matches(ancestor) {
				b, ok := insertNode(a.outputBuffer, input)
				if !ok {
					break
				}
				a.outputBuffer = b
			}
		}
	}
	return a
}
