# Bug: Possible Send on Closed Channel

The analyzer detected a possible send on a closed channel.
Although the send on a closed channel did not occur during the recording, it is possible that it will occur, based on the happens before relation.
Such a send on a closed channel leads to a panic.

## Minimal Example
The following code is a minimal example to visualize the bug type. It is not the code where the bug was found.

```go
func main() {
    c := make(chan int)

    go func() {
        c <- 1          // <-------
    }()

    go func() {
        <- c
    }()

    close(c)            // <-------
}
```

## Test/Program
The bug was found in the following test/program:

- Test/Prog:  Test14
- File:  /home/erik/Uni/HiWi/ADVOCATE/examples/canonicalTests/tests/TP_SendOnClosed3_test.go

## Bug Elements
The full trace of the recording can be found in the `advocateTrace` folder.

The elements involved in the found bug are located at the following positions:

###  Channel: Send
-> /home/erik/Uni/HiWi/ADVOCATE/examples/canonicalTests/tests/TP_SendOnClosed3_test.go:26
```go
15 ...
16 
17 	once := sync.Once{}
18 
19 	go func() {
20 		once.Do(func() {
21 			c <- 1
22 		})
23 	}()
24 
25 	go func() {
26 		time.Sleep(100 * time.Millisecond) // prevent actual send on closed channel           // <-------
27 		once.Do(func() {
28 			// do nothing
29 		})
30 	}()
31 
32 	time.Sleep(200 * time.Millisecond)
33 	close(c)
34 }
35 
```


###  Channel: Close
-> /home/erik/Uni/HiWi/ADVOCATE/examples/canonicalTests/tests/TP_SendOnClosed3_test.go:38


## Replay
The bug is a potential bug.
The analyzer has tries to rewrite the trace in such a way, that the bug will be triggered when replaying the trace.

The rewritten trace can be found in the `rewritten_trace` folder.

**Replaying was successful**.

It exited with the following code: 30

The replay resulted in an expected send on close triggering a panic. The bug was triggered. The replay was therefore able to confirm, that the send on closed can actually occur.

