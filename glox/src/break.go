package main

import "fmt"

type (
	BreakShortCircuit    struct{}
	ContinueShortCircuit struct{}
)

func NewBreakShortCircuit() *BreakShortCircuit {
	return &BreakShortCircuit{}
}

func NewContinueShortCircuit() *ContinueShortCircuit {
	return &ContinueShortCircuit{}
}

func (e *BreakShortCircuit) Error() string {
	return fmt.Sprintf("'break' short circuit statement")
}

func (e *ContinueShortCircuit) Error() string {
	return fmt.Sprintf("'continue' short circuit statement")
}
