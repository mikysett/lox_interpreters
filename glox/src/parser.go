package main

import (
	"fmt"
)

// expression     → commaSeparator ;
// commaSeparator → ternary ( ( "," ) ternary )* ;
// ternary        → ( equality "?" equality ":" )* equality ;
// equality       → comparison ( ( "!=" | "==" ) comparison )* ;
// comparison     → term ( ( ">" | ">=" | "<" | "<=" ) term )* ;
// term           → factor ( ( "-" | "+" ) factor )* ;
// factor         → unary ( ( "/" | "*" ) unary )* ;
// unary          → ( "!" | "-" ) unary | primary ;
// primary        → NUMBER | STRING | "true" | "false" | "nil" | "(" expression ")" ;

type ParseError struct {
	token   *Token
	message string
}

func (e *ParseError) Error() string {
	if e.token.Type == EOF {
		return fmt.Sprintf("Line %v: at end. %v", e.token.Line, e.message)
	}
	return fmt.Sprintf("Line %v: at '%v'. %v", e.token.Line, e.token.Lexeme, e.message)
}

type Parser struct {
	tokens  []Token
	current int
}

func NewParser(tokens []Token) *Parser {
	return &Parser{
		tokens:  tokens,
		current: 0,
	}
}

func (p *Parser) parse() (Expr, error) {
	expr, err := p.expression()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return expr, nil
}

func (p *Parser) expression() (Expr, error) {
	return p.commaSeparator()
}

func (p *Parser) commaSeparator() (expr Expr, firstErr error) {
	expr, err := p.ternary()
	if err != nil {
		firstErr = err
	}

	for p.match(Comma) {
		operator := p.previous()
		right, err := p.ternary()
		if err != nil {
			if firstErr == nil {
				firstErr = err
			}
			break
		}
		expr = NewExprBinary(expr, *operator, right)
	}
	if firstErr != nil {
		return nil, firstErr
	}

	return expr, nil
}

func (p *Parser) ternary() (Expr, error) {
	expr, err := p.equality()
	if err != nil {
		return nil, err
	}

	for p.match(QuestionMark) {
		left, err := p.ternary()
		if err != nil {
			return nil, err
		}
		var right Expr
		if p.match(Colon) {
			right, err = p.ternary()
			if err != nil {
				return nil, err
			}
		} else {
			return nil, NewError(p.peek(), "Expect :.")
		}
		expr = NewExprTernary(expr, left, right)
	}
	return expr, nil
}

func (p *Parser) equality() (expr Expr, firstErr error) {
	expr, err := p.comparison()
	if err != nil {
		firstErr = err
	}

	for p.match(BangEqual, EqualEqual) {
		operator := p.previous()
		right, err := p.comparison()
		if err != nil {
			if firstErr == nil {
				firstErr = err
			}
			break
		}
		expr = NewExprBinary(expr, *operator, right)
	}
	if firstErr != nil {
		return nil, firstErr
	}

	return expr, nil
}

func (p *Parser) comparison() (expr Expr, firstErr error) {
	expr, err := p.term()
	if err != nil {
		firstErr = err
	}

	for p.match(Greater, GreaterEqual, Less, LessEqual) {
		operator := p.previous()
		right, err := p.term()
		if err != nil {
			if firstErr == nil {
				firstErr = err
			}
			break
		}
		expr = NewExprBinary(expr, *operator, right)
	}
	if firstErr != nil {
		return nil, firstErr
	}

	return expr, nil
}

func (p *Parser) term() (expr Expr, firstErr error) {
	expr, err := p.factor()
	if err != nil {
		firstErr = err
	}

	for p.match(Minus, Plus) {
		operator := p.previous()
		right, err := p.factor()
		if err != nil {
			if firstErr == nil {
				firstErr = err
			}
			break
		}
		expr = NewExprBinary(expr, *operator, right)
	}
	if firstErr != nil {
		return nil, firstErr
	}

	return expr, nil
}

func (p *Parser) factor() (expr Expr, firstErr error) {
	expr, err := p.unary()
	if err != nil {
		return nil, err
	}

	for p.match(Slash, Star) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		expr = NewExprBinary(expr, *operator, right)
	}
	return expr, nil
}

func (p *Parser) unary() (Expr, error) {
	if p.match(Bang, Minus) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		return NewExprUnary(*operator, right), nil
	}
	return p.primary()
}

func (p *Parser) primary() (Expr, error) {
	if p.match(False) {
		return NewExprLiteral(false), nil
	} else if p.match(True) {
		return NewExprLiteral(true), nil
	} else if p.match(Nil) {
		return NewExprLiteral(nil), nil
	} else if p.match(Number, String) {
		return NewExprLiteral(p.previous().Literal), nil
	} else if p.match(LeftParen) {
		expr, err := p.expression()
		if err != nil {
			return nil, err
		}
		_, err = p.consume(RightParen, "Expect ')' after expression.")
		if err != nil {
			return nil, err
		}
		return NewExprGrouping(expr), nil
	}
	return nil, NewError(p.peek(), "Expect expression.")
}

func (p *Parser) match(types ...TokenType) bool {
	for _, currType := range types {
		if p.check(currType) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) advance() *Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *Parser) previous() *Token {
	return &p.tokens[p.current-1]
}

func (p *Parser) check(tokenType TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().Type == tokenType
}

func (p *Parser) peek() *Token {
	return &p.tokens[p.current]
}

func (p *Parser) isAtEnd() bool {
	return p.tokens[p.current].Type == EOF
}

func (p *Parser) consume(tokenType TokenType, message string) (*Token, error) {
	if p.check(tokenType) {
		return p.advance(), nil
	}
	return nil, NewError(p.peek(), message)
}

func (p *Parser) synchronize() {
	p.advance()
	for !p.isAtEnd() {
		if p.previous().Type == Semicolon {
			return
		}
		switch p.peek().Type {
		case Class, Fun, Var, For, If, While, Print, Return:
			return
		}

		p.advance()
	}
}

func NewError(token *Token, message string) *ParseError {
	return &ParseError{
		token:   token,
		message: message,
	}
}
