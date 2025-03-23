package main

import (
	"fmt"
	"math"
	"strconv"
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
	globals    map[string]any
	locals     map[Expr]*Position
}

type Position struct {
	depth int
	index int
}

func NewInterpreter() *Interpreter {
	globals := map[string]any{
		"clock": NewProtoCallable(
			func() int { return 0 },
			func(interpreter *Interpreter, arguments []any) (any, error) {
				return float64(time.Now().Unix()), nil
			},
			func() string { return "<native fn>" },
		),
		"len": NewProtoCallable(
			func() int { return 1 },
			func(interpreter *Interpreter, arguments []any) (any, error) {
				val := arguments[0]
				switch val.(type) {
				case *LoxInstance:
					return float64(len(val.(*LoxInstance).fields)), nil
				case []byte:
					return float64(len(val.([]byte))), nil
				default:
					// TODO: improve error message
					dummyToken := NewToken(Identifier, "[PARAMETER]", nil, 0)
					return nil, NewRuntimeError(&dummyToken, "Function call 'len' only valid on object instances, arrays and strings.")
				}
			},
			func() string { return "<native fn>" },
		),
	}

	return &Interpreter{
		globals:    globals,
		enviroment: NewEnvironment(),
		locals:     map[Expr]*Position{},
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

func (i *Interpreter) resolve(expr Expr, depth int, index int) {
	i.locals[expr] = &Position{depth, index}
}

func (i *Interpreter) execute(stmt Stmt) error {
	return stmt.accept(i)
}

func (interpreter *Interpreter) visitBlockStmt(stmt *StmtBlock) error {
	return interpreter.executeBlock(stmt.block, NewEnvironment().WithEnclosing(interpreter.enviroment))
}

func (interpreter *Interpreter) visitClassStmt(stmt *StmtClass) error {
	var superclass *LoxClass
	if stmt.superclass != nil {
		result, err := interpreter.evaluate(stmt.superclass)
		if err != nil {
			return err
		}

		if super, ok := result.(*LoxClass); ok {
			superclass = super
		} else {
			return NewRuntimeError(stmt.superclass.name, "Superclass must be a class.")
		}
	}

	if interpreter.enviroment.enclosing == nil {
		interpreter.globals[stmt.name.Lexeme] = nil
	} else {
		interpreter.enviroment.define(nil)
	}

	if stmt.superclass != nil {
		interpreter.enviroment = NewEnvironment().WithEnclosing(interpreter.enviroment)
		interpreter.enviroment.define("super")
		interpreter.enviroment.assignAtLast(superclass)
	}

	methods := map[string]*Function{}
	for _, method := range stmt.methods {
		isInitializer := false
		if method.name.Lexeme == "init" {
			isInitializer = true
		}
		methods[method.name.Lexeme] = NewFunction(method, interpreter.enviroment, isInitializer)
	}

	staticMethods := map[string]*Function{}
	for _, staticMethod := range stmt.staticMethods {
		staticMethods[staticMethod.name.Lexeme] = NewFunction(staticMethod, interpreter.enviroment, false)
	}

	class := NewLoxClass(
		superclass,
		// metaclass in order to have static methods on class objects [extra feature]
		NewLoxInstance(NewLoxClass(nil, nil, stmt.name.Lexeme, staticMethods)),
		stmt.name.Lexeme,
		methods,
	)

	if stmt.superclass != nil {
		interpreter.enviroment = interpreter.enviroment.enclosing
	}

	if interpreter.enviroment.enclosing == nil {
		interpreter.globals[stmt.name.Lexeme] = class
	} else {
		interpreter.enviroment.assignAtLast(class)
	}

	return nil
}

func (interpreter *Interpreter) visitBreakStmt(stmt *StmtBreak) error {
	return NewBreakShortCircuit()
}

func (interpreter *Interpreter) visitContinueStmt(stmt *StmtContinue) error {
	return NewContinueShortCircuit()
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

	if interpreter.enviroment.enclosing == nil {
		interpreter.globals[stmt.name.Lexeme] = value
	} else {
		interpreter.enviroment.define(value)
	}

	return nil
}

func (interpreter *Interpreter) visitAssignExpr(expr *ExprAssign) (any, error) {
	val, err := interpreter.evaluate(expr.value)
	if err != nil {
		return nil, err
	}

	if position, ok := interpreter.locals[expr]; ok {
		interpreter.enviroment.assignAt(position, val)
		return val, nil
	}

	_, ok := interpreter.globals[expr.name.Lexeme]
	if ok {
		interpreter.globals[expr.name.Lexeme] = val
		return val, nil
	}
	return nil, NewRuntimeError(expr.name, "Undefined variable '"+expr.name.Lexeme+"'.")
}

func (interpreter *Interpreter) visitExpressionStmt(stmt *StmtExpression) error {
	_, err := interpreter.evaluate(stmt.expression)
	return err
}

func (interpreter *Interpreter) visitFunctionStmt(stmt *StmtFunction) error {
	function := NewFunction(stmt, interpreter.enviroment, false)
	if interpreter.enviroment.enclosing == nil {
		interpreter.globals[stmt.name.Lexeme] = function
	} else {
		interpreter.enviroment.define(function)
	}
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

func (interpreter *Interpreter) visitLoopStmt(stmt *StmtLoop) (err error) {
	for eval, err := interpreter.evaluate(stmt.condition); isTruthy(eval); eval, err = interpreter.evaluate(stmt.condition) {
		if err != nil {
			return err
		}

		err = interpreter.execute(stmt.body)
		if err != nil {
			switch err.(type) {
			case *BreakShortCircuit:
				return nil
			case *ContinueShortCircuit:
				if stmt.increment != nil {
					_, err = interpreter.evaluate(stmt.increment)
					if err != nil {
						return err
					}
				}
				continue
			}
			return err
		}
		if stmt.increment != nil {
			_, err = interpreter.evaluate(stmt.increment)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (interpreter *Interpreter) visitPrintStmt(stmt *StmtPrint) error {
	v, err := interpreter.evaluate(stmt.expression)
	if err != nil {
		return err
	}

	fmt.Println(string(stringify(v)))
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
	case Percent:
		err := checkNumberOperands(expr.operator, left, right)
		if err != nil {
			return nil, err
		}
		return float64(int(math.Round(left.(float64))) % int(math.Round(right.(float64)))), nil
	case Plus:
		if err := checkNumberOperands(expr.operator, left, right); err == nil {
			return left.(float64) + right.(float64), nil
		}
		isLeftString, isRightString := isOfType[[]byte](left), isOfType[[]byte](right)
		if (isLeftString && isRightString) ||
			(GlobalConfig.AllowImplicitStringCast && (isLeftString || isRightString)) {
			return append(stringify(left), stringify(right)...), nil
		}
		if GlobalConfig.AllowImplicitStringCast {
			return nil, NewRuntimeError(expr.operator, "Operands must be numbers and/or strings.")
		}
		return nil, NewRuntimeError(expr.operator, "Operands must be two numbers or two strings.")
	case Comma:
		return right, nil
	default:
		return nil, NewRuntimeError(expr.operator, "Unexpected token.")
	}
}

func (interpreter *Interpreter) visitFunctionExpr(expr *ExprFunction) (any, error) {
	return NewFunction(NewStmtFunction(nil, expr), interpreter.enviroment, false), nil
}

func (interpreter *Interpreter) visitArrayExpr(expr *ExprArray) (any, error) {
	array, err := interpreter.evaluate(expr.array)
	if err != nil {
		return nil, err
	}

	index, err := interpreter.evaluate(expr.index)
	if err != nil {
		return nil, err
	}

	var indexStr string
	if isOfType[[]byte](index) {
		indexStr = string(index.([]byte))
	} else {
		indexStr = fmt.Sprintf("%v", index)
	}

	switch array.(type) {
	case *LoxInstance:
		if result, ok := array.(*LoxInstance).fields[indexStr]; ok {
			return result, nil
		}
		return nil, NewRuntimeError(expr.bracket, fmt.Sprintf("Undefined index '%s'.", indexStr))
	case []byte:
		indexInt, err := getIndexForString(array.([]byte), indexStr)
		if err != nil {
			return nil, NewRuntimeError(expr.bracket, err.Error())
		}
		return string(array.([]byte)[indexInt]), nil
	default:
		return nil, NewRuntimeError(expr.bracket, "Can only access array indexes on class instances and strings.")
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

func (interpreter *Interpreter) visitGetExpr(expr *ExprGet) (any, error) {
	object, err := interpreter.evaluate(expr.object)
	if err != nil {
		return nil, err
	}

	if instance, ok := object.(*LoxInstance); ok {
		method, err := instance.Get(expr.name)
		if err != nil {
			return nil, err
		}

		if isOfType[*Function](method) && method.(*Function).IsGetter() {
			return method.(*Function).call(interpreter, nil)
		}
		return method, nil
	}

	if GlobalConfig.AllowStaticMethods {
		if class, ok := object.(*LoxClass); ok {
			return class.metaclass.Get(expr.name)
		}
	}

	return nil, NewRuntimeError(expr.name, "Only instances have properties.")
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

func (interpreter *Interpreter) visitSetExpr(expr *ExprSet) (any, error) {
	object, err := interpreter.evaluate(expr.object)
	if err != nil {
		return nil, err
	}

	if !isOfType[*LoxInstance](object) &&
		(!GlobalConfig.AllowStaticMethods || !isOfType[*LoxClass](object)) {
		return nil, NewRuntimeError(expr.name, "Only instances have fields.")
	}

	value, err := interpreter.evaluate(expr.value)
	if err != nil {
		return nil, err
	}
	if obj, ok := object.(*LoxInstance); ok {
		obj.Set(expr.name, value)
	} else if obj, ok := object.(*LoxClass); ok {
		obj.metaclass.Set(expr.name, value)
	}

	return value, nil
}

func (interpreter *Interpreter) visitSetArrayExpr(expr *ExprSetArray) (any, error) {
	object, err := interpreter.evaluate(expr.object)
	if err != nil {
		return nil, err
	}

	if !isOfType[*LoxInstance](object) &&
		!(GlobalConfig.AllowStaticMethods || isOfType[*LoxClass](object)) &&
		!(GlobalConfig.AllowArrays || isOfType[[]byte](object)) {
		return nil, NewRuntimeError(expr.name, "Only instances have fields.")
	}

	index, err := interpreter.evaluate(expr.index)
	if err != nil {
		return nil, err
	}
	indexStr := fmt.Sprintf("%v", index)

	value, err := interpreter.evaluate(expr.value)
	if err != nil {
		return nil, err
	}
	if obj, ok := object.(*LoxInstance); ok {
		obj.fields[indexStr] = value
	} else if obj, ok := object.(*LoxClass); ok {
		obj.metaclass.fields[indexStr] = value
	} else if str, ok := object.([]byte); ok {
		indexInt, err := getIndexForString(str, indexStr)
		if err != nil {
			return nil, NewRuntimeError(expr.name, err.Error())
		}
		valueStr := fmt.Sprintf("%v", value)
		if len(valueStr) != 1 {
			return nil, NewRuntimeError(expr.name, "Value assigned to a string index must be a single character.")
		}
		str[indexInt] = []byte(valueStr)[0]
		return str, nil
	}

	return value, nil
}

func (interpreter *Interpreter) visitArrayInstanceExpr(expr *ExprArrayInstance) (any, error) {
	array := NewArrayInstance()
	for i, arg := range expr.arguments {
		evaluatedArg, err := interpreter.evaluate(arg)
		if err != nil {
			return nil, err
		}
		array.fields[strconv.Itoa(i)] = evaluatedArg
	}
	return array, nil
}

func (interpreter *Interpreter) visitSuperExpr(expr *ExprSuper) (any, error) {
	superPos := interpreter.locals[expr]
	superclass := interpreter.enviroment.getAt(superPos).(*LoxClass)

	thisPosition := Position{depth: superPos.depth - 1, index: 0}
	object := interpreter.enviroment.getAt(&thisPosition).(*LoxInstance)

	method := superclass.FindMethod(expr.method.Lexeme)
	if method == nil {
		return nil, NewRuntimeError(expr.method, "Undefined property '"+expr.method.Lexeme+"'.")
	}

	return method.Bind(object), nil
}

func (interpreter *Interpreter) visitThisExpr(expr *ExprThis) (any, error) {
	return interpreter.lookUpVariable(expr.keyword, expr)
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
	return interpreter.lookUpVariable(expr.name, expr)
}

func (interpreter *Interpreter) lookUpVariable(name *Token, expr Expr) (any, error) {
	if position, ok := interpreter.locals[expr]; ok {
		return interpreter.enviroment.getAt(position), nil
	}

	value, ok := interpreter.globals[name.Lexeme]
	if !ok {
		return nil, NewRuntimeError(name, "Undefined variable '"+name.Lexeme+"'.")
	}
	if isOfType[Uninitialized](value) {
		if GlobalConfig.ForbidUninitializedVariable {
			return nil, NewRuntimeError(name, "Uninitialized variable '"+name.Lexeme+"'.")
		} else {
			return nil, nil
		}
	}
	return value, nil
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
	if str, ok := left.([]byte); ok {
		left = string(str)
	}
	if str, ok := right.([]byte); ok {
		right = string(str)
	}
	return left == right
}

func stringify(val any) []byte {
	if val == nil {
		return []byte("nil")
	}
	if str, ok := val.([]byte); ok {
		cpy := make([]byte, len(str))
		copy(cpy, str)
		return cpy
	}
	return fmt.Appendf([]byte{}, "%v", val)
}

func getIndexForString(s []byte, indexStr string) (int, error) {
	indexInt, err := strconv.Atoi(indexStr)
	if err != nil {
		return 0, fmt.Errorf("Invalid index: %w.", err)
	}
	if indexInt < 0 || indexInt >= len(s) {
		return 0, fmt.Errorf("Index %d of string out of range.", indexInt)
	}
	return indexInt, nil
}

func isOfType[T any](v any) bool {
	_, ok := v.(T)
	return ok
}
