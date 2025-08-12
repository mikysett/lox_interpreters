use crate::binary_op;
#[cfg(debug)]
use crate::debug::debug::disassemble_instruction;
use crate::domain::{Chunk, OpCode, Value};

const STACK_MAX: usize = 256;
const VALUE_UNKNOWN: Value = Value::Unknown;

pub enum InterpretResult {
    Ok,
    CompileError,
    RuntimeError,
}

pub struct VM {
    chunk: Chunk,
    ip: usize,
    stack: [Value; STACK_MAX],
    stack_top: usize,
}

impl VM {
    pub fn new() -> Self {
        Self {
            chunk: Chunk::new(),
            ip: 0,
            stack: [VALUE_UNKNOWN; STACK_MAX],
            stack_top: 0,
        }
    }

    pub fn free(&mut self) {
        // TODO: may be useless in rust
        self.chunk.free(); // This is implemented in run() after vm.free() by the book
    }

    pub fn interpret(&mut self, chunk: Chunk) -> InterpretResult {
        self.chunk = chunk;
        loop {
            #[cfg(debug)]
            {
                print!("          ");
                for i in 0..self.stack_top {
                    print!("[ {} ]", self.stack[i]);
                }
                println!();
                disassemble_instruction(&self.chunk, self.ip);
            }
            match self.read_byte().into() {
                OpCode::OpConstant => {
                    let constant = self.read_constant().clone();
                    println!("constant: {}", constant);
                    self.push(constant);
                }
                OpCode::OpConstantLong => {
                    let constant = self.read_constant_long().clone();
                    println!("constant: {}", constant);
                    self.push(constant);
                }
                OpCode::OpAdd => binary_op!(self, +),
                OpCode::OpSubtract => binary_op!(self, -),
                OpCode::OpMultiply => binary_op!(self, *),
                OpCode::OpDivide => binary_op!(self, /),
                OpCode::OpNegate => {
                    let last = self.stack_top - 1;
                    if let Value::Double(value) = self.stack[last] {
                        self.stack[last] = Value::Double(-value);
                    } else {
                        return InterpretResult::RuntimeError;
                    }
                }
                OpCode::OpReturn => {
                    println!("{}", self.pop());
                    return InterpretResult::Ok;
                }
                OpCode::Unknown => return InterpretResult::RuntimeError,
            };
        }
    }

    pub fn read_byte(&mut self) -> u8 {
        let result = self.chunk.code[self.ip];
        self.ip += 1;
        result
    }

    pub fn read_constant(&mut self) -> &Value {
        let pointer = self.chunk.code[self.ip];
        let constant = &self.chunk.constants[pointer as usize];
        self.ip += 1;
        constant
    }

    pub fn read_constant_long(&mut self) -> &Value {
        let p1 = self.chunk.code[self.ip];
        let p2 = self.chunk.code[self.ip + 1];
        let p3 = self.chunk.code[self.ip + 2];
        let constant =
            &self.chunk.constants[p1 as usize | (p2 as usize) << 8 | (p3 as usize) << 16];
        self.ip += 3;
        constant
    }

    pub fn reset_stack(&mut self) {
        self.stack_top = 0;
    }

    pub fn push(&mut self, value: Value) {
        self.stack[self.stack_top] = value;
        self.stack_top += 1;
    }

    pub fn pop(&mut self) -> Value {
        self.stack_top -= 1;
        self.stack[self.stack_top].clone()
    }
}

#[macro_export]
macro_rules! binary_op {
    ($vm:ident, $op:tt) => {{
        let right = $vm.pop();
        let left = $vm.pop();
        if let (Value::Double(right), Value::Double(left)) = (right, left) {
            $vm.push(Value::Double(left $op right));
        } else {
            return InterpretResult::RuntimeError;
        }
    }};
}
