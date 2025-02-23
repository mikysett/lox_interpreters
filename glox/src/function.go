package main

import "fmt"

type Function struct {
	declaration *StmtFunction
}

func NewFunction(declaration *StmtFunction) *Function {
	return &Function{
		declaration: declaration,
	}
}

func (f *Function) arity() int {
	return len(f.declaration.params)
}

func (f *Function) call(interpreter *Interpreter, arguments []any) (any, error) {
	env := NewEnvironment().withEnclosing(interpreter.globals)

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
	return fmt.Sprintf("<fn %s>", f.declaration.name.Lexeme)
}
