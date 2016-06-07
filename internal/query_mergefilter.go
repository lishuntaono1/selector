package internal

import "github.com/zhengchun/selector/xpath"

type mergeFilterQuery struct {
	qyInput, child Query
	outputBuffer   []xpath.Navigator
	count          int
}

func (m *mergeFilterQuery) Advance() xpath.Navigator {
	if m.count < len(m.outputBuffer) {
		node := m.outputBuffer[m.count]
		m.count++
		return node
	}
	return nil
}

func (m *mergeFilterQuery) Evaluate(context NodeIterator) interface{} {
	for m.qyInput.Advance() != nil {
		m.child.Evaluate(m.qyInput)
		for {
			node := m.child.Advance()
			if node == nil {
				break
			}
			b, ok := insertNode(m.outputBuffer, node)
			if ok {
				m.outputBuffer = b
			}
		}
	}
	return m
}

func (m *mergeFilterQuery) Count() int {
	return m.count
}

func (m *mergeFilterQuery) Reset() {
	m.count = 0
	m.qyInput.Reset()
}

func (m *mergeFilterQuery) Current() xpath.Navigator {
	if len(m.outputBuffer) == 0 {
		return nil
	}
	return m.outputBuffer[m.count-1]
}
