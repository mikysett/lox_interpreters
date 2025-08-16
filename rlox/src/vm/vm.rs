use crate::compile;
use crate::domain::{Chunk, Value};

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

    // The lifetime is needed to ensure the [source] won't be freed before the VM consumes it.
    // By storing [source] in an unsafe ptr the compiler will not understand it needs to stay live without the lifetime.
    pub fn interpret<'a>(&mut self, source: &'a [u8]) -> InterpretResult {
        compile(source);
        InterpretResult::Ok
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
