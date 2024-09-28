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

- Test/Prog:  Test09
- File:  /home/erik/Uni/HiWi/ADVOCATE/examples/canonicalTests/tests/TP_SendAndRecOnClosed4_test.go

## Bug Elements
The full trace of the recording can be found in the `advocateTrace` folder.

The elements involved in the found bug are located at the following positions:

###  Channel: Send
-> /home/erik/Uni/HiWi/ADVOCATE/examples/canonicalTests/tests/TP_SendAndRecOnClosed4_test.go:30
```go
19 ...
20 
21 	}()
22 
23 	go func() {
24 		select {
25 		case c <- struct{}{}:
26 		default:
27 		}
28 
29 		select {
30 		case <-d:           // <-------
31 		default:
32 		}
33 	}()
34 
35 	d <- struct{}{}
36 	<-c
37 
38 	time.Sleep(1 * time.Second) // prevent termination before receive
39 }
40 
```


###  Channel: Close
-> /home/erik/Uni/HiWi/ADVOCATE/examples/canonicalTests/tests/TP_SendAndRecOnClosed4_test.go:24
```go
13 ...
14 
15 	d := make(chan struct{}, 1)
16 
17 	go func() {
18 		time.Sleep(300 * time.Millisecond) // prevent actual send on closed channel
19 		close(c)
20 		close(d)
21 	}()
22 
23 	go func() {
24 		select {           // <-------
25 		case c <- struct{}{}:
26 		default:
27 		}
28 
29 		select {
30 		case <-d:
31 		default:
32 		}
33 	}()
34 
35 
36 ...
```


## Replay
The bug is a potential bug.
The analyzer has tries to rewrite the trace in such a way, that the bug will be triggered when replaying the trace.

The rewritten trace can be found in the `rewritten_trace` folder.

**Replaying was successful**.

It exited with the following code: 30

The replay resulted in an expected send on close triggering a panic. The bug was triggered. The replay was therefore able to confirm, that the send on closed can actually occur.

