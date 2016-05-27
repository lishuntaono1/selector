package internal

import (
	"fmt"
	"reflect"

	"github.com/zhengchun/selector/xpath"
)

type XPathParser struct {
	scanner *XPathScanner
	depth   int
}

//
func ParseXPathExpression(expr string) AstNode {
	p := newXPathParser(expr)
	if n := p.parseExpression(nil); p.scanner.kind != LexEof {
		panic(fmt.Sprintf("%s has an invalid token.", p.scanner.expr))
	} else {
		return n
	}
}

func ParseXPathPattern(pattern string) AstNode {
	p := newXPathParser(pattern)
	if n := p.parsePattern(); p.scanner.kind != LexEof {
		panic(fmt.Sprintf("%s has an invalid token.", p.scanner.expr))
	} else {
		return n
	}
}

func newXPathParser(q string) *XPathParser {
	s := &XPathScanner{expr: q}
	s.NextChar()
	s.NextLex()
	return &XPathParser{s, 0}
}

// --------------- Expression Parsing ----------------------

func (p *XPathParser) parseExpression(qyInput AstNode) AstNode {
	//ParseOrExpr->ParseAndExpr->ParseEqualityExpr->ParseRelationalExpr...->ParseFilterExpr->ParsePredicate->ParseExpression
	//So put 200 limitation here will max cause about 2000~3000 depth stack.
	if p.depth = p.depth + 1; p.depth > 200 {
		panic("The xpath query is too complex.")
	}
	n := p.parseOrExpr(qyInput)
	p.depth--
	return n
}

// >> OrExpr ::= ( OrExpr 'or' )? AndExpr
func (p *XPathParser) parseOrExpr(qyInput AstNode) AstNode {
	opnd := p.parseAndExpr(qyInput)
	for {
		if !p.testOp("or") {
			return opnd
		}
		p.nextLex()
		opnd = &Operator{OpOR, opnd, p.parseAndExpr(qyInput)}
	}
}

//>> AndExpr ::= ( AndExpr 'and' )? EqualityExpr
func (p *XPathParser) parseAndExpr(qyInput AstNode) AstNode {
	opnd := p.parseEqualityExpr(qyInput)
	for {
		if !p.testOp("and") {
			return opnd
		}
		p.nextLex()
		opnd = &Operator{OpAND, opnd, p.parseEqualityExpr(qyInput)}
	}
}

//>> EqualityOp ::= '=' | '!='
//>> EqualityExpr    ::= ( EqualityExpr EqualityOp )? RelationalExpr
func (p *XPathParser) parseEqualityExpr(qyInput AstNode) AstNode {
	opnd := p.parseRelationalExpr(qyInput)
	for {
		var op OpType = OpINVALID

		switch p.scanner.kind {
		case LexEq:
			op = OpEQ
		case LexNe:
			op = OpNE
		}
		if op == OpINVALID {
			return opnd
		}
		p.nextLex()
		opnd = &Operator{op, opnd, p.parseRelationalExpr(qyInput)}
	}
}

//>> RelationalOp ::= '<' | '>' | '<=' | '>='
//>> RelationalExpr    ::= ( RelationalExpr RelationalOp )? AdditiveExpr
func (p *XPathParser) parseRelationalExpr(qyInput AstNode) AstNode {
	opnd := p.parseAdditiveExpr(qyInput)
	for {
		var op OpType = OpINVALID

		switch p.scanner.kind {
		case LexLt:
			op = OpLT
		case LexLe:
			op = OpLE
		case LexGt:
			op = OpGT
		case LexGe:
			op = OpGE
		}
		if op == OpINVALID {
			return opnd
		}
		p.nextLex()
		opnd = &Operator{op, opnd, p.parseAdditiveExpr(qyInput)}
	}
}

//>> AdditiveOp   ::= '+' | '-'
//>> AdditiveExpr ::= ( AdditiveExpr AdditiveOp )? MultiplicativeExpr
func (p *XPathParser) parseAdditiveExpr(qyInput AstNode) AstNode {
	opnd := p.parseMultiplicativeExpr(qyInput)
	for {
		var op OpType = OpINVALID

		if p.scanner.kind == LexPlus {
			op = OpPLUS
		} else if p.scanner.kind == LexMinus {
			op = OpMINUS
		}
		if op == OpINVALID {
			return opnd
		}
		p.nextLex()
		opnd = &Operator{op, opnd, p.parseMultiplicativeExpr(qyInput)}
	}
}

//>> MultiplicativeOp   ::= '*' | 'div' | 'mod'
//>> MultiplicativeExpr ::= ( MultiplicativeExpr MultiplicativeOp )? UnaryExpr
func (p *XPathParser) parseMultiplicativeExpr(qyInput AstNode) AstNode {
	opnd := p.parseUnaryExpr(qyInput)
	for {
		var op OpType = OpINVALID
		if p.scanner.kind == LexStar {
			op = OpMUL
		} else {
			if p.testOp("div") {
				op = OpDIV
			} else if p.testOp("mod") {
				op = OpMOD
			}
		}
		if op == OpINVALID {
			return opnd
		}
		p.nextLex()
		opnd = &Operator{op, opnd, p.parseUnaryExpr(qyInput)}
	}
}

//>> UnaryExpr    ::= UnionExpr | '-' UnaryExpr
func (p *XPathParser) parseUnaryExpr(qyInput AstNode) AstNode {
	minus := false
	for p.scanner.kind == LexMinus {
		p.nextLex()
		minus = !minus
	}
	if minus {
		return &Operator{OpMINUS, p.parseUnionExpr(qyInput), &Operand{NumberType, -1}}
	} else {
		return p.parseUnionExpr(qyInput)
	}
}

//>> UnionExpr ::= ( UnionExpr '|' )? PathExpr
func (p *XPathParser) parseUnionExpr(qyInput AstNode) AstNode {
	opnd := p.parsePathExpr(qyInput)
	for {
		if p.scanner.kind != LexUnion {
			return opnd
		}
		p.nextLex()
		opnd2 := p.parsePathExpr(qyInput)
		checkNodeSet(opnd.ReturnType())
		checkNodeSet(opnd2.ReturnType())
		opnd = &Operator{OpUNION, opnd, opnd2}
	}
}

//>> PathOp   ::= '/' | '//'
//>> PathExpr ::= LocationPath |
//>>              FilterExpr ( PathOp  RelativeLocationPath )?
func (p *XPathParser) parsePathExpr(qyInput AstNode) AstNode {
	var opnd AstNode
	if isPrimaryExpr(p.scanner) {
		opnd = p.parseFilterExpr(qyInput)
		if p.scanner.kind == LexSlash {
			p.nextLex()
			opnd = p.parseRelativeLocationPath(opnd)
		} else if p.scanner.kind == LexSlashSlash {
			p.nextLex()
			opnd = p.parseRelativeLocationPath(&Axis{
				abbr:      true,
				input:     opnd,
				axis_type: AxisDescendantOrSelf,
				node_type: xpath.AllNode,
			})
		}
	} else {
		opnd = p.parseLocationPath(nil)
	}
	return opnd
}

//>> FilterExpr ::= PrimaryExpr | FilterExpr Predicate
func (p *XPathParser) parseFilterExpr(qyInput AstNode) AstNode {
	opnd := p.parsePrimaryExpr(qyInput)
	for p.scanner.kind == LexLBracket {
		// opnd must be a query
		opnd = &Filter{opnd, p.parsePredicate(opnd)}
	}
	return opnd
}

//>> Predicate ::= '[' Expr ']'
func (p *XPathParser) parsePredicate(qyInput AstNode) AstNode {
	var opnd AstNode
	// we have predicates. Check that input type is NodeSet
	checkNodeSet(qyInput.ReturnType())
	p.passToken(LexLBracket)
	opnd = p.parseExpression(qyInput)
	p.passToken(LexRBracket)
	return opnd
}

//>> LocationPath ::= RelativeLocationPath | AbsoluteLocationPath
func (p *XPathParser) parseLocationPath(qyInput AstNode) AstNode {
	if p.scanner.kind == LexSlash {
		p.nextLex()
		var opnd AstNode
		opnd = &Root{}
		if isStep(p.scanner.kind) {
			opnd = p.parseRelativeLocationPath(opnd)
		}
		return opnd
	} else if p.scanner.kind == LexSlashSlash {
		p.nextLex()
		return p.parseRelativeLocationPath(&Axis{
			abbr:      true,
			input:     &Root{},
			axis_type: AxisDescendantOrSelf,
			node_type: xpath.AllNode,
		})
	} else {
		return p.parseRelativeLocationPath(qyInput)
	}
}

//>> PathOp   ::= '/' | '//'
//>> RelativeLocationPath ::= ( RelativeLocationPath PathOp )? Step
func (p *XPathParser) parseRelativeLocationPath(qyInput AstNode) AstNode {
	opnd := qyInput
	for {
		opnd = p.parseStep(opnd)
		if p.scanner.kind == LexSlashSlash {
			p.nextLex()
			opnd = &Axis{
				abbr:      true,
				input:     opnd,
				axis_type: AxisDescendantOrSelf,
				node_type: xpath.AllNode,
			}
		} else if p.scanner.kind == LexSlash {
			p.nextLex()
		} else {
			break
		}
	}
	return opnd
}

//>> Step ::= '.' | '..' | ( AxisName '::' | '@' )? NodeTest Predicate*
func (p *XPathParser) parseStep(qyInput AstNode) AstNode {
	var opnd AstNode
	if p.scanner.kind == LexDot { //>> '.'
		p.nextLex()
		opnd = &Axis{
			abbr:      true,
			input:     qyInput,
			axis_type: AxisSelf,
			node_type: xpath.AllNode,
		}
	} else if p.scanner.kind == LexDotDot { //>> '..'
		p.nextLex()
		opnd = &Axis{
			abbr:      true,
			input:     qyInput,
			axis_type: AxisParent,
			node_type: xpath.AllNode,
		}
	} else { //>> ( AxisName '::' | '@' )? NodeTest Predicate*
		axisType := AxisChild
		switch p.scanner.kind {
		case LexAt: //>> '@'
			axisType = AxisAttribute
			p.nextLex()
		case LexAxe: //>> AxisName '::'
			axisType = p.getAxis()
			p.nextLex()
		}
		nodeType := xpath.ElementNode
		if axisType == AxisAttribute {
			nodeType = xpath.AttributeNode
		}

		opnd = p.parseNodeTest(qyInput, axisType, nodeType)
		for p.scanner.kind == LexLBracket {
			opnd = &Filter{opnd, p.parsePredicate(opnd)}
		}
	}
	return opnd
}

//>> NodeTest ::= NameTest | 'comment ()' | 'text ()' | 'node ()' | 'processing-instruction ('  Literal ? ')'
func (p *XPathParser) parseNodeTest(qyInput AstNode, axisType AxisType, nodeType xpath.NodeType) AstNode {
	var nodeName, nodePrefix string
	switch p.scanner.kind {
	case LexName:
		{
			if p.scanner.canBeFunction && isNodeType(p.scanner) {
				nodePrefix = ""
				nodeName = ""
				switch p.scanner.name {
				case "comment":
					nodeType = xpath.CommentNode
				case "text":
					nodeType = xpath.TextNode
				case "node":
					nodeType = xpath.AllNode
				case "processing-instruction":
					nodeType = xpath.ProcessingInstructionNode
				default:
					nodeType = xpath.RootNode
				}
				p.nextLex()
				p.passToken(LexLParens)
				if nodeType == xpath.ProcessingInstructionNode {
					if p.scanner.kind != LexRParens {
						checkToken(p.scanner, LexString)
						nodeName = p.scanner.strval
						p.nextLex()
					}
				}
				p.passToken(LexRParens)
			} else {
				nodePrefix = p.scanner.prefix
				nodeName = p.scanner.name
				p.nextLex()
				if nodeName == "*" {
					nodeName = ""
				}
			}
		}
	case LexStar:
		nodePrefix = ""
		nodeName = ""
		p.nextLex()
	default:
		panic("Expression must evaluate to a node-set.")
	}
	return &Axis{
		input:     qyInput,
		prefix:    nodePrefix,
		name:      nodeName,
		axis_type: axisType,
		node_type: nodeType,
	}
}

//>> PrimaryExpr ::= Literal | Number | VariableReference | '(' Expr ')' | FunctionCall
func (p *XPathParser) parsePrimaryExpr(qyInput AstNode) AstNode {
	var opnd AstNode
	switch p.scanner.kind {
	case LexString:
		opnd = &Operand{StringType, p.scanner.strval}
		p.nextLex()
	case LexNumber:
		opnd = &Operand{NumberType, p.scanner.numval}
		p.nextLex()
	case LexDollar:
		p.nextLex()
		checkToken(p.scanner, LexName)
		opnd = &Variable{p.scanner.name, p.scanner.prefix}
		p.nextLex()
	case LexLParens:
		p.nextLex()
		opnd = p.parseExpression(qyInput)
		if opnd.Type() != ConstantOperandAst {
			opnd = &Group{opnd}
		}
		p.passToken(LexRParens)
	case LexName:
		if p.scanner.canBeFunction && !isNodeType(p.scanner) {
			opnd = p.parseMethod(nil)
		}
	}
	return opnd
}

func (p *XPathParser) parseMethod(qyInput AstNode) AstNode {
	argList := make([]AstNode, 0)
	name := p.scanner.name
	prefix := p.scanner.prefix
	p.passToken(LexName)
	p.passToken(LexLParens)
	if p.scanner.kind != LexRParens {
		for {
			argList = append(argList, p.parseExpression(qyInput))
			if p.scanner.kind == LexRParens {
				break
			}
			p.passToken(LexComma)
		}
	}
	p.passToken(LexRParens)
	if prefix == "" {
		if pi, ok := functionTable[name]; ok {
			argCount := len(argList)
			if argCount < pi.minargs {
				panic(fmt.Sprintf("Function %s in %s has an invalid number of arguments.", name, p.scanner.expr))
			}
			if pi.functype == FuncConcat {
				for i, arg := range argList {
					if arg.ReturnType() != StringType {
						arg = &Function{functype: FuncString, argument: []AstNode{arg}}
					}
					argList[i] = arg
				}
			} else {
				if pi.maxargs < argCount {
					panic(fmt.Sprintf("Function %s in %s has an invalid number of arguments.", name, p.scanner.expr))
				}
				if len(pi.argTypes) < argCount {
					argCount = len(pi.argTypes) // argument we have the type specified (can be < pi.Minargs)
				}
				for i := 0; i < argCount; i++ {
					arg := argList[i]
					if pi.argTypes[i] != AnyType && pi.argTypes[i] != arg.ReturnType() {
						switch pi.argTypes[i] {
						case NodeSetType:
							if !(reflect.TypeOf(arg) == reflect.TypeOf((*Variable)(nil))) &&
								!(reflect.TypeOf(arg) == reflect.TypeOf((*Function)(nil)) && arg.ReturnType() == AnyType) {
								panic(fmt.Sprintf("The argument to function %s in %s cannot be converted to a node-set.", name, p.scanner.expr))
							}
						case StringType:
							arg = &Function{functype: FuncString, argument: []AstNode{arg}}
						case NumberType:
							arg = &Function{functype: FuncNumber, argument: []AstNode{arg}}
						case BooleanType:
							arg = &Function{functype: FuncBoolean, argument: []AstNode{arg}}
						}
						argList[i] = arg
					}
				}
			}
			return &Function{functype: pi.functype, argument: argList}
		}
	}
	return &Function{functype: FuncUserDefined, argument: []AstNode{}, prefix: prefix, name: name}
}

// --------------- Pattern Parsing ----------------------

//>> Pattern ::= ( Pattern '|' )? LocationPathPattern
func (p *XPathParser) parsePattern() AstNode {
	opnd := p.parseLocationPathPattern()
	for {
		if p.scanner.kind != LexUnion {
			return opnd
		}
		p.nextLex()
		opnd = &Operator{OpUNION, opnd, p.parseLocationPathPattern()}
	}
}

//>> LocationPathPattern ::= '/' | RelativePathPattern | '//' RelativePathPattern  |  '/' RelativePathPattern
//>>                       | IdKeyPattern (('/' | '//') RelativePathPattern)?
func (p *XPathParser) parseLocationPathPattern() AstNode {
	var opnd AstNode
	switch p.scanner.kind {
	case LexSlash:
		p.nextLex()
		opnd = &Root{}
		if p.scanner.kind == LexEof || p.scanner.kind == LexUnion {
			return opnd
		}
	case LexSlashSlash:
		p.nextLex()
		opnd = &Axis{
			input:     &Root{},
			abbr:      true,
			axis_type: AxisDescendantOrSelf,
			node_type: xpath.AllNode,
		}
	case LexName:
		if p.scanner.canBeFunction {
			opnd = p.parseIdKeyPattern()
			if opnd != nil {
				switch p.scanner.kind {
				case LexSlash:
					p.nextLex()
				case LexSlashSlash:
					p.nextLex()
					opnd = &Axis{
						abbr:      true,
						input:     opnd,
						axis_type: AxisDescendantOrSelf,
						node_type: xpath.AllNode,
					}
				default:
					return opnd
				}
			}
		}
	}
	return p.parseRelativePathPattern(opnd)
}

//>> IdKeyPattern ::= 'id' '(' Literal ')' | 'key' '(' Literal ',' Literal ')'
func (p *XPathParser) parseIdKeyPattern() AstNode {
	argList := make([]AstNode, 0)
	if len(p.scanner.prefix) == 0 {
		if p.scanner.name == "id" {
			pi, _ := functionTable["id"]
			p.nextLex()
			p.passToken(LexLParens)
			checkToken(p.scanner, LexString)
			argList = append(argList, &Operand{StringType, p.scanner.strval})
			p.nextLex()
			p.passToken(LexRParens)
			return &Function{functype: pi.functype, argument: argList}
		}
		if p.scanner.name == "key" {
			p.nextLex()
			p.passToken(LexLParens)
			checkToken(p.scanner, LexString)
			argList = append(argList, &Operand{StringType, p.scanner.strval})
			p.nextLex()
			p.passToken(LexComma)
			checkToken(p.scanner, LexString)
			argList = append(argList, &Operand{StringType, p.scanner.strval})
			p.nextLex()
			p.passToken(LexRParens)
			return &Function{functype: FuncUserDefined, argument: argList, name: "key"}
		}
	}
	return nil
}

//>> PathOp   ::= '/' | '//'
//>> RelativePathPattern ::= ( RelativePathPattern PathOp )? StepPattern
func (p *XPathParser) parseRelativePathPattern(qyInput AstNode) AstNode {
	opnd := p.parseStepPattern(qyInput)
	if p.scanner.kind == LexSlashSlash {
		p.nextLex()
		opnd = p.parseRelativePathPattern(&Axis{
			abbr:      true,
			input:     opnd,
			node_type: xpath.AllNode,
			axis_type: AxisDescendantOrSelf,
		})
	} else if p.scanner.kind == LexSlash {
		p.nextLex()
		opnd = p.parseRelativePathPattern(opnd)
	}
	return opnd
}

//>> StepPattern    ::=    ChildOrAttributeAxisSpecifier NodeTest Predicate*
//>> ChildOrAttributeAxisSpecifier    ::=    @ ? | ('child' | 'attribute') '::'
func (p *XPathParser) parseStepPattern(qyInput AstNode) AstNode {
	var opnd AstNode
	axisType := AxisChild
	switch p.scanner.kind {
	case LexAt: //>> '@'
		axisType = AxisAttribute
		p.nextLex()
	case LexAxe: //>> AxisName '::'
		axisType = p.getAxis()
		if axisType != AxisChild && axisType != AxisAttribute {
			panic(fmt.Sprintf("%s has an invalid token.", p.scanner.expr))
		}
		p.nextLex()
	}
	nodeType := xpath.ElementNode
	if axisType == AxisAttribute {
		nodeType = xpath.AttributeNode
	}

	opnd = p.parseNodeTest(qyInput, axisType, nodeType)

	for p.scanner.kind == LexLBracket {
		opnd = &Filter{opnd, p.parsePredicate(opnd)}
	}
	return opnd
}

func (p *XPathParser) getAxis() AxisType {
	t, ok := axesTable[p.scanner.name]
	if !ok {
		panic(fmt.Sprintf("%s has an invalid token.", p.scanner.expr))
	}
	return t
}

func (p *XPathParser) testOp(op string) bool {
	return p.scanner.kind == LexName && len(p.scanner.prefix) == 0 && p.scanner.name == op
}

func (p *XPathParser) nextLex() {
	p.scanner.NextLex()
}

func (p *XPathParser) passToken(t LexKind) {
	checkToken(p.scanner, t)
	p.nextLex()
}

func checkNodeSet(t ResultType) {
	if t != NodeSetType && t != AnyType {
		panic("Expression must evaluate to a node-set.")
	}
}

func isPrimaryExpr(s *XPathScanner) bool {
	return s.kind == LexString ||
		s.kind == LexNumber ||
		s.kind == LexDollar ||
		s.kind == LexLParens ||
		s.kind == LexName && (s.canBeFunction && !isNodeType(s))
}

func isNodeType(s *XPathScanner) bool {
	return s.prefix == "" && (s.name == "node" ||
		s.name == "text" ||
		s.name == "processing-instruction" ||
		s.name == "comment")
}

func isStep(kind LexKind) bool {
	return kind == LexDot ||
		kind == LexDotDot ||
		kind == LexAt ||
		kind == LexAxe ||
		kind == LexStar ||
		kind == LexName
}

func checkToken(s *XPathScanner, t LexKind) {
	if s.kind != t {
		panic(fmt.Sprintf("%s has an invalid token.", s.expr))
	}
}

var (
	temparray1 []ResultType = []ResultType{}
	temparray2 []ResultType = []ResultType{NodeSetType}
	temparray3 []ResultType = []ResultType{AnyType}
	temparray4 []ResultType = []ResultType{StringType}
	temparray5 []ResultType = []ResultType{StringType, StringType}
	temparray6 []ResultType = []ResultType{StringType, NumberType, NumberType}
	temparray7 []ResultType = []ResultType{StringType, StringType, StringType}
	temparray8 []ResultType = []ResultType{BooleanType}
	temparray9 []ResultType = []ResultType{NumberType}
)

var functionTable = map[string]*ParamInfo{
	"last":             &ParamInfo{FuncLast, 0, 0, temparray1},
	"position":         &ParamInfo{FuncPosition, 0, 0, temparray1},
	"name":             &ParamInfo{FuncName, 0, 1, temparray2},
	"namespace-uri":    &ParamInfo{FuncNameSpaceUri, 0, 1, temparray2},
	"local-name":       &ParamInfo{FuncLocalName, 0, 1, temparray2},
	"count":            &ParamInfo{FuncCount, 1, 1, temparray2},
	"id":               &ParamInfo{FuncID, 1, 1, temparray3},
	"string":           &ParamInfo{FuncString, 0, 1, temparray3},
	"concat":           &ParamInfo{FuncConcat, 2, 100, temparray4},
	"starts-with":      &ParamInfo{FuncStartsWith, 2, 2, temparray5},
	"contains":         &ParamInfo{FuncContains, 2, 2, temparray5},
	"substring-before": &ParamInfo{FuncSubstringBefore, 2, 2, temparray5},
	"substring-after":  &ParamInfo{FuncSubstringAfter, 2, 2, temparray5},
	"substring":        &ParamInfo{FuncSubstring, 2, 3, temparray6},
	"string-length":    &ParamInfo{FuncStringLength, 0, 1, temparray4},
	"normalize-space":  &ParamInfo{FuncNormalize, 0, 1, temparray4},
	"translate":        &ParamInfo{FuncTranslate, 3, 3, temparray7},
	"boolean":          &ParamInfo{FuncBoolean, 1, 1, temparray3},
	"not":              &ParamInfo{FuncNot, 1, 1, temparray8},
	"true":             &ParamInfo{FuncTrue, 0, 0, temparray8},
	"false":            &ParamInfo{FuncFalse, 0, 0, temparray8},
	"lang":             &ParamInfo{FuncLang, 1, 1, temparray4},
	"number":           &ParamInfo{FuncNumber, 0, 1, temparray3},
	"sum":              &ParamInfo{FuncSum, 1, 1, temparray2},
	"floor":            &ParamInfo{FuncFloor, 1, 1, temparray9},
	"ceiling":          &ParamInfo{FuncCeiling, 1, 1, temparray9},
	"round":            &ParamInfo{FuncRound, 1, 1, temparray9},
}

var axesTable map[string]AxisType = map[string]AxisType{
	"ancestor":           AxisAncestor,
	"estor-or-self":      AxisAncestorOrSelf,
	"attribute":          AxisAttribute,
	"child":              AxisChild,
	"descendant":         AxisDescendant,
	"descendant-or-self": AxisDescendantOrSelf,
	"following":          AxisFollowing,
	"following-sibling":  AxisFollowingSibling,
	"namespace":          AxisNamespace,
	"parent":             AxisParent,
	"preceding":          AxisPreceding,
	"preceding-sibling":  AxisPrecedingSibling,
	"self":               AxisSelf,
}
