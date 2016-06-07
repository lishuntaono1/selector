package internal

import (
	"fmt"
	"reflect"

	"github.com/zhengchun/selector/xpath"
)

type filterQuery struct {
	qyInput  Query
	cond     Query
	currnode xpath.Navigator
}

func (f *filterQuery) Advance() xpath.Navigator {
	for {
		nav := f.qyInput.Advance()
		if nav == nil {
			return nil
		}
		f.currnode = nav
		if f.EvaluatePredicate() {
			return nav
		}
	}
}

func (f *filterQuery) Evaluate(iter NodeIterator) interface{} {
	f.qyInput.Evaluate(iter)
	return f
}

func (f *filterQuery) EvaluatePredicate() bool {
	var x = f.cond.Evaluate(f.qyInput)
	v := reflect.ValueOf(x)
	switch v.Kind() {
	case reflect.Bool:
		return v.Bool()
	case reflect.String:
		return len(v.String()) > 0
	case reflect.Float64:
		fmt.Println(v.Float())
		fmt.Println(iteratorPosition(f.qyInput))
		return int(v.Float()) == iteratorPosition(f.qyInput)
	default:
		if reflect.TypeOf(x).Implements(reflect.TypeOf((*NodeIterator)(nil)).Elem()) {
			return f.cond.Advance() != nil
		}
	}
	return false
}

func (f *filterQuery) MoveNext() bool {
	return f.Advance() != nil
}

func (f *filterQuery) Reset() {
	f.cond.Reset()
	f.qyInput.Reset()
}

func (f *filterQuery) Count() int {
	clone := *f
	clone.Reset()
	var count int
	for clone.MoveNext() {
		count++
	}
	return count
}

func (f *filterQuery) Current() xpath.Navigator {
	return f.currnode
}
