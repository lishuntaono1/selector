package internal

import "github.com/zhengchun/selector/xpath"

// XPath query interface.
//
type Query interface {
	Matches(xpath.Navigator) bool
}

type selector func(xpath.Navigator) bool

func (s selector) Matches(nav xpath.Navigator) bool {
	return s(nav)
}

type Query2 interface {
	Advance() xpath.Navigator
	Evaluate() interface{}
	Value(xpath.Navigator) interface{}
}
