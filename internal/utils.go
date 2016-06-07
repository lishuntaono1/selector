package internal

import (
	"strings"

	"github.com/zhengchun/selector/xpath"
)

func iteratorCount(iter NodeIterator) int {
	type Counter interface {
		Count() int
	}
	count, ok := iter.(Counter)
	if !ok {
		return 0
	}
	return count.Count()
}

func iteratorPosition(iter NodeIterator) int {
	type Position interface {
		CurrentPosition() int
	}
	position, ok := iter.(Position)
	if !ok {
		return 0
	}
	return position.CurrentPosition()
}

func insertNode(b []xpath.Navigator, n xpath.Navigator) ([]xpath.Navigator, bool) {
	var l int
	r := len(b)
	if r != 0 {
		switch compareNodes(b[r-1], n) {
		case xpath.XmlNodeOrderSame:
			return nil, false
		case xpath.XmlNodeOrderBefore:
			b = append(b, n.Clone())
			return b, true
		default:
			r--
		}
	}
	for l < r {
		m := (l + r) >> 1
		switch compareNodes(b[m], n) {
		case xpath.XmlNodeOrderSame:
			return nil, false
		case xpath.XmlNodeOrderBefore:
			l = m + 1
		default:
			r = m
		}
	}
	b = append(b, nil)
	copy(b[l+1:], b[l:])
	b[l] = n.Clone()
	return b, true
}

func compareNodes(l, r xpath.Navigator) xpath.XmlNodeOrder {
	cmp := l.ComparePosition(r)
	if cmp == xpath.XmlNodeOrderUnknown {
		copy := l.Clone()
		copy.MoveToRoot()
		baseUriL := copy.BaseURI()
		if !copy.MoveTo(r) {
			copy = r.Clone()
		}
		copy.MoveToRoot()
		baseUriR := copy.BaseURI()
		cmpBase := strings.Compare(baseUriL, baseUriR)

		if cmpBase < 0 {
			cmp = xpath.XmlNodeOrderBefore
		} else if cmpBase > 0 {
			cmp = xpath.XmlNodeOrderAfter
		}
		cmp = xpath.XmlNodeOrderUnknown
	}
	return cmp
}
