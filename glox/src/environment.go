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

func (env *Environment) assign(name *Token, value any) error {
	_, ok := env.values[name.Lexeme]
	if !ok {
		return NewRuntimeError(name, "Undefined variable '"+name.Lexeme+"'.")
	}
	env.values[name.Lexeme] = value
	return nil
}

func (env *Environment) get(name *Token) (any, error) {
	v, ok := env.values[name.Lexeme]
	if !ok {
		return nil, NewRuntimeError(name, "Undefined variable '"+name.Lexeme+"'.")
	}
	return v, nil
}
