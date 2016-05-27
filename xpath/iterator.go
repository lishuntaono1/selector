package xpath

type NodeIterator interface {
	Current() Navigator
	MoveNext() bool
}
