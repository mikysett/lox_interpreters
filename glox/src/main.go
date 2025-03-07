package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
)

// Check [sysexits.h](https://man.freebsd.org/cgi/man.cgi?query=sysexits&apropos=0&sektion=0&manpath=FreeBSD+4.3-RELEASE&format=html)
const (
	exUsage      = 64
	exDataErr    = 65
	exRuntimeErr = 70
)

// To prevent interpreter execution on errors not triggering parser panic mode
var hadError = false

// When `true` expressions will be evaluated in the REPL instead of throwing an error
// For example: `3 < 2` will print `false` in the REPL and throw an error in a file.
var isReplMode = false

var (
	memprofile    = flag.String("memprofile", "", "write memory profile to `file`")
	disableExtras = flag.Bool("disable-extras", false, "exclude extra features (`false` by default)")
)

func main() {
	flag.Parse()
	if *disableExtras {
		GlobalConfig = BasicConfig
	}

	if len(flag.Args()) > 1 {
		println("Usage: glox [script]")
		os.Exit(exUsage)
	} else if flag.Arg(0) != "" {
		err := runFile(flag.Arg(0))
		if err != nil {
			switch err.(type) {
			case *RuntimeError:
				os.Exit(exRuntimeErr)
			case *ParseError:
				os.Exit(exDataErr)
			// Scanner
			default:
				os.Exit(exDataErr)
			}
		}
	} else {
		err := runPrompt()
		if err != nil {
			os.Exit(exDataErr)
		}
	}
	if *memprofile != "" {
		saveMemProfile(*memprofile)
	}
}

func saveMemProfile(fileName string) {
	f, err := os.Create(fileName)
	if err != nil {
		log.Fatal("could not create memory profile: ", err)
	}
	defer f.Close()
	runtime.GC()
	if err := pprof.Lookup("allocs").WriteTo(f, 0); err != nil {
		log.Fatal("could not write memory profile: ", err)
	}
}

func runFile(filePath string) error {
	bytes, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	return run(string(bytes), NewInterpreter())
}

func runPrompt() error {
	isReplMode = true
	interpreter := NewInterpreter()
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		line, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		_ = run(line, interpreter)
	}
}

func run(source string, interpreter *Interpreter) (err error) {
	scanner := NewScanner(source)
	err = scanner.scanTokens()
	if err != nil {
		return err
	}

	parser := NewParser(scanner.Tokens)
	stmts, err := parser.parse()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return err
	}

	// For errors not propagated to the `parse()` return
	if hadError {
		return NewParserError(parser.peek(), "Don't run interpreter due to previous errors.")
	}

	resolver := NewResolver(interpreter)
	// The resolver never returns errors so we can safely skip the check
	resolver.resolveStmts(stmts)

	if hadError {
		return NewParserError(parser.peek(), "Don't run interpreter due to previous errors.")
	}

	err = interpreter.interpret(stmts)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return err
	}
	return nil
}
