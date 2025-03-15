package main

type LoxClass struct {
	name    string
	methods map[string]*Function
}

func NewLoxClass(name string, methods map[string]*Function) *LoxClass {
	return &LoxClass{
		name:    name,
		methods: methods,
	}
}

func (c *LoxClass) arity() int {
	if initializer := c.FindMethod("init"); initializer != nil {
		return initializer.arity()
	}

	return 0
}

func (c *LoxClass) call(interpreter *Interpreter, arguments []any) (any, error) {
	instance := NewLoxInstance(c)

	if initializer := c.FindMethod("init"); initializer != nil {
		initializer.Bind(instance).call(interpreter, arguments)
	}

	return instance, nil
}

func (c *LoxClass) FindMethod(name string) *Function {
	return c.methods[name]
}

func (c *LoxClass) String() string {
	return c.name
}
