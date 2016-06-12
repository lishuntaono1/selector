package internal

type ReversePositionQuery struct {
	ForwardPositionQuery
}

func (r *ReversePositionQuery) CurrentPosition() int {
	return len(r.buff) - r.count + 1
}
