package internal

import (
	"fmt"
	"reflect"

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
	query string
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
						qyGrandInput = &contextQuery{}
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
		qyInput = &contextQuery{}
		*props = noneProp
	}
	switch root.axistype {
	case AxisChild:
		result = childrenSelector(qyInput, matches)
	case AxisAncestor:
		//result = ancestorSelector(qyInput, false, matches)
	case AxisAncestorOrSelf:
		//result = ancestorSelector(qyInput, true, matches)
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

	return filterSelector(qyInput, cond, propsCond&hasPositionProp == 0, nil)
}

func (builder *QueryBuilder) processOperator(root *Operator, flags Flags, props *Props) Query {
	var props1, props2 Props

	op1 := builder.processNode(root.opnd1, noneFlag, &props1)
	op2 := builder.processNode(root.opnd2, noneFlag, &props2)
	*props = props1 | props2
	switch root.op {
	case OpPLUS, OpMINUS, OpMUL, OpMOD, OpDIV:
		return &numericExpr{root.op, op1, op2}
	case OpLT, OpGT, OpLE, OpGE, OpEQ, OpNE:
		return &logicalExpr{root.op, op1, op2}
	case OpOR, OpAND:
		if root.opnd1.ReturnType() != BooleanType {
			//BooleanFunctions
		}
		if root.opnd2.ReturnType() != BooleanType {
			// /BooleanFunctions
		}
		return &booleanExpr{root.op == OpOR, op1, op2}
	case OpUNION:
	}
	return nil
}

// nodeFunctions is ValueQuery interface implemented.
type nodeFunctions struct {
	ft  FunctionType
	arg Query
}

func (f *nodeFunctions) Advance() xpath.Navigator {
	return nil
}

func (f *nodeFunctions) Evaluate(iter NodeIterator) interface{} {
	switch f.ft {
	case FuncPosition:
		return float64(iter.CurrentPosition())
	case FuncLast:
		fmt.Println(reflect.TypeOf(iter))
		return iter.Count()
		//default:
		//	panic("sorry,this feature not supported yet.")
	}
	return float64(iter.CurrentPosition())
}

func (f *nodeFunctions) Current() xpath.Navigator {
	return nil
}

func (f *nodeFunctions) MoveNext() bool {
	return false
}

func (f *nodeFunctions) Reset() {
}

func (f *nodeFunctions) Count() int {
	return 0
}

func (f *nodeFunctions) CurrentPosition() int {
	return 0
}

func (builder *QueryBuilder) processFunction(root *Function, flags Flags, props *Props) Query {
	*props = noneProp
	var qy Query
	switch root.functype {
	case FuncPosition:
		qy = &nodeFunctions{root.functype, nil}
		*props |= hasPositionProp
	case FuncLast:
		qy = &nodeFunctions{root.functype, nil}
		*props |= hasLastProp
	default:
		panic(fmt.Sprintf("The XPath query %s is not supported.", builder.query))
	}
	return qy
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
	case FunctionAst:
		result = builder.processFunction(root.(*Function), flags, props)
	case RootAst:
		result = &absoluteQuery{}
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

type filterQuery struct {
	qyInput  Query
	cond     Query
	matches  func(xpath.Navigator) bool
	currnode xpath.Navigator
}

func (f *filterQuery) Advance() xpath.Navigator {
	for {
		nav := f.qyInput.Advance()
		if nav == nil {
			return nil
		}
		f.currnode = nav
		if f.EvaluatePredicate() {
			return nav
		}
	}
}

func (f *filterQuery) Evaluate(iter NodeIterator) interface{} {
	f.qyInput.Evaluate(iter)
	return f
}

func (f *filterQuery) EvaluatePredicate() bool {
	var x = f.cond.Evaluate(f.qyInput)
	v := reflect.ValueOf(x)
	switch v.Kind() {
	case reflect.Bool:
		return v.Bool()
	case reflect.String:
		return len(v.String()) > 0
	case reflect.Float64:
		return int(v.Float()) == f.qyInput.CurrentPosition()
	default:
		if reflect.TypeOf(x).Implements(reflect.TypeOf((*NodeIterator)(nil)).Elem()) {
			return f.cond.Advance() != nil
		}
	}
	return false
}

func (f *filterQuery) MoveNext() bool {
	return f.Advance() != nil
}

func (f *filterQuery) Reset() {
	f.cond.Reset()
	f.qyInput.Reset()
}

func (f *filterQuery) Count() int {
	clone := *f
	clone.Reset()
	var count int
	for clone.MoveNext() {
		count++
	}
	return count
}

func (f *filterQuery) Current() xpath.Navigator {
	return f.currnode
}

func (f *filterQuery) CurrentPosition() int {
	return 0
}

func filterSelector(qyInput Query, cond Query, noPosition bool, matches func(xpath.Navigator) bool) Query {
	return &filterQuery{qyInput: qyInput, cond: cond, matches: matches}
}

type childrenQuery struct {
	qyInput  Query
	position int
	movenext func() bool
	matches  func(xpath.Navigator) bool
	currnode xpath.Navigator
}

func (c *childrenQuery) Advance() xpath.Navigator {
	for {
		if c.movenext == nil {
			nav := c.qyInput.Advance()

			if nav == nil {
				return nil
			}
			c.currnode = nav.Clone()
			c.movenext = func() bool {
				for {
					if c.position == 0 && !c.currnode.MoveToFirstChild() {
						return false
					} else if c.position > 0 && !c.currnode.MoveToNext() {
						c.currnode.MoveToParent()
						return false
					}
					c.position++
					if c.matches(c.currnode) {
						return true
					}
				}
			}
		}
		if c.movenext() {
			return c.currnode
		} else {
			c.position = 0
			c.movenext = nil
		}
	}
}

func (c *childrenQuery) Evaluate(iter NodeIterator) interface{} {
	c.qyInput.Evaluate(iter)
	return c
}

func (c *childrenQuery) MoveNext() bool {
	return c.Advance() != nil
}

func (c *childrenQuery) Current() xpath.Navigator {
	return c.currnode
}

func (c *childrenQuery) CurrentPosition() int {
	return c.position
}

func (c *childrenQuery) Reset() {
	c.currnode = nil
	c.position = 0
	c.movenext = nil
	c.qyInput.Reset()
}

func (c *childrenQuery) Count() int {
	clone := *c
	clone.Reset()
	var count int
	for clone.MoveNext() {
		count++
	}
	return count
}

func childrenSelector(qyInput Query, matches func(xpath.Navigator) bool) Query {
	return &childrenQuery{qyInput: qyInput, matches: matches}
}

type attributeQuery struct {
	qyInput  Query
	position int
	matches  func(xpath.Navigator) bool

	onAttr   bool
	currNode xpath.Navigator
}

func (a *attributeQuery) Advance() xpath.Navigator {

	for {
		if !a.onAttr {
			nav := a.qyInput.Advance()
			if nav == nil {
				return nil
			}
			a.position = 0
			a.currNode = nav.Clone()
			a.onAttr = a.currNode.MoveToFirstAttribute()
		} else {
			a.onAttr = a.currNode.MoveToNextAttribute()
		}
		if a.onAttr {
			if a.matches(a.currNode) {
				a.position++
				return a.currNode
			}
		}
	}
}

func (a *attributeQuery) Evaluate(iter NodeIterator) interface{} {
	a.qyInput.Evaluate(iter)
	return a
}

func (a *attributeQuery) Current() xpath.Navigator {
	return a.currNode
}

func (a *attributeQuery) MoveNext() bool {
	return a.Advance() != nil
}

func (a *attributeQuery) Reset() {
	a.currNode = nil
	a.position = 0
	a.qyInput.Reset()
}

func (a *attributeQuery) Count() int {
	clone := *a
	clone.Reset()
	var count int
	for clone.MoveNext() {
		count++
	}
	return count
}

func (a *attributeQuery) CurrentPosition() int {
	return a.position
}

func attributeSelector(qyInput Query, matches func(xpath.Navigator) bool) Query {
	return &attributeQuery{qyInput: qyInput, matches: matches}
}

type descendantQuery struct {
	qyInput   Query
	matchSelf bool
	matches   func(xpath.Navigator) bool
	position  int

	movenext func() bool
	currnode xpath.Navigator
}

func (d *descendantQuery) Advance() xpath.Navigator {
	for {
		if d.movenext == nil {
			var first bool = true
			var level int

			if nav := d.qyInput.Advance(); nav == nil {
				return nil
			} else {
				d.currnode = nav.Clone()
			}
			d.movenext = func() bool {
				if first {
					first = false
					if d.matchSelf && d.matches(d.currnode) {
						d.position = 1
						return true
					}
				}
				for {
					if d.currnode.MoveToFirstChild() {
						level++
					} else {
						for {
							if level == 0 {
								return false
							}
							if d.currnode.MoveToNext() {
								break
							}
							d.currnode.MoveToParent()
							level--
						}
					}
					if d.matches(d.currnode) {
						return true
					}
				}
			}
		}
		if d.movenext() {
			d.position++
			return d.currnode
		} else {
			d.movenext = nil
		}
	}
}

func (d *descendantQuery) Evaluate(ctx NodeIterator) interface{} {
	d.qyInput.Evaluate(ctx)
	return d
}

func (d *descendantQuery) Current() xpath.Navigator {
	return d.currnode
}

func (d *descendantQuery) MoveNext() bool {
	return d.Advance() != nil
}

func (d *descendantQuery) CurrentPosition() int {
	return d.position
}

func (d *descendantQuery) Reset() {
	d.position = 0
	d.currnode = nil
	d.movenext = nil
	d.qyInput.Reset()
}

func (d *descendantQuery) Count() int {
	clone := *d
	clone.Reset()
	var count int
	for !clone.MoveNext() {
		count++
	}
	return count
}

func descendantSelector(qyInput Query, matchSelf bool, matches func(xpath.Navigator) bool) Query {
	return &descendantQuery{qyInput: qyInput, matchSelf: matchSelf, matches: matches}
}

type operandQuery struct {
	val interface{}
}

func (o *operandQuery) Advance() xpath.Navigator {
	return nil
}

func (o *operandQuery) Evaluate(iter NodeIterator) interface{} {
	return o.val
}

func (o *operandQuery) Current() xpath.Navigator {
	return nil
}

func (o *operandQuery) MoveNext() bool {
	return true
}

func (o *operandQuery) Count() int {
	return 0
}

func (o *operandQuery) Reset() {
}

func (o *operandQuery) CurrentPosition() int {
	return 0
}

func getXPathType(v interface{}) ResultType {
	vt := reflect.ValueOf(v)
	switch vt.Kind() {
	case reflect.Bool:
		return BooleanType
	case reflect.Float64:
		return NumberType
	case reflect.String:
		return StringType
	default:
		if _, ok := v.(NodeIterator); ok {
			return NodeSetType
		}
	}
	return StringType
}
