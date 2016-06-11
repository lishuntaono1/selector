package internal

import "github.com/zhengchun/selector/xpath"

// CacheAxisQuery is the axis query with cache feature.
type CacheAxisQuery struct {
	BaseAxisQuery
	buff  []xpath.Navigator
	count int
}

func (c *CacheAxisQuery) Advance() xpath.Navigator {
	if c.count < len(c.buff) {
		c.count++
		return c.buff[c.count-1]
	}
	return nil
}

func (c *CacheAxisQuery) Evaluate(context NodeIterator) interface{} {
	c.buff = c.buff[:0]
	c.BaseAxisQuery.Evaluate(context)
	return c
}

func (c *CacheAxisQuery) Count() int {
	return c.count
}

func (c *CacheAxisQuery) Reset() {
	c.count = 0
}

func (c *CacheAxisQuery) Current() xpath.Navigator {
	if c.count == 0 {
		return nil
	}
	return c.buff[c.count-1]
}
