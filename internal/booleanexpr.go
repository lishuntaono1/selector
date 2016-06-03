package internal

import "github.com/zhengchun/selector/xpath"

type booleanExpr struct {
	isOr  bool
	opnd1 Query
	opnd2 Query
}

func (e *booleanExpr) Advance() xpath.Navigator {
	return nil
}

func (e *booleanExpr) Evaluate(iter NodeIterator) interface{} {
	n1 := e.opnd1.Evaluate(iter)
	if n1.(bool) == e.isOr {
		return n1
	}
	return e.opnd2.Evaluate(iter)
}

func (e *booleanExpr) Current() xpath.Navigator {
	return nil
}

func (e *booleanExpr) MoveNext() bool {
	return false
}

func (e *booleanExpr) CurrentPosition() int {
	return 0
}
