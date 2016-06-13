package internal

type ParentQuery struct {
	CacheAxisQuery
}

func (p *ParentQuery) Evaluate(context NodeIterator) interface{} {
	p.CacheAxisQuery.Evaluate(context)

	for {
		input := p.qyInput.Advance()
		if input == nil {
			break
		}
		input = input.Clone()
		if input.MoveToParent() {
			if p.matches(input) {
				if b, ok := insertNode(p.buff[0:], input); ok {
					p.buff = b
				}
			}
		}
	}
	return p
}
