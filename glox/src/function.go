package main

import "fmt"

type Function struct {
	declaration   *StmtFunction
	closure       *Environment
	isInitializer bool
}

func NewFunction(declaration *StmtFunction, closure *Environment, isInitializer bool) *Function {
	return &Function{
		declaration:   declaration,
		closure:       closure,
		isInitializer: isInitializer,
	}
}

func (f *Function) arity() int {
	return len(f.declaration.function.params)
}

func (f *Function) call(interpreter *Interpreter, arguments []any) (any, error) {
	env := NewEnvironment().WithEnclosing(f.closure)

	for i := range f.declaration.function.params {
		env.define(arguments[i])
	}

	err := interpreter.executeBlock(f.declaration.function.body, env)
	if res, ok := err.(*ReturnShortCircuit); ok {
		// In `init` method an empty `return` will always return `this` implicitely
		if f.isInitializer {
			return f.closure.getAt(&Position{depth: 0, index: 0}), nil
		}
		return res.value, nil
	}
	if err != nil {
		return nil, err
	}

	// `this` is always the first element declared in the stack of the instance
	if f.isInitializer {
		return f.closure.getAt(&Position{depth: 0, index: 0}), nil
	}
	return nil, nil
}

func (f *Function) Bind(instance *LoxInstance) *Function {
	env := NewEnvironment().WithEnclosing(f.closure)
	// The first element of the array will always be `this` for object methods
	env.define(instance)
	return NewFunction(f.declaration, env, f.isInitializer)
}

func (f *Function) String() string {
	if f.declaration.name == nil {
		return "<fn>"
	}
	return fmt.Sprintf("<fn %s>", f.declaration.name.Lexeme)
}
