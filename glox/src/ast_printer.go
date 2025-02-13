package main

import (
	"fmt"
	"strings"
)

type AstPrinter struct{}

func NewAstPrinter() *AstPrinter {
	return &AstPrinter{}
}

func (ast *AstPrinter) print(expr Expr) any {
	return expr.accept(ast)
}

func (ast *AstPrinter) visitBinaryExpr(expr *ExprBinary) any {
	return ast.parenthesize(expr.operator.Lexeme, expr.left, expr.right)
}

func (ast *AstPrinter) visitTernaryExpr(expr *ExprTernary) any {
	return ast.parenthesize("?:", expr.condition, expr.left, expr.right)
}

func (ast *AstPrinter) visitGroupingExpr(expr *ExprGrouping) any {
	return ast.parenthesize("group", expr.expression)
}

func (ast *AstPrinter) visitLiteralExpr(expr *ExprLiteral) any {
	if expr.value == nil {
		return "nil"
	}
	switch expr.value.(type) {
	case string:
		return fmt.Sprintf("\"%v\"", expr.value)
	case float64:
		return fmt.Sprintf("%.1f", expr.value)
	default:
		return fmt.Sprintf("%v", expr.value)
	}
}

func (ast *AstPrinter) visitUnaryExpr(expr *ExprUnary) any {
	return ast.parenthesize(expr.operator.Lexeme, expr.right)
}

func (ast *AstPrinter) parenthesize(name string, exprs ...Expr) string {
	var builder strings.Builder
	builder.WriteByte('(')
	builder.WriteString(name)
	for _, expr := range exprs {
		builder.WriteByte(' ')
		if v, ok := expr.accept(ast).(string); ok {
			builder.WriteString(v)
		} else {
			panic("Unreachable non string return in AstPrinter!")
		}
	}
	builder.WriteByte(')')
	return builder.String()
}
