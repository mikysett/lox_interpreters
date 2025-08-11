#[cfg(debug)]
pub mod debug {
    use crate::domain::Chunk;
    use crate::domain::OpCode;

    pub fn disassemble_instruction(chunk: &Chunk, offset: usize) -> usize {
        print!("{:04} ", offset);

        let current_line = chunk.get_line(offset);
        if offset > 0 && chunk.get_line(offset - 1) == current_line {
            print!("   | ");
        } else {
            print!("{:4} ", current_line);
        }

        match OpCode::from(chunk.code[offset]) {
            OpCode::OpReturn => simple_instruction("OP_RETURN", offset),
            OpCode::OpConstant => constant_instruction("OP_CONSTANT", offset, chunk),
            OpCode::OpConstantLong => constant_long_instruction("OP_CONSTANT_LONG", offset, chunk),
            OpCode::OpNegate => simple_instruction("OP_NEGATE", offset),
            OpCode::Unknown => {
                println!("Unknown opcode {} at offset {}", chunk.code[offset], offset);
                offset + 1
            }
        }
    }

    fn simple_instruction(name: &str, offset: usize) -> usize {
        println!("{name}");
        offset + 1
    }

    fn constant_instruction(name: &str, offset: usize, chunk: &Chunk) -> usize {
        let pointer = chunk.code[offset + 1] as usize;
        println!("{name:-16} {:4} '{}'", pointer, chunk.constants[pointer]);
        offset + 2
    }

    fn constant_long_instruction(name: &str, offset: usize, chunk: &Chunk) -> usize {
        let pointer = chunk.code[offset + 1] as usize
            | (chunk.code[offset + 2] as usize) << 8
            | (chunk.code[offset + 3] as usize) << 16;
        println!("{name:-16} {:4} '{}'", pointer, chunk.constants[pointer]);
        offset + 4
    }
}
