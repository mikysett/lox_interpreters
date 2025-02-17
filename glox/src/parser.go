package main

import (
	"fmt"
)

// program        → declaration* EOF ;
// declaration    → varDecl | statement ;
// varDecl        → "var" IDENTIFIER ( "=" expression )? ";" ;
// statement      → exprStmt
//                | ifStmt
//                | printStmt
//                | block ;

// if             → if "(" expression ")" statement ;
//                ( "else" statement )+ ;
// block          → "{" declaration* "}" ;
// exprStmt       → expression ";" ;
// printStmt      → "print" expression ";" ;
// expression     → commaSeparator ;
// commaSeparator → assignment ( ( "," ) assignment )* ;
// assignment     → IDENTIFIER "=" assignment
//                | ternary ;
// ternary        → ( logic_or "?" logic_or ":" )* logic_or ;
// logic_or       → logic_and ( "or" logic_and )* ;
// logic_and      → equality ( "and" equality )* ;
// equality       → comparison ( ( "!=" | "==" ) comparison )* ;
// comparison     → term ( ( ">" | ">=" | "<" | "<=" ) term )* ;
// term           → factor ( ( "-" | "+" ) factor )* ;
// factor         → unary ( ( "/" | "*" ) unary )* ;
// unary          → ( "!" | "-" ) unary | primary ;
// primary        → "true" | "false" | "nil"
//                | NUMBER | STRING
//                | "(" expression ")"
//                | IDENTIFIER ;

type ParseError struct {
	token   *Token
	message string
}

func (e *ParseError) Error() string {
	if e.token.Type == EOF {
		return fmt.Sprintf("[line %v] Error at end: %v", e.token.Line, e.message)
	}
	return fmt.Sprintf("[line %v] Error at '%v': %v", e.token.Line, e.token.Lexeme, e.message)
}

type Parser struct {
	tokens  []*Token
	current int
}

func NewParser(tokens []*Token) *Parser {
	return &Parser{
		tokens:  tokens,
		current: 0,
	}
}

func (p *Parser) parse() ([]Stmt, error) {
	statements := []Stmt{}
	for !p.isAtEnd() {
		stmt, err := p.declaration()
		if err != nil {
			return nil, err
		}
		statements = append(statements, stmt)
	}
	return statements, nil
}

func (p *Parser) declaration() (stmt Stmt, err error) {
	defer func() {
		// In case of error parser moves to end of statement
		// So it can catch further errors in one pass
		if err != nil {
			p.synchronize()
		}
	}()
	if p.match(Var) {
		return p.varDeclaration()
	}
	return p.statement()
}

func (p *Parser) varDeclaration() (stmt Stmt, err error) {
	name, err := p.consume(Identifier, "Expect variable name.")
	if err != nil {
		return nil, err
	}
	var initializer Expr
	if p.match(Equal) {
		initializer, err = p.expression()
		if err != nil {
			return nil, err
		}
	}
	_, err = p.consume(Semicolon, "Expect ';' after variable declaration.")
	if err != nil {
		return nil, err
	}

	return NewStmtVar(name, initializer), nil
}

func (p *Parser) statement() (Stmt, error) {
	if p.match(If) {
		return p.ifStatement()
	}
	if p.match(Print) {
		return p.printStatement()
	}
	if p.match(LeftBrace) {
		return p.blockStatement()
	}
	return p.expressionStatement()
}

func (p *Parser) ifStatement() (Stmt, error) {
	_, err := p.consume(LeftParen, "Expect '(' after 'if'.")
	if err != nil {
		return nil, err
	}
	condition, err := p.expression()
	if err != nil {
		return nil, err
	}
	_, err = p.consume(RightParen, "Expect ')' after if condition.")
	if err != nil {
		return nil, err
	}
	thenBranch, err := p.statement()
	if err != nil {
		return nil, err
	}
	var elseBranch Stmt
	if p.match(Else) {
		elseBranch, err = p.statement()
		if err != nil {
			return nil, err
		}
	}
	return NewStmtIf(condition, thenBranch, elseBranch), nil
}

func (p *Parser) blockStatement() (Stmt, error) {
	statements := []Stmt{}
	for !p.check(RightBrace) && !p.isAtEnd() {
		stmt, err := p.declaration()
		if err != nil {
			return nil, err
		}
		statements = append(statements, stmt)
	}
	_, err := p.consume(RightBrace, "Expect '}' after block.")
	if err != nil {
		return nil, err
	}
	return NewStmtBlock(statements), nil
}

func (p *Parser) printStatement() (Stmt, error) {
	value, err := p.expression()
	if err != nil {
		return nil, err
	}
	_, err = p.consume(Semicolon, "Expect ';' after value.")
	if err != nil {
		return nil, err
	}
	return NewStmtPrint(value), nil
}

func (p *Parser) expressionStatement() (Stmt, error) {
	value, err := p.expression()
	if err != nil {
		return nil, err
	}
	_, err = p.consume(Semicolon, "Expect ';' after value.")
	if err == nil {
		return NewStmtExpression(value), nil
	} else if isReplMode && p.isAtEnd() {
		// Mimic last expression evaluation in the REPL when no `;` is found
		return NewStmtPrint(value), nil
	}
	return nil, err
}

func (p *Parser) expression() (Expr, error) {
	return p.commaSeparator()
}

func (p *Parser) commaSeparator() (expr Expr, firstErr error) {
	expr, err := p.assignment()
	if err != nil {
		firstErr = err
	}

	for p.match(Comma) {
		operator := p.previous()
		right, err := p.assignment()
		if err != nil {
			if firstErr == nil {
				firstErr = err
			}
			break
		}
		expr = NewExprBinary(expr, operator, right)
	}
	if firstErr != nil {
		return nil, firstErr
	}

	return expr, nil
}

func (p *Parser) assignment() (Expr, error) {
	expr, err := p.ternary()
	if err != nil {
		return nil, err
	}

	if p.match(Equal) {
		equals := p.previous()
		value, err := p.assignment()
		if err != nil {
			return nil, err
		}

		if _, ok := expr.(*ExprVariable); !ok {
			return nil, NewError(equals, "Invalid assignment target.")
		}
		return NewExprAssign(expr.(*ExprVariable).name, value), nil
	}
	return expr, nil
}

func (p *Parser) ternary() (Expr, error) {
	expr, err := p.logical_or()
	if err != nil {
		return nil, err
	}

	for p.match(QuestionMark) {
		operator := p.previous()
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
		expr = NewExprTernary(operator, expr, left, right)
	}
	return expr, nil
}

func (p *Parser) logical_or() (Expr, error) {
	expr, err := p.logical_and()
	if err != nil {
		return nil, err
	}

	for p.match(Or) {
		operator := p.previous()
		right, err := p.logical_and()
		if err != nil {
			return nil, err
		}
		expr = NewExprLogical(expr, operator, right)
	}
	return expr, nil
}

func (p *Parser) logical_and() (Expr, error) {
	expr, err := p.equality()
	if err != nil {
		return nil, err
	}

	for p.match(And) {
		operator := p.previous()
		right, err := p.equality()
		if err != nil {
			return nil, err
		}
		expr = NewExprLogical(expr, operator, right)
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
		expr = NewExprBinary(expr, operator, right)
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
		expr = NewExprBinary(expr, operator, right)
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
		expr = NewExprBinary(expr, operator, right)
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
		expr = NewExprBinary(expr, operator, right)
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
		return NewExprUnary(operator, right), nil
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
	} else if p.match(Identifier) {
		return NewExprVariable(p.previous()), nil
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
	return p.tokens[p.current-1]
}

func (p *Parser) check(tokenType TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().Type == tokenType
}

func (p *Parser) peek() *Token {
	return p.tokens[p.current]
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
