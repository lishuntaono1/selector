package internal

type FunctionType uint

const (
	FuncLast FunctionType = iota
	FuncPosition
	FuncCount
	FuncID
	FuncLocalName
	FuncNameSpaceUri
	FuncName
	FuncString
	FuncBoolean
	FuncNumber
	FuncTrue
	FuncFalse
	FuncNot
	FuncConcat
	FuncStartsWith
	FuncContains
	FuncSubstringBefore
	FuncSubstringAfter
	FuncSubstring
	FuncStringLength
	FuncNormalize
	FuncTranslate
	FuncLang
	FuncSum
	FuncFloor
	FuncCeiling
	FuncRound
	FuncUserDefined
)

type ParamInfo struct {
	funcType FunctionType
	minargs  int
	maxargs  int
	argTypes []ResultType
}

var funcTypes = []ResultType{
	/* FunctionType.FuncLast            */ NumberType,
	/* FunctionType.FuncPosition        */ NumberType,
	/* FunctionType.FuncCount           */ NumberType,
	/* FunctionType.FuncID              */ NodeSetType,
	/* FunctionType.FuncLocalName       */ StringType,
	/* FunctionType.FuncNameSpaceUri    */ StringType,
	/* FunctionType.FuncName            */ StringType,
	/* FunctionType.FuncString          */ StringType,
	/* FunctionType.FuncBoolean         */ BooleanType,
	/* FunctionType.FuncNumber          */ NumberType,
	/* FunctionType.FuncTrue            */ BooleanType,
	/* FunctionType.FuncFalse           */ BooleanType,
	/* FunctionType.FuncNot             */ BooleanType,
	/* FunctionType.FuncConcat          */ StringType,
	/* FunctionType.FuncStartsWith      */ BooleanType,
	/* FunctionType.FuncContains        */ BooleanType,
	/* FunctionType.FuncSubstringBefore */ StringType,
	/* FunctionType.FuncSubstringAfter  */ StringType,
	/* FunctionType.FuncSubstring       */ StringType,
	/* FunctionType.FuncStringLength    */ NumberType,
	/* FunctionType.FuncNormalize       */ StringType,
	/* FunctionType.FuncTranslate       */ StringType,
	/* FunctionType.FuncLang            */ BooleanType,
	/* FunctionType.FuncSum             */ NumberType,
	/* FunctionType.FuncFloor           */ NumberType,
	/* FunctionType.FuncCeiling         */ NumberType,
	/* FunctionType.FuncRound           */ NumberType,
	/* FunctionType.FuncUserDefined     */ AnyType,
}

type Function struct {
	func_type FunctionType
	argument  []AstNode

	prefix, name string
}

func (f *Function) Type() AstType {
	return AstFunction
}

func (f *Function) ReturnType() ResultType {
	return funcTypes[f.func_type]
}
