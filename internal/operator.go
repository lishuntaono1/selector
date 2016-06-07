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

var invertOp = []OpType{
	/*INVALID*/ OpINVALID,
	/*OR     */ OpINVALID,
	/*END    */ OpINVALID,
	/*EQ     */ OpEQ,
	/*NE     */ OpNE,
	/*LT     */ OpGT,
	/*LE     */ OpGE,
	/*GT     */ OpLT,
	/*GE     */ OpLE,
}

func invertOperator(op OpType) OpType {
	return invertOp[op]
}

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
