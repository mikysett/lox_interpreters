package main

type Stmt interface {
	accept(StmtVisitor) error
}

type StmtVisitor interface {
	visitExpressionStmt(*StmtExpression) error
	visitFunctionStmt(*StmtFunction) error
	visitIfStmt(*StmtIf) error
	visitPrintStmt(*StmtPrint) error
	visitReturnStmt(*StmtReturn) error
	visitVarStmt(*StmtVar) error
	visitBlockStmt(*StmtBlock) error
	visitWhileStmt(*StmtWhile) error
	visitBreakStmt(*StmtBreak) error
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

// Function   : Token name, List<Token> params, List<Stmt> body
type StmtFunction struct {
	name   *Token
	params []*Token
	body   []Stmt
}

func NewStmtFunction(name *Token, params []*Token, body []Stmt) *StmtFunction {
	return &StmtFunction{
		name:   name,
		params: params,
		body:   body,
	}
}

func (stmt *StmtFunction) accept(v StmtVisitor) error {
	return v.visitFunctionStmt(stmt)
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

func (expr *StmtPrint) accept(v StmtVisitor) error {
	return v.visitPrintStmt(expr)
}

// Return     : Token keyword, Expr value
type StmtReturn struct {
	keyword    *Token
	expression Expr
}

func NewStmtReturn(keyword *Token, expression Expr) *StmtReturn {
	return &StmtReturn{
		keyword:    keyword,
		expression: expression,
	}
}

func (expr *StmtReturn) accept(v StmtVisitor) error {
	return v.visitReturnStmt(expr)
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

// While      : Expr condition, Stmt body
type StmtWhile struct {
	condition Expr
	body      Stmt
}

func NewStmtWhile(condition Expr, body Stmt) *StmtWhile {
	return &StmtWhile{
		condition: condition,
		body:      body,
	}
}

func (stmt *StmtWhile) accept(v StmtVisitor) error {
	return v.visitWhileStmt(stmt)
}

// Break      :
type StmtBreak struct{}

func NewStmtBreak() *StmtBreak {
	return &StmtBreak{}
}

func (stmt *StmtBreak) accept(v StmtVisitor) error {
	return v.visitBreakStmt(stmt)
}
