#[cfg(debug)]
use crate::debug::debug::disassemble_chunk;
use crate::domain::Chunk;
use crate::domain::OpCode;
use crate::domain::Value;
use crate::scanner::Scanner;
use crate::scanner::Token;
use crate::scanner::TokenType;
use crate::vm::InterpretError;

pub struct Parser {
    current: Token,
    previous: Token,
    scanner: Scanner,
    had_error: bool,
    panic_mode: bool,
}

pub enum ParserToken {
    Current,
    Previous,
}

#[derive(Debug, Clone, Copy, PartialEq, PartialOrd)]
#[repr(u8)]
pub enum Precedence {
    None = 0,
    Assignment, // =
    Or,         // or
    And,        // and
    Equality,   // == !=
    Comparison, // < > <= >=
    Term,       // + -
    Factor,     // * /
    Unary,      // ! -
    Call,       // . ()
    Primary,
}

type ParseFn = fn(&mut Compiler);

pub struct ParserRule {
    prefix: Option<ParseFn>,
    infix: Option<ParseFn>,
    precedence: Precedence,
}

impl ParserRule {
    pub const fn new(
        prefix: Option<ParseFn>,
        infix: Option<ParseFn>,
        precedence: Precedence,
    ) -> Self {
        Self {
            prefix,
            infix,
            precedence,
        }
    }
}

impl Precedence {
    pub fn from_unchecked(value: u8) -> Self {
        unsafe { std::mem::transmute(value) }
    }
}

impl Parser {
    pub fn new(scanner: Scanner, dummy_token: Token) -> Self {
        Self {
            current: dummy_token,
            previous: dummy_token,
            scanner,
            had_error: false,
            panic_mode: false,
        }
    }

    pub fn advance(&mut self) {
        self.previous = self.current;

        loop {
            self.current = self.scanner.scan_token();
            if self.current.kind != TokenType::Error {
                break;
            }
            self.error_at_current(self.current.as_str().to_owned().as_str());
        }
    }

    pub fn consume(&mut self, token_type: TokenType, message: &str) {
        if self.current.kind == token_type {
            self.advance();
        } else {
            self.error_at_current(message);
        }
    }

    pub fn error_at_current(&mut self, message: &str) {
        self.error_at(ParserToken::Current, message);
    }

    pub fn error(&mut self, message: &str) {
        self.error_at(ParserToken::Previous, message);
    }

    pub fn error_at(&mut self, parser_token: ParserToken, message: &str) {
        if self.panic_mode {
            return;
        }

        self.had_error = true;
        self.panic_mode = true;

        let token = match parser_token {
            ParserToken::Current => self.current,
            ParserToken::Previous => self.previous,
        };
        print!("[line {}] Error", token.line);

        if token.kind == TokenType::Eof {
            eprint!(" at end");
        } else if token.kind == TokenType::Error {
            // Nothing.
        } else {
            eprint!(" at '{}'", token.as_str());
        }

        eprintln!(": {}", message);
    }
}

pub struct Compiler {
    parser: Parser,
    chunk: Chunk,
}

impl Compiler {
    pub fn new(parser: Parser, chunk: Chunk) -> Self {
        Self { parser, chunk }
    }

    pub fn emit_bytes(&mut self, byte1: u8, byte2: u8) {
        self.emit_byte(byte1);
        self.emit_byte(byte2);
    }

    pub fn emit_byte(&mut self, byte: u8) {
        let line = self.parser.previous.line as usize;
        self.current_chunk().write(byte, line);
    }

    pub fn current_chunk(&mut self) -> &mut Chunk {
        &mut self.chunk
    }

    pub fn end_compiler(&mut self) {
        self.emit_return();
        #[cfg(debug)]
        {
            if !self.parser.had_error {
                disassemble_chunk(&self.current_chunk(), "code");
            }
        }
    }

    pub fn grouping(&mut self) {
        self.expression();
        self.parser
            .consume(TokenType::RightParen, "Expect ')' after expression.");
    }

    pub fn expression(&mut self) {
        self.parse_precedence(Precedence::Assignment);
    }

    pub fn number(&mut self) {
        let constant = Value::Double(self.parser.previous.as_str().parse::<f64>().unwrap());
        self.emit_constant(constant);
    }

    pub fn unary(&mut self) {
        let operator_type = self.parser.previous.kind;

        self.parse_precedence(Precedence::Unary);
        match operator_type {
            TokenType::Minus => self.emit_byte(OpCode::OpNegate as u8),
            _ => unreachable!("Operator type can't be anything but Minus for unary expressions."),
        }
    }

    pub fn binary(&mut self) {
        let operator_type = self.parser.previous.kind;
        let rule = Compiler::get_rule(operator_type);

        self.parse_precedence(Precedence::from_unchecked(rule.precedence as u8 + 1));
        match operator_type {
        TokenType::Plus => self.emit_byte(OpCode::OpAdd as u8),
        TokenType::Minus => self.emit_byte(OpCode::OpSubtract as u8),
        TokenType::Star => self.emit_byte(OpCode::OpMultiply as u8),
        TokenType::Slash => self.emit_byte(OpCode::OpDivide as u8),
        _ => unreachable!("Operator type can't be anything but Plus, Minus, Star, or Slash for binary expressions."),
      }
    }

    pub fn parse_precedence(&mut self, precedence: Precedence) {
        self.parser.advance();

        let prefix_rule = Compiler::get_rule(self.parser.previous.kind).prefix;
        match prefix_rule {
            Some(prefix_rule) => {
                prefix_rule(self);

                while precedence <= Compiler::get_rule(self.parser.current.kind).precedence {
                    self.parser.advance();

                    if let Some(infix_rule) = Compiler::get_rule(self.parser.previous.kind).infix {
                        infix_rule(self);
                    }
                }
            }
            None => self.parser.error("Expect expression."),
        }
    }

    pub fn get_rule(token_type: TokenType) -> &'static ParserRule {
        &RULES[token_type as usize]
    }

    pub fn emit_constant(&mut self, constant: Value) {
        let constant_index = self.make_constant(constant);
        self.emit_bytes(OpCode::OpConstant as u8, constant_index);
    }

    pub fn make_constant(&mut self, constant: Value) -> u8 {
        let constant_index = self.current_chunk().add_constant(constant);
        if constant_index > u8::MAX as usize {
            self.parser
                .error_at_current("Too many constants in one chunk.");
            return 0;
        }
        constant_index as u8
    }

    pub fn emit_return(&mut self) {
        self.emit_byte(OpCode::OpReturn as u8);
    }
}

pub fn compile(source: &[u8]) -> Result<Chunk, InterpretError> {
    let dummy_token = Token::new(TokenType::Eof, 0, "");
    let parser = Parser::new(Scanner::new(source), dummy_token);

    let mut compiler = Compiler::new(parser, Chunk::new());

    compiler.parser.advance();
    compiler.expression();
    compiler
        .parser
        .consume(TokenType::Eof, "Expect end of expression.");
    compiler.end_compiler();

    if compiler.parser.had_error {
        return Result::Err(InterpretError::CompileError);
    }
    Result::Ok(compiler.chunk)
}

#[rustfmt::skip]
const RULES: [ParserRule; 40] = [
    ParserRule::new(Some(Compiler::grouping), None, Precedence::None), // TokenType::LeftParen
    ParserRule::new(None, None, Precedence::None),                     // TokenType::RightParen
    ParserRule::new(None, None, Precedence::None),                     // TokenType::LeftBrace
    ParserRule::new(None, None, Precedence::None),                     // TokenType::RightBrace
    ParserRule::new(None, None, Precedence::None),                     // TokenType::Comma
    ParserRule::new(None, None, Precedence::None),                     // TokenType::Dot
    ParserRule::new(Some(Compiler::unary),Some(Compiler::binary),Precedence::Term,), // TokenType::Minus
    ParserRule::new(None, Some(Compiler::binary), Precedence::Term),   // TokenType::Plus
    ParserRule::new(None, None, Precedence::None),                     // TokenType::Semicolon
    ParserRule::new(None, Some(Compiler::binary), Precedence::Factor), // TokenType::Slash
    ParserRule::new(None, Some(Compiler::binary), Precedence::Factor), // TokenType::Star
    ParserRule::new(None, None, Precedence::None),                     // TokenType::Bang
    ParserRule::new(None, None, Precedence::None),                     // TokenType::BangEqual
    ParserRule::new(None, None, Precedence::None),                     // TokenType::Equal
    ParserRule::new(None, None, Precedence::None),                     // TokenType::EqualEqual
    ParserRule::new(None, None, Precedence::None),                     // TokenType::Greater
    ParserRule::new(None, None, Precedence::None),                     // TokenType::GreaterEqual
    ParserRule::new(None, None, Precedence::None),                     // TokenType::Less
    ParserRule::new(None, None, Precedence::None),                     // TokenType::LessEqual
    ParserRule::new(None, None, Precedence::None),                     // TokenType::Identifier
    ParserRule::new(None, None, Precedence::None),                     // TokenType::String
    ParserRule::new(Some(Compiler::number), None, Precedence::None),   // TokenType::Number
    ParserRule::new(None, None, Precedence::None),                     // TokenType::And
    ParserRule::new(None, None, Precedence::None),                     // TokenType::Class
    ParserRule::new(None, None, Precedence::None),                     // TokenType::Else
    ParserRule::new(None, None, Precedence::None),                     // TokenType::False
    ParserRule::new(None, None, Precedence::None),                     // TokenType::For
    ParserRule::new(None, None, Precedence::None),                     // TokenType::Fun
    ParserRule::new(None, None, Precedence::None),                     // TokenType::If
    ParserRule::new(None, None, Precedence::None),                     // TokenType::Nil
    ParserRule::new(None, None, Precedence::None),                     // TokenType::Or
    ParserRule::new(None, None, Precedence::None),                     // TokenType::Print
    ParserRule::new(None, None, Precedence::None),                     // TokenType::Return
    ParserRule::new(None, None, Precedence::None),                     // TokenType::Super
    ParserRule::new(None, None, Precedence::None),                     // TokenType::This
    ParserRule::new(None, None, Precedence::None),                     // TokenType::True
    ParserRule::new(None, None, Precedence::None),                     // TokenType::Var
    ParserRule::new(None, None, Precedence::None),                     // TokenType::While
    ParserRule::new(None, None, Precedence::None),                     // TokenType::Error
    ParserRule::new(None, None, Precedence::None),                     // TokenType::Eof
];
