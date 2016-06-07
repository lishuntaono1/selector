package internal

import "github.com/zhengchun/selector/xpath"

// XPath query interface.
//
type Query interface {
	NodeIterator

	Advance() xpath.Navigator
	Evaluate(NodeIterator) interface{}
}

type QueryProps uint

const (
	NoneQueryProp QueryProps = iota << 1
	PositionQueryProp
	CountQuery
	CachedQuery
	ReverseQueryProp
	MergeQueryProp
)
