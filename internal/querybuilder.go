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
	depth      int
	query      string
	firstInput interface{}
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
						qyGrandInput = &ContextQuery{}
						*props = noneProp
					}
					result = &DescendantQuery{BaseAxisQuery: BaseAxisQuery{qyInput: qyGrandInput, matches: matches}}
					if *props&nonFlatProp != 0 {
						result = &DocumentOrderQuery{CacheOutputQuery: CacheOutputQuery{qyInput: result}}
					}
					*props |= nonFlatProp
					return result
				}
			}
			if root.axistype == AxisDescendant || root.axistype == AxisDescendantOrSelf {
				inputFlags |= smartDescFlag
			}
		}
		qyInput = builder.processNode(root.input, flags, props)
	} else {
		qyInput = &ContextQuery{}
		*props = noneProp
	}
	switch root.axistype {
	case AxisChild:
		if *props&nonFlatProp != 0 {
			result = &CacheChildrenQuery{
				ChildrenQuery: ChildrenQuery{BaseAxisQuery: BaseAxisQuery{qyInput: qyInput, matches: matches}},
				elementStk:    &Stack{},
				positionStk:   &Stack{},
				needInput:     true,
			}
		} else {
			result = &ChildrenQuery{BaseAxisQuery: BaseAxisQuery{qyInput: qyInput, matches: matches}}
		}
	case AxisAncestor:
		result = &AncestorQuery{
			ParentQuery: ParentQuery{
				CacheAxisQuery: CacheAxisQuery{BaseAxisQuery: BaseAxisQuery{qyInput: qyInput, matches: matches},
					buff: make([]xpath.Navigator, 0)},
			},
		}
		*props |= nonFlatProp
	case AxisAncestorOrSelf:
		result = &AncestorQuery{
			ParentQuery: ParentQuery{
				CacheAxisQuery: CacheAxisQuery{BaseAxisQuery: BaseAxisQuery{qyInput: qyInput, matches: matches},
					buff: make([]xpath.Navigator, 0)},
			},
			self: true,
		}
		*props |= nonFlatProp
	case AxisAttribute:
		result = &AttributeQuery{
			BaseAxisQuery: BaseAxisQuery{qyInput: qyInput, matches: matches}, // AxisQuery
		}
	case AxisDescendant:
		if flags&smartDescFlag != 0 {
			// result = new DescendantOverDescendantQuery(qyInput, false, root.Name, root.Prefix, root.NodeType, /*abbrAxis:*/false);
		} else {
			result = &DescendantQuery{
				BaseAxisQuery: BaseAxisQuery{qyInput: qyInput, matches: matches},
			}
			if *props&nonFlatProp != 0 {
				result = &DocumentOrderQuery{CacheOutputQuery: CacheOutputQuery{qyInput: result}}
			}
		}

		*props |= nonFlatProp
	case AxisDescendantOrSelf:
		result = &DescendantQuery{
			BaseAxisQuery: BaseAxisQuery{qyInput: qyInput, matches: matches}, self: true}
		if *props&nonFlatProp != 0 {
			result = &DocumentOrderQuery{CacheOutputQuery: CacheOutputQuery{qyInput: result}}
		}
		*props |= nonFlatProp
	case AxisParent:
		result = &ParentQuery{
			CacheAxisQuery: CacheAxisQuery{BaseAxisQuery: BaseAxisQuery{qyInput: qyInput, matches: matches},
				buff: make([]xpath.Navigator, 0)},
		}
	case AxisSelf:
		result = &SelfQuery{BaseAxisQuery: BaseAxisQuery{qyInput: qyInput, matches: matches}}
	default:
		panic("axis type not supported.")
	}
	return result
}

func (builder *QueryBuilder) processFilter(root *Filter, flags Flags, props *Props) Query {
	first := flags&filterFlag == 0

	var propsCond Props
	cond := builder.processNode(root.condition, noneFlag, &propsCond)
	if canBeNumber(root.condition.ReturnType()) && (propsCond&(hasPositionProp|hasLastProp)) != 0 {
		propsCond |= hasPositionProp
		flags |= posFilterFlag
	}
	flags &= ^smartDescFlag

	qyInput := builder.processNode(root.input, flags|filterFlag, props)

	if root.input.Type() != FilterAst {
		*props &= ^posFilterProp
	}

	if propsCond&hasPositionProp != 0 {
		*props |= posFilterProp
	}

	qyFilter, ok := qyInput.(*FilterQuery)
	if ok && (propsCond&hasPositionProp) == 0 && getQueryReturnType(qyFilter.cond) != AnyType {
		prevCond := qyFilter.cond
		if getQueryReturnType(prevCond) == NumberType {
			prevCond = &logicalExpr{op: OpEQ, opnd1: nil, opnd2: prevCond}
		}
		cond = newBooleanExpr(OpAND, prevCond, cond)
		qyInput = qyFilter.qyInput
	}

	if *props&posFilterProp != 0 {
		input, ok := qyInput.(*DocumentOrderQuery)
		if ok {
			qyInput = input.qyInput
		}
	}
	var first_qyInput Query

	if builder.firstInput == nil {
		//BaseAxisQuery in reflect
		val := reflect.Indirect(reflect.ValueOf(qyInput))
		for {
			if val.Type().Field(0).Name == "BaseAxisQuery" {
				first_qyInput = qyInput
				builder.firstInput = val.FieldByName("BaseAxisQuery").Interface()
				break
			} else if val.NumField() > 0 {
				val = val.Field(0)
			} else {
				break
			}
		}

	}

	merge := getQueryProperties(qyInput)&MergeQueryProp != 0
	reverse := false // getQueryProperties(qyInput)&ReverseQueryProp != 0
	if propsCond&hasPositionProp != 0 {
		if reverse {
			qyInput = &ReversePositionQuery{
				ForwardPositionQuery: ForwardPositionQuery{CacheOutputQuery: CacheOutputQuery{qyInput: qyInput, buff: make([]xpath.Navigator, 0)}},
			}
		} else if propsCond&hasLastProp != 0 {
			qyInput = &ForwardPositionQuery{CacheOutputQuery: CacheOutputQuery{qyInput: qyInput, buff: make([]xpath.Navigator, 0)}}
		}
	}

	if first && builder.firstInput != nil {
		if merge && (*props&posFilterProp) != 0 {
			qyInput = &FilterQuery{BaseAxisQuery: BaseAxisQuery{qyInput: qyInput}, cond: cond}
			parent := builder.firstInput.(BaseAxisQuery).qyInput
			if reflect.TypeOf(parent) != reflect.TypeOf((*ContextQuery)(nil)) {
				if ca, ok := first_qyInput.(*CacheChildrenQuery); ok {
					ca.qyInput = &ContextQuery{}
				}

				builder.firstInput = nil
				return &MergeFilterQuery{CacheOutputQuery: CacheOutputQuery{qyInput: parent}, child: qyInput}
			}
			builder.firstInput = nil
			return qyInput
		}
		builder.firstInput = nil
	}
	return &FilterQuery{BaseAxisQuery: BaseAxisQuery{qyInput: qyInput}, cond: cond}
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
		return newBooleanExpr(root.op, op1, op2)
	case OpUNION:
		*props |= nonFlatProp
		return &unionExpr{qy1: op1, qy2: op2}
	}
	return nil
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
	case FuncCount:
		return &nodeFunctions{root.functype, builder.processNode(root.argument[0], flags, props)}
	case FuncName, FuncNameSpaceUri, FuncLocalName:
		if len(root.argument) > 0 {
			return &nodeFunctions{root.functype, builder.processNode(root.argument[0], flags, props)}
		} else {
			return &nodeFunctions{root.functype, nil}
		}
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
		result = &OperandQuery{(root.(*Operand).val)}
	case FunctionAst:
		result = builder.processFunction(root.(*Function), flags, props)
	case RootAst:
		result = &AbsoluteQuery{}
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
	builder.query = xpath
	return builder.build(ParseXPathExpression(xpath))
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

var queryReturnTypes = map[string]ResultType{}

func getQueryReturnType(q Query) ResultType {
	name := reflect.TypeOf(q).Elem().Name()
	if val, ok := queryReturnTypes[name]; ok {
		return val
	}
	return AnyType
}

var queryProperties = map[string]QueryProps{
	"MergeFilterQuery": MergeQueryProp | CachedQueryProp | PositionQueryProp | CountQueryProp,
}

func getQueryProperties(q Query) QueryProps {
	name := reflect.TypeOf(q).Elem().Name()
	if prop, ok := queryProperties[name]; ok {
		return prop
	}
	return MergeQueryProp
}
