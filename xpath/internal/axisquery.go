package internal

import "github.com/zhengchun/selector/xpath"

type AxisQuery struct {
	position    int
	currentNode xpath.Navigator
	count       int

	qyInput      Query
	typetest     xpath.NodeType
	name, prefix string
}

func (q *AxisQuery) MatchNode(context xpath.Navigator) xpath.Navigator {
	panic("MatchNode method NOT implemented.")
}

/*
func (q *AxisQuery) Reset() {
	q.position = 0
	q.currentNode = nil
	q.qyInput.Reset()
}

func (q *AxisQuery) Advance() xpath.Navigator {
	panic(fmt.Sprintf("%v Advance method NOT implemented.", q))
}

func (q *AxisQuery) Evaluate(nodeIterator xpath.NodeIterator) interface{} {
	q.count = -1
	q.Reset()
	q.qyInput.Evaluate(nodeIterator)
	return q
}

func (q *AxisQuery) MatchNode(current xpath.Navigator) xpath.Navigator {
	panic("'{0}' is an invalid XSLT pattern.")
}

func (q *AxisQuery) Clone() xpath.NodeIterator {
	panic(fmt.Sprintf("%v Clone method NOT implemented.", q))
}

func (q *AxisQuery) Current() xpath.Navigator {
	panic(fmt.Sprintf("%v Current method NOT implemented.", q))
}

func (q *AxisQuery) matches(e xpath.Navigator) bool {
	if q.typetest == e.NodeType() || q.typetest == xpath.AllNode ||
		q.typetest == xpath.TextNode && (e.NodeType() == xpath.WhitespaceNode || e.NodeType() == xpath.SignificantWhitespaceNode) {
		if len(q.prefix) > 0 || len(q.name) > 0 {
			if q.name == e.LocalName() || q.name == "" {
				// Namespace checking?
				return true
			}
		} else {
			return true
		}
	}
	return false
}
*/

func newAxisQuery(qyInput Query, name, prefix string, nodetype xpath.NodeType) *AxisQuery {
	return &AxisQuery{
		qyInput:  qyInput,
		name:     name,
		prefix:   prefix,
		typetest: nodetype,
	}
}
