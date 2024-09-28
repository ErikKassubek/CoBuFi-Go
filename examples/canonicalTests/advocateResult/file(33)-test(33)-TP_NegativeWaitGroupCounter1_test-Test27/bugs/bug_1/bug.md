# Bug: Possible Negative WaitGroup cCounter

The analyzer detected a possible negative WaitGroup counter.
Although the negative counter did not occur during the recording, it is possible that it will occur, based on the happens before relation.
A negative counter will lead to a panic.

## Minimal Example
The following code is a minimal example to visualize the bug type. It is not the code where the bug was found.

```go
func main() {
    var wg sync.WaitGroup

    go func() {
        wg.Add(1)       // <-------
    }()

    go func() {
        wg.Done()       // <-------
    }()

    wg.Wait()
}
```

## Test/Program
The bug was found in the following test/program:

- Test/Prog:  Test27
- File:  /home/erik/Uni/HiWi/ADVOCATE/examples/canonicalTests/tests/TP_NegativeWaitGroupCounter1_test.go

## Bug Elements
The full trace of the recording can be found in the `advocateTrace` folder.

The elements involved in the found bug are located at the following positions:

###  Waitgroup: Done
-> /home/erik/Uni/HiWi/ADVOCATE/examples/canonicalTests/tests/TP_NegativeWaitGroupCounter1_test.go:44


###  Waitgroup: Add
-> /home/erik/Uni/HiWi/ADVOCATE/examples/canonicalTests/tests/TP_NegativeWaitGroupCounter1_test.go:29
```go
18 ...
19 
20 	}()
21 
22 	go func() {
23 		wg.Add(1)
24 		wg.Add(1)
25 		wg.Done()
26 		// d <- 1
27 	}()
28 
29 	go func() {           // <-------
30 		wg.Add(1)
31 		// <-d
32 		c <- 1
33 	}()
34 
35 	<-c
36 
37 	time.Sleep(100 * time.Millisecond) // prevent negative wait counter
38 	wg.Done()
39 	wg.Done()
40 
41 ...
```


## Replay
The bug is a potential bug.
The analyzer has tries to rewrite the trace in such a way, that the bug will be triggered when replaying the trace.

The rewritten trace can be found in the `rewritten_trace` folder.

**Replaying was successful**.

It exited with the following code: 32

The replay resulted in an expected negative wait group triggering a panic. The bug was triggered. The replay was therefore able to confirm, that the negative wait group can actually occur.

