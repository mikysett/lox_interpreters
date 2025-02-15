package main

type Stmt interface {
	accept(StmtVisitor) error
}

type StmtVisitor interface {
	visitExpressionStmt(*StmtExpression) error
	visitPrintStmt(*StmtPrint) error
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

func (expr *StmtPrint) accept(v StmtVisitor) error {
	return v.visitPrintStmt(expr)
}
