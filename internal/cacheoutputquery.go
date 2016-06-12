package internal

import "github.com/zhengchun/selector/xpath"

type CacheOutputQuery struct {
	qyInput Query
	buff    []xpath.Navigator
	count   int
}

func (c *CacheOutputQuery) Advance() xpath.Navigator {
	if c.count < len(c.buff) {
		c.count++
		return c.buff[c.count-1]
	}
	return nil
}

func (c *CacheOutputQuery) Evaluate(context NodeIterator) interface{} {
	c.buff = c.buff[:0]
	c.count = 0
	return c.qyInput.Evaluate(context)
}

func (c *CacheOutputQuery) Reset() {
	c.count = 0
}

func (c *CacheOutputQuery) Current() xpath.Navigator {
	if c.count == 0 {
		return nil
	}
	return c.buff[c.count-1]
}

func (c *CacheOutputQuery) CurrentPosition() int {
	return c.count
}

func (c *CacheOutputQuery) Count() int {
	return len(c.buff)
}
