package internal

import "github.com/zhengchun/selector/xpath"

// nodeFunctions is ValueQuery interface implemented.
type nodeFunctions struct {
	ft  FunctionType
	arg Query
}

func (f *nodeFunctions) evaluateArg(context NodeIterator) xpath.Navigator {
	if f.arg == nil {
		return context.Current()
	}
	f.arg.Evaluate(context)
	return f.arg.Advance()
}

func (f *nodeFunctions) Advance() xpath.Navigator {
	return nil
}

func (f *nodeFunctions) Evaluate(context NodeIterator) interface{} {
	switch f.ft {
	case FuncPosition:
		return float64(iteratorPosition(context))
	case FuncLast:
		return float64(iteratorCount(context))
	case FuncName:
		argVal := f.evaluateArg(context)
		if argVal != nil {
			return argVal.LocalName()
		}
	case FuncCount:
		f.arg.Evaluate(context)
		var count int
		for f.arg.Advance() != nil {
			count++
		}
		return float64(count)
		//default:
		//	panic("sorry,this feature not supported yet.")
	}
	return ""
}

func (f *nodeFunctions) Current() xpath.Navigator {
	return nil
}

func (f *nodeFunctions) Reset() {
}
