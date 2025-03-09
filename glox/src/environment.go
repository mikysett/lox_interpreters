package main

type Environment struct {
	// Parent-pointer tree (cactus stack)
	enclosing    *Environment
	globalValues map[string]any
	localValues  []any
}

type Uninitialized struct{}

func NewGlobalEnvironment() *Environment {
	return &Environment{
		enclosing:    nil,
		localValues:  nil,
		globalValues: map[string]any{},
	}
}

func NewLocalEnvironment() *Environment {
	return &Environment{
		enclosing:    nil,
		localValues:  []any{},
		globalValues: nil,
	}
}

func (env *Environment) WithEnclosing(enclosing *Environment) *Environment {
	env.enclosing = enclosing
	return env
}

func (env *Environment) defineGlobal(k string, v any) {
	env.globalValues[k] = v
}

func (env *Environment) define(v any) {
	env.localValues = append(env.localValues, v)
}

func (env *Environment) assign(name *Token, value any) error {
	_, ok := env.globalValues[name.Lexeme]
	if ok {
		env.globalValues[name.Lexeme] = value
		return nil
	}
	return NewRuntimeError(name, "Undefined variable '"+name.Lexeme+"'.")
}

func (env *Environment) assignAt(position *Position, value any) {
	env.ancestor(position.depth).localValues[position.index] = value
}

func (env *Environment) get(name *Token) (any, error) {
	value, ok := env.globalValues[name.Lexeme]
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

func (env *Environment) getAt(position *Position) (any, error) {
	return env.ancestor(position.depth).localValues[position.index], nil
}

func (env *Environment) ancestor(distance int) *Environment {
	for range distance {
		env = env.enclosing
	}
	return env
}
