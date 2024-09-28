# Leak: Leak on sync.Mutex

The analyzer detected a leak on a sync.Mutex.
A leak on a sync.Mutex is a situation, where a sync.Mutex lock operations is still blocking at the end of the program.
A sync.Mutex lock operation is a operation, which is blocking, because the lock is already acquired.

## Minimal Example
The following code is a minimal example to visualize the bug type. It is not the code where the bug was found.

```go
func main() {
    var m sync.Mutex

    go func() {
        m.Lock()        // <------- Leak
    }()

    m.Lock()            // <------- Lock, no unlock
}
```

## Test/Program
The bug was found in the following test/program:

- Test/Prog:  Test45
- File:  /home/erik/Uni/HiWi/ADVOCATE/examples/canonicalTests/tests/TP_GoRoutineLeak7_test.go

## Bug Elements
The full trace of the recording can be found in the `advocateTrace` folder.

The elements involved in the found leak are located at the following positions:

###  Mutex: Lock
-> /home/erik/Uni/HiWi/ADVOCATE/examples/canonicalTests/tests/TP_GoRoutineLeak7_test.go:23


###  Mutex: Lock
-> /home/erik/Uni/HiWi/ADVOCATE/examples/canonicalTests/tests/TP_GoRoutineLeak7_test.go:22
```go
11 ...
12 
13 func n45() {
14 	m := sync.Mutex{}
15 
16 	go func() {
17 		m.Lock()
18 		m.Lock()
19 	}()
20 
21 	time.Sleep(100 * time.Millisecond)
22 }           // <-------
23 
```


## Replay
The analyzer found a leak in the recorded trace.
The analyzer found a way to resolve the leak, meaning the leak should not reappear in the rewritten trace.

The rewritten trace can be found in the `rewritten_trace` folder.

**Replaying failed**.

It exited with the following code: 

No replay info available. Could not find trace number in output.log

