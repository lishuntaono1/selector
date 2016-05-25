package internal

import "github.com/zhengchun/selector/xpath"

type AxisQuery struct {
	position    int
	currentNode xpath.Navigator

	qyInput      Query
	nodetype     xpath.NodeType
	name, prefix string
}

func (q *AxisQuery) MoveNext() bool {
	return false
}
