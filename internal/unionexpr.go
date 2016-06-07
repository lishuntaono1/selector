package internal

import "github.com/zhengchun/selector/xpath"

type unionExpr struct {
	qy1, qy2              Query
	advance1, advance2    bool
	currentNode, nextNode xpath.Navigator
}

func (e *unionExpr) Advance() xpath.Navigator {
	var res1, res2 xpath.Navigator
	if e.advance1 {
		res1 = e.qy1.Advance()
	} else {
		res1 = e.nextNode
	}
	if e.advance2 {
		res2 = e.qy2.Advance()
	} else {
		res2 = e.nextNode
	}

	var order = xpath.XmlNodeOrderBefore

	if res1 != nil && res2 != nil {
		order = compareNodes(res1, res2)
	} else if res2 == nil {
		e.advance1 = true
		e.advance2 = false
		e.currentNode = res1
		e.nextNode = nil
		return res1
	} else {
		e.advance1 = false
		e.advance2 = true
		e.currentNode = res2
		e.nextNode = nil
		return res2
	}
	switch order {
	case xpath.XmlNodeOrderBefore:
		e.nextNode = res2
		e.advance2 = false
		e.advance1 = true
		e.currentNode = res1
		return res1
	case xpath.XmlNodeOrderAfter:
		e.nextNode = res1
		e.advance1 = false
		e.advance2 = true
		e.currentNode = res2
		return res2
	}
	e.advance1 = true
	e.advance2 = true
	return res1
}

func (e *unionExpr) Evaluate(context NodeIterator) interface{} {
	e.qy1.Evaluate(context)
	e.qy2.Evaluate(context)
	e.advance1 = true
	e.advance2 = true
	e.nextNode = nil
	return e
}

func (e *unionExpr) Reset() {
	e.qy1.Reset()
	e.qy2.Reset()
	e.advance1 = true
	e.advance2 = true
	e.nextNode = nil
}

func (e *unionExpr) Current() xpath.Navigator {
	return e.currentNode
}
