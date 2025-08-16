use std::env;
use std::io::stdin;
use std::io::stdout;
use std::io::Write;

use crate::vm::InterpretResult;
use crate::vm::VM;

pub fn run() {
    let args: Vec<String> = env::args().collect();

    let mut vm = VM::new();
    match args.len() {
        1 => repl(&mut vm),
        2 => run_file(&mut vm, &args[1]),
        _ => {
            println!("Usage: rlox [path]");
            std::process::exit(64);
        }
    }

    vm.free();
}

fn repl(vm: &mut VM) {
    let mut buf = String::new();

    loop {
        print!("> ");
        stdout().flush().unwrap_or_else(|err| {
            eprintln!("error: Failed to flush stdout: {}", err);
            std::process::exit(65);
        });
        stdin().read_line(&mut buf).unwrap_or_else(|err| {
            eprintln!("error: Failed to read line: {}", err);
            std::process::exit(65);
        });
        vm.interpret(buf.as_bytes());
    }
}

fn run_file(vm: &mut VM, path: &str) {
    let source = std::fs::read(path).unwrap_or_else(|err| {
        eprintln!("error: Failed to read file: {}", err);
        std::process::exit(66);
    });

    match vm.interpret(&source) {
        InterpretResult::Ok => {}
        InterpretResult::CompileError => std::process::exit(65),
        InterpretResult::RuntimeError => std::process::exit(70),
    }
}
