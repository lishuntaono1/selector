package internal

import "github.com/zhengchun/selector/xpath"

type EmptyQuery struct {
}

func (q *EmptyQuery) MatchNode(context xpath.Navigator) xpath.Navigator {
	panic("MatchNode method NOT implemented.")
}
