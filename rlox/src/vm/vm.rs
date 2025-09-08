use crate::binary_op;
use crate::compile;
use crate::domain::{Chunk, OpCode, Value};
use crate::runtime_error;

const STACK_MAX: usize = 256;
const VALUE_UNKNOWN: Value = Value::Unknown;

#[cfg(speedtest)]
const SPEED_TEST_RUNS: i32 = 10_000_000;

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

            let result;
            #[cfg(speedtest)]
            {
                let mut last_result = Result::Ok(());
                for _ in 0..SPEED_TEST_RUNS {
                    self.ip = self.chunk.code.as_ptr();
                    last_result = self.run();
                }
                result = last_result;
            }
            #[cfg(not(speedtest))]
            {
                self.ip = self.chunk.code.as_ptr();
                result = self.run();
            }

            self.chunk.free();
            result
        })
    }

    pub fn run(&mut self) -> Result<(), InterpretError> {
        loop {
            match self.read_byte().into() {
                OpCode::Constant => {
                    let constant = self.read_constant().clone();
                    self.push(constant);
                }
                OpCode::True => self.push(Value::Bool(true)),
                OpCode::False => self.push(Value::Bool(false)),
                OpCode::Equal => {
                    let right = self.pop();
                    let left = self.pop();
                    self.push(Value::Bool(left.equal(&right)));
                }
                #[cfg(optimize)]
                OpCode::EqualZero => {
                    let value = self.pop();
                    self.push(Value::Bool(value.is_zero()));
                }
                OpCode::Greater => binary_op!(self, Value::Bool, >),
                OpCode::Less => binary_op!(self, Value::Bool, <),
                OpCode::Nil => self.push(Value::Nil),
                OpCode::ConstantLong => {
                    let constant = self.read_constant_long().clone();
                    self.push(constant);
                }
                OpCode::Add => binary_op!(self, Value::Double, +),
                OpCode::Subtract => binary_op!(self, Value::Double, -),
                OpCode::Multiply => binary_op!(self, Value::Double, *),
                OpCode::Divide => binary_op!(self, Value::Double, /),
                OpCode::Negate => {
                    let last = self.peek(0);
                    if let Value::Double(value) = *last {
                        *last = Value::Double(-value);
                    } else {
                        runtime_error!(self, "Operand must be a number.");
                        return Err(InterpretError::RuntimeError);
                    }
                }
                OpCode::Not => {
                    let last = self.peek(0);
                    *last = Value::Bool(last.is_falsey());
                }
                OpCode::Return => {
                    // Writing to stdout makes speedtests slower with no benefit
                    #[cfg(not(speedtest))]
                    println!("{}", self.pop());

                    #[cfg(speedtest)]
                    self.pop();
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

    pub fn peek(&mut self, distance: usize) -> &mut Value {
        &mut self.stack[self.stack_top - 1 - distance]
    }
}

#[macro_export]
macro_rules! binary_op {
    ($vm:ident, $value_type:path, $op:tt) => {{
        let right = $vm.pop();
        let left = $vm.pop();
        if let (Value::Double(right), Value::Double(left)) = (right, left) {
            $vm.push($value_type(left $op right));
        } else {
            runtime_error!($vm, "Operands must be numbers.");
            return Err(InterpretError::RuntimeError);
        }
    }};
}

#[macro_export]
macro_rules! runtime_error {
    ($vm:ident, $format:expr $(, $($arg:expr),*)?) => {{
        eprintln!($format, $($arg),*);

        let line = $vm.chunk.get_line($vm.ip as usize - 1 - $vm.chunk.code.as_ptr() as usize);
        println!("[line {}] in script", line);
        $vm.reset_stack();
    }};
}
