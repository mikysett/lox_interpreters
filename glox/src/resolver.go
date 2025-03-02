package main

type FunctionType int

const (
	None FunctionType = iota
	Func
)

type Resolver struct {
	interpreter     *Interpreter
	scopes          *Scopes
	currentFunction FunctionType
}

type Scopes []Scope

type Scope map[string]*LocalVariable

type LocalVariable struct {
	declaration   *Token
	isInitialized bool
	isUsed        bool
	// TODO: Implement index to speed lookup in interpreter
	// scopedIndex   int
}

func NewScopes() *Scopes {
	return &Scopes{}
}

func NewScope() Scope {
	return Scope{}
}

func NewLocalVariable(declaration *Token) *LocalVariable {
	return &LocalVariable{
		declaration: declaration,
	}
}

func (s *Scopes) push() {
	*s = append(*s, NewScope())
}

func (s *Scopes) pop() {
	*s = (*s)[:len(*s)-1]
}

func (s *Scopes) peek() Scope {
	return (*s)[len(*s)-1]
}

func (s *Scopes) isEmpty() bool {
	return len(*s) == 0
}

func NewResolver(interpreter *Interpreter) *Resolver {
	return &Resolver{
		interpreter:     interpreter,
		scopes:          NewScopes(),
		currentFunction: None,
	}
}

func (resolver *Resolver) resolveStmts(stmts []Stmt) error {
	for _, stmt := range stmts {
		resolver.resolveStmt(stmt)
	}
	return nil
}

func (resolver *Resolver) resolveStmt(stmt Stmt) error {
	return stmt.accept(resolver)
}

func (resolver *Resolver) resolveExpr(expr Expr) error {
	expr.accept(resolver)
	return nil
}

func (resolver *Resolver) visitBlockStmt(stmt *StmtBlock) error {
	resolver.beginScope()
	resolver.resolveStmts(stmt.block)
	resolver.endScope()
	return nil
}

func (resolver *Resolver) beginScope() {
	resolver.scopes.push()
}

func (resolver *Resolver) endScope() {
	for _, localVar := range resolver.scopes.peek() {
		if !localVar.isUsed {
			printError(localVar.declaration, "Variable declared but never read")
		}
	}
	resolver.scopes.pop()
}

func (resolver *Resolver) declare(name *Token) {
	if resolver.scopes.isEmpty() {
		return
	}
	scope := resolver.scopes.peek()
	if _, ok := scope[name.Lexeme]; ok {
		printError(name, "Already a variable with this name in this scope.")
	}

	scope[name.Lexeme] = NewLocalVariable(name)
}

func (resolver *Resolver) define(name *Token) {
	if resolver.scopes.isEmpty() {
		return
	}
	scope := resolver.scopes.peek()
	scope[name.Lexeme].isInitialized = true
}

func (resolver *Resolver) visitBreakStmt(stmt *StmtBreak) error {
	return nil
}

func (resolver *Resolver) visitVarStmt(stmt *StmtVar) (err error) {
	resolver.declare(stmt.name)
	if stmt.initializer != nil {
		resolver.resolveExpr(stmt.initializer)
	}
	resolver.define(stmt.name)
	return nil
}

func (resolver *Resolver) visitAssignExpr(expr *ExprAssign) (any, error) {
	resolver.resolveExpr(expr.value)
	resolver.resolveLocal(expr, expr.name, false)
	return nil, nil
}

func (resolver *Resolver) visitExpressionStmt(stmt *StmtExpression) error {
	return resolver.resolveExpr(stmt.expression)
}

func (resolver *Resolver) visitFunctionStmt(stmt *StmtFunction) error {
	resolver.declare(stmt.name)
	resolver.define(stmt.name)
	resolver.resolveFunction(stmt.function, Func)
	return nil
}

func (resolver *Resolver) visitIfStmt(stmt *StmtIf) (err error) {
	resolver.resolveExpr(stmt.condition)
	resolver.resolveStmt(stmt.thenBranch)
	if stmt.elseBranch != nil {
		resolver.resolveStmt(stmt.elseBranch)
	}
	return nil
}

func (resolver *Resolver) visitWhileStmt(stmt *StmtWhile) (err error) {
	resolver.resolveExpr(stmt.condition)
	resolver.resolveStmt(stmt.body)
	return nil
}

func (resolver *Resolver) visitPrintStmt(stmt *StmtPrint) error {
	return resolver.resolveExpr(stmt.expression)
}

func (resolver *Resolver) visitReturnStmt(stmt *StmtReturn) (err error) {
	if resolver.currentFunction == None {
		printError(stmt.keyword, "Can't return from top-level code.")
	}
	if stmt.expression != nil {
		resolver.resolveExpr(stmt.expression)
	}
	return nil
}

func (resolver *Resolver) visitBinaryExpr(expr *ExprBinary) (any, error) {
	resolver.resolveExpr(expr.left)
	resolver.resolveExpr(expr.right)
	return nil, nil
}

func (resolver *Resolver) visitFunctionExpr(expr *ExprFunction) (any, error) {
	resolver.resolveFunction(expr, Func)
	return nil, nil
}

func (resolver *Resolver) visitCallExpr(expr *ExprCall) (any, error) {
	resolver.resolveExpr(expr.callee)

	for _, argumentExpr := range expr.arguments {
		resolver.resolveExpr(argumentExpr)
	}
	return nil, nil
}

func (resolver *Resolver) visitTernaryExpr(expr *ExprTernary) (any, error) {
	resolver.resolveExpr(expr.condition)
	resolver.resolveExpr(expr.left)
	resolver.resolveExpr(expr.right)
	return nil, nil
}

func (resolver *Resolver) visitLogicalExpr(expr *ExprLogical) (any, error) {
	resolver.resolveExpr(expr.left)
	resolver.resolveExpr(expr.right)
	return nil, nil
}

func (resolver *Resolver) visitGroupingExpr(expr *ExprGrouping) (any, error) {
	resolver.resolveExpr(expr.expression)
	return nil, nil
}

func (resolver *Resolver) visitLiteralExpr(expr *ExprLiteral) (any, error) {
	return nil, nil
}

func (resolver *Resolver) visitUnaryExpr(expr *ExprUnary) (any, error) {
	resolver.resolveExpr(expr.right)
	return nil, nil
}

func (resolver *Resolver) visitVariableExpr(expr *ExprVariable) (any, error) {
	if !resolver.scopes.isEmpty() {
		if localVar, ok := resolver.scopes.peek()[expr.name.Lexeme]; ok && !localVar.isInitialized {
			printError(expr.name, "Can't read local variable in its own initializer.")
		}
	}
	resolver.resolveLocal(expr, expr.name, true)
	return nil, nil
}

func (resolver *Resolver) resolveLocal(expr Expr, name *Token, isRead bool) {
	for i := len(*resolver.scopes) - 1; i >= 0; i-- {
		if localVar, ok := (*resolver.scopes)[i][name.Lexeme]; ok {
			if isRead {
				localVar.isUsed = true
			}
			resolver.interpreter.resolve(expr, len(*resolver.scopes)-1-i)
			return
		}
	}
}

func (resolver *Resolver) resolveFunction(expr *ExprFunction, funcType FunctionType) {
	enclosingFunction := resolver.currentFunction
	resolver.currentFunction = funcType

	resolver.beginScope()
	for _, param := range expr.params {
		resolver.declare(param)
		resolver.define(param)
	}
	resolver.resolveStmts(expr.body)
	resolver.endScope()

	resolver.currentFunction = enclosingFunction
}
