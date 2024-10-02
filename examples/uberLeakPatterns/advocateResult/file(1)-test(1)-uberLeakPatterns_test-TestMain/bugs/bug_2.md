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
-> /home/erik/Uni/HiWi/ADVOCATE/examples/uberLeakPatterns/uberLeakPatterns_test.go:101
```go
90 ...
91 
92    Fix is to have multiple receivers, increase buffer space of data channel, ...
93 
94 */
95 
96 func communicationContention() {
97 	n := 5
98 	data := make(chan int)
99 	for i := 0; i < n; i++ {
100 		go func() {
101 			data <- i           // <-------
102 		}()
103 
104 	}
105 
106 	<-data
107 
108 }
109 
110 // Channel Iteration Misuse
111 
112 
113 ...
```


###  Channel: Receive
-> /home/erik/Uni/HiWi/ADVOCATE/examples/uberLeakPatterns/uberLeakPatterns_test.go:106
```go
95 ...
96 
97 	n := 5
98 	data := make(chan int)
99 	for i := 0; i < n; i++ {
100 		go func() {
101 			data <- i
102 		}()
103 
104 	}
105 
106 	<-data           // <-------
107 
108 }
109 
110 // Channel Iteration Misuse
111 
112 /*
113 
114 1. Once no messages are sent the consumer will block.
115 2. There's a single blocked receiver and all potential senders are not stuck.
116 3. To resolve the leak, the solution is to close the channel.
117 
118 ...
```


## Replay
The analyzer found a leak in the recorded trace.
The analyzer found a way to resolve the leak, meaning the leak should not reappear in the rewritten trace.

**Replaying was successful**.

It exited with the following code: 20

The replay was able to get the leaking unbuffered channel or select unstuck.

