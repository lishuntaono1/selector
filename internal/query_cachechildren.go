package internal

import "github.com/zhengchun/selector/xpath"

type CacheChildrenQuery struct {
	ChildrenQuery

	elementStk  *Stack
	positionStk *Stack
	nextInput   xpath.Navigator
	needInput   bool
}

func (c *CacheChildrenQuery) getNextInput() xpath.Navigator {
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

func (c *CacheChildrenQuery) decideNextNode() bool {
	c.nextInput = c.getNextInput()

	if c.nextInput != nil {
		if compareNodes(c.currNode, c.nextInput) == xpath.XmlNodeOrderAfter {
			c.elementStk.Push(c.currNode)
			c.positionStk.Push(c.position)
			c.currNode = c.nextInput
			c.nextInput = nil
			if !c.currNode.MoveToFirstChild() {
				return false
			}
			c.position = 0
		}
	}
	return true
}

func (c *CacheChildrenQuery) Advance() xpath.Navigator {
	for {
		if c.needInput {
			if c.elementStk.count == 0 {
				c.currNode = c.getNextInput()
				if c.currNode == nil {
					return nil
				}
				if !c.currNode.MoveToFirstChild() {
					continue
				}
				c.position = 0
			} else {
				c.currNode = c.elementStk.Pop().(xpath.Navigator)
				c.position = c.positionStk.Pop().(int)
				if !c.decideNextNode() {
					continue
				}
			}
			c.needInput = false
		} else {
			if !c.currNode.MoveToNext() || !c.decideNextNode() {
				c.needInput = true
				continue
			}
		}

		if c.matches(c.currNode) {
			c.position++
			return c.currNode
		}
	}
}

func (c *CacheChildrenQuery) Reset() {
	c.elementStk.Clear()
	c.positionStk.Clear()
	c.needInput = true
	c.nextInput = nil
	c.ChildrenQuery.Reset()
}
