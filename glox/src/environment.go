package main

type Environment struct {
	values map[string]any
}

func NewEnvironment() *Environment {
	return &Environment{
		values: map[string]any{},
	}
}

func (env *Environment) define(k string, v any) {
	env.values[k] = v
}

func (env *Environment) get(name *Token) (any, error) {
	v, ok := env.values[name.Lexeme]
	if !ok {
		return nil, NewRuntimeError(name, "Undefined variable '"+name.Lexeme+"'.")
	}
	return v, nil
}
