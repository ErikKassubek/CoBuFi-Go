# overview Stats

## Program
| Info | Value |
| - | - |
| Number of go files | 1 |
| Number of lines | 942 |
| Number of non-empty lines | 678 |


## Trace
| Info | Value |
| - | - |
| Number of routines | 18 |
| Number of spawns | 4 |
| Number of atomics | 0 |
| Number of atomic operations | 0 |
| Number of channels | 1 |
| Number of channel operations | 3 |
| Number of selects | 0 |
| Number of select cases | 0 |
| Number of select channel operations | 0 |
| Number of select default operations | 0 |
| Number of mutexes | 0 |
| Number of mutex operations | 0 |
| Number of wait groups | 0 |
| Number of wait group operations | 0 |
| Number of cond vars | 0 |
| Number of cond var operations | 0 |
| Number of once | 0| 
| Number of once operations | 0 |


## Times
| Info | Value |
| - | - |
| Time for run without ADVOCATE | 0.303932 s |
| Time for run with ADVOCATE | 0.318076 s |
| Overhead of ADVOCATE | 4.653673 % |
| Replay without changes | 0.319061 s |
| Overhead of Replay | 4.977758 % s |
| Analysis | 0.044729 s |


## Results
==================== Summary ====================

-------------------- Critical -------------------
1 Possible send on closed channel:
	close: /home/erikkassubek/Uni/HiWi/ADVOCATE/examples/constructed/potentialBugs.go:133@32
	send : /home/erikkassubek/Uni/HiWi/ADVOCATE/examples/constructed/potentialBugs.go:125@29
-------------------- Warning --------------------
2 Possible receive on closed channel:
	close: /home/erikkassubek/Uni/HiWi/ADVOCATE/examples/constructed/potentialBugs.go:133@32
	recv : /home/erikkassubek/Uni/HiWi/ADVOCATE/examples/constructed/potentialBugs.go:129@28