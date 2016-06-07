package internal

import "github.com/zhengchun/selector/xpath"

type NodeIterator interface {
	Current() xpath.Navigator
	Reset()
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

func (x *XPathSelectionIterator) Count() int {
	return 0
}

func (x *XPathSelectionIterator) Reset() {
}

func (x *XPathSelectionIterator) CurrentPosition() int {
	return 0
}

func (x *XPathSelectionIterator) Current() xpath.Navigator {
	return x.nav
}
