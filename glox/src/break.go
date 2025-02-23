package main

import "fmt"

type BreakShortCircuit struct {
	value any
}

func NewBreakShortCircuit() *BreakShortCircuit {
	return &BreakShortCircuit{}
}

func (e *BreakShortCircuit) Error() string {
	return fmt.Sprintf("'break' statement outside of loop scope")
}
