package internal

import (
	"fmt"

	"github.com/zhengchun/selector/xpath"
)

type Flags uint

const (
	noneFlag Flags = iota
	smartDescFlag
	posFilterFlag
	filterFlag
)

type Props uint

const (
	noneProp Props = iota
	posFilterProp
	hasPositionProp
	hasLastProp
	nonFlatProp
)

type QueryBuilder struct {
	depth int
}

var canBeNumber = func(t ResultType) bool {
	return t == AnyType || t == NumberType
}

func (builder *QueryBuilder) processAxis(root *Axis, flags Flags, props *Props) Query {
	var matches func(xpath.Navigator) bool
	matches = func(nav xpath.Navigator) bool {
		if root.nodetype == nav.NodeType() ||
			root.nodetype == xpath.AllNode ||
			(root.nodetype == xpath.TextNode &&
				(nav.NodeType() == xpath.WhitespaceNode || nav.NodeType() == xpath.SignificantWhitespaceNode)) {
			if root.name != "" || root.prefix != "" {
				if root.name == nav.LocalName() || root.name == "" {
					return true
				}
			} else {
				return true
			}
		}
		return false
	}

	var result, qyInput Query
	if root.input != nil {
		inputFlags := noneFlag
		if flags&posFilterFlag == 0 {
			input, ok := root.input.(*Axis)
			if ok {
				if root.axistype == AxisChild && input.axistype == AxisDescendantOrSelf && input.nodetype == xpath.AllNode {
					var qyGrandInput Query
					if input.input != nil {
						qyGrandInput = builder.processNode(input.input, smartDescFlag, props)
					} else {
						qyGrandInput = contextSelector()
						*props = noneProp
					}
					result = descendantSelector(qyGrandInput, false, matches)
					return result
				}
			}
			if root.axistype == AxisDescendant || root.axistype == AxisDescendantOrSelf {
				inputFlags |= smartDescFlag
			}
		}
		qyInput = builder.processNode(root.input, flags, props)
	} else {
		qyInput = contextSelector()
		*props = noneProp
	}
	switch root.axistype {
	case AxisChild:
		result = childrenSelector(qyInput, matches)
	case AxisAncestor:
		result = ancestorSelector(qyInput, false, matches)
	case AxisAncestorOrSelf:
		result = ancestorSelector(qyInput, true, matches)
	case AxisAttribute:
		result = attributeSelector(qyInput, matches)
	case AxisDescendant:
		result = descendantSelector(qyInput, false, matches)
	case AxisDescendantOrSelf:
		result = descendantSelector(qyInput, true, matches)
	default:
		panic("axis type not supported.")
	}
	return result
}

func (builder *QueryBuilder) processFilter(root *Filter, flags Flags, props *Props) Query {
	//first := flags&filterFlag == 0
	var propsCond Props
	cond := builder.processNode(root.condition, noneFlag, &propsCond)
	if canBeNumber(root.condition.ReturnType()) && (propsCond&(hasPositionProp|hasLastProp)) != 0 {
		propsCond |= hasPositionProp
		flags |= posFilterFlag
	}
	flags &= ^smartDescFlag
	qyInput := builder.processNode(root.input, flags|filterFlag, props)
	if root.input.Type() != FilterAst {
		// Props.PosFilter is for nested filters only.
		// We clean it here to avoid cleaning it in all other ast nodes.
		*props &= ^posFilterProp
	}
	if propsCond&hasPositionProp != 0 {
		// this condition is positional rightmost filter should be avare of this.
		*props |= posFilterProp
	}
	/*merging predicates*/
	return filterSelector(qyInput, cond, propsCond&hasPositionProp == 0, nil)
}

func (builder *QueryBuilder) processOperator(root *Operator, flags Flags, props *Props) Query {
	var props1, props2 Props
	op1 := builder.processNode(root.opnd1, noneFlag, &props1)
	op2 := builder.processNode(root.opnd2, noneFlag, &props2)

	return logicalSelector(root.op, op1, op2)

	/*
		*props = props1 | props2
		switch root.op {
		case OpPLUS, OpMINUS, OpMUL, OpMOD, OpDIV:
			//Numeric
		case OpLT, OpGT, OpLE, OpGE, OpEQ, OpNE:
			//Logical
		case OpOR, OpAND:
			//Boolean
		case OpUNION:
			*props |= nonFlatProp
			//return new UnionExpr(op1, op2);
		default:
			return nil
		}
	*/

}

func (builder *QueryBuilder) processNode(root AstNode, flags Flags, props *Props) Query {
	if builder.depth = builder.depth + 1; builder.depth > 1024 {
		panic("The xpath query is too complex.")
	}
	var result Query
	switch root.Type() {
	case AxisAst:
		result = builder.processAxis(root.(*Axis), flags, props)
	case FilterAst:
		result = builder.processFilter(root.(*Filter), flags, props)
	case OperatorAst:
		result = builder.processOperator(root.(*Operator), flags, props)
	case ConstantOperandAst:
		result = &operandQuery{(root.(*Operand).val)}
	case RootAst:
		result = absoluteSelector()
	default:
		panic("Unknown QueryType encountered!!")
	}
	return result
}

func (builder *QueryBuilder) build(root AstNode) Query {
	props := noneProp
	return builder.processNode(root, noneFlag, &props)
}

func (builder *QueryBuilder) Build(xpath string) Query {
	return builder.build(ParseXPathExpression(xpath))
}

func contextSelector() selector {
	var count = 0
	return func(nav xpath.Navigator) bool {
		if count == 0 {
			count = 1
			return true
		}
		return false
	}
}

func absoluteSelector() selector {
	return contextSelector()
}

func logicalSelector(op OpType, opnd1, opnd2 Query) selector {
	return func(nav xpath.Navigator) bool {

		if !(opnd1.Matches(nav) && opnd2.Matches(nav)) {
			return false
		}
		type Evaluate interface {
			Evaluate(xpath.Navigator) interface{}
		}
		var val1, val2 interface{}
		if eval, ok := opnd1.(Evaluate); ok {
			val1 = eval.Evaluate(nav)
		} else {
			val1 = nav.Value() // string
		}
		if eval, ok := opnd2.(Evaluate); ok {
			val2 = eval.Evaluate(nav)
		} else {
			val2 = nav.Value() // string
		}

		fmt.Println(val1.(string))

		switch op {
		case OpEQ: // `=`
			return val1.(string) == val2.(string)
		}
		return false
	}
}

func filterSelector(qyInput Query, cond Query, noPosition bool, matches func(xpath.Navigator) bool) selector {
	return func(nav xpath.Navigator) bool {
		for {
			if !qyInput.Matches(nav) {
				return false
			}
			if cond.Matches(nav) {
				return true
			}
		}
	}
}

func childrenSelector(qyInput Query, matches func(xpath.Navigator) bool) selector {
	var movenext func() bool
	var position = 0
	return func(nav xpath.Navigator) bool {
		for {
			if movenext == nil {
				if !qyInput.Matches(nav) {
					return false
				}
				movenext = func() bool {
					fmt.Println(nav.LocalName())
					for {
						if position == 0 && !nav.MoveToFirstChild() {
							return false
						} else if position > 0 && !nav.MoveToNext() {
							nav.MoveToParent()
							return false
						}
						position++
						if matches(nav) {
							return true
						}
					}
				}
			}
			if movenext() {
				return true
			} else {
				position = 0
				movenext = nil
			}
		}
	}
}

func descendantSelector(qyInput Query, matchSelf bool, matches func(xpath.Navigator) bool) selector {
	var movenext func() bool

	return func(nav xpath.Navigator) bool {
		for {
			if movenext == nil {
				var first bool = true
				var level int

				if !qyInput.Matches(nav) {
					return false
				}
				movenext = func() bool {
					if first {
						first = false
						if matchSelf && matches(nav) {
							return true
						}
					}
					for {
						if nav.MoveToFirstChild() {
							level++
						} else {
							for {
								if level == 0 {
									return false
								}
								if nav.MoveToNext() {
									break
								}
								nav.MoveToParent()
								level--
							}
						}
						if matches(nav) {
							return true
						}
					}
				}
			}

			if movenext() {
				return true
			} else {
				movenext = nil
			}
		}

	}
}

func attributeSelector(qyInput Query, matches func(xpath.Navigator) bool) selector {
	return func(nav xpath.Navigator) bool {
		var onAttr bool
		var currNode xpath.Navigator
		for {
			if !onAttr {
				if !qyInput.Matches(nav) {
					return false
				}
				currNode = nav.Clone()
				onAttr = currNode.MoveToFirstAttribute()
			} else {
				onAttr = currNode.MoveToNextAttribute()
			}
			if onAttr {
				if matches(currNode) {
					nav.MoveTo(currNode)
					return true
				}
			}
		}
	}
}

func ancestorSelector(qyInput Query, matchSelf bool, matches func(xpath.Navigator) bool) selector {
	var ancestor xpath.Navigator
	return func(nav xpath.Navigator) bool {
		for {
			if !qyInput.Matches(nav) {
				return false
			}
			if matchSelf {

			}
			if ancestor == nil {
				ancestor = nav.Clone()
			}
			for ancestor.MoveToParent() {
				if matches(ancestor) {
					nav.MoveTo(ancestor)
					return true
				}
			}
			ancestor = nil
		}
	}
}

type operandQuery struct {
	val interface{}
}

func (o *operandQuery) Matches(nav xpath.Navigator) bool {
	return true
}

func (o *operandQuery) Evaluate(nav xpath.Navigator) interface{} {
	return o.val
}
