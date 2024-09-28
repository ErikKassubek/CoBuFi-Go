# Diagnostics: Concurrent Receive

During the execution of the program, a channel waited to receive at multiple positions at the same time.
In this case, the actual receiver of a send message is chosen randomly.
This can lead to nondeterministic behavior.

## Minimal Example
The following code is a minimal example to visualize the bug type. It is not the code where the bug was found.

```go
func main() {
    c := make(chan int, 1)

    go func() {
        <-c             // <-------
    }()

    go func() {
        <-c             // <-------
    }()

    c <- 1
}
```

## Test/Program
The bug was found in the following test/program:

- Test/Prog:  Test21
- File:  /home/erik/Uni/HiWi/ADVOCATE/examples/canonicalTests/tests/TP_ConcurrentRecSameChannel_test.go

## Bug Elements
The full trace of the recording can be found in the `advocateTrace` folder.

The elements involved in the found diagnostics are located at the following positions:

###  Channel: Receive
-> /home/erik/Uni/HiWi/ADVOCATE/examples/canonicalTests/tests/TP_ConcurrentRecSameChannel_test.go:23
```go
12 ...
13 
14 func n21() {
15 	x := make(chan int)
16 
17 	go func() {
18 		<-x
19 	}()
20 
21 	go func() {
22 		<-x
23 	}()           // <-------
24 
25 	time.Sleep(100 * time.Millisecond)
26 
27 	x <- 1
28 	x <- 1
29 
30 	time.Sleep(300 * time.Millisecond)
31 }
32 
```


###  Channel: Receive
-> /home/erik/Uni/HiWi/ADVOCATE/examples/canonicalTests/tests/TP_ConcurrentRecSameChannel_test.go:27
```go
16 ...
17 
18 		<-x
19 	}()
20 
21 	go func() {
22 		<-x
23 	}()
24 
25 	time.Sleep(100 * time.Millisecond)
26 
27 	x <- 1           // <-------
28 	x <- 1
29 
30 	time.Sleep(300 * time.Millisecond)
31 }
32 
```


## Replay
The bug is an actual bug. Therefore no rewrite is possibel.

**Replaying was not possible**.

