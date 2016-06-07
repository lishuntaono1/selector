package internal

import "github.com/zhengchun/selector/xpath"

type contextQuery struct {
	count int
	nav   xpath.Navigator
}

func (c *contextQuery) Advance() xpath.Navigator {
	if c.count == 0 {
		c.count = 1
		return c.nav
	}
	return nil
}

func (c *contextQuery) Evaluate(iter NodeIterator) interface{} {
	c.nav = iter.Current()
	c.count = 0
	return c
}

func (c *contextQuery) MoveNext() bool {
	return c.Advance() != nil
}

func (c *contextQuery) Reset() {
	c.count = 0
}

func (c *contextQuery) Current() xpath.Navigator {
	return c.nav
}
