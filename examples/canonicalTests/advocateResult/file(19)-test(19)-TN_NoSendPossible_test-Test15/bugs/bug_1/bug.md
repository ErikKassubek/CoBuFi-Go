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

- Test/Prog:  Test15
- File:  /home/erik/Uni/HiWi/ADVOCATE/examples/canonicalTests/tests/TN_NoSendPossible_test.go

## Bug Elements
The full trace of the recording can be found in the `advocateTrace` folder.

The elements involved in the found bug are located at the following positions:

###  Channel: Send
-> /home/erik/Uni/HiWi/ADVOCATE/examples/canonicalTests/tests/TN_NoSendPossible_test.go:26
```go
15 ...
16 
17 
18 	go func() {
19 		t := m.TryLock()
20 		if t {
21 			c <- 1
22 			m.Unlock()
23 		}
24 	}()
25 
26 	go func() {           // <-------
27 		time.Sleep(100 * time.Millisecond)
28 		t := m.TryLock()
29 		if t {
30 			m.Unlock()
31 		}
32 		<-c
33 	}()
34 
35 	time.Sleep(1000 * time.Millisecond)
36 	close(c)
37 
38 ...
```


###  Channel: Close
-> /home/erik/Uni/HiWi/ADVOCATE/examples/canonicalTests/tests/TN_NoSendPossible_test.go:41


## Replay
The bug is a potential bug.
The analyzer has tries to rewrite the trace in such a way, that the bug will be triggered when replaying the trace.

The rewritten trace can be found in the `rewritten_trace` folder.

**Replaying was successful**.

It exited with the following code: 30

The replay resulted in an expected send on close triggering a panic. The bug was triggered. The replay was therefore able to confirm, that the send on closed can actually occur.

