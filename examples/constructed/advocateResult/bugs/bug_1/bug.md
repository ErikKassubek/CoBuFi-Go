# Bug: Possible unlock of not locked mutex

The analyzer detected a possible unlock on a not locked mutex.
Although the unlock of a not locked mutex did not occur during the recording, it is possible that it will occur, based on the happens before relation.
A unlock of a not locked mutex will result in a panic.

## Minimal Example
The following code is a minimal example to visualize the bug type. It is not the code where the bug was found.

```go
func main() {
    var m sync.Mutex

    go func() {
        m.Lock()       // <-------
    }()

    go func() {
        m.Unlock()     // <-------
    }()

}
```

## Test/Program
The bug was found in the following test/program:

- Test: unknown
- File: unknown

## Bug Elements
The full trace of the recording can be found in the `advocateTrace` folder.

The elements involved in the found bug are located at the following positions:

###  Mutex: RUnlock
-> /home/erik/Uni/HiWi/ADVOCATE/examples/constructed/constructed.go:1125
```go
1114 ...
1115 
1116 	}()
1117 
1118 	go func() {
1119 		m.RLock()
1120 	}()
1121 
1122 	go func() {
1123 		time.Sleep(100 * time.Millisecond)
1124 		m.RUnlock()
1125 	}()           // <-------
1126 
1127 	time.Sleep(100 * time.Millisecond)
1128 	m.RUnlock()
1129 
1130 	time.Sleep(200 * time.Millisecond)
1131 
1132 }
1133 
1134 // =============== use for testing ===============
1135 // MARK: FOR TESTING
1136 
1137 ...
```


-> /home/erik/Uni/HiWi/ADVOCATE/examples/constructed/constructed.go:1129
```go
1118 ...
1119 
1120 	}()
1121 
1122 	go func() {
1123 		time.Sleep(100 * time.Millisecond)
1124 		m.RUnlock()
1125 	}()
1126 
1127 	time.Sleep(100 * time.Millisecond)
1128 	m.RUnlock()
1129            // <-------
1130 	time.Sleep(200 * time.Millisecond)
1131 
1132 }
1133 
1134 // =============== use for testing ===============
1135 // MARK: FOR TESTING
1136 // leak because of wait group
1137 func nTest() {
1138 	c := make(chan int, 0)
1139 	m := sync.Mutex{}
1140 
1141 ...
```


###  Mutex: RLock
-> /home/erik/Uni/HiWi/ADVOCATE/examples/constructed/constructed.go:1116
```go
1105 ...
1106 
1107 	m.Unlock()
1108 }
1109 
1110 // possible unlock of locked
1111 func n59() {
1112 	var m sync.RWMutex
1113 
1114 	go func() {
1115 		m.RLock()
1116 	}()           // <-------
1117 
1118 	go func() {
1119 		m.RLock()
1120 	}()
1121 
1122 	go func() {
1123 		time.Sleep(100 * time.Millisecond)
1124 		m.RUnlock()
1125 	}()
1126 
1127 
1128 ...
```


-> /home/erik/Uni/HiWi/ADVOCATE/examples/constructed/constructed.go:1120
```go
1109 ...
1110 
1111 func n59() {
1112 	var m sync.RWMutex
1113 
1114 	go func() {
1115 		m.RLock()
1116 	}()
1117 
1118 	go func() {
1119 		m.RLock()
1120 	}()           // <-------
1121 
1122 	go func() {
1123 		time.Sleep(100 * time.Millisecond)
1124 		m.RUnlock()
1125 	}()
1126 
1127 	time.Sleep(100 * time.Millisecond)
1128 	m.RUnlock()
1129 
1130 	time.Sleep(200 * time.Millisecond)
1131 
1132 ...
```


## Replay
The bug is a potential bug.
The analyzer has tries to rewrite the trace in such a way, that the bug will be triggered when replaying the trace.

The rewritten trace can be found in the `rewritten_trace` folder.

**Replaying was successful**.

It exited with the following code: 33

The replay resulted in an expected lock of an unlocked mutex triggering a panic. The bug was triggered. The replay was therefore able to confirm, that the unlock of a not locked mutex can actually occur.

