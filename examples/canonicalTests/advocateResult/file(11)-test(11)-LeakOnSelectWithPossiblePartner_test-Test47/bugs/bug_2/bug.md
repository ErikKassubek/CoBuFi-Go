# Leak: Leak of unbuffered Channel with possible partner

The analyzer detected a leak of an unbuffered channel with a possible partner.
A leak of an unbuffered channel is a situation, where a unbuffered channel is still blocking at the end of the program.
The partner is a corresponding send or receive operation, which communicated with another operation, but could communicated with the stuck operation instead, resolving the deadlock.

## Minimal Example
The following code is a minimal example to visualize the bug type. It is not the code where the bug was found.

```go
func main() {
    c := make(chan int)

    go func() {
        c <- 1          // <------- Communicates
    }()

    go func() {
        <- c            // <------- Communicates, possible partner
    }()

    go func() {
        c <- 1          // <------- Leak
    }()
}
```

## Test/Program
The bug was found in the following test/program:

- Test/Prog:  Test47
- File:  /home/erik/Uni/HiWi/ADVOCATE/examples/canonicalTests/tests/LeakOnSelectWithPossiblePartner_test.go

## Bug Elements
The full trace of the recording can be found in the `advocateTrace` folder.

The elements involved in the found leak are located at the following positions:

###  Channel: Send
-> /home/erik/Uni/HiWi/ADVOCATE/examples/canonicalTests/tests/LeakOnSelectWithPossiblePartner_test.go:27
```go
16 ...
17 
18 	}()
19 
20 	go func() {
21 		time.Sleep(100 * time.Millisecond)
22 		d <- 1
23 	}()
24 
25 	select {
26 	case <-c:
27 	case <-d:           // <-------
28 	}
29 
30 	close(c)
31 
32 	time.Sleep(300 * time.Millisecond)
33 }
34 
```


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


## Replay
The analyzer found a leak in the recorded trace.
The analyzer found a way to resolve the leak, meaning the leak should not reappear in the rewritten trace.

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

