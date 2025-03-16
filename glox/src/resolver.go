package main

type FunctionType int

const (
	FunctionTypeNone FunctionType = iota
	FunctionTypeFunc
	FunctionTypeInitializer
	FunctionTypeMethod
)

type ClassType int

const (
	ClassTypeNone ClassType = iota
	ClassTypeClass
)

type Resolver struct {
	interpreter     *Interpreter
	scopes          *Scopes
	currentFunction FunctionType
	currentClass    ClassType
}

type Scopes []*Scope

type Scope struct {
	variables       map[string]*LocalVariable
	currentIndex    int
	unusedVariables map[*Token]bool
}

type LocalVariable struct {
	declaration   *Token
	isInitialized bool
	scopedIndex   int
}

func NewScopes() *Scopes {
	return &Scopes{}
}

func NewScope() *Scope {
	return &Scope{
		variables:       map[string]*LocalVariable{},
		currentIndex:    0,
		unusedVariables: map[*Token]bool{},
	}
}

func (s *Scope) NewLocalVariable(declaration *Token) {
	s.variables[declaration.Lexeme] = &LocalVariable{
		declaration: declaration,
		scopedIndex: s.currentIndex,
	}
	s.currentIndex++
}

func (s *Scopes) push() {
	*s = append(*s, NewScope())
}

func (s *Scopes) pop() {
	*s = (*s)[:len(*s)-1]
}

func (s *Scopes) peek() *Scope {
	return (*s)[len(*s)-1]
}

func (s *Scopes) isEmpty() bool {
	return len(*s) == 0
}

func NewResolver(interpreter *Interpreter) *Resolver {
	return &Resolver{
		interpreter:     interpreter,
		scopes:          NewScopes(),
		currentFunction: FunctionTypeNone,
		currentClass:    ClassTypeNone,
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

func (resolver *Resolver) visitClassStmt(stmt *StmtClass) error {
	enclosingClass := resolver.currentClass
	resolver.currentClass = ClassTypeClass

	resolver.declare(stmt.name)
	resolver.define(stmt.name)

	resolver.beginScope()
	thisToken := NewToken(This, "this", True, 0)
	resolver.scopes.peek().NewLocalVariable(&thisToken)

	for _, method := range stmt.methods {
		var declaration FunctionType
		if method.name.Lexeme == "init" {
			declaration = FunctionTypeInitializer
		} else {
			declaration = FunctionTypeMethod
		}
		resolver.resolveFunction(method.function, declaration)
	}

	for _, staticMethod := range stmt.staticMethods {
		resolver.resolveFunction(staticMethod.function, FunctionTypeMethod)
	}

	resolver.endScope()
	resolver.currentClass = enclosingClass
	return nil
}

func (resolver *Resolver) beginScope() {
	resolver.scopes.push()
}

func (resolver *Resolver) endScope() {
	if GlobalConfig.ForbidUnusedVariable {
		for varDeclaration := range resolver.scopes.peek().unusedVariables {
			printError(varDeclaration, "Variable declared but never read")
		}
	}
	resolver.scopes.pop()
}

func (resolver *Resolver) declare(name *Token) {
	if resolver.scopes.isEmpty() {
		return
	}
	scope := resolver.scopes.peek()
	if _, ok := scope.variables[name.Lexeme]; ok {
		printError(name, "Already a variable with this name in this scope.")
	}

	scope.NewLocalVariable(name)
	if GlobalConfig.ForbidUnusedVariable {
		scope.unusedVariables[name] = true
	}
}

func (resolver *Resolver) define(name *Token) {
	if resolver.scopes.isEmpty() {
		return
	}
	variables := resolver.scopes.peek().variables
	variables[name.Lexeme].isInitialized = true
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
	resolver.resolveFunction(stmt.function, FunctionTypeFunc)
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
	if resolver.currentFunction == FunctionTypeNone {
		printError(stmt.keyword, "Can't return from top-level code.")
	}
	if stmt.expression != nil {
		if resolver.currentFunction == FunctionTypeInitializer {
			printError(stmt.keyword, "Can't return a value from an initializer.")
		}
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
	resolver.resolveFunction(expr, FunctionTypeFunc)
	return nil, nil
}

func (resolver *Resolver) visitCallExpr(expr *ExprCall) (any, error) {
	resolver.resolveExpr(expr.callee)

	for _, argumentExpr := range expr.arguments {
		resolver.resolveExpr(argumentExpr)
	}
	return nil, nil
}

func (resolver *Resolver) visitGetExpr(expr *ExprGet) (any, error) {
	resolver.resolveExpr(expr.object)
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

func (resolver *Resolver) visitSetExpr(expr *ExprSet) (any, error) {
	resolver.resolveExpr(expr.value)
	resolver.resolveExpr(expr.object)
	return nil, nil
}

func (resolver *Resolver) visitThisExpr(expr *ExprThis) (any, error) {
	if resolver.currentClass != ClassTypeClass {
		printError(expr.keyword, "Can't use 'this' outside of a class.")
		return nil, nil
	}
	resolver.resolveLocal(expr, expr.keyword, true)
	return nil, nil
}

func (resolver *Resolver) visitUnaryExpr(expr *ExprUnary) (any, error) {
	resolver.resolveExpr(expr.right)
	return nil, nil
}

func (resolver *Resolver) visitVariableExpr(expr *ExprVariable) (any, error) {
	if !resolver.scopes.isEmpty() {
		if localVar, ok := resolver.scopes.peek().variables[expr.name.Lexeme]; ok && !localVar.isInitialized {
			printError(expr.name, "Can't read local variable in its own initializer.")
		}
	}
	resolver.resolveLocal(expr, expr.name, true)
	return nil, nil
}

func (resolver *Resolver) resolveLocal(expr Expr, name *Token, isRead bool) {
	for i := len(*resolver.scopes) - 1; i >= 0; i-- {
		if localVar, ok := (*resolver.scopes)[i].variables[name.Lexeme]; ok {
			if isRead && GlobalConfig.ForbidUnusedVariable {
				delete((*resolver.scopes)[i].unusedVariables, localVar.declaration)
			}
			resolver.interpreter.resolve(expr, len(*resolver.scopes)-1-i, localVar.scopedIndex)
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
