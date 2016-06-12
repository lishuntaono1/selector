package internal

type ForwardPositionQuery struct {
	CacheOutputQuery
}

func (f *ForwardPositionQuery) Evaluate(context NodeIterator) interface{} {
	f.CacheOutputQuery.Evaluate(context)
	for {
		node := f.qyInput.Advance()
		if node == nil {
			break
		}
		f.buff = append(f.buff, node.Clone())
	}
	return f
}
