package main

type Expr interface {
	accept(Visitor) any
}

type Visitor interface {
	visitBinaryExpr(*ExprBinary)
	visitGroupingExpr(*ExprGrouping)
	visitLiteralExpr(*ExprLiteral)
	visitUnaryExpr(*ExprUnary)
}

// Binary   : Expr left, Token operator, Expr right
type ExprBinary struct {
	left     Expr
	operator Token
	right    Expr
}

func NewExprBinary(left Expr, operator Token, right Expr) ExprBinary {
	return ExprBinary{
		left:     left,
		operator: operator,
		right:    right,
	}
}

func (expr *ExprBinary) accept(v Visitor) any {
	v.visitBinaryExpr(expr)
	return nil
}

// Grouping : Expr expression
type ExprGrouping struct {
	expression Expr
}

func NewExprGrouping(expression Expr) ExprGrouping {
	return ExprGrouping{expression: expression}
}

func (expr *ExprGrouping) accept(v Visitor) any {
	v.visitGroupingExpr(expr)
	return nil
}

// Literal  : Object value
type ExprLiteral struct {
	value any
}

func NewExprLiteral(value any) ExprLiteral {
	return ExprLiteral{value: value}
}

func (expr *ExprLiteral) accept(v Visitor) any {
	v.visitLiteralExpr(expr)
	return nil
}

// Unary    : Token operator, Expr right
type ExprUnary struct {
	operator Token
	right    Expr
}

func NewExprUnary(operator Token, right Expr) ExprUnary {
	return ExprUnary{
		operator: operator,
		right:    right,
	}
}

func (expr *ExprUnary) accept(v Visitor) any {
	v.visitUnaryExpr(expr)
	return nil
}
