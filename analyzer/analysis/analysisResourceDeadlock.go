// Copyrigth (c) 2024 Erik Kassubek
//
// File: analysisResourceDeadlock.go
// Brief: Alternative analysis for cyclic mutex deadlocks.
//
// Author: Sebastian Pohsner
// Created: 2025-01-01
//
// License: BSD-3-Clause

package analysis

import (
	"analyzer/clock"
	"fmt"
	"log"
	"strconv"
	"strings"
)

// Computation of "abstract" lock dependencies
// Lock dependencies are 3-tuples of the following form:
//    (ThreadID, Lock, LockSet)
// Lock dependencies are computed thread-local.
// For each thread there might be several (acquire) events that lead to "lock" acquired under some "lockset".
//
// Each acquire event carries its own vector clock.
// We wish to make use of vector clocks to elimiante infeasible replay candidates.
//
// This means that lock dependencies are 4-tuples of the following form:
//    (ThreadID, Lock, LockSet, []Event)

////////////////////////
// DATA STRUCTURES

// Lock dependencies are computed thread-local. We make use of the following structures.

type Thread struct {
	thread_id ThreadId
	ls        Lockset // The thread's current lockset.
	ldeps     map[LockId][]Dep
}

// Unfortunately, we can't use double-indexed map of the following form in Go.
// type Deps map[Lock]map[Lockset][]Event
// Hence, we introduce some intermediate structure.

type Dep struct {
	ls       Lockset
	requests []Event
}

// Representation of vector clocks, events, threads, lock and lockset.

type Event struct {
	thread_id ThreadId
	vc        clock.VectorClock
}

type ThreadId int
type LockId int
type Lockset map[LockId]struct{}

///////////////////////////////
// ALGORITHM
//
// There are two phases.
//  1. Recording of lock dependencies.
//  2. Checking if lock dependencies imply a cycle.

type State struct {
	threads map[ThreadId]Thread // Recording lock dependencies in phase 1
	cycles  []Cycle             // Computing cycles in phase 2
}

// TODO this isn't very pretty and will never get reset..
var currentState = State{
	threads: make(map[ThreadId]Thread),
	cycles:  make([]Cycle, 0),
}

// Algorithm phase 1

// We show the event processing functions for acquire and release.

func acquire(s *State, lockId LockId, e Event) {
	if _, ok := s.threads[e.thread_id]; !ok {
		s.threads[e.thread_id] = Thread{
			thread_id: e.thread_id,
			ls:        make(Lockset),
			ldeps:     make(map[LockId][]Dep),
		}
	}

	ls := s.threads[e.thread_id].ls
	if !ls.empty() {
		deps := s.threads[e.thread_id].ldeps
		deps[lockId] = insert(deps[lockId], ls, e)
		// In an actual implementation we would record e's vector clock.
		// In fact, we record the vector clock of the associated request.

	}
	s.threads[e.thread_id].ls.add(lockId)
}

func release(s *State, x LockId, e Event) {
	s.threads[e.thread_id].ls.remove(x)
}

// Insert a new lock dependency for a given thread and lock x.
// We assume that event e acquired lock x.
// We might have already an entry that shares the same lock and lockset!
func insert(ds []Dep, ls Lockset, e Event) []Dep {
	for i, v := range ds {
		if v.ls.equal(ls) {
			ds[i].requests = append(ds[i].requests, e)
			return ds
		}
	}
	return append(ds, Dep{ls.Clone(), []Event{e}})
}

// The above insert function records all requests that share the same dependency (tid,l,ls).
// In case of loops, we may end up with many request entrys.
// For performance reasons, we may want to reduce their size.
//
// Eviction strategy.
// Insert variant where we evict event an already stored event f by e,
// if in between f and e no intra-thread synchronization took place.
// This can be checked via helper function equalModuloTID.
// Assumption: Vector clocks underapproximate the must happen-before relation.
func insert2(ds []Dep, ls Lockset, e Event) []Dep {
	// Helper function.
	// Assumes that vc1 and vc2 are connected to two events that are from the same thread tid.
	// Yields true if vc1[k] == vc2[k] for all threads k but tid.
	equalModuloTID := func(tid ThreadId, vc1 clock.VectorClock, vc2 clock.VectorClock) bool {
		if vc1.GetSize() != vc2.GetSize() {
			return false
		}

		for i := 1; i <= vc1.GetSize(); i++ {
			if i == int(tid) {
				continue
			}

			if vc1.GetClock()[i] != vc2.GetClock()[i] {
				return false
			}
		}

		return true
	}

	for i, v := range ds {
		if v.ls.equal(ls) {
			add_vc := true

			for j, f := range ds[i].requests {
				if equalModuloTID(e.thread_id, e.vc, f.vc) {
					ds[i].requests[j] = e
					add_vc = false
				}

			}

			if add_vc {
				ds[i].requests = append(ds[i].requests, e)
			}

			// TODO
			// We also might want to impose a limit on the size of "[]Dep
			// A good strategy is to evict by "replacing".
			return ds
		}
	}
	return append(ds, Dep{ls.Clone(), []Event{e}})
}

// Algorithm phase 2

// Based on lock dependencies we can check for cycles.
// A cycle involves n threads and results from some n lock dependencies.
// For each thread we record the requests that might block.

type Entry struct {
	thread_id ThreadId
	requests  []Event
}

type Cycle []Entry

func report(s *State, c Cycle) {
	s.cycles = append(s.cycles, c)

}

// After phase 1, the following function yields all cyclice lock depndencies.
// TODO kinda unnecessary?
func get_cycles(s *State) []Cycle {
	s.cycles = []Cycle{}
	find_cycles(s)

	return s.cycles
}

type LockDependency struct {
	thread_id ThreadId
	l         LockId
	ls        Lockset
	reqs      []Event
}

// The implementation below follows the algorithm used in UNDEAD (https://github.com/UTSASRG/UnDead/blob/master/analyzer.hh)
func find_cycles(s *State) {
	thread_ids := make(map[ThreadId]bool)
	is_traversed := make(map[ThreadId]bool)
	for tid := range s.threads {
		thread_ids[tid] = true
		is_traversed[tid] = false
	}

	var chain_stack []LockDependency
	for thread_id := range thread_ids {
		if len(s.threads[thread_id].ldeps) == 0 {
			continue
		}
		visiting := thread_id
		is_traversed[thread_id] = true
		for l, ds := range s.threads[thread_id].ldeps {
			for _, d := range ds {
				chain_stack = append(chain_stack, LockDependency{thread_id, l, d.ls, d.requests}) // push
			}
			dfs(s, &chain_stack, visiting, is_traversed, thread_ids)
			chain_stack = chain_stack[:len(chain_stack)-1] // pop
		}

	}

}

// Check if adding dependency to chain will still be a chain.
func isChain(chain_stack *[]LockDependency, dependency *LockDependency) bool {

	for _, d := range *chain_stack {
		// Exit early. No two deps can hold the same lock.
		if d.l == dependency.l {
			return false
		}
		// Check (LD-1) LS(ls_j) cap LS(ls_i+1) for j in {1,..,i}
		if !d.ls.disjoint(dependency.ls) {
			return false
		}
	}
	// Check (LD-2) l_i in ls_i+1
	for l := range dependency.ls {

		if (*chain_stack)[len(*chain_stack)-1].l == l {
			return true
		}

	}

	return false
}

// Check (LD-3) l_n in ls_1
func isCycleChain(chain_stack *[]LockDependency, dependency *LockDependency) bool {
	for l := range (*chain_stack)[0].ls {
		if l == dependency.l {
			return true
		}

	}
	return false
}

func dfs(s *State, chain_stack *[]LockDependency, visiting ThreadId, is_traversed map[ThreadId]bool, thread_ids map[ThreadId]bool) {
	for tid := range thread_ids {
		// if tid <= visiting {
		// 	continue
		// }
		if len(s.threads[tid].ldeps) == 0 {
			continue
		}

		if !is_traversed[tid] {
			for l, l_d := range s.threads[tid].ldeps {
				for _, l_ls_d := range l_d {
					ld := LockDependency{tid, l, l_ls_d.ls, l_ls_d.requests}
					if isChain(chain_stack, &ld) {
						if isCycleChain(chain_stack, &ld) {
							var c Cycle = make([]Entry, len(*chain_stack)+1)
							for i, d := range *chain_stack {
								c[i].requests = d.reqs
								c[i].thread_id = d.thread_id

							}
							c[len(*chain_stack)].requests = ld.reqs
							c[len(*chain_stack)].thread_id = ld.thread_id
							report(s, c)

						} else {
							is_traversed[tid] = true
							*chain_stack = append(*chain_stack, ld) // push
							dfs(s, chain_stack, visiting, is_traversed, thread_ids)
							*chain_stack = (*chain_stack)[:len(*chain_stack)-1] // pop
							is_traversed[tid] = false

						}

					}

				}

			}

		}

	}

}

// ////////////////////////////////
// High level functions for integration with Advocate
func checkForResourceDeadlock() {
	log.Println("Checking for Resource Deadlocks, current state is:" + fmt.Sprint(currentState))

	get_cycles(&currentState)

	for _, cycle := range currentState.cycles {
		// var cycleElements []results.ResultElem
		log.Println("Found cycle:", cycle)
		for i := 0; i < len(cycle); i++ {
			// file, line, tPre, err := infoFromTID(cycle[i].thread_id)
			// if err != nil {
			// 	log.Print(err.Error())
			// 	continue
			// }

			// cycleElements = append(cycleElements, results.TraceElementResult{
			// 	RoutineID: int(cycle[i].thread_id),
			// 	ObjID:     0, //TODO in our current approach we loose the Object Ids?
			// 	TPre:      tPre,
			// 	ObjType:   "DC",
			// 	File:      file,
			// 	Line:      line,
			// })
			log.Println("Cycle element in routine", cycle[i].thread_id, "with requests", cycle[i].requests)
		}

		//		results.Result(results.CRITICAL, results.PCyclicDeadlock, "head", []results.ResultElem{cycleElements[0]}, "tail", cycleElements)

	}
}

func handleMutexEventForRessourceDeadlock(element TraceElementMutex) {
	e := Event{
		thread_id: ThreadId(element.GetRoutine()),
		vc:        element.GetVC(),
	}

	switch element.opM {
	case LockOp:
		acquire(&currentState, LockId(element.GetID()), e)
	case TryLockOp:
	case RLockOp:
	case UnlockOp:
		release(&currentState, LockId(element.GetID()), e)
	case RUnlockOp:
	}
}

/////////////////////////////////
// Auxiliary functions.

// Lockset methods.

func (ls Lockset) empty() bool {
	return len(ls) == 0

}

func (ls Lockset) add(x LockId) {
	ls[x] = struct{}{}
}

func (ls Lockset) remove(x LockId) {
	delete(ls, x)
}

func (ls Lockset) Clone() Lockset {
	clone := make(Lockset, 0)
	for l := range ls {
		clone[l] = ls[l]
	}
	return clone
}

func (ls Lockset) String() string {
	b := strings.Builder{}
	b.WriteString("Lockset{")
	for l := range ls {
		b.WriteString(strconv.Itoa(int(l)))
	}
	b.WriteString("}")
	return b.String()
}

func (ls Lockset) equal(ls2 Lockset) bool {
	if len(ls) != len(ls2) {
		return false
	}

	for l, _ := range ls {
		if _, contains := ls2[l]; !contains {
			return false
		}
	}
	return true
}

func (ls Lockset) disjoint(ls2 Lockset) bool {
	for l, _ := range ls {
		if _, contains := ls2[l]; contains {
			return false
		}
	}
	return true
}

// Further notes.
//
// If possible we would like to use a double-indexed map of the following form.
//
// type Deps map[Lock]map[Lockset][]Event
//
// Unfortunately, this is not possible in Go because keys must be comparable (but slices, maps, ... are not comparable).
// This is not an issue in Haskell or C++ where we can extend the set of comparable types (but providing additional definitions for "==",...)
//
// Hence, we use single-indexed (by Lock) map.
