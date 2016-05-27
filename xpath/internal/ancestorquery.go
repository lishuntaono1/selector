package internal

type AncestorQuery struct {
	*AxisQuery
	matchSelf bool
}
