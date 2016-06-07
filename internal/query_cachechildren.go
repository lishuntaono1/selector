package internal

import "github.com/zhengchun/selector/xpath"

type cacheChildrenQuery struct {
	childrenQuery

	elementStk  *Stack
	positionStk *Stack
	nextInput   xpath.Navigator
	needInput   bool
}

func (c *cacheChildrenQuery) getNextInput() xpath.Navigator {
	var result xpath.Navigator
	if c.nextInput != nil {
		result = c.nextInput
		c.nextInput = nil
	} else {
		result = c.qyInput.Advance()
		if result != nil {
			result = result.Clone()
		}
	}
	return result
}

func (c *cacheChildrenQuery) decideNextNode() bool {
	c.nextInput = c.getNextInput()

	if c.nextInput != nil {
		if compareNodes(c.currnode, c.nextInput) == xpath.XmlNodeOrderAfter {
			c.elementStk.Push(c.currnode)
			c.positionStk.Push(c.position)
			c.currnode = c.nextInput
			c.nextInput = nil
			if !c.currnode.MoveToFirstChild() {
				return false
			}
			c.position = 0
		}
	}
	return true
}

func (c *cacheChildrenQuery) Advance() xpath.Navigator {
	for {
		if c.needInput {
			if c.elementStk.count == 0 {
				c.currnode = c.getNextInput()
				if c.currnode == nil {
					return nil
				}
				if !c.currnode.MoveToFirstChild() {
					continue
				}
				c.position = 0
			} else {
				c.currnode = c.elementStk.Pop().(xpath.Navigator)
				c.position = c.positionStk.Pop().(int)
				if !c.decideNextNode() {
					continue
				}
			}
			c.needInput = false
		} else {
			if !c.currnode.MoveToNext() || !c.decideNextNode() {
				c.needInput = true
				continue
			}
		}

		if c.matches(c.currnode) {
			c.position++
			return c.currnode
		}
	}
}

func (c *cacheChildrenQuery) Reset() {
	c.elementStk.Clear()
	c.positionStk.Clear()
	c.needInput = true
	c.nextInput = nil
	//base reset
	c.childrenQuery.Reset()
}

func (c *cacheChildrenQuery) CurrentPosition() int {
	return c.position
}
