package main

type Stmt interface {
	Run(lox *Lox)
}

/*----------  Print Stmt  ----------*/

type StmtPrint struct {
	expr Expr
}

func NewStmtPrint(expr Expr) *StmtPrint {
	return &StmtPrint{expr}
}

/*----------  Expression Stmt  ----------*/

type StmtExpression struct {
	expr Expr
}

func NewStmtExpression(expr Expr) *StmtExpression {
	return &StmtExpression{expr}
}

/*----------  Var Decl Stmt  ----------*/
type StmtVarDecl struct {
	name  *Token
	value Expr
}

func NewStmtVarDecl(name *Token, value Expr) *StmtVarDecl {
	return &StmtVarDecl{name, value}
}