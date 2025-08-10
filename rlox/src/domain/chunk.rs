use std::u8;

use crate::domain::value::Value;

#[derive(Debug)]
pub enum OpCode {
    OpReturn,
    OpConstant,
    Byte(u8),
    OpConstantLong,
    FourBytes(u32),
}

pub struct Line {
    code_offset: usize,
    line_number: usize,
}

impl Line {
    pub fn new(code_offset: usize, line_number: usize) -> Self {
        Self {
            code_offset,
            line_number,
        }
    }
}

pub struct Chunk {
    pub code: Vec<OpCode>,
    pub constants: Vec<Value>,
    pub lines: Vec<Line>,
}

impl Chunk {
    pub fn new() -> Self {
        Chunk {
            code: Vec::new(),
            constants: Vec::new(),
            lines: Vec::new(),
        }
    }

    pub fn write(&mut self, byte: OpCode, line: usize) {
        self.code.push(byte);

        if self.lines.is_empty() {
            self.lines.push(Line::new(self.code.len() - 1, line));
        } else if line != self.lines.last().unwrap().line_number {
            self.lines.push(Line::new(self.code.len() - 1, line));
        }
    }

    pub fn free(&mut self) {
        self.code.clear();
        self.code.shrink_to_fit();
        self.constants.clear();
        self.constants.shrink_to_fit();
        self.lines.clear();
        self.lines.shrink_to_fit();
    }

    pub fn write_constant(&mut self, value: Value, line: usize) {
        let index = self.add_constant(value);

        if index >= u8::MAX as usize {
            self.write(OpCode::OpConstantLong, line);
            self.write(OpCode::FourBytes(index as u32), line);
        } else {
            self.write(OpCode::OpConstant, line);
            self.write(OpCode::Byte(index as u8), line);
        }
    }

    fn add_constant(&mut self, value: Value) -> usize {
        self.constants.push(value);
        self.constants.len() - 1
    }

    pub fn get_line(&self, code_offset: usize) -> usize {
        if self.lines.is_empty() {
            return 0;
        }

        for i in 0..self.lines.len() - 1 {
            if self.lines[i + 1].code_offset > code_offset {
                return self.lines[i].line_number;
            }
        }
        self.lines.last().map(|line| line.line_number).unwrap_or(0)
    }
}
