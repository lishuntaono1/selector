package internal

import "github.com/zhengchun/selector/xpath"

type operandQuery struct {
	val interface{}
}

func (o *operandQuery) Advance() xpath.Navigator {
	return nil
}

func (o *operandQuery) Evaluate(iter NodeIterator) interface{} {
	return o.val
}

func (o *operandQuery) Current() xpath.Navigator {
	return nil
}

func (o *operandQuery) Reset() {
}
