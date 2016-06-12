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

func (e *booleanExpr) Reset() {

}

func newBooleanExpr(op OpType, opnd1, opnd2 Query) *booleanExpr {
	if getQueryReturnType(opnd1) != BooleanType {
		// opnd1 =  BooleanFunctions(Function.FunctionType.FuncBoolean, opnd1);
	}
	if getQueryReturnType(opnd2) != BooleanType {
		// opnd2 = new BooleanFunctions(Function.FunctionType.FuncBoolean, opnd2);
	}
	expr := &booleanExpr{}
	expr.opnd1 = opnd1
	expr.opnd2 = opnd2
	expr.isOr = op == OpOR
	return expr
}
