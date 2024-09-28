# Diagnostic: Possible Receive on Closed Channel

The analyzer detected a possible receive on a closed channel.
Although the receive on a closed channel did not occur during the recording, it is possible that it will occur, based on the happens before relation.This is not necessarily a bug, but it can be an indication of a bug.

## Minimal Example
The following code is a minimal example to visualize the bug type. It is not the code where the bug was found.

```go
func main() {
    c := make(chan int)

    go func() {
        c <- 1
    }()

    go func() {
        <- c            // <-------
    }()

    close(c)            // <-------
}
```

## Test/Program
The bug was found in the following test/program:

- Test/Prog:  Test16
- File:  /home/erik/Uni/HiWi/ADVOCATE/examples/canonicalTests/tests/TP_SendOnClosed4_test.go

## Bug Elements
The full trace of the recording can be found in the `advocateTrace` folder.

The elements involved in the found diagnostic are located at the following positions:

###  Channel: Receive
-> /home/erik/Uni/HiWi/ADVOCATE/examples/canonicalTests/tests/TP_SendOnClosed4_test.go:33
```go
22 ...
23 
24 		}
25 	}()
26 
27 	go func() {
28 		<-c
29 	}()
30 
31 	time.Sleep(100 * time.Millisecond)
32 	close(c)
33 }           // <-------
34 
```


###  Channel: Close
-> /home/erik/Uni/HiWi/ADVOCATE/examples/canonicalTests/tests/TP_SendOnClosed4_test.go:37


## Replay
The bug is a potential bug.
The analyzer has tries to rewrite the trace in such a way, that the bug will be triggered when replaying the trace.

The rewritten trace can be found in the `rewritten_trace` folder.

**Replaying failed**.

It exited with the following code: 3

During the replay, the program panicked unexpectedly.
This can be expected behavior, e.g. if the program tries to replay a recv on closed but the recv on closed is necessarily preceded by a send on closed.

