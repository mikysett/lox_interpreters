package main

import (
	"fmt"
	"time"
)

type RuntimeError struct {
	token   *Token
	message string
}

func (e *RuntimeError) Error() string {
	return fmt.Sprintf("%v\n[line %v]", e.message, e.token.Line)
}

func NewRuntimeError(token *Token, message string) *RuntimeError {
	return &RuntimeError{
		token:   token,
		message: message,
	}
}

type Interpreter struct {
	enviroment *Environment
	globals    *Environment
}

func NewInterpreter() *Interpreter {
	globals := NewEnvironment()
	globals.define("clock", NewProtoCallable(
		func() int { return 0 },
		func(interpreter *Interpreter, arguments []any) (any, error) {
			return float64(time.Now().Unix()), nil
		},
		func() string { return "<native fn>" },
	))

	return &Interpreter{
		globals:    globals,
		enviroment: globals,
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

func (interpreter *Interpreter) visitBlockStmt(stmt *StmtBlock) error {
	return interpreter.executeBlock(stmt.block, NewEnvironment().withEnclosing(interpreter.enviroment))
}

func (interpreter *Interpreter) visitBreakStmt(stmt *StmtBreak) error {
	return NewBreakShortCircuit()
}

func (interpreter *Interpreter) executeBlock(stmts []Stmt, env *Environment) error {
	enclosingEnv := interpreter.enviroment
	interpreter.enviroment = env
	defer func() {
		interpreter.enviroment = enclosingEnv
	}()

	for _, stmt := range stmts {
		err := interpreter.execute(stmt)
		if err != nil {
			return err
		}
	}
	return nil
}

func (interpreter *Interpreter) visitVarStmt(stmt *StmtVar) (err error) {
	var value any
	if stmt.initializer != nil {
		value, err = interpreter.evaluate(stmt.initializer)
		if err != nil {
			return err
		}
	} else {
		value = Uninitialized{}
	}
	interpreter.enviroment.define(stmt.name.Lexeme, value)
	return nil
}

func (interpreter *Interpreter) visitAssignExpr(expr *ExprAssign) (any, error) {
	val, err := interpreter.evaluate(expr.value)
	if err != nil {
		return nil, err
	}
	err = interpreter.enviroment.assign(expr.name, val)
	if err != nil {
		return nil, err
	}
	return val, nil
}

func (interpreter *Interpreter) visitExpressionStmt(stmt *StmtExpression) error {
	_, err := interpreter.evaluate(stmt.expression)
	return err
}

func (interpreter *Interpreter) visitFunctionStmt(stmt *StmtFunction) error {
	function := NewFunction(stmt)
	interpreter.enviroment.define(stmt.name.Lexeme, function)
	return nil
}

func (interpreter *Interpreter) visitIfStmt(stmt *StmtIf) (err error) {
	eval, err := interpreter.evaluate(stmt.condition)
	if err != nil {
		return err
	}

	if isTruthy(eval) {
		return interpreter.execute(stmt.thenBranch)
	}
	if stmt.elseBranch != nil {
		return interpreter.execute(stmt.elseBranch)
	}
	return nil
}

func (interpreter *Interpreter) visitWhileStmt(stmt *StmtWhile) (err error) {
	for eval, err := interpreter.evaluate(stmt.condition); isTruthy(eval); eval, err = interpreter.evaluate(stmt.condition) {
		if err != nil {
			return err
		}

		err = interpreter.execute(stmt.body)
		if err != nil {
			if _, ok := err.(*BreakShortCircuit); ok {
				return nil
			}
			return err
		}
	}
	return nil
}

func (interpreter *Interpreter) visitPrintStmt(stmt *StmtPrint) error {
	v, err := interpreter.evaluate(stmt.expression)
	if err != nil {
		return err
	}

	fmt.Println(stringify(v))
	return nil
}

func (interpreter *Interpreter) visitReturnStmt(stmt *StmtReturn) (err error) {
	var value any
	if stmt.expression != nil {
		value, err = interpreter.evaluate(stmt.expression)
		if err != nil {
			return err
		}
	}
	return NewReturnShortCircuit(value)
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
		// This error message is not 100% correct (see case above)
		// But it needs to be phrased like that to pass the standard tests
		return nil, NewRuntimeError(expr.operator, "Operands must be two numbers or two strings.")
	case Comma:
		return right, nil
	default:
		return nil, NewRuntimeError(expr.operator, "Unexpected token.")
	}
}

func (interpreter *Interpreter) visitCallExpr(expr *ExprCall) (any, error) {
	callee, err := interpreter.evaluate(expr.callee)
	if err != nil {
		return nil, err
	}

	arguments := []any{}
	for _, argumentExpr := range expr.arguments {
		argument, err := interpreter.evaluate(argumentExpr)
		if err != nil {
			return nil, err
		}
		arguments = append(arguments, argument)
	}

	function, ok := callee.(Callable)
	if !ok {
		return nil, NewRuntimeError(expr.paren, "Can only call functions and classes.")
	}

	if function.arity() != len(arguments) {
		return nil, NewRuntimeError(expr.paren, fmt.Sprintf("Expected %d arguments but got %d.", function.arity(), len(arguments)))
	}

	result, err := function.call(interpreter, arguments)
	if err != nil {
		return nil, err
	}

	return result, nil
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

func (interpreter *Interpreter) visitLogicalExpr(expr *ExprLogical) (any, error) {
	left, err := interpreter.evaluate(expr.left)
	if err != nil {
		return nil, err
	}

	if expr.operator.Type == And {
		if !isTruthy(left) {
			return left, nil
		}
	} else if isTruthy(left) { // Logical `Or`
		return left, nil
	}

	return interpreter.evaluate(expr.right)
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
