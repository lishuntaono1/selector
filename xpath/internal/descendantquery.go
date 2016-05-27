package internal

type DescendantQuery struct {
	*AxisQuery

	matchSelf bool
	abbrAxis  bool
}

func (q *DescendantQuery) MoveNext() bool {
	return false
}
