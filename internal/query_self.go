package internal

import "github.com/zhengchun/selector/xpath"

type selfQuery struct {
	qyInput  Query
	matches  func(xpath.Navigator) bool
	currnode xpath.Navigator
	position int
}

func (s *selfQuery) Advance() xpath.Navigator {
	for {
		nav := s.qyInput.Advance()
		if nav == nil {
			break
		}
		if s.matches(nav) {
			s.currnode = nav
			s.position = 1
			return nav
		}
	}
	return nil
}

func (s *selfQuery) Evaluate(iter NodeIterator) interface{} {
	s.qyInput.Evaluate(iter)
	return s
}

func (s *selfQuery) Reset() {
	s.currnode = nil
	s.position = 0
	s.qyInput.Reset()
}

func (s *selfQuery) Current() xpath.Navigator {
	return s.currnode
}
