package internal

import "github.com/zhengchun/selector/xpath"

type descendantQuery struct {
	input Query
}

func (dq *descendantQuery) Advance() xpath.Navigator {
	return nil
}

func (dq *descendantQuery) Evaluate(ctx xpath.NodeIterator) Query {
	return dq.input.Evaluate(ctx)
}

type descendantIterator struct {
	currNode xpath.Navigator
}

func (iter *descendantIterator) MoveNext() bool {
	return false
}
