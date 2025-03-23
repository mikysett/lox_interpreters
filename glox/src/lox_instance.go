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

func NewArrayInstance() *LoxInstance {
	return NewLoxInstance(NewLoxClass(nil, nil, "Array", map[string]*Function{}))
}

func (i *LoxInstance) Get(name *Token) (any, error) {
	if val, ok := i.fields[name.Lexeme]; ok {
		return val, nil
	}

	method := i.class.FindMethod(name.Lexeme)
	if method != nil {
		return method.Bind(i), nil
	}

	return nil, NewRuntimeError(name, "Undefined property '"+name.Lexeme+"'.")
}

func (i *LoxInstance) Set(name *Token, value any) {
	i.fields[name.Lexeme] = value
}

func (i *LoxInstance) String() string {
	return i.class.name + " instance"
}
