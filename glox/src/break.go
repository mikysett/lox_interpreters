package main

import "fmt"

type BreakShortCircuit struct{}

func NewBreakShortCircuit() *BreakShortCircuit {
	return &BreakShortCircuit{}
}

func (e *BreakShortCircuit) Error() string {
	return fmt.Sprintf("'break' short circuit statement")
}
