# Diagnostics: Select Case without Partner

During the execution of the program, a select was executed, where, based on the happens-before relation, at least one case could never be triggered.
This can be a desired behavior, especially considering, that only executed operations are considered, but it can also be an hint of an unnecessary select case.

## Minimal Example
The following code is a minimal example to visualize the bug type. It is not the code where the bug was found.

```go
func main() {
    c := make(chan int)
    d := make(chan int)
    go func() {
        <-c
    }()

    select{
    case c1 := <- c:
        print(c1)
    case d <- 1:      // <-------
        print("d")
    }

```

## Test/Program
The bug was found in the following test/program:

- Test/Prog:  Test48
- File:  /home/erik/Uni/HiWi/ADVOCATE/examples/canonicalTests/tests/LeakOnBufferedChannelWithPossiblePartner_test.go

## Bug Elements
The full trace of the recording can be found in the `advocateTrace` folder.

The elements involved in the found diagnostics are located at the following positions:

###  Select:
-> /home/erik/Uni/HiWi/ADVOCATE/examples/canonicalTests/tests/LeakOnBufferedChannelWithPossiblePartner_test.go:24


###  
## Replay
The bug is an actual bug. Therefore no rewrite is possibel.

**Replaying was not possible**.

