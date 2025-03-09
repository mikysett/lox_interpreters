package main

import "fmt"

type Function struct {
	declaration *StmtFunction
	closure     *Environment
}

func NewFunction(declaration *StmtFunction, closure *Environment) *Function {
	return &Function{
		declaration: declaration,
		closure:     closure,
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
		return res.value, nil
	}
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (f *Function) String() string {
	if f.declaration.name == nil {
		return "<fn>"
	}
	return fmt.Sprintf("<fn %s>", f.declaration.name.Lexeme)
}
