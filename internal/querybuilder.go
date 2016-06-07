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
	firstInput Query
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
					result = &descendantQuery{qyInput: qyGrandInput, matchSelf: false, matches: matches}
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
		if *props&nonFlatProp != 0 {
			result = &cacheChildrenQuery{
				childrenQuery: childrenQuery{qyInput: qyInput, matches: matches},
				elementStk:    &Stack{},
				positionStk:   &Stack{},
				needInput:     true,
			}
		} else {
			result = &childrenQuery{qyInput: qyInput, matches: matches}
		}

	case AxisAncestor:
		result = &ancestorQuery{parentQuery: parentQuery{qyInput: qyInput, matches: matches}, matchSelf: false}
		*props |= nonFlatProp
	case AxisAncestorOrSelf:
		result = &ancestorQuery{parentQuery: parentQuery{qyInput: qyInput, matches: matches}, matchSelf: true}
		*props |= nonFlatProp
	case AxisAttribute:
		result = &attributeQuery{qyInput: qyInput, matches: matches}
	case AxisDescendant:
		result = &descendantQuery{qyInput: qyInput, matchSelf: false, matches: matches}
		*props |= nonFlatProp
	case AxisDescendantOrSelf:
		result = &descendantQuery{qyInput: qyInput, matchSelf: true, matches: matches}
		*props |= nonFlatProp
	case AxisParent:
		result = &parentQuery{qyInput: qyInput, matches: matches}
	case AxisSelf:
		result = &selfQuery{qyInput: qyInput, matches: matches}
	default:
		panic("axis type not supported.")
	}
	return result
}

func (builder *QueryBuilder) processFilter(root *Filter, flags Flags, props *Props) Query {
	var properties QueryProps
	properties = NoneQueryProp

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

	if builder.firstInput == nil {
		type BaseAxisQuery struct {
		}
		//firstInput, ok := qyInput.(BaseAxisQuery)
		//if ok {
		//builder.firstInput = firstInput
		//}
	}

	//merge := properties&MergeQueryProp != 0
	reverse := properties&ReverseQueryProp != 0

	if propsCond&hasPositionProp != 0 {
		if reverse {
			//qyInput = new ReversePositionQuery(qyInput)
			panic("ReversePositionQuery not implemented.")
		} else if propsCond&hasLastProp != 0 {
			qyInput = &forwardPositionQuery{qyInput: qyInput}
		}
	}
	return &filterQuery{qyInput: qyInput, cond: cond}
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
