package main

type Stmt interface {
	accept(StmtVisitor) error
}

type StmtVisitor interface {
	visitExpressionStmt(*StmtExpression) error
	visitPrintStmt(*StmtPrint) error
	visitVarStmt(*StmtVar) error
}

// Expression : Expr expression
type StmtExpression struct {
	expression Expr
}

func NewStmtExpression(expression Expr) *StmtExpression {
	return &StmtExpression{
		expression: expression,
	}
}

func (expr *StmtExpression) accept(v StmtVisitor) error {
	return v.visitExpressionStmt(expr)
}

// Print      : Expr expression
type StmtPrint struct {
	expression Expr
}

func NewStmtPrint(expression Expr) *StmtPrint {
	return &StmtPrint{
		expression: expression,
	}
}

// Var        : Token name, Expr initializer
type StmtVar struct {
	name        *Token
	initializer Expr
}

func NewStmtVar(name *Token, initializer Expr) *StmtVar {
	return &StmtVar{
		name:        name,
		initializer: initializer,
	}
}

func (stmt *StmtVar) accept(v StmtVisitor) error {
	return v.visitVarStmt(stmt)
}

func (expr *StmtPrint) accept(v StmtVisitor) error {
	return v.visitPrintStmt(expr)
}
