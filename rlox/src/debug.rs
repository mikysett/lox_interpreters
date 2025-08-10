use crate::domain::Chunk;
use crate::domain::OpCode;

pub fn disassemble_chunk(chunk: &Chunk, name: &str) {
    println!("== {name} ==");

    let mut offset = 0;
    while offset < chunk.code.len() {
        offset = disassemble_instruction(chunk, offset);
    }
}

fn disassemble_instruction(chunk: &Chunk, offset: usize) -> usize {
    print!("{:04} ", offset);

    let current_line = chunk.get_line(offset);
    if offset > 0 && chunk.get_line(offset - 1) == current_line {
        print!("   | ");
    } else {
        print!("{:4} ", current_line);
    }

    match &chunk.code[offset] {
        OpCode::OpReturn => simple_instruction("OP_RETURN", offset),
        OpCode::OpConstant => constant_instruction("OP_CONSTANT", offset, chunk),
        OpCode::OpConstantLong => constant_long_instruction("OP_CONSTANT_LONG", offset, chunk),
        instruction => {
            println!("Unknown opcode {instruction:?}");
            offset + 1
        }
    }
}

fn simple_instruction(name: &str, offset: usize) -> usize {
    println!("{name}");
    offset + 1
}

fn constant_instruction(name: &str, offset: usize, chunk: &Chunk) -> usize {
    let OpCode::Byte(pointer) = chunk.code[offset + 1] else {
        panic!("error: invalid pointer at offset {offset} for constant instruction.");
    };
    print!("{name:-16} {:4} '", pointer);
    chunk.constants[pointer as usize].print();
    println!("'");
    offset + 2
}

fn constant_long_instruction(name: &str, offset: usize, chunk: &Chunk) -> usize {
    let OpCode::FourBytes(pointer) = chunk.code[offset + 1] else {
        panic!("error: invalid pointer at offset {offset} for constant long instruction.");
    };
    print!("{name:-16} {:4} '", pointer);
    chunk.constants[pointer as usize].print();
    println!("'");
    offset + 2
}
