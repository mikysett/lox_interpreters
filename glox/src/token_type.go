package main

type TokenType int

const (
	// Single-character tokens.
	LeftParen    = iota // Char: '('
	RightParen          // Char: ')'
	LeftBrace           // Char: '{'
	RightBrace          // Char: '}'
	Comma               // Char: ','
	Dot                 // Char: '.'
	Minus               // Char: '-'
	Plus                // Char: '+'
	Semicolon           // Char: ";"
	Slash               // Char: "/"
	Star                // Char: '*'
	QuestionMark        // Char: '?'
	Colon               // Char: ':'

	// One or two character tokens.
	Bang         // Char: '!'
	BangEqual    // Char: '!='
	Equal        // Char: '='
	EqualEqual   // Char: '=='
	Greater      // Char: '>'
	GreaterEqual // Char: '>='
	Less         // Char: '<'
	LessEqual    // Char: '<='

	// Literals.
	Identifier
	String
	Number

	// Keywords.
	And    // Char: 'and'
	Class  // Char: 'class'
	Else   // Char: 'else'
	False  // Char: 'false'
	Fun    // Char: 'fun'
	For    // Char: 'for'
	If     // Char: 'if'
	Nil    // Char: 'nil'
	Or     // Char: 'or'
	Print  // Char: 'print'
	Return // Char: 'return'
	Super  // Char: 'super'
	This   // Char: 'this'
	True   // Char: 'true'
	Var    // Char: 'var'
	While  // Char: 'while'

	EOF
)

func (t *TokenType) toString() string {
	switch *t {
	case LeftParen:
		return "LeftParen"
	case RightParen:
		return "RightParen"
	case LeftBrace:
		return "LeftBrace"
	case RightBrace:
		return "RightBrace"
	case Comma:
		return "Comma"
	case Dot:
		return "Dot"
	case Minus:
		return "Minus"
	case Plus:
		return "Plus"
	case Semicolon:
		return "Semicolon"
	case Slash:
		return "Slash"
	case Star:
		return "Star"
	case Bang:
		return "Bang"
	case BangEqual:
		return "BangEqual"
	case Equal:
		return "Equal"
	case EqualEqual:
		return "EqualEqual"
	case Greater:
		return "Greater"
	case GreaterEqual:
		return "GreaterEqual"
	case Less:
		return "Less"
	case LessEqual:
		return "LessEqual"
	case Identifier:
		return "Identifier"
	case String:
		return "String"
	case Number:
		return "Number"
	case And:
		return "And"
	case Class:
		return "Class"
	case Else:
		return "Else"
	case False:
		return "False"
	case Fun:
		return "Fun"
	case For:
		return "For"
	case If:
		return "If"
	case Nil:
		return "Nil"
	case Or:
		return "Or"
	case Print:
		return "Print"
	case Return:
		return "Return"
	case Super:
		return "Super"
	case This:
		return "This"
	case True:
		return "True"
	case Var:
		return "Var"
	case While:
		return "While"
	case EOF:
		return "EOF"
	default:
		return "UNKNOWN"
	}
}
