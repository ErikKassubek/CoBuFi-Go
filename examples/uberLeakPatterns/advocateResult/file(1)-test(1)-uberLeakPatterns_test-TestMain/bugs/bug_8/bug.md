# Leak: Leak on routine without blocking operation

The analyzer detected a leak on a routine without a blocking operations.
This means that the routine was terminated because of a panic in another routine or because the main routine terminated while this routine was still running.
This can be a desired behavior, but it can also be a signal for a not otherwise detected block.

## Minimal Example
The following code is a minimal example to visualize the bug type. It is not the code where the bug was found.

```go
func main() {
    go func() {
        time.Sleep(time.Second)          // <------- Is still running when main routine terminates
    }()

}
```

## Test/Program
The bug was found in the following test/program:

- Test/Prog: /home/erik/Uni/HiWi/ADVOCATE/examples/uberLeakPatterns/uberLeakPatterns_test.go
- File: /home/erik/Uni/HiWi/ADVOCATE/examples/uberLeakPatterns/uberLeakPatterns_test.go

## Bug Elements
The full trace of the recording can be found in the `advocateTrace` folder.

The elements involved in the found leak are located at the following positions:

###  Routine
-> /home/erik/Uni/HiWi/ADVOCATE/examples/uberLeakPatterns/uberLeakPatterns_test.go:69
```go
58 ...
59 
60 5. This is useful information to lead to the fix
61    where we use buffered (done) channel.
62 
63 */
64 
65 func timeOutBug() {
66 
67 	done := make(chan int)
68 
69 	go func() {           // <-------
70 		time.Sleep(2 * time.Second)
71 		done <- 1
72 		fmt.Printf("done")
73 	}()
74 
75 	select {
76 	case <-done:
77 	case <-time.After(1 * time.Second):
78 
79 	}
80 
81 ...
```


## Replay


The rewritten trace can be found in the `rewritten_trace` folder.

**Replaying was not run**.

The analyzer was not able to rewrite the bug.
This can be because the bug is an actual bug, because the bug is a leak without a possible partner or blocking operations or because the analyzer was not able to rewrite the trace for other reasons.

