package main

import (
	"fmt"
	"strconv"
)

type Scanner struct {
	Source   string
	Tokens   []Token
	keywords map[string]TokenType
	start    int
	current  int
	line     int
}

var keywords = map[string]TokenType{
	"and":    And,
	"class":  Class,
	"else":   Else,
	"false":  False,
	"for":    For,
	"fun":    Fun,
	"if":     If,
	"nil":    Nil,
	"or":     Or,
	"print":  Print,
	"return": Return,
	"super":  Super,
	"this":   This,
	"true":   True,
	"var":    Var,
	"while":  While,
}

func NewScanner(source string) Scanner {
	return Scanner{
		Source:   source,
		Tokens:   []Token{},
		keywords: keywords,
		line:     1,
	}
}

func (scanner *Scanner) scanTokens() (err error) {
	for !scanner.isAtEnd() {
		scanner.start = scanner.current
		tokenErr := scanner.scanToken()
		if tokenErr != nil {
			err = tokenErr
		}
	}

	scanner.Tokens = append(scanner.Tokens, Token{
		Type:    EOF,
		Lexeme:  "",
		Literal: nil,
		Line:    scanner.line,
	})
	return err
}

func (scanner *Scanner) isAtEnd() bool {
	return scanner.current >= len(scanner.Source)
}

func (scanner *Scanner) scanToken() (err error) {
	c := scanner.advance()
	switch c {
	case '(':
		scanner.addToken(LeftParen)
	case ')':
		scanner.addToken(RightParen)
	case '{':
		scanner.addToken(LeftBrace)
	case '}':
		scanner.addToken(RightBrace)
	case ',':
		scanner.addToken(Comma)
	case '.':
		scanner.addToken(Dot)
	case '-':
		scanner.addToken(Minus)
	case '+':
		scanner.addToken(Plus)
	case ';':
		scanner.addToken(Semicolon)
	case '*':
		scanner.addToken(Star)
	case '!':
		if scanner.match('=') {
			scanner.addToken(BangEqual)
		} else {
			scanner.addToken(Bang)
		}
	case '=':
		if scanner.match('=') {
			scanner.addToken(EqualEqual)
		} else {
			scanner.addToken(Equal)
		}
	case '<':
		if scanner.match('=') {
			scanner.addToken(LessEqual)
		} else {
			scanner.addToken(Less)
		}
	case '>':
		if scanner.match('=') {
			scanner.addToken(GreaterEqual)
		} else {
			scanner.addToken(Greater)
		}
	case '/':
		if scanner.match('/') {
			for char, err := scanner.peek(); err == nil && char != '\n'; char, err = scanner.peek() {
				scanner.advance()
			}
		} else if scanner.match('*') {
			scanner.blockComment()
		} else {
			scanner.addToken(Slash)
		}
	case '"':
		err = scanner.stringLiteral()
	case ' ':
	case '\r':
	case '\t':
	case '\n':
		scanner.line++
	default:
		if IsDigit(c) {
			err = scanner.numberLiteral()
		} else if IsAlpha(c) {
			err = scanner.identifier()
		} else {
			err = fmt.Errorf("Unexpected character.")
		}
	}

	if err != nil {
		report(scanner.line, "", err.Error())
	}
	return err
}

func (scanner *Scanner) addToken(tokenType TokenType) {
	scanner.addTokenWithLiteral(tokenType, nil)
}

func (scanner *Scanner) advance() byte {
	char := scanner.Source[scanner.current]
	scanner.current += 1
	return char
}

func (scanner *Scanner) match(char byte) bool {
	if !scanner.isAtEnd() &&
		scanner.Source[scanner.current] == char {
		scanner.current += 1
		return true
	}
	return false
}

func (scanner *Scanner) peek() (byte, error) {
	if scanner.isAtEnd() {
		return '\n', fmt.Errorf("EOF")
	}
	return scanner.Source[scanner.current], nil
}

func (scanner *Scanner) peekNext() (byte, error) {
	if scanner.current+1 >= len(scanner.Source) {
		return '\n', fmt.Errorf("EOF")
	}
	return scanner.Source[scanner.current+1], nil
}

func (scanner *Scanner) stringLiteral() error {
	for char, err := scanner.peek(); err == nil && char != '"'; char, err = scanner.peek() {
		if char == '\n' {
			scanner.line++
		}
		scanner.advance()
	}
	if scanner.isAtEnd() {
		return fmt.Errorf("Unterminated string.")
	}

	// Consuming the closing '"'
	scanner.advance()

	scanner.addTokenWithLiteral(String, scanner.Source[scanner.start+1:scanner.current-1])
	return nil
}

func (scanner *Scanner) numberLiteral() error {
	var char byte
	var err error
	for char, err = scanner.peek(); err == nil && IsDigit(char); char, err = scanner.peek() {
		scanner.advance()
	}

	nextChar, nextErr := scanner.peekNext()
	if nextErr == nil && err == nil &&
		char == '.' && IsDigit(nextChar) {
		// Consume the '.'
		scanner.advance()
		for char, err = scanner.peek(); err == nil && IsDigit(char); char, err = scanner.peek() {
			scanner.advance()
		}
	}
	number, err := strconv.ParseFloat(scanner.Source[scanner.start:scanner.current], 64)
	if err != nil {
		return err
	}

	scanner.addTokenWithLiteral(Number, number)
	return nil
}

func (scanner *Scanner) identifier() error {
	var identifierType TokenType = Identifier

	for char, err := scanner.peek(); err == nil && IsAlphaNumeric(char); char, err = scanner.peek() {
		scanner.advance()
	}

	text := scanner.Source[scanner.start:scanner.current]
	standardType, ok := scanner.keywords[text]
	if ok {
		identifierType = standardType
	}

	scanner.addToken(identifierType)
	return nil
}

func (scanner *Scanner) blockComment() error {
	for char, err := scanner.peek(); err == nil; char, err = scanner.peek() {
		if char == '*' {
			nextChar, err := scanner.peekNext()
			if err == nil && nextChar == '/' {
				break
			}
		} else if char == '/' {
			nextChar, err := scanner.peekNext()
			if err == nil && nextChar == '*' {
				scanner.advance()
				scanner.advance()
				scanner.blockComment()
			}
		} else {
			if char == '\n' {
				scanner.line++
			}
			scanner.advance()
		}
	}
	if scanner.isAtEnd() {
		return fmt.Errorf("Unterminated block comment.")
	}

	// Consuming the closing '*/'
	scanner.advance()
	scanner.advance()

	return nil
}

func (scanner *Scanner) addTokenWithLiteral(tokenType TokenType, literal any) {
	scanner.Tokens = append(scanner.Tokens, Token{
		Type:    tokenType,
		Literal: literal,
		Lexeme:  scanner.Source[scanner.start:scanner.current],
		Line:    scanner.line,
	})
}
