package internal

type MergeFilterQuery struct {
	CacheOutputQuery
	child Query
}

func (m *MergeFilterQuery) Evaluate(context NodeIterator) interface{} {
	m.CacheOutputQuery.Evaluate(context)

	for m.qyInput.Advance() != nil {
		m.child.Evaluate(m.qyInput)
		for {
			node := m.child.Advance()
			if node == nil {
				break
			}
			b, ok := insertNode(m.buff, node)
			if ok {
				m.buff = b
			}
		}
	}
	return m
}
