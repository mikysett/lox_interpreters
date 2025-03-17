package main

type Expr interface {
	accept(ExprVisitor) (any, error)
}

type ExprVisitor interface {
	visitBinaryExpr(*ExprBinary) (any, error)
	visitFunctionExpr(*ExprFunction) (any, error)
	visitCallExpr(*ExprCall) (any, error)
	visitGetExpr(*ExprGet) (any, error)
	visitTernaryExpr(*ExprTernary) (any, error)
	visitGroupingExpr(*ExprGrouping) (any, error)
	visitLiteralExpr(*ExprLiteral) (any, error)
	visitSetExpr(*ExprSet) (any, error)
	visitSuperExpr(*ExprSuper) (any, error)
	visitThisExpr(*ExprThis) (any, error)
	visitUnaryExpr(*ExprUnary) (any, error)
	visitVariableExpr(*ExprVariable) (any, error)
	visitAssignExpr(*ExprAssign) (any, error)
	visitLogicalExpr(*ExprLogical) (any, error)
}

// Assign   : Token name, Expr value
type ExprAssign struct {
	name  *Token
	value Expr
}

func NewExprAssign(name *Token, value Expr) *ExprAssign {
	return &ExprAssign{
		name:  name,
		value: value,
	}
}

func (expr *ExprAssign) accept(v ExprVisitor) (any, error) {
	return v.visitAssignExpr(expr)
}

// Binary   : Expr left, Token operator, Expr right
type ExprBinary struct {
	left     Expr
	operator *Token
	right    Expr
}

func NewExprBinary(left Expr, operator *Token, right Expr) *ExprBinary {
	return &ExprBinary{
		left:     left,
		operator: operator,
		right:    right,
	}
}

func (expr *ExprBinary) accept(v ExprVisitor) (any, error) {
	return v.visitBinaryExpr(expr)
}

// ExprFunction   : List<Token> params, List<Stmt> body
type ExprFunction struct {
	params []*Token
	body   []Stmt
}

func NewExprFunction(params []*Token, body []Stmt) *ExprFunction {
	return &ExprFunction{
		params: params,
		body:   body,
	}
}

func (stmt *ExprFunction) accept(v ExprVisitor) (any, error) {
	return v.visitFunctionExpr(stmt)
}

// Call     : Expr callee, Token paren, List<Expr> arguments
type ExprCall struct {
	callee    Expr
	paren     *Token
	arguments []Expr
}

func NewExprCall(callee Expr, paren *Token, arguments []Expr) *ExprCall {
	return &ExprCall{
		callee:    callee,
		paren:     paren,
		arguments: arguments,
	}
}

func (expr *ExprCall) accept(v ExprVisitor) (any, error) {
	return v.visitCallExpr(expr)
}

// Get     : Expr object, Token name
type ExprGet struct {
	object Expr
	name   *Token
}

func NewExprGet(object Expr, name *Token) *ExprGet {
	return &ExprGet{
		object: object,
		name:   name,
	}
}

func (expr *ExprGet) accept(v ExprVisitor) (any, error) {
	return v.visitGetExpr(expr)
}

// Ternary   : Expr condition, Token operator, Expr left, Expr right
type ExprTernary struct {
	operator  *Token
	condition Expr
	left      Expr
	right     Expr
}

func NewExprTernary(operator *Token, condition Expr, left Expr, right Expr) *ExprTernary {
	return &ExprTernary{
		operator:  operator,
		condition: condition,
		left:      left,
		right:     right,
	}
}

func (expr *ExprTernary) accept(v ExprVisitor) (any, error) {
	return v.visitTernaryExpr(expr)
}

// Grouping : Expr expression
type ExprGrouping struct {
	expression Expr
}

func NewExprGrouping(expression Expr) *ExprGrouping {
	return &ExprGrouping{expression: expression}
}

func (expr *ExprGrouping) accept(v ExprVisitor) (any, error) {
	return v.visitGroupingExpr(expr)
}

// Literal  : Object value
type ExprLiteral struct {
	value any
}

func NewExprLiteral(value any) *ExprLiteral {
	return &ExprLiteral{value: value}
}

func (expr *ExprLiteral) accept(v ExprVisitor) (any, error) {
	return v.visitLiteralExpr(expr)
}

// Logical  : Expr left, Token operator, Expr right

type ExprLogical struct {
	left     Expr
	operator *Token
	right    Expr
}

func NewExprLogical(left Expr, operator *Token, right Expr) *ExprLogical {
	return &ExprLogical{
		left:     left,
		operator: operator,
		right:    right,
	}
}

func (expr *ExprLogical) accept(v ExprVisitor) (any, error) {
	return v.visitLogicalExpr(expr)
}

// Set    : Expr object, Token name, Expr value
type ExprSet struct {
	object Expr
	name   *Token
	value  Expr
}

func NewExprSet(object Expr, name *Token, value Expr) *ExprSet {
	return &ExprSet{
		object: object,
		name:   name,
		value:  value,
	}
}

func (expr *ExprSet) accept(v ExprVisitor) (any, error) {
	return v.visitSetExpr(expr)
}

// Super    : Token keyword, Token method
type ExprSuper struct {
	keyword *Token
	method  *Token
}

func NewExprSuper(keyword, method *Token) *ExprSuper {
	return &ExprSuper{
		keyword: keyword,
		method:  method,
	}
}

func (expr *ExprSuper) accept(v ExprVisitor) (any, error) {
	return v.visitSuperExpr(expr)
}

// This    : Token keyword
type ExprThis struct {
	keyword *Token
}

func NewExprThis(keyword *Token) *ExprThis {
	return &ExprThis{
		keyword: keyword,
	}
}

func (expr *ExprThis) accept(v ExprVisitor) (any, error) {
	return v.visitThisExpr(expr)
}

// Unary    : Token operator, Expr right
type ExprUnary struct {
	operator *Token
	right    Expr
}

func NewExprUnary(operator *Token, right Expr) *ExprUnary {
	return &ExprUnary{
		operator: operator,
		right:    right,
	}
}

func (expr *ExprUnary) accept(v ExprVisitor) (any, error) {
	return v.visitUnaryExpr(expr)
}

// Variable : Token name
type ExprVariable struct {
	name *Token
}

func NewExprVariable(name *Token) *ExprVariable {
	return &ExprVariable{
		name: name,
	}
}

func (expr *ExprVariable) accept(v ExprVisitor) (any, error) {
	return v.visitVariableExpr(expr)
}
