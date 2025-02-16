package main

type Stmt interface {
	accept(StmtVisitor) error
}

type StmtVisitor interface {
	visitExpressionStmt(*StmtExpression) error
	visitIfStmt(*StmtIf) error
	visitPrintStmt(*StmtPrint) error
	visitVarStmt(*StmtVar) error
	visitBlockStmt(*StmtBlock) error
}

// Block      : List<Stmt> statements
type StmtBlock struct {
	block []Stmt
}

func NewStmtBlock(block []Stmt) *StmtBlock {
	return &StmtBlock{
		block: block,
	}
}

func (stmt *StmtBlock) accept(v StmtVisitor) error {
	return v.visitBlockStmt(stmt)
}

// If         : Expr condition, Stmt thenBranch, Stmt elseBranch,
type StmtIf struct {
	condition  Expr
	thenBranch Stmt
	elseBranch Stmt
}

func NewStmtIf(condition Expr, thenBranch Stmt, elseBranch Stmt) *StmtIf {
	return &StmtIf{
		condition:  condition,
		thenBranch: thenBranch,
		elseBranch: elseBranch,
	}
}

func (stmt *StmtIf) accept(v StmtVisitor) error {
	return v.visitIfStmt(stmt)
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
