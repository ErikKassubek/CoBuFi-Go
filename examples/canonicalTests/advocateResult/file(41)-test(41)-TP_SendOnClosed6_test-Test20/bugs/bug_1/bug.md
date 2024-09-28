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

- Test/Prog:  Test20
- File:  /home/erik/Uni/HiWi/ADVOCATE/examples/canonicalTests/tests/TP_SendOnClosed6_test.go

## Bug Elements
The full trace of the recording can be found in the `advocateTrace` folder.

The elements involved in the found bug are located at the following positions:

###  Channel: Send
-> /home/erik/Uni/HiWi/ADVOCATE/examples/canonicalTests/tests/TP_SendOnClosed6_test.go:22
```go
11 ...
12 
13 func n20() {
14 	ch := make(chan int, 2)
15 
16 	f := func() {
17 		ch <- 1
18 	}
19 
20 	go func() {
21 		f()
22 	}()           // <-------
23 
24 	time.Sleep(200 * time.Millisecond)
25 	close(ch)
26 }
27 
```


###  Channel: Close
-> /home/erik/Uni/HiWi/ADVOCATE/examples/canonicalTests/tests/TP_SendOnClosed6_test.go:30


## Replay
The bug is a potential bug.
The analyzer has tries to rewrite the trace in such a way, that the bug will be triggered when replaying the trace.

The rewritten trace can be found in the `rewritten_trace` folder.

**Replaying was successful**.

It exited with the following code: 30

The replay resulted in an expected send on close triggering a panic. The bug was triggered. The replay was therefore able to confirm, that the send on closed can actually occur.

