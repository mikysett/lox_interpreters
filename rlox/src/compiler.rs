use crate::scanner::Scanner;
use crate::scanner::TokenType;

pub fn compile(source: &[u8]) {
    let mut scanner = Scanner::new(source);

    let mut line = -1;
    loop {
        let token = scanner.scan_token();
        if token.line != line {
            line = token.line;
            print!("{:4} ", line);
        } else {
            print!("   | ");
        }

        println!(" {:2} '{}'", token.kind as u8, token.as_str());
        if token.kind == TokenType::Eof {
            break;
        }
    }
}
