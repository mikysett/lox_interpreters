use crate::binary_op;
use crate::compile;
use crate::domain::{Chunk, OpCode, Value};

const STACK_MAX: usize = 256;
const VALUE_UNKNOWN: Value = Value::Unknown;

#[derive(Debug)]
pub enum InterpretError {
    CompileError,
    RuntimeError,
}

pub struct VM {
    chunk: Chunk,
    ip: *const u8,
    stack: [Value; STACK_MAX],
    stack_top: usize,
}

impl VM {
    pub fn new() -> Self {
        Self {
            chunk: Chunk::new(),
            ip: std::ptr::null(),
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
    pub fn interpret<'a>(&mut self, source: &'a [u8]) -> Result<(), InterpretError> {
        compile(source).and_then(|chunk| {
            self.chunk = chunk;
            self.ip = self.chunk.code.as_ptr();
            let result = self.run();
            self.chunk.free();
            result
        })
    }

    pub fn run(&mut self) -> Result<(), InterpretError> {
        loop {
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
                        return Err(InterpretError::RuntimeError);
                    }
                }
                OpCode::OpReturn => {
                    println!("{}", self.pop());
                    return Ok(());
                }
                OpCode::Unknown => return Err(InterpretError::RuntimeError),
            }
        }
    }

    pub fn read_byte(&mut self) -> u8 {
        unsafe {
            let result = *self.ip;
            self.ip = self.ip.add(1);
            result
        }
    }

    pub fn read_constant(&mut self) -> &Value {
        unsafe {
            let pointer = *self.ip;
            let constant = &self.chunk.constants[pointer as usize];
            self.ip = self.ip.add(1);
            constant
        }
    }

    pub fn read_constant_long(&mut self) -> &Value {
        unsafe {
            let p1 = *self.ip;
            let p2 = *self.ip.add(1);
            let p3 = *self.ip.add(2);
            let constant =
                &self.chunk.constants[p1 as usize | (p2 as usize) << 8 | (p3 as usize) << 16];
            self.ip = self.ip.add(3);
            constant
        }
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
            return Err(InterpretError::RuntimeError);
        }
    }};
}
