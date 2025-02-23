package main

import (
	"fmt"
	"os"
)

// program        → declaration* EOF ;
// declaration    → funDecl
//                | varDecl
//                | statement ;
// funDecl        → "fun" function ;
// function       → IDENTIFIER "(" parameters? ")" block ;
// parameters     → IDENTIFIER ( "," IDENTIFIER )* ;
// varDecl        → "var" IDENTIFIER ( "=" commaOperator )? ";" ;
// statement      → exprStmt
//                | ifStmt
//                | printStmt
//                | returnStmt
//                | whileStmt
//                | forStmt
//                | block
//                | "break" ; // This is not context free as it is only valid in `while` and `for` loops

// returnStmt     → "return" expression? ";" ;
// while          → while "(" commaOperator ")" statement ;
// for            → for "(" ( varDecl | exprStmt ";" )
//                  commaOperator? ";"
//                  commaOperator? ")" statement ;
// if             → if "(" commaOperator ")" statement ;
//                ( "else" statement )+ ;
// block          → "{" declaration* "}" ;
// exprStmt       → commaOperator ";" ;
// printStmt      → "print" commaOperator ";" ;
// commaOperator → expression ( ( "," ) expression )* ;
// expression     → assignment ;
// assignment     → IDENTIFIER "=" assignment
//                | ternary ;
// ternary        → ( logic_or "?" logic_or ":" )* logic_or ;
// logic_or       → logic_and ( "or" logic_and )* ;
// logic_and      → equality ( "and" equality )* ;
// equality       → comparison ( ( "!=" | "==" ) comparison )* ;
// comparison     → term ( ( ">" | ">=" | "<" | "<=" ) term )* ;
// term           → factor ( ( "-" | "+" ) factor )* ;
// factor         → unary ( ( "/" | "*" ) unary )* ;
// unary          → ( "!" | "-" ) unary | call ;
// call           → primary ( "(" arguments? ")" )* ;
// arguments      → expression ( "," expression )* ;
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
	tokens           []*Token
	current          int
	acceptBreakCount int
	// To prevent interpreter execution on errors not triggering parser panic mode
	hadError bool
}

func NewParser(tokens []*Token) *Parser {
	return &Parser{
		tokens:  tokens,
		current: 0,
	}
}

func (p *Parser) parse() (statements []Stmt, err error) {
	// To ensure the parser returns all errors concatenated in order
	errors := []error{}

	for !p.isAtEnd() {
		stmt, declarationErr := p.declaration()
		if declarationErr != nil {
			errors = append(errors, declarationErr)
		}
		statements = append(statements, stmt)
	}

	if len(errors) == 0 {
		return statements, nil
	}

	err = errors[0]
	for i := 1; i < len(errors); i++ {
		err = fmt.Errorf("%w\n%w", err, errors[i])
	}
	return nil, err
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
	} else if p.match(Fun) {
		return p.function("function")
	}

	return p.statement()
}

func (p *Parser) function(kind string) (stmt Stmt, err error) {
	name, err := p.consume(Identifier, "Expect "+kind+" name.")
	if err != nil {
		return nil, err
	}

	_, err = p.consume(LeftParen, "Expect '(' after "+kind+" name.")
	if err != nil {
		return nil, err
	}

	parameters := []*Token{}
	for !p.check(RightParen) {
		if len(parameters) != 0 {
			_, err = p.consume(Comma, "Expect ',' between parameters.")
			if err != nil {
				return nil, err
			}
		}
		if len(parameters) >= 255 {
			// Error here is just shown but doesn't stop parser execution as the parser is not in panic mode
			p.hadError = true
			fmt.Fprintln(os.Stderr, NewParserError(p.peek(), "Can't have more than 255 arguments."))
		}

		param, err := p.consume(Identifier, "Expect parameter name.")
		if err != nil {
			return nil, err
		}
		parameters = append(parameters, param)
	}
	_, err = p.consume(RightParen, "Expect ')' after parameters.")
	if err != nil {
		return nil, err
	}

	_, err = p.consume(LeftBrace, "Expect '{' before "+kind+" body.")
	if err != nil {
		return nil, err
	}

	body, err := p.block()
	if err != nil {
		return nil, err
	}

	return NewStmtFunction(name, parameters, body), nil
}

func (p *Parser) block() (stmts []Stmt, err error) {
	for !p.check(RightBrace) && !p.isAtEnd() {
		stmt, err := p.declaration()
		if err != nil {
			return nil, err
		}
		stmts = append(stmts, stmt)
	}
	_, err = p.consume(RightBrace, "Expect '}' after block.")
	if err != nil {
		return nil, err
	}

	return stmts, nil
}

func (p *Parser) varDeclaration() (stmt Stmt, err error) {
	name, err := p.consume(Identifier, "Expect variable name.")
	if err != nil {
		return nil, err
	}
	var initializer Expr
	if p.match(Equal) {
		initializer, err = p.commaOperator()
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
	if p.match(While) {
		return p.whileStatement()
	}
	if p.match(For) {
		return p.forStatement()
	}
	if p.match(Print) {
		return p.printStatement()
	}
	if p.match(Return) {
		return p.returnStatement()
	}
	if p.match(LeftBrace) {
		return p.blockStatement()
	}
	if p.match(Break) {
		return p.breakStatement()
	}
	return p.expressionStatement()
}

func (p *Parser) ifStatement() (Stmt, error) {
	_, err := p.consume(LeftParen, "Expect '(' after 'if'.")
	if err != nil {
		return nil, err
	}
	condition, err := p.commaOperator()
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

func (p *Parser) whileStatement() (Stmt, error) {
	_, err := p.consume(LeftParen, "Expect '(' after 'while'.")
	if err != nil {
		return nil, err
	}
	condition, err := p.commaOperator()
	if err != nil {
		return nil, err
	}
	_, err = p.consume(RightParen, "Expect ')' after while condition.")
	if err != nil {
		return nil, err
	}

	p.acceptBreakCount += 1
	defer func() {
		p.acceptBreakCount -= 1
	}()

	body, err := p.statement()
	if err != nil {
		return nil, err
	}
	return NewStmtWhile(condition, body), nil
}

func (p *Parser) forStatement() (Stmt, error) {
	_, err := p.consume(LeftParen, "Expect '(' after 'for'.")
	if err != nil {
		return nil, err
	}

	var initializer Stmt
	if p.match(Semicolon) {
		initializer = nil
	} else if p.match(Var) {
		initializer, err = p.varDeclaration()
	} else {
		initializer, err = p.expressionStatement()
	}
	if err != nil {
		return nil, err
	}

	var condition Expr
	if !p.check(Semicolon) {
		condition, err = p.commaOperator()
		if err != nil {
			return nil, err
		}
	}

	_, err = p.consume(Semicolon, "Expect ';' after loop condition.")
	if err != nil {
		return nil, err
	}

	var increment Expr
	if !p.check(RightParen) {
		increment, err = p.commaOperator()
	}
	if err != nil {
		return nil, err
	}

	_, err = p.consume(RightParen, "Expect ')' after loop condition.")
	if err != nil {
		return nil, err
	}

	p.acceptBreakCount += 1
	defer func() {
		p.acceptBreakCount -= 1
	}()

	body, err := p.statement()
	if err != nil {
		return nil, err
	}

	if increment != nil {
		body = NewStmtBlock([]Stmt{
			body,
			NewStmtExpression(increment),
		})
	}

	if condition == nil {
		condition = NewExprLiteral(true)
	}
	body = NewStmtWhile(condition, body)

	if initializer != nil {
		body = NewStmtBlock([]Stmt{
			initializer,
			body,
		})
	}

	return body, nil
}

func (p *Parser) blockStatement() (Stmt, error) {
	statements, err := p.block()
	if err != nil {
		return nil, err
	}

	return NewStmtBlock(statements), nil
}

func (p *Parser) breakStatement() (Stmt, error) {
	breakToken := p.previous()
	_, err := p.consume(Semicolon, "Expect ';' after break.")
	if err != nil {
		return nil, err
	}

	if p.acceptBreakCount <= 0 {
		return nil, NewParserError(breakToken, "Only valid in 'while' and 'for' loops.")
	}
	return NewStmtBreak(breakToken), nil
}

func (p *Parser) printStatement() (Stmt, error) {
	value, err := p.commaOperator()
	if err != nil {
		return nil, err
	}
	_, err = p.consume(Semicolon, "Expect ';' after value.")
	if err != nil {
		return nil, err
	}
	return NewStmtPrint(value), nil
}

func (p *Parser) returnStatement() (stmt Stmt, err error) {
	keyword := p.previous()

	var value Expr
	if !p.check(Semicolon) {
		value, err = p.commaOperator()
		if err != nil {
			return nil, err
		}
	}

	_, err = p.consume(Semicolon, "Expect ';' after value.")
	if err != nil {
		return nil, err
	}
	return NewStmtReturn(keyword, value), nil
}

func (p *Parser) expressionStatement() (Stmt, error) {
	value, err := p.commaOperator()
	if err != nil {
		return nil, err
	}
	_, err = p.consume(Semicolon, "Expect ';' after expression.")
	if err == nil {
		return NewStmtExpression(value), nil
	} else if isReplMode && p.isAtEnd() {
		// Mimic last expression evaluation in the REPL when no `;` is found
		return NewStmtPrint(value), nil
	}
	return nil, err
}

func (p *Parser) commaOperator() (expr Expr, firstErr error) {
	expr, err := p.expression()
	if err != nil {
		firstErr = err
	}

	for p.match(Comma) {
		operator := p.previous()
		right, err := p.expression()
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

func (p *Parser) expression() (Expr, error) {
	return p.assignment()
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
			return nil, NewParserError(equals, "Invalid assignment target.")
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
			return nil, NewParserError(p.peek(), "Expect :.")
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
	return p.call()
}

func (p *Parser) call() (Expr, error) {
	expr, err := p.primary()
	if err != nil {
		return nil, err
	}

	for {
		if p.match(LeftParen) {
			expr, err = p.finishCall(expr)
			if err != nil {
				return nil, err
			}
		} else {
			break
		}
	}

	return expr, nil
}

func (p *Parser) finishCall(expr Expr) (Expr, error) {
	arguments := []Expr{}
	if !p.check(RightParen) {
		for {
			if len(arguments) >= 255 {
				// Error here is just shown but doesn't stop parser execution as the parser is not in panic mode
				p.hadError = true
				fmt.Fprintln(os.Stderr, NewParserError(p.peek(), "Can't have more than 255 arguments."))
			}

			arg, err := p.expression()
			if err != nil {
				return nil, err
			}
			arguments = append(arguments, arg)

			if !p.match(Comma) {
				break
			}
		}
	}

	paren, err := p.consume(RightParen, "Expect ')' after arguments.")
	if err != nil {
		return nil, err
	}

	return NewExprCall(expr, paren, arguments), nil
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
	return nil, NewParserError(p.peek(), "Expect expression.")
}

// If the current Token.type matches one of the given types returns `true` and advance the parser's cursor
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

// Return `true` if the current token == [tokenType] but do not move the parser's cursor
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
	return nil, NewParserError(p.peek(), message)
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

func NewParserError(token *Token, message string) *ParseError {
	return &ParseError{
		token:   token,
		message: message,
	}
}
