package main

type Callable interface {
	arity() int
	call(interpreter *Interpreter, arguments []any) (any, error)
}
