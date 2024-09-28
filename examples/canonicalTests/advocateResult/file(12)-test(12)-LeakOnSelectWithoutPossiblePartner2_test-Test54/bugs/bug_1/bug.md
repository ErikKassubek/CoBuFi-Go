# Leak: Leak on select without possible partner

The analyzer detected a leak of a select without a possible partner.
A leak of a select is a situation, where a select is still blocking at the end of the program.
The analyzer could not find a partner for the stuck operation, which would resolve the leak.

## Minimal Example
The following code is a minimal example to visualize the bug type. It is not the code where the bug was found.

```go
func main() {
    c := make(chan int)

    go func() {
        select {        // <------- Leak, no possible partner
        case c <- 1:
        }
    }()
}
```

## Test/Program
The bug was found in the following test/program:

- Test/Prog:  Test54
- File:  /home/erik/Uni/HiWi/ADVOCATE/examples/canonicalTests/tests/LeakOnSelectWithoutPossiblePartner2_test.go

## Bug Elements
The full trace of the recording can be found in the `advocateTrace` folder.

The elements involved in the found leak are located at the following positions:

###  Select:
-> /home/erik/Uni/HiWi/ADVOCATE/examples/canonicalTests/tests/LeakOnSelectWithoutPossiblePartner2_test.go:23
```go
12 ...
13 
14 	c := make(chan int, 0)
15 	d := make(chan int, 0)
16 
17 	go func() {
18 		select {
19 		case c <- 1:
20 		case d <- 1:
21 		}
22 	}()
23            // <-------
24 	time.Sleep(200 * time.Millisecond)
25 }
26 
```


## Replay
The analyzer found a leak in the recorded trace.
The analyzer could not find a way to resolve the leak.No rewritten trace was created. This does not need to mean, that the leak can not be resolved, especially because the analyzer is only aware of executed operations.

**Replaying was not possible**.

