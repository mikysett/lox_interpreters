use crate::domain::value::Value;

#[derive(Debug)]
#[repr(u8)]
pub enum OpCode {
    Return = 0,
    Constant = 1,
    True = 2,
    False = 3,
    Nil = 4,
    ConstantLong = 5,
    Add = 6,
    Subtract = 7,
    Multiply = 8,
    Divide = 9,
    Negate = 10,
    Not = 11,
    Equal = 12,
    #[cfg(feature = "optimize")]
    EqualZero = 13,
    Greater = 14,
    Less = 15,
    Unknown = 0xff,
}

// #[cfg(debug)]
impl From<u8> for OpCode {
    fn from(value: u8) -> Self {
        match value {
            0 => OpCode::Return,
            1 => OpCode::Constant,
            2 => OpCode::True,
            3 => OpCode::False,
            4 => OpCode::Nil,
            5 => OpCode::ConstantLong,
            6 => OpCode::Add,
            7 => OpCode::Subtract,
            8 => OpCode::Multiply,
            9 => OpCode::Divide,
            10 => OpCode::Negate,
            11 => OpCode::Not,
            12 => OpCode::Equal,
            #[cfg(feature = "optimize")]
            13 => OpCode::EqualZero,
            14 => OpCode::Greater,
            15 => OpCode::Less,
            _ => OpCode::Unknown,
        }
    }
}

impl Into<u8> for OpCode {
    fn into(self) -> u8 {
        self as u8
    }
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
    pub code: Vec<u8>,
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

    pub fn write(&mut self, byte: u8, line: usize) {
        self.code.push(byte);

        if self.lines.is_empty() || line != self.lines.last().unwrap().line_number {
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
            self.write(OpCode::ConstantLong as u8, line);
            self.write((index & 0xff) as u8, line);
            self.write(((index >> 8) & 0xff) as u8, line);
            self.write(((index >> 16) & 0xff) as u8, line);
        } else {
            self.write(OpCode::Constant as u8, line);
            self.write(index as u8, line);
        }
    }

    pub fn add_constant(&mut self, value: Value) -> usize {
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
        self.lines.last().map(|line| line.line_number).unwrap()
    }

    #[cfg(feature = "optimize")]
    pub fn optimize(&mut self) {
        let permutations = vec![Permutation::new(3, |chunk, off| {
            if chunk.code.len() > off + 3
                && chunk.code[off] == OpCode::Constant as u8
                && chunk.code[off + 2] == OpCode::Equal as u8
            {
                let pointer = chunk.code[off + 1] as usize;
                if chunk.constants[pointer].is_zero() {
                    return Some(vec![OpCode::EqualZero as u8]);
                }
            }
            None
        })];

        loop {
            let mut new_chunk = Chunk::new();
            let mut last_written = 0;
            self.code.iter().enumerate().for_each(|(offset, _)| {
                if last_written > offset {
                    return;
                }

                for permutation in &permutations {
                    if let Some(new_opcodes) = (permutation.try_apply)(self, offset) {
                        while last_written != offset {
                            new_chunk.write(self.code[last_written], self.get_line(last_written));
                            last_written += 1;
                        }

                        let permutated_line = self.get_line(last_written);
                        for new_opcode in new_opcodes {
                            new_chunk.write(new_opcode, permutated_line);
                        }
                        last_written += permutation.len;
                    }
                }
            });

            if last_written != 0 {
                while last_written < self.code.len() {
                    new_chunk.write(self.code[last_written], self.get_line(last_written));
                    last_written += 1;
                }

                self.lines = new_chunk.lines;
                self.code = new_chunk.code;
            } else {
                break;
            }
        }
    }
}

#[cfg(feature = "optimize")]
#[derive(Debug)]
struct Permutation {
    len: usize,
    try_apply: fn(&Chunk, usize) -> Option<Vec<u8>>,
}

#[cfg(feature = "optimize")]
impl Permutation {
    pub fn new(len: usize, try_apply: fn(&Chunk, usize) -> Option<Vec<u8>>) -> Self {
        Permutation { len, try_apply }
    }
}
