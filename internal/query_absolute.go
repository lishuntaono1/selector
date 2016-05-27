package internal

import "github.com/zhengchun/selector/xpath"

type absoluteQuery struct {
	name, prefix string
}

func (aq *absoluteQuery) Advance(ctx xpath.Context) xpath.Navigator {
	return nil
}
