type Line = isize;

#[derive(Debug, Clone, Copy, PartialEq)]
#[repr(u8)]
#[rustfmt::skip]
pub enum TokenType {
    // Single-character tokens.
    LeftParen = 0, RightParen, LeftBrace, RightBrace, Comma, Dot, Minus, Plus, Semicolon, Slash, Star,
    // One or two character tokens.
    Bang, BangEqual, Equal, EqualEqual, Greater, GreaterEqual, Less, LessEqual,
    // Literals.
    Identifier, String, Number,
    // Keywords.
    And, Class, Else, False, For, Fun, If, Nil, Or, Print, Return, Super, This, True, Var, While,

    Error, Eof,
}

pub struct Scanner {
    start: *const u8,
    // `end` is not a valid pointer, used only to check when the string ends (start == end).
    end: *const u8,
    current: *const u8,
    line: Line,
}

#[derive(Debug, Clone, Copy)]
pub struct Token {
    pub kind: TokenType,
    pub line: Line,
    pub start: *const u8,
    pub length: usize,
}

impl Scanner {
    pub fn new(source: &[u8]) -> Self {
        let start = source.as_ptr();
        let current = start;
        // This is the first position after the end of the string, it should never be dereferenced, used only to check for the end of the string.
        let end = unsafe { start.add(source.len()) };
        Self {
            start,
            current,
            end,
            line: 1,
        }
    }

    pub fn scan_token(&mut self) -> Token {
        self.skip_whitespace();
        self.start = self.current;

        if self.is_at_end() {
            return Token::new(TokenType::Eof, self.line, "");
        }

        let c = self.advance();
        return match c as char {
            '(' => Token::new(TokenType::LeftParen, self.line, self.get_sub_str()),
            ')' => Token::new(TokenType::RightParen, self.line, self.get_sub_str()),
            '{' => Token::new(TokenType::LeftBrace, self.line, self.get_sub_str()),
            '}' => Token::new(TokenType::RightBrace, self.line, self.get_sub_str()),
            ';' => Token::new(TokenType::Semicolon, self.line, self.get_sub_str()),
            ',' => Token::new(TokenType::Comma, self.line, self.get_sub_str()),
            '.' => Token::new(TokenType::Dot, self.line, self.get_sub_str()),
            '-' => Token::new(TokenType::Minus, self.line, self.get_sub_str()),
            '+' => Token::new(TokenType::Plus, self.line, self.get_sub_str()),
            '/' => Token::new(TokenType::Slash, self.line, self.get_sub_str()),
            '*' => Token::new(TokenType::Star, self.line, self.get_sub_str()),
            '!' if self.match_char('=') => {
                Token::new(TokenType::BangEqual, self.line, self.get_sub_str())
            }
            '=' if self.match_char('=') => {
                Token::new(TokenType::EqualEqual, self.line, self.get_sub_str())
            }
            '<' if self.match_char('=') => {
                Token::new(TokenType::LessEqual, self.line, self.get_sub_str())
            }
            '>' if self.match_char('=') => {
                Token::new(TokenType::GreaterEqual, self.line, self.get_sub_str())
            }
            '!' => Token::new(TokenType::Bang, self.line, self.get_sub_str()),
            '=' => Token::new(TokenType::Equal, self.line, self.get_sub_str()),
            '<' => Token::new(TokenType::Less, self.line, self.get_sub_str()),
            '>' => Token::new(TokenType::Greater, self.line, self.get_sub_str()),
            '"' => self.string_token(),
            c if c.is_ascii_digit() => self.number_token(),
            c if is_alpha(c) => self.identifier_token(),
            _ => Token::error_token("Unexpected character.", self.line),
        };
    }

    // [advance] should only be used after checking the position is valid with [is_at_end].
    fn advance(&mut self) -> u8 {
        let c = unsafe { self.current.read() };
        unsafe {
            self.current = self.current.add(1);
        }
        c
    }

    // If [c] matches the current token return true and advance the cursor.
    fn match_char(&mut self, c: char) -> bool {
        // this ensures self.current is a valid pointer.
        if self.is_at_end() {
            false
        } else if unsafe { *self.current } == c as u8 {
            unsafe {
                self.current = self.current.add(1);
            }
            true
        } else {
            false
        }
    }

    fn peek(&self) -> char {
        unsafe { *self.current as char }
    }

    fn peek_next(&self) -> char {
        if self.is_at_end() {
            '\0'
        } else {
            unsafe { *self.current.add(1) as char }
        }
    }

    fn get_sub_str(&self) -> &str {
        unsafe {
            std::str::from_utf8_unchecked(std::slice::from_raw_parts(
                self.start,
                self.current as usize - self.start as usize,
            ))
        }
    }

    fn is_at_end(&self) -> bool {
        self.current == self.end as *mut u8
    }

    fn skip_whitespace(&mut self) {
        loop {
            match self.peek() {
                ' ' | '\r' | '\t' => {
                    self.advance();
                }
                '\n' => {
                    self.line += 1;
                    self.advance();
                }
                '/' if self.peek_next() == '/' => {
                    self.advance();
                    self.advance();
                    while !self.is_at_end() && self.peek() != '\n' {
                        self.advance();
                    }
                }
                _ => return,
            }
        }
    }

    fn string_token(&mut self) -> Token {
        while !self.is_at_end() && self.peek() != '"' {
            if self.peek() == '\n' {
                self.line += 1;
            }
            self.advance();
        }

        if self.is_at_end() {
            return Token::error_token("Unterminated string.", self.line);
        }

        // Match closing quote
        self.advance();

        Token::new(TokenType::String, self.line, self.get_sub_str())
    }

    fn number_token(&mut self) -> Token {
        while !self.is_at_end() && self.peek().is_ascii_digit() {
            self.advance();
        }

        if self.peek() == '.' && self.peek_next().is_ascii_digit() {
            self.advance();
            while !self.is_at_end() && self.peek().is_ascii_digit() {
                self.advance();
            }
        }

        Token::new(TokenType::Number, self.line, self.get_sub_str())
    }

    fn identifier_token(&mut self) -> Token {
        // TODO: Maybe we need is_at_end() check here
        while is_alpha(self.peek()) || self.peek().is_ascii_digit() {
            self.advance();
        }

        Token::new(self.identifier_type(), self.line, self.get_sub_str())
    }

    fn identifier_type(&mut self) -> TokenType {
        let identifier = self.get_sub_str();

        match identifier.as_bytes()[0] as char {
            'a' => check_keyword(identifier, 1, 2, "nd", TokenType::And),
            'c' => check_keyword(identifier, 1, 4, "lass", TokenType::Class),
            'e' => check_keyword(identifier, 1, 3, "lse", TokenType::Else),
            'i' => check_keyword(identifier, 1, 1, "f", TokenType::If),
            'n' => check_keyword(identifier, 1, 2, "il", TokenType::Nil),
            'o' => check_keyword(identifier, 1, 1, "r", TokenType::Or),
            'p' => check_keyword(identifier, 1, 4, "rint", TokenType::Print),
            'r' => check_keyword(identifier, 1, 5, "eturn", TokenType::Return),
            's' => check_keyword(identifier, 1, 4, "uper", TokenType::Super),
            'v' => check_keyword(identifier, 1, 2, "ar", TokenType::Var),
            'w' => check_keyword(identifier, 1, 4, "hile", TokenType::While),
            'f' if identifier.as_bytes().len() > 1 => match identifier.as_bytes()[1] as char {
                'a' => check_keyword(identifier, 2, 3, "lse", TokenType::False),
                'o' => check_keyword(identifier, 2, 1, "r", TokenType::For),
                'u' => check_keyword(identifier, 2, 1, "n", TokenType::Fun),
                _ => TokenType::Identifier,
            },
            't' if identifier.as_bytes().len() > 1 => match identifier.as_bytes()[1] as char {
                'h' => check_keyword(identifier, 2, 2, "is", TokenType::This),
                'r' => check_keyword(identifier, 2, 2, "ue", TokenType::True),
                _ => TokenType::Identifier,
            },
            _ => TokenType::Identifier,
        }
    }
}

impl Token {
    pub fn new(kind: TokenType, line: Line, value: &str) -> Self {
        Self {
            kind,
            line,
            start: value.as_ptr(),
            length: value.len(),
        }
    }

    pub fn error_token(message: &str, line: Line) -> Self {
        Self::new(TokenType::Error, line, message)
    }

    pub fn as_str(&self) -> &str {
        unsafe {
            std::str::from_utf8_unchecked(std::slice::from_raw_parts(self.start, self.length))
        }
    }
}

fn is_alpha(c: char) -> bool {
    c.is_ascii_alphabetic() || c == '_'
}

fn check_keyword(
    identifier: &str,
    start: usize,
    length: usize,
    rest: &str,
    kind: TokenType,
) -> TokenType {
    if identifier.len() == start + length && &identifier[start..] == rest {
        kind
    } else {
        TokenType::Identifier
    }
}
