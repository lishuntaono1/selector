package internal

import "github.com/zhengchun/selector/xpath"

type AbsoluteQuery struct {
	*ContextQuery
}

func (q *AbsoluteQuery) Evaluate(context xpath.NodeIterator) interface{} {

	return q
}

func (q *AbsoluteQuery) MatchNode(context xpath.Navigator) xpath.Navigator {
	if context != nil && context.NodeType() == xpath.RootNode {
		return context
	}
	return nil
}
