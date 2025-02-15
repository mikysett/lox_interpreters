package main

import (
	"bufio"
	"fmt"
	"os"
)

// Check [sysexits.h](https://man.freebsd.org/cgi/man.cgi?query=sysexits&apropos=0&sektion=0&manpath=FreeBSD+4.3-RELEASE&format=html)
const (
	exUsage   = 64
	exDataErr = 65
)

func main() {
	args := os.Args
	if len(args) > 2 {
		println("Usage: glox [script]")
		os.Exit(exUsage)
	} else if len(args) == 2 {
		err := runFile(args[1])
		if err != nil {
			os.Exit(exDataErr)
		}
	} else {
		err := runPrompt()
		if err != nil {
			os.Exit(exDataErr)
		}
	}
}

func runFile(filePath string) error {
	bytes, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	return run(string(bytes))
}

func runPrompt() error {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		line, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		err = run(line)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func run(source string) (err error) {
	scanner := NewScanner(source)
	err = scanner.scanTokens()
	if err != nil {
		return err
	}

	parser := NewParser(scanner.Tokens)
	expr, err := parser.parse()
	if err != nil {
		return err
	}

	result, err := NewInterpreter().evaluate(expr)
	if err != nil {
		return err
	}
	fmt.Println(result)

	return nil
}
