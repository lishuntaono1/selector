package xpath

type NodeIterator interface {
	MoveNext() bool
	Current() Navigator
	Clone() NodeIterator
}
