# Leak: Leak on sync.Cond

The analyzer detected a leak on a sync.Cond.
A leak on a sync.Cond is a situation, where a sync.Cond wait is still blocking at the end of the program.
A sync.Cond wait is blocking, because the condition is not met.

## Minimal Example
The following code is a minimal example to visualize the bug type. It is not the code where the bug was found.

```go
func main() {
    var c sync.Cond

    c.Wait()            // <------- Leak, no signal/broadcast
}
```

## Test/Program
The bug was found in the following test/program:

- Test/Prog:  Test56
- File:  /home/erik/Uni/HiWi/ADVOCATE/examples/canonicalTests/tests/LeakConditionalVariable_test.go

## Bug Elements
The full trace of the recording can be found in the `advocateTrace` folder.

The elements involved in the found leak are located at the following positions:

###  Conditional Variable: Wait
-> /home/erik/Uni/HiWi/ADVOCATE/examples/canonicalTests/tests/LeakConditionalVariable_test.go:25
```go
14 ...
15 
16 
17 	// wait for signal
18 	go func() {
19 		c.L.Lock()
20 		c.Wait()
21 		c.L.Unlock()
22 	}()
23 
24 	time.Sleep(200 * time.Millisecond)
25 }           // <-------
26 
```


## Replay
The analyzer found a leak in the recorded trace.
The analyzer found a way to resolve the leak, meaning the leak should not reappear in the rewritten trace.

The rewritten trace can be found in the `rewritten_trace` folder.

**Replaying failed**.

It exited with the following code: 

No replay info available. Could not find trace number in output.log

