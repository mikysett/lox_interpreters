package main

import "fmt"

type ContinueShortCircuit struct{}

func NewContinueShortCircuit() *ContinueShortCircuit {
	return &ContinueShortCircuit{}
}

func (e *ContinueShortCircuit) Error() string {
	return fmt.Sprintf("'continue' short circuit statement")
}
