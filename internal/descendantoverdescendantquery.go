package internal

import "github.com/zhengchun/selector/xpath"

type DescendantOverDescendantQuery struct {
	BaseAxisQuery
	self, abbr bool
	level      int
}

func (d *DescendantOverDescendantQuery) Advance() xpath.Navigator {
	/*for {
		if d.level == 0 {

		}
	}*/
	return nil
}

func (d *DescendantOverDescendantQuery) Reset() {
	d.level = 0
	d.BaseAxisQuery.Reset()
}
