use crate::disassemble_chunk;
use crate::domain::Chunk;
use crate::domain::OpCode;
use crate::domain::Value;

pub fn run() {
    let mut chunk = Chunk::new();

    chunk.write_constant(Value::Double(1.2), 123);
    chunk.write(OpCode::OpReturn, 123);

    disassemble_chunk(&chunk, "test chunk");
    chunk.free();
}
