use crate::domain::Chunk;
use crate::domain::OpCode;
use crate::domain::Value;
use crate::vm::VM;

pub fn run() {
    let mut vm = VM::new();

    let mut chunk = Chunk::new();

    chunk.write_constant(Value::Double(1.2), 123);

    chunk.write_constant(Value::Double(3.4), 123);
    chunk.write(OpCode::OpAdd as u8, 123);

    chunk.write_constant(Value::Double(5.6), 123);
    chunk.write(OpCode::OpDivide as u8, 123);

    chunk.write(OpCode::OpNegate as u8, 123);
    chunk.write(OpCode::OpReturn as u8, 123);

    vm.interpret(chunk);
    vm.free();
}
