# Leak: Leak on unbuffered Channel without possible partner

The analyzer detected a leak of an unbuffered channel without a possible partner.
A leak of an unbuffered channel is a situation, where a unbuffered channel is still blocking at the end of the program.
The analyzer could not find a partner for the stuck operation, which would resolve the leak.

## Minimal Example
The following code is a minimal example to visualize the bug type. It is not the code where the bug was found.

```go
func main() {
    c := make(chan int)

    go func() {
        c <- 1          // <------- Leak, no possible partner
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

###  Channel: Send
-> /home/erik/Uni/HiWi/ADVOCATE/examples/uberLeakPatterns/uberLeakPatterns_test.go:38
```go
27 ...
28 
29 
30 */
31 
32 func unbalancedConditionalCommunication() {
33 	c := make(chan int)
34 
35 	go func() {
36 		err := rand.Intn(2) == 1
37 		if err {
38 			c <- 1           // <-------
39 		} else {
40 			c <- 0
41 		}
42 
43 	}()
44 
45 	return // early
46 
47 	<-c
48 
49 
50 ...
```


## Replay
The analyzer found a leak in the recorded trace.
The analyzer could not find a way to resolve the leak.No rewritten trace was created. This does not need to mean, that the leak can not be resolved, especially because the analyzer is only aware of executed operations.

The rewritten trace can be found in the `rewritten_trace` folder.

**Replaying was not run**.

The analyzer was not able to rewrite the bug.
This can be because the bug is an actual bug, because the bug is a leak without a possible partner or blocking operations or because the analyzer was not able to rewrite the trace for other reasons.

