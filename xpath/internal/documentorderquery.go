package internal

type DocumentOrderQuery struct {
	input Query
}

func (q *DocumentOrderQuery) MoveNext() bool { return false }
