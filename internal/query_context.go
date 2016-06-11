package internal

import "github.com/zhengchun/selector/xpath"

type ContextQuery struct {
	count    int
	currNode xpath.Navigator
}

func (c *ContextQuery) Advance() xpath.Navigator {
	if c.count == 0 {
		c.count = 1
		return c.currNode
	}
	return nil
}

func (c *ContextQuery) Evaluate(iter NodeIterator) interface{} {
	c.currNode = iter.Current()
	c.count = 0
	return c
}

func (c *ContextQuery) Reset() {
	c.count = 0
}

func (c *ContextQuery) Current() xpath.Navigator {
	return c.currNode
}
