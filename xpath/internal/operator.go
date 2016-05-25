package internal

type OpType uint

const (
	// order is aligned with XPathOperator
	OpINVALID OpType = iota
	/*Logical   */
	OpOR
	OpAND
	/*Equality  */
	OpEQ
	OpNE
	/*Relational*/
	OpLT
	OpLE
	OpGT
	OpGE
	/*Arithmetic*/
	OpPLUS
	OpMINUS
	OpMUL
	OpDIV
	OpMOD
	/*Union     */
	OpUNION
)

type Operator struct {
	op    OpType
	opnd1 AstNode
	opnd2 AstNode
}

func (o *Operator) Type() AstNodeType {
	return OperatorAst
}

func (o *Operator) ReturnType() ResultType {
	if o.op <= OpGE {
		return BooleanType
	}
	if o.op <= OpMOD {
		return NumberType
	}
	return NodeSetType
}
