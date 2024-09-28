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

- Test/Prog:  Test51
- File:  /home/erik/Uni/HiWi/ADVOCATE/examples/canonicalTests/tests/LeakChannelWithPossiblePartner_test.go

## Bug Elements
The full trace of the recording can be found in the `advocateTrace` folder.

The elements involved in the found leak are located at the following positions:

###  Channel: Send
-> /home/erik/Uni/HiWi/ADVOCATE/examples/canonicalTests/tests/LeakChannelWithPossiblePartner_test.go:29
```go
18 ...
19 
20 		println(1)
21 	}()
22 
23 	go func() {
24 		c <- 1
25 		println(2)
26 	}()
27 
28 	<-c
29 	time.Sleep(200 * time.Millisecond)           // <-------
30 }
31 
```


###  Channel: Receive
-> /home/erik/Uni/HiWi/ADVOCATE/examples/canonicalTests/tests/LeakChannelWithPossiblePartner_test.go:33


## Replay
The analyzer found a leak in the recorded trace.
The analyzer found a way to resolve the leak, meaning the leak should not reappear in the rewritten trace.

The rewritten trace can be found in the `rewritten_trace` folder.

**Replaying was successful**.

It exited with the following code: 20

The replay was able to get the leaking unbuffered channel or select unstuck.

