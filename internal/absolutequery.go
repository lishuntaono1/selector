package internal

type AbsoluteQuery struct {
	ContextQuery
}

func (a *AbsoluteQuery) Evaluate(context NodeIterator) interface{} {
	a.currNode = context.Current().Clone()
	a.currNode.MoveToRoot()
	a.count = 0
	return a
}
