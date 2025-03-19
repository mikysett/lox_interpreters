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
	visitClassStmt(*StmtClass) error
	visitWhileStmt(*StmtWhile) error
	visitBreakStmt(*StmtBreak) error
	visitContinueStmt(*StmtContinue) error
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

// Class      : Token name, ExprVariable supercleass, List<StmtFunction> methods, List<StmtFunction> staticMethods
type StmtClass struct {
	name          *Token
	superclass    *ExprVariable
	methods       []*StmtFunction
	staticMethods []*StmtFunction
}

func NewStmtClass(name *Token, superclass *ExprVariable, methods, staticMethods []*StmtFunction) *StmtClass {
	return &StmtClass{
		name:          name,
		superclass:    superclass,
		methods:       methods,
		staticMethods: staticMethods,
	}
}

func (stmt *StmtClass) accept(v StmtVisitor) error {
	return v.visitClassStmt(stmt)
}

// Function   : Token name, ExprFunction body
type StmtFunction struct {
	name     *Token
	function *ExprFunction
}

func NewStmtFunction(name *Token, function *ExprFunction) *StmtFunction {
	return &StmtFunction{
		name:     name,
		function: function,
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

// Continue      :
type StmtContinue struct{}

func NewStmtContinue() *StmtContinue {
	return &StmtContinue{}
}

func (stmt *StmtContinue) accept(v StmtVisitor) error {
	return v.visitContinueStmt(stmt)
}
