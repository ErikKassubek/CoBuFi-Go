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

- Test/Prog:  Test47
- File:  /home/erik/Uni/HiWi/ADVOCATE/examples/canonicalTests/tests/LeakOnSelectWithPossiblePartner_test.go

## Bug Elements
The full trace of the recording can be found in the `advocateTrace` folder.

The elements involved in the found diagnostic are located at the following positions:

###  Channel: Receive
-> /home/erik/Uni/HiWi/ADVOCATE/examples/canonicalTests/tests/LeakOnSelectWithPossiblePartner_test.go:30
```go
19 ...
20 
21 		time.Sleep(100 * time.Millisecond)
22 		d <- 1
23 	}()
24 
25 	select {
26 	case <-c:
27 	case <-d:
28 	}
29 
30 	close(c)           // <-------
31 
32 	time.Sleep(300 * time.Millisecond)
33 }
34 
```


###  Channel: Close
-> /home/erik/Uni/HiWi/ADVOCATE/examples/canonicalTests/tests/LeakOnSelectWithPossiblePartner_test.go:35


## Replay
The bug is a potential bug.
The analyzer has tries to rewrite the trace in such a way, that the bug will be triggered when replaying the trace.

The rewritten trace can be found in the `rewritten_trace` folder.

**Replaying failed**.

It exited with the following code: 12

The replay got stuck during the execution.
No trace element was executed for a long tim.
This can be caused by a stuck replay.
Possible causes are:
    - The program was altered between recording and replay
    - The program execution path is not deterministic, e.g. its execution path is determined by a random number
    - The program execution path depends on the order of not tracked operations
    - The program execution depends on outside input, that was not exactly reproduced

