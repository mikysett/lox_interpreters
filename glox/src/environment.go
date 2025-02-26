package main

type Environment struct {
	// Parent-pointer tree (cactus stack)
	enclosing *Environment
	values    map[string]any
}

type Uninitialized struct{}

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

func (env *Environment) assignAt(distance int, name *Token, value any) {
	env.ancestor(distance).values[name.Lexeme] = value
}

func (env *Environment) get(name *Token) (any, error) {
	value, ok := env.values[name.Lexeme]
	if !ok && env.enclosing != nil {
		return env.enclosing.get(name)
	} else if ok {
		if isOfType[Uninitialized](value) {
			return nil, NewRuntimeError(name, "Uninitialized variable '"+name.Lexeme+"'.")
		}
		return value, nil
	}
	return nil, NewRuntimeError(name, "Undefined variable '"+name.Lexeme+"'.")
}

func (env *Environment) getAt(distance int, name *Token) any {
	return env.ancestor(distance).values[name.Lexeme]
}

func (env *Environment) ancestor(distance int) *Environment {
	for i := 0; i < distance; i++ {
		env = env.enclosing
	}
	return env
}
