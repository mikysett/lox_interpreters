package main

type Environment struct {
	// Parent-pointer tree (cactus stack)
	enclosing   *Environment
	localValues []any
}

type Uninitialized struct{}

func NewEnvironment() *Environment {
	return &Environment{
		enclosing:   nil,
		localValues: []any{},
	}
}

func (env *Environment) WithEnclosing(enclosing *Environment) *Environment {
	env.enclosing = enclosing
	return env
}

func (env *Environment) define(v any) {
	env.localValues = append(env.localValues, v)
}

func (env *Environment) assignAt(position *Position, value any) {
	env.ancestor(position.depth).localValues[position.index] = value
}

func (env *Environment) getAt(position *Position) any {
	return env.ancestor(position.depth).localValues[position.index]
}

func (env *Environment) ancestor(distance int) *Environment {
	for range distance {
		env = env.enclosing
	}
	return env
}
