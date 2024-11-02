# Toolchain

This toolchain is intended to applay the analyzer onto tests and programs.

## Preparation
Make sure to first build the analyzer and the runtime.

## Usage
The script can be run with
```
./toolchain main [args]
```
to run a full program with a main function or with
```
./toolchain tests [args]
```
to run the unit tests.
The following args are required:

- `-a [path]`: path to the ADVOCATE directory
- `-f [path]`: path to the program containing the test files

For main, the following arg is also required

- `-E [name]`: only for main, name of the executable of the program

For test, the following arg can be set to run only one test. If it is not set, all tests will be run

- `-n [name]`: name of the test to run

The following arguments can be set:

- `-t`: if set, the toolchain will measure the runtime of the runs and analysis. It will also run the tests/the program without any recording or replay to measure a base time
- `-m`: if set, the toolchain check if there are relevant operations in the program, that have never been executed in the runs
- `-s`: create a file containing statistics about the program runs

If either `-t` or `-s` is set, the following arg must be set:

- `-N` [name]: name of the analyzed program



Its result and additional information (rewritten traces, logs, etc) will be written to `advocateResult`.
