package internal

/*
type QueryBuilder struct {
	query       string
	depth       int
	needContext bool
	firstInput  Query

	allowVar, allowKey, allowCurrent bool
}

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

func (b *QueryBuilder) processNode(root AstNode, flags Flags, props *Props) Query {
	if b.depth = b.depth + 1; b.depth > 1024 {
		panic("The xpath query is too complex.")
	}
	var result Query
	switch root.Type() {
	case AxisAst:
		result = b.processAxis(root.(*Axis), flags, props)
	case RootAst:
		result = &AbsoluteQuery{}
	}
	b.depth--
	return result
}

func (b *QueryBuilder) processAxis(root *Axis, flags Flags, props *Props) Query {
	var result Query
	if len(root.prefix) > 0 {
		b.needContext = true
	}
	b.firstInput = nil
	var qyInput Query
	if root.input != nil {
		inputFlags := noneFlag
		if flags&posFilterFlag == 0 {
			input, ok := root.input.(*Axis)
			if ok {
				if root.axis_type == AxisChild && input.axis_type == AxisDescendantOrSelf && input.node_type == xpath.AllNode {
					var qyGrandInput Query
					if input.input != nil {
						qyGrandInput = b.processNode(input.input, smartDescFlag, props)
					} else {
						qyGrandInput = &ContextQuery{}
						*props = noneProp
					}
					result = &DescendantQuery{
						AxisQuery: newAxisQuery(qyGrandInput, root.name, root.prefix, root.node_type),
						abbrAxis:  input.abbr,
					}
					if *props&nonFlatProp != 0 {
						result = &DocumentOrderQuery{}
					}
					*props |= nonFlatProp
					return result
				}
			}
			if root.axis_type == AxisDescendant || root.axis_type == AxisDescendantOrSelf {
				inputFlags |= smartDescFlag
			}
		}
		qyInput = b.processNode(root.input, inputFlags, props)
	} else {
		qyInput = &ContextQuery{}
		*props = noneProp
	}
	switch root.axis_type {
	case AxisAncestor:
		result = &AncestorQuery{AxisQuery: newAxisQuery(qyInput, root.name, root.prefix, root.node_type)}
		*props |= nonFlatProp
	case AxisAncestorOrSelf:
		result = &AncestorQuery{AxisQuery: newAxisQuery(qyInput, root.name, root.prefix, root.node_type), matchSelf: true}
		*props |= nonFlatProp
	case AxisChild:
		result = &ChildrenQuery{newAxisQuery(qyInput, root.name, root.prefix, root.node_type)}
	case AxisParent:
		result = &ParentQuery{newAxisQuery(qyInput, root.name, root.prefix, root.node_type)}
	case AxisDescendant:
		if flags&smartDescFlag != 0 {
			result = &DescendantOverDescendantQuery{AxisQuery: newAxisQuery(qyInput, root.name, root.prefix, root.node_type), matchSelf: false}
		} else {
			result = &DescendantQuery{AxisQuery: newAxisQuery(qyInput, root.name, root.prefix, root.node_type),
				matchSelf: false,
				abbrAxis:  false,
			}
			if *props&nonFlatProp != 0 {
				result = &DocumentOrderQuery{result}
			}
		}
	case AxisDescendantOrSelf:
		if flags&smartDescFlag != 0 {
			result = &DescendantOverDescendantQuery{AxisQuery: newAxisQuery(qyInput, root.name, root.prefix, root.node_type), matchSelf: true}
		} else {
			result = &DescendantQuery{AxisQuery: newAxisQuery(qyInput, root.name, root.prefix, root.node_type),
				matchSelf: true,
				abbrAxis:  root.abbr,
			}
			if *props&nonFlatProp != 0 {
				result = &DocumentOrderQuery{result}
			}
		}
	case AxisPreceding:
		result = &PrecedingQuery{newAxisQuery(qyInput, root.name, root.prefix, root.node_type)}
		*props |= nonFlatProp
	case AxisFollowing:
		result = &FollowingQuery{newAxisQuery(qyInput, root.name, root.prefix, root.node_type)}
		*props |= nonFlatProp
	case AxisFollowingSibling:
		result = &FollSiblingQuery{newAxisQuery(qyInput, root.name, root.prefix, root.node_type)}
		*props |= nonFlatProp
	case AxisPrecedingSibling:
		result = &PreSiblingQuery{newAxisQuery(qyInput, root.name, root.prefix, root.node_type)}
		*props |= nonFlatProp
	case AxisAttribute:
		result = &AttributeQuery{newAxisQuery(qyInput, root.name, root.prefix, root.node_type)}
		*props |= nonFlatProp
	case AxisSelf:
		result = &XPathSelfQuery{newAxisQuery(qyInput, root.name, root.prefix, root.node_type)}
		*props |= nonFlatProp
	case AxisNamespace:
		if (root.node_type == xpath.AllNode || root.node_type == xpath.ElementNode || root.node_type == xpath.AttributeNode) && len(root.prefix) == 0 {
			result = &XPathSelfQuery{newAxisQuery(qyInput, root.name, root.prefix, root.node_type)}
		} else {
			result = &EmptyQuery{}
		}
	default:
		panic(fmt.Sprintf("The XPath query '%s' is not supported.", b.query))
	}
	return result
}
*/
