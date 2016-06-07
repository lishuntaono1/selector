package internal

import "github.com/zhengchun/selector/xpath"

type descendantQuery struct {
	qyInput   Query
	matchSelf bool
	matches   func(xpath.Navigator) bool
	position  int

	movenext func() bool
	currnode xpath.Navigator
}

func (d *descendantQuery) Advance() xpath.Navigator {
	for {
		if d.movenext == nil {
			var first bool = true
			var level int

			if nav := d.qyInput.Advance(); nav == nil {
				return nil
			} else {
				d.currnode = nav.Clone()
			}
			d.movenext = func() bool {
				if first {
					first = false
					if d.matchSelf && d.matches(d.currnode) {
						d.position = 1
						return true
					}
				}
				for {
					if d.currnode.MoveToFirstChild() {
						level++
					} else {
						for {
							if level == 0 {
								return false
							}
							if d.currnode.MoveToNext() {
								break
							}
							d.currnode.MoveToParent()
							level--
						}
					}
					if d.matches(d.currnode) {
						return true
					}
				}
			}
		}
		if d.movenext() {
			d.position++
			return d.currnode
		} else {
			d.movenext = nil
		}
	}
}

func (d *descendantQuery) Evaluate(ctx NodeIterator) interface{} {
	d.qyInput.Evaluate(ctx)
	return d
}

func (d *descendantQuery) Current() xpath.Navigator {
	return d.currnode
}

func (d *descendantQuery) MoveNext() bool {
	return d.Advance() != nil
}

func (d *descendantQuery) CurrentPosition() int {
	return d.position
}

func (d *descendantQuery) Reset() {
	d.position = 0
	d.currnode = nil
	d.movenext = nil
	d.qyInput.Reset()
}

func (d *descendantQuery) Count() int {
	clone := *d
	clone.Reset()
	var count int
	for !clone.MoveNext() {
		count++
	}
	return count
}
