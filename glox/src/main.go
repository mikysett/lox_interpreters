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

var memprofile = flag.String("memprofile", "", "write memory profile to `file`")

func main() {
	flag.Parse()
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

	err = interpreter.interpret(stmts)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return err
	}
	return nil
}
