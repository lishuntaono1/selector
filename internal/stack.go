package internal

type Stack struct {
	nodes []interface{}
	count int
}

func (s *Stack) Push(n interface{}) {
	if s.nodes == nil {
		s.nodes = make([]interface{}, 5)
	}
	if s.count >= len(s.nodes) {
		nodes := make([]interface{}, len(s.nodes)*2)
		copy(nodes, s.nodes)
		s.nodes = nodes
	}
	s.nodes[s.count] = n
	s.count++
}

func (s *Stack) Pop() interface{} {
	if s.count == 0 {
		return nil
	}
	node := s.nodes[s.count-1]
	s.count--
	return node
}

func (s *Stack) Clear() {
	s.nodes = s.nodes[:0]
	s.count = 0
}
