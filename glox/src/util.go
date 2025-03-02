package main

import (
	"fmt"
	"os"
)

// For a more convenient wrapper use [printError]
// Set global [hadError] to true and writes the error to stderr
func report(line int, where, message string) {
	hadError = true
	fmt.Fprintf(os.Stderr, "[line %v] Error%v: %v\n", line, where, message)
}

// Set global [hadError] to true and writes the error to stderr
func printError(token *Token, message string) {
	if token.Type == EOF {
		report(token.Line, " at end", message)
	}
	report(token.Line, " at '"+token.Lexeme+"'", message)
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
