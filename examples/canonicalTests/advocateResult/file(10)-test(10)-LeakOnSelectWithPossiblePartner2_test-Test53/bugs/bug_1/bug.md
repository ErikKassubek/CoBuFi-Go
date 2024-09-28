# Leak: Leak of select with possible partner

The analyzer detected a leak of a select with a possible partner.
A leak of a select is a situation, where a select is still blocking at the end of the program.
The partner is a corresponding send or receive operation, which communicated with another operation, but could communicated with the stuck operation instead, resolving the leak.

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
        select {        // <------- Leak
        case c <- 1:    // <------- Possible partner
        }
    }()
}
```

## Test/Program
The bug was found in the following test/program:

- Test/Prog:  Test53
- File:  /home/erik/Uni/HiWi/ADVOCATE/examples/canonicalTests/tests/LeakOnSelectWithPossiblePartner2_test.go

## Bug Elements
The full trace of the recording can be found in the `advocateTrace` folder.

The elements involved in the found leak are located at the following positions:

###  Select:
-> /home/erik/Uni/HiWi/ADVOCATE/examples/canonicalTests/tests/LeakOnSelectWithPossiblePartner2_test.go:33
```go
22 ...
23 
24 
25 	go func() {
26 		time.Sleep(300 * time.Millisecond)
27 
28 		select {
29 		case c <- 1:
30 		case d <- 1:
31 		}
32 	}()
33            // <-------
34 	c <- 1
35 	d <- 1
36 
37 	time.Sleep(800 * time.Millisecond)
38 }
39 
```


###  Channel: Receive
-> /home/erik/Uni/HiWi/ADVOCATE/examples/canonicalTests/tests/LeakOnSelectWithPossiblePartner2_test.go:27
```go
16 ...
17 
18 		<-d
19 	}()
20 
21 	go func() {
22 		<-c
23 	}()
24 
25 	go func() {
26 		time.Sleep(300 * time.Millisecond)
27            // <-------
28 		select {
29 		case c <- 1:
30 		case d <- 1:
31 		}
32 	}()
33 
34 	c <- 1
35 	d <- 1
36 
37 	time.Sleep(800 * time.Millisecond)
38 
39 ...
```


## Replay
The analyzer found a leak in the recorded trace.
The analyzer found a way to resolve the leak, meaning the leak should not reappear in the rewritten trace.

The rewritten trace can be found in the `rewritten_trace` folder.

**Replaying was successful**.

It exited with the following code: 20

The replay was able to get the leaking unbuffered channel or select unstuck.

