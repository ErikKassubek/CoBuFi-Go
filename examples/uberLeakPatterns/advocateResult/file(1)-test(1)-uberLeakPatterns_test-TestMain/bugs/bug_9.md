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
-> /home/erik/Uni/HiWi/ADVOCATE/examples/uberLeakPatterns/uberLeakPatterns_test.go:100
```go
89 ...
90 
91  (message contention)
92    Fix is to have multiple receivers, increase buffer space of data channel, ...
93 
94 */
95 
96 func communicationContention() {
97 	n := 5
98 	data := make(chan int)
99 	for i := 0; i < n; i++ {
100 		go func() {           // <-------
101 			data <- i
102 		}()
103 
104 	}
105 
106 	<-data
107 
108 }
109 
110 // Channel Iteration Misuse
111 
112 ...
```


## Replay


**Replaying was not run**.

The analyzer was not able to rewrite the bug.
This can be because the bug is an actual bug, because the bug is a leak without a possible partner or blocking operations or because the analyzer was not able to rewrite the trace for other reasons.

