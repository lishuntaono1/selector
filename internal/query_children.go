package internal

import "github.com/zhengchun/selector/xpath"

type childrenQuery struct {
	qyInput  Query
	position int
	movenext func() (xpath.Navigator, bool)
	matches  func(xpath.Navigator) bool
	currnode xpath.Navigator

	firstChild bool
}

func (c *childrenQuery) Advance() xpath.Navigator {
	/*
		for {
			nav := c.qyInput.Advance()
			if nav == nil {
				return nil
			}
			var firstChild = false
			for {
				if !firstChild {
					firstChild = true
					if !nav.MoveToFirstChild() {
						break
					}
				} else {
					if !nav.MoveToNext() {
						break
					}
				}
				if c.matches(nav) {
					c.position++
					return nav
				}
			}
			c.position = 0
		}

	*/

	for {
		if c.movenext == nil {
			c.position = 0
			c.firstChild = false
			nav := c.qyInput.Advance()
			if nav == nil {
				return nil
			}
			//nav = nav.Clone()
			c.movenext = func() (xpath.Navigator, bool) {
				for {
					if !c.firstChild {
						c.firstChild = true
						if !nav.MoveToFirstChild() {
							return nil, false
						}
					} else {
						if !nav.MoveToNext() {
							return nil, false
						}
					}
					if c.matches(nav) {
						c.position++
						c.currnode = nav
						return c.currnode, true
					}
				}
			}
		}
		if nav, ok := c.movenext(); ok {
			c.position++
			c.currnode = nav
			return nav
		} else {
			c.movenext = nil
		}
	}
}

func (c *childrenQuery) Evaluate(iter NodeIterator) interface{} {
	c.qyInput.Evaluate(iter)
	return c
}

func (c *childrenQuery) MoveNext() bool {
	return c.Advance() != nil
}

func (c *childrenQuery) Current() xpath.Navigator {
	return c.currnode
}

func (c *childrenQuery) CurrentPosition() int {
	return c.position
}

func (c *childrenQuery) Reset() {
	c.currnode = nil
	c.movenext = nil
	c.firstChild = false
	c.qyInput.Reset()
}

func (c *childrenQuery) Count() int {
	clone := *c
	clone.Reset()
	var count int
	for clone.MoveNext() {
		count++
	}
	return count
}
