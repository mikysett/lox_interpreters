package main

type LoxClass struct {
	superclass *LoxClass
	metaclass  *LoxInstance
	name       string
	methods    map[string]*Function
}

func NewLoxClass(superclass *LoxClass, metaclass *LoxInstance, name string, methods map[string]*Function) *LoxClass {
	return &LoxClass{
		superclass: superclass,
		metaclass:  metaclass,
		name:       name,
		methods:    methods,
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
	if val, ok := c.methods[name]; ok {
		return val
	}
	if c.superclass != nil {
		return c.superclass.FindMethod(name)
	}
	return nil
}

func (c *LoxClass) String() string {
	return c.name
}
