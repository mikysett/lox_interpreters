package main

type Expr interface {
	accept(Visitor) (any, error)
}

type Visitor interface {
	visitBinaryExpr(*ExprBinary) (any, error)
	visitTernaryExpr(*ExprTernary) (any, error)
	visitGroupingExpr(*ExprGrouping) (any, error)
	visitLiteralExpr(*ExprLiteral) (any, error)
	visitUnaryExpr(*ExprUnary) (any, error)
}

// Binary   : Expr left, Token operator, Expr right
type ExprBinary struct {
	left     Expr
	operator Token
	right    Expr
}

func NewExprBinary(left Expr, operator Token, right Expr) *ExprBinary {
	return &ExprBinary{
		left:     left,
		operator: operator,
		right:    right,
	}
}

func (expr *ExprBinary) accept(v Visitor) (any, error) {
	return v.visitBinaryExpr(expr)
}

// Binary   : Expr left, Token operator, Expr right
type ExprTernary struct {
	operator  Token
	condition Expr
	left      Expr
	right     Expr
}

func NewExprTernary(operator Token, condition Expr, left Expr, right Expr) *ExprTernary {
	return &ExprTernary{
		operator:  operator,
		condition: condition,
		left:      left,
		right:     right,
	}
}

func (expr *ExprTernary) accept(v Visitor) (any, error) {
	return v.visitTernaryExpr(expr)
}

// Grouping : Expr expression
type ExprGrouping struct {
	expression Expr
}

func NewExprGrouping(expression Expr) *ExprGrouping {
	return &ExprGrouping{expression: expression}
}

func (expr *ExprGrouping) accept(v Visitor) (any, error) {
	return v.visitGroupingExpr(expr)
}

// Literal  : Object value
type ExprLiteral struct {
	value any
}

func NewExprLiteral(value any) *ExprLiteral {
	return &ExprLiteral{value: value}
}

func (expr *ExprLiteral) accept(v Visitor) (any, error) {
	return v.visitLiteralExpr(expr)
}

// Unary    : Token operator, Expr right
type ExprUnary struct {
	operator Token
	right    Expr
}

func NewExprUnary(operator Token, right Expr) *ExprUnary {
	return &ExprUnary{
		operator: operator,
		right:    right,
	}
}

func (expr *ExprUnary) accept(v Visitor) (any, error) {
	return v.visitUnaryExpr(expr)
}
