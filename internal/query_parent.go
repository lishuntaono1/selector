package internal

import "github.com/zhengchun/selector/xpath"

// parentQuery is the axis query with cache.
type parentQuery struct {
	qyInput      Query
	outputBuffer []xpath.Navigator
	matches      func(xpath.Navigator) bool
	count        int
}

func (p *parentQuery) Advance() xpath.Navigator {
	if p.count < len(p.outputBuffer) {
		nav := p.outputBuffer[p.count]
		p.count++
		return nav
	}
	return nil
}

func (p *parentQuery) Evaluate(context NodeIterator) interface{} {
	p.Reset()
	p.qyInput.Evaluate(context)

	if p.outputBuffer == nil {
		p.outputBuffer = make([]xpath.Navigator, 0)
	} else {
		p.outputBuffer = p.outputBuffer[:0]
	}

	for {
		input := p.qyInput.Advance()
		if input == nil {
			break
		}
		input = input.Clone()
		if input.MoveToParent() {
			if p.matches(input) {
				// insert
				if b, ok := insertNode(p.outputBuffer[0:], input); ok {
					p.outputBuffer = b
				}
			}
		}
	}
	return p
}

func (p *parentQuery) Reset() {
	p.count = 0
	p.qyInput.Reset()
}

func (p *parentQuery) Current() xpath.Navigator {
	if p.count == 0 {
		return nil
	}
	return p.outputBuffer[p.count-1]
}

func (p *parentQuery) Count() int {
	return len(p.outputBuffer)
}

func (p *parentQuery) CurrentPosition() int {
	return p.count
}
