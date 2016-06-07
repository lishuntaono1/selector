package internal

import "github.com/zhengchun/selector/xpath"

func Express(nav xpath.Navigator, xpath string) xpath.NodeIterator {
	builder := &QueryBuilder{}
	q := builder.Build(xpath)
	iter := &XPathSelectionIterator{nav: nav, query: q}
	q.Evaluate(iter)
	return iter
}
