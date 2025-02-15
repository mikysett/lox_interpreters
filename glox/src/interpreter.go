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

type Interpreter struct{}

func NewInterpreter() *Interpreter {
	return &Interpreter{}
}

func (i *Interpreter) evaluate(expr Expr) (any, error) {
	return expr.accept(i)
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
		err := checkNumberOperands(&expr.operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.(float64) > right.(float64), nil
	case GreaterEqual:
		err := checkNumberOperands(&expr.operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.(float64) >= right.(float64), nil
	case Less:
		err := checkNumberOperands(&expr.operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.(float64) < right.(float64), nil
	case LessEqual:
		err := checkNumberOperands(&expr.operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.(float64) <= right.(float64), nil
	case Minus:
		err := checkNumberOperands(&expr.operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.(float64) - right.(float64), nil
	case Slash:
		err := checkNumberOperands(&expr.operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.(float64) / right.(float64), nil
	case Star:
		err := checkNumberOperands(&expr.operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.(float64) * right.(float64), nil
	case Plus:
		err := checkNumberOperands(&expr.operator, left, right)
		if err == nil {
			return left.(float64) + right.(float64), nil
		}
		leftStr, leftOk := left.(string)
		rightStr, rightOk := right.(string)
		if !rightOk || !leftOk {
			return nil, NewRuntimeError(&expr.operator, "Operands must be numbers or strings.")
		}
		return leftStr + rightStr, nil
	case Comma:
		return right, nil
	default:
		return nil, NewRuntimeError(&expr.operator, "Unexpected token.")
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
		return nil, NewRuntimeError(&expr.operator, "Condition must evaluate to boolean.")
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
		err := checkNumberOperand(&expr.operator, right)
		if err != nil {
			return nil, err
		}
		return -right.(float64), nil
	}
	panic("Unreachable!")
}

func checkNumberOperand(operator *Token, operand any) error {
	if _, ok := operand.(float64); !ok {
		return NewRuntimeError(operator, "Operand must be a number.")
	}
	return nil
}

func checkNumberOperands(operator *Token, left, right any) error {
	_, okLeft := left.(float64)
	_, okRight := right.(float64)
	if !okLeft || !okRight {
		return NewRuntimeError(operator, "Operands must be numbers.")
	}
	return nil
}

func isTruthy(input any) bool {
	if input == nil {
		return false
	}
	if v, ok := input.(bool); ok {
		return v
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
