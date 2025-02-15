package main

import "fmt"

type RuntimeError struct {
	token   *Token
	message string
}

func (e *RuntimeError) Error() string {
	if e.token.Type == EOF {
		return fmt.Sprintf("Line %v: at end. %v", e.token.Line, e.message)
	}
	return fmt.Sprintf("Line %v: at '%v'. %v", e.token.Line, e.token.Lexeme, e.message)
}

func NewRuntimeError(token *Token, message string) *RuntimeError {
	return &RuntimeError{
		token:   token,
		message: message,
	}
}

type Interpreter struct {
	enviroment *Environment
}

func NewInterpreter() *Interpreter {
	return &Interpreter{
		enviroment: NewEnvironment(),
	}
}

func (i *Interpreter) interpret(stmts []Stmt) error {
	for _, stmt := range stmts {
		err := i.execute(stmt)
		if err != nil {
			return err
		}
	}
	return nil
}

func (i *Interpreter) evaluate(expr Expr) (any, error) {
	return expr.accept(i)
}

func (i *Interpreter) execute(stmt Stmt) error {
	return stmt.accept(i)
}

func (interpreter *Interpreter) visitVarStmt(stmt *StmtVar) error {
	if stmt.initializer != nil {
		value, err := interpreter.evaluate(stmt.initializer)
		if err != nil {
			return err
		}
		interpreter.enviroment.define(stmt.name.Lexeme, value)
	}
	return nil
}

func (interpreter *Interpreter) visitExpressionStmt(stmt *StmtExpression) error {
	_, err := interpreter.evaluate(stmt.expression)
	return err
}

func (interpreter *Interpreter) visitPrintStmt(stmt *StmtPrint) error {
	v, err := interpreter.evaluate(stmt.expression)
	if err != nil {
		return err
	}

	fmt.Println(stringify(v))
	return nil
}

func (interpreter *Interpreter) visitBinaryExpr(expr *ExprBinary) (any, error) {
	left, err := interpreter.evaluate(expr.left)
	if err != nil {
		return nil, err
	}
	right, err := interpreter.evaluate(expr.right)
	if err != nil {
		return nil, err
	}

	switch expr.operator.Type {
	case BangEqual:
		return !isEqual(left, right), nil
	case EqualEqual:
		return isEqual(left, right), nil
	case Greater:
		err := checkNumberOperands(expr.operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.(float64) > right.(float64), nil
	case GreaterEqual:
		err := checkNumberOperands(expr.operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.(float64) >= right.(float64), nil
	case Less:
		err := checkNumberOperands(expr.operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.(float64) < right.(float64), nil
	case LessEqual:
		err := checkNumberOperands(expr.operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.(float64) <= right.(float64), nil
	case Minus:
		err := checkNumberOperands(expr.operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.(float64) - right.(float64), nil
	case Slash:
		err := checkNumberOperands(expr.operator, left, right)
		if err != nil {
			return nil, err
		}
		if right.(float64) == 0 {
			return nil, NewRuntimeError(expr.operator, "Division by 0.")
		}
		return left.(float64) / right.(float64), nil
	case Star:
		err := checkNumberOperands(expr.operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.(float64) * right.(float64), nil
	case Plus:
		err := checkNumberOperands(expr.operator, left, right)
		if err == nil {
			return left.(float64) + right.(float64), nil
		}
		if isOfType[string](left) || isOfType[string](right) {
			return stringify(left) + stringify(right), nil
		}
		return nil, NewRuntimeError(expr.operator, "Operands must be numbers or strings.")
	case Comma:
		return right, nil
	default:
		return nil, NewRuntimeError(expr.operator, "Unexpected token.")
	}
}

func (interpreter *Interpreter) visitTernaryExpr(expr *ExprTernary) (any, error) {
	condition, err := interpreter.evaluate(expr.condition)
	if err != nil {
		return nil, err
	}
	left, err := interpreter.evaluate(expr.left)
	if err != nil {
		return nil, err
	}
	right, err := interpreter.evaluate(expr.right)
	if err != nil {
		return nil, err
	}

	conditionBool, ok := condition.(bool)
	if !ok {
		return nil, NewRuntimeError(expr.operator, "Condition must evaluate to boolean.")
	}
	if conditionBool {
		return left, nil
	}
	return right, nil
}

func (interpreter *Interpreter) visitGroupingExpr(expr *ExprGrouping) (any, error) {
	return interpreter.evaluate(expr.expression)
}

func (interpreter *Interpreter) visitLiteralExpr(expr *ExprLiteral) (any, error) {
	return expr.value, nil
}

func (interpreter *Interpreter) visitUnaryExpr(expr *ExprUnary) (any, error) {
	right, err := interpreter.evaluate(expr.right)
	if err != nil {
		return nil, err
	}

	switch expr.operator.Type {
	case Bang:
		return !isTruthy(right), nil
	case Minus:
		err := checkNumberOperand(expr.operator, right)
		if err != nil {
			return nil, err
		}
		return -right.(float64), nil
	}
	panic("Unreachable!")
}

func (interpreter *Interpreter) visitVariableExpr(expr *ExprVariable) (any, error) {
	return interpreter.enviroment.get(expr.name)
}

func checkNumberOperand(operator *Token, operand any) error {
	if !isOfType[float64](operand) {
		return NewRuntimeError(operator, "Operand must be a number.")
	}
	return nil
}

func checkNumberOperands(operator *Token, left, right any) error {
	if !isOfType[float64](left) || !isOfType[float64](right) {
		return NewRuntimeError(operator, "Operands must be numbers.")
	}
	return nil
}

func isTruthy(input any) bool {
	if input == nil {
		return false
	}
	if isOfType[bool](input) {
		return input.(bool)
	}
	return true
}

func isEqual(left any, right any) bool {
	if left == nil && right == nil {
		return true
	}
	if left == nil {
		return false
	}
	return left == right
}

func stringify(val any) string {
	if val == nil {
		return "nil"
	}
	return fmt.Sprintf("%v", val)
}

func isOfType[T any](v any) bool {
	_, ok := v.(T)
	return ok
}
