package html

import (
	"github.com/zhengchun/selector/internal"
	"github.com/zhengchun/selector/xpath"
)

type NodeIterator struct {
	cur   xpath.Navigator
	query internal.Query
}

func (iter *NodeIterator) MoveNext() bool {
	if iter.query.Matches(iter.cur) {
		return true
	}
	return false
}

func (iter *NodeIterator) Current() xpath.Navigator {
	return iter.cur
}
