package main

import (
	"fmt"
	"strings"
)

type AstPrinter struct{}

func NewAstPrinter() *AstPrinter {
	return &AstPrinter{}
}

// TODO: Update the printer to accept also Statements
func (ast *AstPrinter) print(expr Expr) any {
	str, _ := expr.accept(ast)
	return str
}

func (ast *AstPrinter) visitAssignExpr(expr *ExprAssign) (any, error) {
	return ast.parenthesize("= "+expr.name.Lexeme, expr.value)
}

func (ast *AstPrinter) visitBinaryExpr(expr *ExprBinary) (any, error) {
	return ast.parenthesize(expr.operator.Lexeme, expr.left, expr.right)
}

func (ast *AstPrinter) visitTernaryExpr(expr *ExprTernary) (any, error) {
	return ast.parenthesize("?:", expr.condition, expr.left, expr.right)
}

func (ast *AstPrinter) visitGroupingExpr(expr *ExprGrouping) (any, error) {
	return ast.parenthesize("group", expr.expression)
}

func (ast *AstPrinter) visitLiteralExpr(expr *ExprLiteral) (any, error) {
	if expr.value == nil {
		return "nil", nil
	}
	switch expr.value.(type) {
	case string:
		return fmt.Sprintf("\"%v\"", expr.value), nil
	case float64:
		return fmt.Sprintf("%.1f", expr.value), nil
	default:
		return fmt.Sprintf("%v", expr.value), nil
	}
}

func (ast *AstPrinter) visitUnaryExpr(expr *ExprUnary) (any, error) {
	return ast.parenthesize(expr.operator.Lexeme, expr.right)
}

func (ast *AstPrinter) visitVariableExpr(expr *ExprVariable) (any, error) {
	return expr.name.Lexeme, nil
}

func (ast *AstPrinter) parenthesize(name string, exprs ...Expr) (string, error) {
	var builder strings.Builder
	builder.WriteByte('(')
	builder.WriteString(name)
	for _, expr := range exprs {
		builder.WriteByte(' ')
		astResult, _ := expr.accept(ast)
		if v, ok := astResult.(string); ok {
			builder.WriteString(v)
		} else {
			panic("Unreachable non string return in AstPrinter!")
		}
	}
	builder.WriteByte(')')
	return builder.String(), nil
}
