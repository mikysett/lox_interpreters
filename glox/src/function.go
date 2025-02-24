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
	return len(f.declaration.params)
}

func (f *Function) call(interpreter *Interpreter, arguments []any) (any, error) {
	env := NewEnvironment().withEnclosing(f.closure)

	for i, param := range f.declaration.params {
		env.define(param.Lexeme, arguments[i])
	}

	err := interpreter.executeBlock(f.declaration.body, env)
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
