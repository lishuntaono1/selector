package internal

import "github.com/zhengchun/selector/xpath"

type SelfQuery struct {
	BaseAxisQuery
}

func (s *SelfQuery) Advance() xpath.Navigator {
	for {
		nav := s.qyInput.Advance()
		if nav == nil {
			break
		}
		if s.matches(nav) {
			s.currNode = nav
			s.position = 1
			return nav
		}
	}
	return nil
}
