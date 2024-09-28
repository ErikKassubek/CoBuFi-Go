# Diagnostics: Actual Receive on Closed Channel

During the execution of the program, a receive on a closed channel occurred.


## Minimal Example
The following code is a minimal example to visualize the bug type. It is not the code where the bug was found.

```go
func main() {
    c := make(chan int)
    close(c)          // <-------
    <-c               // <-------
}
```

## Test/Program
The bug was found in the following test/program:

- Test/Prog:  Test08
- File:  /home/erik/Uni/HiWi/ADVOCATE/examples/canonicalTests/tests/TP_RecOnClosed_test.go

## Bug Elements
The full trace of the recording can be found in the `advocateTrace` folder.

The elements involved in the found diagnostics are located at the following positions:

###  Channel: Receive
-> /home/erik/Uni/HiWi/ADVOCATE/examples/canonicalTests/tests/TP_RecOnClosed_test.go:24
```go
13 ...
14 
15 	c := make(chan int)
16 
17 	go func() {
18 		time.Sleep(300 * time.Millisecond) // force actual recv on closed channel
19 		<-c
20 	}()
21 
22 	close(c)
23 	time.Sleep(1 * time.Second) // prevent termination before receive
24 }           // <-------
25 
```


###  Channel: Close
-> /home/erik/Uni/HiWi/ADVOCATE/examples/canonicalTests/tests/TP_RecOnClosed_test.go:27


## Replay
The bug is an actual bug. Therefore no rewrite is possibel.

**Replaying was not possible**.

