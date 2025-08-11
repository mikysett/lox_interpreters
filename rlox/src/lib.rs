pub mod debug;
pub mod domain;
pub mod run;
pub mod vm;

#[cfg(debug)]
pub use debug::debug::disassemble_instruction;
pub use run::run;
pub use vm::VM;
