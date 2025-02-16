package main

import (
	"fmt"
	"os"
)

func report(line int, where, message string) {
	fmt.Fprintf(os.Stderr, "[line %v] Error%v: %v\n", line, where, message)
}

func IsDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func IsAlpha(c byte) bool {
	return (c >= 'a' && c <= 'z') ||
		(c >= 'A' && c <= 'Z') ||
		c == '_'
}

func IsAlphaNumeric(c byte) bool {
	return IsAlpha(c) || IsDigit(c)
}
