package internal

import "github.com/zhengchun/selector/xpath"

type ContextQuery struct {
	count int

	contextNode xpath.Navigator
}

func (q *ContextQuery) MatchNode(context xpath.Navigator) xpath.Navigator {
	return nil
}

func (q *ContextQuery) Clone() *ContextQuery {
	return &ContextQuery{
		count:       0,
		contextNode: q.contextNode,
	}
}

/*

func (q *ContextQuery) Reset() {
	q.count = 0
}

func (q *ContextQuery) Evaluate(context xpath.NodeIterator) interface{} {
	q.count = 0
	q.contextNode = context.Current()
	return q
}

func (q *ContextQuery) Advance() xpath.Navigator {
	if q.count == 0 {
		q.count = 1
		return q.contextNode
	}
	return nil
}

func (q *ContextQuery) MatchNode(current xpath.Navigator) xpath.Navigator {
	return q.contextNode
}



func (q *ContextQuery) Current() xpath.Navigator {
	return q.contextNode
}
*/
