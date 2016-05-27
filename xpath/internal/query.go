package internal

import "github.com/zhengchun/selector/xpath"

type Query interface {
	MatchNode(xpath.Navigator) xpath.Navigator
	/*
		MoveNext() bool
		Advance() xpath.Navigator
		Evaluate(xpath.NodeIterator) interface{}
		MatchNode(xpath.Navigator) xpath.Navigator
		Current() xpath.Navigator
	*/
}
