package internal

import "github.com/zhengchun/selector/xpath"

type childrenQuery struct {
	qyInput  Query
	position int
	movenext func() bool
	matches  func(xpath.Navigator) bool
	currnode xpath.Navigator
}

func (c *childrenQuery) Advance() xpath.Navigator {
	for {
		if c.movenext == nil {
			nav := c.qyInput.Advance()

			if nav == nil {
				return nil
			}
			c.currnode = nav.Clone()
			c.movenext = func() bool {
				for {
					if c.position == 0 && !c.currnode.MoveToFirstChild() {
						return false
					} else if c.position > 0 && !c.currnode.MoveToNext() {
						c.currnode.MoveToParent()
						return false
					}
					c.position++
					if c.matches(c.currnode) {
						return true
					}
				}
			}
		}
		if c.movenext() {
			return c.currnode
		} else {
			c.position = 0
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
	c.position = 0
	c.movenext = nil
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
