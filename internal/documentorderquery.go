package internal

type DocumentOrderQuery struct {
	CacheOutputQuery
}

func (d *DocumentOrderQuery) Evaluate(context NodeIterator) interface{} {
	d.CacheOutputQuery.Evaluate(context)
	for {
		node := d.CacheOutputQuery.qyInput.Advance()
		if node == nil {
			break
		}
		if buff, ok := insertNode(d.buff[0:], node); ok {
			d.buff = buff
		}
	}
	return d
}
