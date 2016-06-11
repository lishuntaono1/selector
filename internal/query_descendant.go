package internal

import "github.com/zhengchun/selector/xpath"

type DescendantQuery struct {
	BaseAxisQuery

	self     bool
	iterator func() (xpath.Navigator, bool)
}

func (d *DescendantQuery) Advance() xpath.Navigator {
	for {
		if d.iterator == nil {
			var first bool = true
			var level int

			nav := d.qyInput.Advance()
			if nav == nil {
				return nil
			}
			nav = nav.Clone()
			d.iterator = func() (xpath.Navigator, bool) {
				if first {
					first = false
					if d.self && d.matches(nav) {
						d.position = 1
						return nav, true
					}
				}
				for {
					if nav.MoveToFirstChild() {
						level++
					} else {
						for {
							if level == 0 {
								return nil, false
							}
							if nav.MoveToNext() {
								break
							}
							nav.MoveToParent()
							level--
						}
					}
					if d.matches(nav) {
						return nav, true
					}
				}
			}
		}
		if nav, ok := d.iterator(); ok {
			d.position++
			d.currNode = nav
			return nav
		} else {
			d.iterator = nil
		}
	}
}

func (d *DescendantQuery) Reset() {
	d.iterator = nil
	d.BaseAxisQuery.Reset()
}

func (d *DescendantQuery) moveNext() bool {
	return d.Advance() != nil
}

func (d *DescendantQuery) Count() int {
	clone := *d
	clone.Reset()
	var count int
	for clone.moveNext() {
		count++
	}
	return count
}
