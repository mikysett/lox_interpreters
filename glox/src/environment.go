package main

type Environment struct {
	// Parent-pointer tree (cactus stack)
	enclosing *Environment
	values    []any
}

type Uninitialized struct{}

func NewEnvironment() *Environment {
	return &Environment{
		enclosing: nil,
		values:    []any{},
	}
}

func (env *Environment) withEnclosing(enclosing *Environment) *Environment {
	env.enclosing = enclosing
	return env
}

func (env *Environment) define(v any) {
	env.values = append(env.values, v)
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

func (env *Environment) assignAt(position Position, value any) {
	env.ancestor(position.depth).values[position.index] = value
}

func (env *Environment) get(name *Token) (any, error) {
	value, ok := env.values[name.Lexeme]
	if !ok && env.enclosing != nil {
		return env.enclosing.get(name)
	} else if !ok {
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

func (env *Environment) getAt(position Position) any {
	return env.ancestor(position.depth).values[position.index]
}

func (env *Environment) ancestor(distance int) *Environment {
	for i := 0; i < distance; i++ {
		env = env.enclosing
	}
	return env
}
