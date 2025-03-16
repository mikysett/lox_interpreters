package main

type LoxInstance struct {
	class  *LoxClass
	fields map[string]any
}

func NewLoxInstance(class *LoxClass) *LoxInstance {
	return &LoxInstance{
		class:  class,
		fields: map[string]any{},
	}
}

func (i *LoxInstance) Get(name *Token) (any, error) {
	getter := i.class.FindGetter(name.Lexeme)
	if getter != nil {
		return getter.Bind(i), nil
	}

	if val, ok := i.fields[name.Lexeme]; ok {
		return val, nil
	}

	method := i.class.FindMethod(name.Lexeme)
	if method != nil {
		return method.Bind(i), nil
	}

	return nil, NewRuntimeError(name, "Undefined property '"+name.Lexeme+"'.")
}

func (i *LoxInstance) Set(name *Token, value any) error {
	if i.class.FindGetter(name.Lexeme) != nil {
		return NewRuntimeError(name, "Can't override a getter: '"+name.Lexeme+"'.")
	}
	i.fields[name.Lexeme] = value
	return nil
}

func (i *LoxInstance) String() string {
	return i.class.name + " instance"
}
