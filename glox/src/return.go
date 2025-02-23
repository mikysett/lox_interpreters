package main

import "fmt"

type ReturnShortCircuit struct {
	value any
}

func NewReturnShortCircuit(value any) *ReturnShortCircuit {
	return &ReturnShortCircuit{
		value: value,
	}
}

func (e *ReturnShortCircuit) Error() string {
	return fmt.Sprintf("'return' statement outside of function scope, value not catched: %v", e.value)
}
