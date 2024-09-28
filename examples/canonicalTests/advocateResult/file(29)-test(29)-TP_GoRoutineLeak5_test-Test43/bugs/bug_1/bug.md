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

- Test/Prog:  Test43
- File:  /home/erik/Uni/HiWi/ADVOCATE/examples/canonicalTests/tests/TP_GoRoutineLeak5_test.go

## Bug Elements
The full trace of the recording can be found in the `advocateTrace` folder.

The elements involved in the found leak are located at the following positions:

###  Channel: Send
-> /home/erik/Uni/HiWi/ADVOCATE/examples/canonicalTests/tests/TP_GoRoutineLeak5_test.go:24
```go
13 ...
14 
15 		<-c
16 	}()
17 
18 	go func() {
19 		c <- 1
20 	}()
21 
22 	go func() {
23 		c <- 1
24 	}()           // <-------
25 
26 	time.Sleep(100 * time.Millisecond)
27 }
28 
```


###  Channel: Receive
-> /home/erik/Uni/HiWi/ADVOCATE/examples/canonicalTests/tests/TP_GoRoutineLeak5_test.go:20
```go
9 ...
10 
11 func n43() {
12 	c := make(chan int, 0)
13 
14 	go func() {
15 		<-c
16 	}()
17 
18 	go func() {
19 		c <- 1
20 	}()           // <-------
21 
22 	go func() {
23 		c <- 1
24 	}()
25 
26 	time.Sleep(100 * time.Millisecond)
27 }
28 
```


## Replay
The analyzer found a leak in the recorded trace.
The analyzer found a way to resolve the leak, meaning the leak should not reappear in the rewritten trace.

The rewritten trace can be found in the `rewritten_trace` folder.

**Replaying was successful**.

It exited with the following code: 20

The replay was able to get the leaking unbuffered channel or select unstuck.

