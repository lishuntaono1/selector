package internal

type DescendantOverDescendantQuery struct {
	*AxisQuery
	matchSelf bool
}
