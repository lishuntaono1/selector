package internal

import "github.com/zhengchun/selector/xpath"

type cmpXslt func(OpType, interface{}, interface{}) bool

type logicalExpr struct {
	op    OpType
	opnd1 Query
	opnd2 Query
}

func (l *logicalExpr) Advance() xpath.Navigator {
	return nil
}

func (l *logicalExpr) Evaluate(iter NodeIterator) interface{} {
	var val1 = l.opnd1.Evaluate(iter)
	var val2 = l.opnd2.Evaluate(iter)

	var type1 = getXPathType(val1)
	var type2 = getXPathType(val2)

	if type1 < type2 {
		l.op = invertOperator(l.op)
		valTemp := val1
		val1 = val2
		val2 = valTemp
		typeTmp := type1
		type1 = type2
		type2 = typeTmp
	}
	if l.op == OpEQ || l.op == OpNE {
		return s_CompXsltE[type1][type2](l.op, val1, val2)
	}
	return s_CompXsltO[type1][type2](l.op, val1, val2)
}

func (l *logicalExpr) Current() xpath.Navigator {
	return nil
}

func (l *logicalExpr) Reset() {
}

var s_CompXsltE = [][]cmpXslt{
	[]cmpXslt{cmpNumberNumber, nil, nil, nil, nil},
	[]cmpXslt{nil, nil, nil, nil, nil},
	[]cmpXslt{nil, nil, nil, nil, nil},
	[]cmpXslt{nil, cmpQueryStringE, nil, nil, nil},
	[]cmpXslt{nil, nil, nil, nil, nil},
}

var s_CompXsltO = [][]cmpXslt{
	[]cmpXslt{cmpNumberNumber, nil, nil, nil, nil},
	[]cmpXslt{nil, nil, nil, nil, nil},
	[]cmpXslt{nil, nil, nil, nil, nil},
	[]cmpXslt{nil, nil, nil, nil, nil},
	[]cmpXslt{nil, nil, nil, nil, nil},
}

func cmpStringStringE(op OpType, n1, n2 string) bool {
	return (op == OpEQ) == (n1 == n2)
}

func cmpQueryStringE(op OpType, val1, val2 interface{}) bool {
	n1 := &nodeSet{val1.(Query), nil}
	n2 := val2.(string)
	for n1.MoveNext() {
		if cmpStringStringE(op, n1.Value(), n2) {
			return true
		}
	}
	return false
}

func cmpNumberNumberF(op OpType, n1, n2 float64) bool {
	switch op {
	case OpLT:
		return n1 < n2
	case OpGT:
		return n1 > n2
	case OpLE:
		return n1 <= n2
	case OpGE:
		return n1 >= n2
	case OpEQ:
		return n1 == n2
	case OpNE:
		return n1 != n2
	}
	return false
}

func cmpNumberNumber(op OpType, val1, val2 interface{}) bool {
	return cmpNumberNumberF(op, val1.(float64), val2.(float64))
}

type nodeSet struct {
	opnd Query
	curr xpath.Navigator
}

func (s *nodeSet) MoveNext() bool {
	s.curr = s.opnd.Advance()
	return s.curr != nil
}

func (s *nodeSet) Value() string {
	return s.curr.Value()
}
