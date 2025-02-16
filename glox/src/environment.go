package main

type Environment struct {
	// Parent-pointer tree (cactus stack)
	enclosing *Environment
	values    map[string]any
}

func NewEnvironment() *Environment {
	return &Environment{
		enclosing: nil,
		values:    map[string]any{},
	}
}

func (env *Environment) withEnclosing(enclosing *Environment) *Environment {
	env.enclosing = enclosing
	return env
}

func (env *Environment) define(k string, v any) {
	env.values[k] = v
}

func (env *Environment) assign(name *Token, value any) error {
	_, ok := env.values[name.Lexeme]
	if ok {
		env.values[name.Lexeme] = value
		return nil
	}
	if env.enclosing != nil {
		return env.enclosing.assign(name, value)
	}
	return NewRuntimeError(name, "Undefined variable '"+name.Lexeme+"'.")
}

func (env *Environment) get(name *Token) (any, error) {
	v, ok := env.values[name.Lexeme]
	if ok {
		return v, nil
	}
	if env.enclosing != nil {
		return env.enclosing.get(name)
	}
	return nil, NewRuntimeError(name, "Undefined variable '"+name.Lexeme+"'.")
}
