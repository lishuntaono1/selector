package internal

import "github.com/zhengchun/selector/xpath"

type NodeIterator interface {
	xpath.NodeIterator

	CurrentPosition() int
}

type XPathSelectionIterator struct {
	nav   xpath.Navigator
	query Query
}

func (x *XPathSelectionIterator) MoveNext() bool {
	n := x.query.Advance()
	if n != nil {
		if !x.nav.MoveTo(n) {
			x.nav = n.Clone()
		}
		return true
	}
	return false
}

func (x *XPathSelectionIterator) Current() xpath.Navigator {
	return x.nav
}

func (x *XPathSelectionIterator) CurrentPosition() int {
	return 0
}
