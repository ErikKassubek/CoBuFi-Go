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

- Test/Prog: /home/erik/Uni/HiWi/ADVOCATE/examples/uberLeakPatterns/uberLeakPatterns_test.go
- File: /home/erik/Uni/HiWi/ADVOCATE/examples/uberLeakPatterns/uberLeakPatterns_test.go

## Bug Elements
The full trace of the recording can be found in the `advocateTrace` folder.

The elements involved in the found leak are located at the following positions:

###  Channel: Send
-> /home/erik/Uni/HiWi/ADVOCATE/examples/uberLeakPatterns/uberLeakPatterns_test.go:71
```go
60 ...
61 
62 
63 */
64 
65 func timeOutBug() {
66 
67 	done := make(chan int)
68 
69 	go func() {
70 		time.Sleep(2 * time.Second)
71 		done <- 1           // <-------
72 		fmt.Printf("done")
73 	}()
74 
75 	select {
76 	case <-done:
77 	case <-time.After(1 * time.Second):
78 
79 	}
80 
81 }
82 
83 ...
```


###  Channel: Receive
-> /home/erik/Uni/HiWi/ADVOCATE/examples/uberLeakPatterns/uberLeakPatterns_test.go:75
```go
64 ...
65 
66 
67 	done := make(chan int)
68 
69 	go func() {
70 		time.Sleep(2 * time.Second)
71 		done <- 1
72 		fmt.Printf("done")
73 	}()
74 
75 	select {           // <-------
76 	case <-done:
77 	case <-time.After(1 * time.Second):
78 
79 	}
80 
81 }
82 
83 // The NCast Leak
84 
85 /*
86 
87 ...
```


## Replay
The analyzer found a leak in the recorded trace.
The analyzer found a way to resolve the leak, meaning the leak should not reappear in the rewritten trace.

The rewritten trace can be found in the `rewritten_trace` folder.

**Replaying was successful**.

It exited with the following code: 20

The replay was able to get the leaking unbuffered channel or select unstuck.

