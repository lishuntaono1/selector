package internal

import "github.com/zhengchun/selector/xpath"

type numericExpr struct {
	op    OpType
	opnd1 Query
	opnd2 Query
}

func (e *numericExpr) Advance() xpath.Navigator {
	return nil
}

func (e *numericExpr) Evaluate(iter NodeIterator) interface{} {
	n1 := convertToXPathDouble(e.opnd1.Evaluate(iter))
	n2 := convertToXPathDouble(e.opnd2.Evaluate(iter))
	switch e.op {
	case OpPLUS:
		return n1 + n2
	case OpMINUS:
		return n1 - n2
	case OpMOD:
		return int(n1) % int(n2)
	case OpDIV:
		return n1 / n2
	case OpMUL:
		return n1 * n2
	}
	return 0
}

func (e *numericExpr) Current() xpath.Navigator {
	return nil
}

func (e *numericExpr) MoveNext() bool {
	return false
}

func (e *numericExpr) CurrentPosition() int {
	return 0
}
