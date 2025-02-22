package main

type Callable struct {
	arity func() int
	call  func(interpreter *Interpreter, arguments []any) (any, error)
}

func NewCallable(arity func() int, call func(interpreter *Interpreter, arguments []any) (any, error)) *Callable {
	return &Callable{
		arity: arity,
		call:  call,
	}
}
