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

- Test/Prog:  Test11
- File:  /home/erik/Uni/HiWi/ADVOCATE/examples/canonicalTests/tests/TP_SendOnClosed_test.go

## Bug Elements
The full trace of the recording can be found in the `advocateTrace` folder.

The elements involved in the found bug are located at the following positions:

###  Channel: Send
-> /home/erik/Uni/HiWi/ADVOCATE/examples/canonicalTests/tests/TP_SendOnClosed_test.go:27
```go
16 ...
17 
18 	}()
19 
20 	go func() {
21 		select {
22 		case c <- struct{}{}:
23 		default:
24 		}
25 	}()
26 
27 	go func() {           // <-------
28 		time.Sleep(100 * time.Millisecond)
29 		select {
30 		case <-c:
31 		default:
32 		}
33 	}()
34 
35 	time.Sleep(300 * time.Millisecond) // make sure, that the default values are taken
36 }
37 
```


###  Channel: Close
-> /home/erik/Uni/HiWi/ADVOCATE/examples/canonicalTests/tests/TP_SendOnClosed_test.go:22
```go
11 ...
12 
13 	c := make(chan struct{}, 0)
14 
15 	go func() {
16 		time.Sleep(200 * time.Millisecond) // prevent actual send on closed channel
17 		close(c)
18 	}()
19 
20 	go func() {
21 		select {
22 		case c <- struct{}{}:           // <-------
23 		default:
24 		}
25 	}()
26 
27 	go func() {
28 		time.Sleep(100 * time.Millisecond)
29 		select {
30 		case <-c:
31 		default:
32 		}
33 
34 ...
```


## Replay
The bug is a potential bug.
The analyzer has tries to rewrite the trace in such a way, that the bug will be triggered when replaying the trace.

The rewritten trace can be found in the `rewritten_trace` folder.

**Replaying was successful**.

It exited with the following code: 30

The replay resulted in an expected send on close triggering a panic. The bug was triggered. The replay was therefore able to confirm, that the send on closed can actually occur.

