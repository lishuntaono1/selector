package internal

import (
	"reflect"

	"github.com/zhengchun/selector/xpath"
)

type FilterQuery struct {
	BaseAxisQuery

	cond Query
}

func (f *FilterQuery) Advance() xpath.Navigator {
	for {
		nav := f.qyInput.Advance()
		if nav == nil {
			return nil
		}
		f.currNode = nav
		if f.EvaluatePredicate() {
			return nav
		}
	}
}

func (f *FilterQuery) EvaluatePredicate() bool {
	var x = f.cond.Evaluate(f.qyInput)
	v := reflect.ValueOf(x)
	switch v.Kind() {
	case reflect.Bool:
		return v.Bool()
	case reflect.String:
		return len(v.String()) > 0
	case reflect.Float64:
		return int(v.Float()) == iteratorPosition(f.qyInput)
	default:
		if reflect.TypeOf(x).Implements(reflect.TypeOf((*NodeIterator)(nil)).Elem()) {
			return f.cond.Advance() != nil
		}
	}
	return false
}

func (f *FilterQuery) Reset() {
	f.cond.Reset()
	f.BaseAxisQuery.Reset()
}
