# Crafting Interpreters: In Go and Rust

Those are my interpreters for `lox`, the language used as an example in the fantastic book [Crafting Interpreters](https://craftinginterpreters.com/).

## glox: The Go interpreter [DONE]

> In the book this corresponds to `jlox`, a Java interpreter

In `glox` directory you can run:
- `make` to build the binary in `bin/glox`
- `make install` to install the interpreter globally
- Use `disable-extras=true` flag to have a canonical implementation without optional improvements

## rlox: The Rust interpreter [TODO]

> In the book this corresponds to `clox`, a C compiler to bytecode with a VM

In `rlox` directory you can run:
- `make` to build the binary in `bin/rlox`
- To build an optimized version `OPT=true make`
- To run `hyperfine` benchmark use `make speedtest`

## Example of Lox file

```lox
// fibonacci.lox
print "This program will output fibonacci numbers up to 10!";

fun fibonacci(first, second, count) {
    if (count <= 0) {
        return;
    }
    print first + second;
    fibonacci(second, first + second, count - 1);
}

fibonacci(0, 1, 10);
```

## To test those implementations with the offical tests

- `git clone` the [offical repo](https://github.com/munificent/craftinginterpreters) in the same directory of this one
- Follow the readme to install all dependencies and prepare for testing
- For example for running the tests for chapter 6 for `glox` use:
```bash
dart tool/bin/test.dart chap06_parsing -i [path/to/bin/glox] --arguments --disable-extras=true
```
