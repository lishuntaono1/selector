package internal

import "github.com/zhengchun/selector/xpath"

type ChildrenQuery struct {
	BaseAxisQuery

	iterator func() (xpath.Navigator, bool)
	first    bool
}

func (c *ChildrenQuery) Advance() xpath.Navigator {
	for {
		if c.iterator == nil {
			c.first = true
			c.position = 0
			nav := c.qyInput.Advance()
			if nav == nil {
				return nil
			}
			nav = nav.Clone()
			c.iterator = func() (xpath.Navigator, bool) {
				for {
					if c.first {
						if !nav.MoveToFirstChild() {
							return nil, false
						}
					} else {
						if !nav.MoveToNext() {
							return nil, false
						}
					}
					c.first = false
					if c.matches(nav) {
						return nav, true
					}
				}
			}
		}
		if nav, ok := c.iterator(); ok {
			c.position++
			c.currNode = nav
			return nav
		} else {
			c.iterator = nil
		}
	}
}

func (c *ChildrenQuery) Evaluate(iter NodeIterator) interface{} {
	c.qyInput.Evaluate(iter)
	return c
}

func (c *ChildrenQuery) Reset() {
	c.iterator = nil
	c.BaseAxisQuery.Reset()
}
