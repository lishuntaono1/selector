package internal

import "github.com/zhengchun/selector/xpath"

type forwardPositionQuery struct {
	qyInput      Query
	outputBuffer []xpath.Navigator
	count        int
}

func (f *forwardPositionQuery) Advance() xpath.Navigator {
	if f.count < len(f.outputBuffer) {
		node := f.outputBuffer[f.count]
		f.count++
		return node
	}
	return nil
}

func (f *forwardPositionQuery) Evaluate(context NodeIterator) interface{} {
	//
	if f.outputBuffer == nil {
		f.outputBuffer = make([]xpath.Navigator, 0)
	} else {
		f.outputBuffer = f.outputBuffer[:0]
	}
	// base evaluate ??
	f.Reset()
	f.qyInput.Evaluate(context)

	for {
		node := f.qyInput.Advance()
		if node == nil {
			break
		}
		f.outputBuffer = append(f.outputBuffer, node.Clone())
	}
	return f
}

func (f *forwardPositionQuery) Reset() {
	f.count = 0
	f.qyInput.Reset()
}

func (f *forwardPositionQuery) Current() xpath.Navigator {
	if f.count == 0 {
		return nil
	}
	return f.outputBuffer[f.count-1]
}

func (f *forwardPositionQuery) Count() int {
	return f.count
}

func (f *forwardPositionQuery) CurrentPosition() int {
	return f.count
}
