package internal

import "github.com/zhengchun/selector/xpath"

type OperandQuery struct {
	val interface{}
}

func (o *OperandQuery) Advance() xpath.Navigator {
	return nil
}

func (o *OperandQuery) Evaluate(context NodeIterator) interface{} {
	return o.val
}

func (o *OperandQuery) Current() xpath.Navigator {
	return nil
}

func (o *OperandQuery) Reset() {
}
