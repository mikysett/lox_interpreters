package main

type Callable interface {
	arity() int
	call(interpreter *Interpreter, arguments []any) (any, error)
	String() string
}

type ProtoCallable struct {
	arityFunc    func() int
	callFunc     func(interpreter *Interpreter, arguments []any) (any, error)
	toStringFunc func() string
}

func NewProtoCallable(arity func() int, call func(interpreter *Interpreter, arguments []any) (any, error), toString func() string) *ProtoCallable {
	return &ProtoCallable{
		arityFunc:    arity,
		callFunc:     call,
		toStringFunc: toString,
	}
}

func (pc *ProtoCallable) arity() int {
	return pc.arityFunc()
}

func (pc *ProtoCallable) call(interpreter *Interpreter, arguments []any) (any, error) {
	return pc.callFunc(interpreter, arguments)
}

func (pc *ProtoCallable) String() string {
	return pc.toStringFunc()
}
