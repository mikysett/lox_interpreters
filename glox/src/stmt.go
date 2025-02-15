package main

type Stmt interface {
	accept(StmtVisitor) (any, error)
}

type StmtVisitor interface {
	visitExpressionStmt(*StmtExpression) (any, error)
	visitPrintStmt(*StmtPrint) (any, error)
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

func (expr *StmtExpression) accept(v StmtVisitor) (any, error) {
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

func (expr *StmtPrint) accept(v StmtVisitor) (any, error) {
	return v.visitPrintStmt(expr)
}
