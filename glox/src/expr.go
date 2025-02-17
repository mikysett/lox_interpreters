package main

type Expr interface {
	accept(ExprVisitor) (any, error)
}

type ExprVisitor interface {
	visitBinaryExpr(*ExprBinary) (any, error)
	visitTernaryExpr(*ExprTernary) (any, error)
	visitGroupingExpr(*ExprGrouping) (any, error)
	visitLiteralExpr(*ExprLiteral) (any, error)
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

// "Variable : Token name"
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
