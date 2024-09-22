// Copyrigth (c) 2024 Erik Kassubek
//
// File: vcMutex.go
// Brief: Update functions for vector clocks from mutex operation
//        Some of the functions start analysis functions
//
// Author: Erik Kassubek <kassubek.erik@gmail.com>
// Created: 2023-07-25
// LastChange: 2024-09-01
//
// License: BSD-3-Clause

package analysis

import "analyzer/clock"

/*
 * Create a new relW and relR if needed
 * Args:
 *   index (int): The id of the atomic variable
 *   nRout (int): The number of routines in the trace
 */
func newRel(index int, nRout int) {
	if _, ok := relW[index]; !ok {
		relW[index] = clock.NewVectorClock(nRout)
	}
	if _, ok := relR[index]; !ok {
		relR[index] = clock.NewVectorClock(nRout)
	}
}

/*
 * Update and calculate the vector clocks given a lock operation
 * Args:
 *   mu (*TraceElementMutex): The trace element
 *   vc (map[int]VectorClock): The current vector clocks
 *   wVc (map[int]VectorClock): The current weak vector clocks
 */
func Lock(mu *TraceElementMutex, vc map[int]clock.VectorClock, wVc map[int]clock.VectorClock) {
	if mu.tPost == 0 {
		vc[mu.routine] = vc[mu.routine].Inc(mu.routine)
		return
	}

	newRel(mu.id, vc[mu.routine].GetSize())
	vc[mu.routine] = vc[mu.routine].Sync(relW[mu.id])
	vc[mu.routine] = vc[mu.routine].Sync(relR[mu.id])
	vc[mu.routine] = vc[mu.routine].Inc(mu.routine)

	if analysisCases["leak"] {
		addMostRecentAcquireTotal(mu.routine, mu.id, mu.tID, vc[mu.routine], 0)
	}

	if analysisCases["mixedDeadlock"] {
		lockSetAddLock(mu.routine, mu.id, mu.tID, wVc[mu.routine])
	}
}

/*
 * Update and calculate the vector clocks given a unlock operation
 * Args:
 *   mu (*TraceElementMutex): The trace element
 *   vc (map[int]VectorClock): The current vector clocks
 */
func Unlock(mu *TraceElementMutex, vc map[int]clock.VectorClock) {
	if mu.tPost == 0 {
		return
	}

	newRel(mu.id, vc[mu.routine].GetSize())
	relW[mu.id] = vc[mu.routine].Copy()
	relR[mu.id] = vc[mu.routine].Copy()
	vc[mu.routine] = vc[mu.routine].Inc(mu.routine)

	if analysisCases["mixedDeadlock"] {
		lockSetRemoveLock(mu.routine, mu.id)
	}
}

/*
 * Update and calculate the vector clocks given a rlock operation
 * Args:
 *   mu (*TraceElementMutex): The trace element
 *   vc (map[int]VectorClock): The current vector clocks
 *   wVc (map[int]VectorClock): The current weak vector clocks
 * Returns:
 *   (vectorClock): The new vector clock
 */
func RLock(mu *TraceElementMutex, vc map[int]clock.VectorClock, wVc map[int]clock.VectorClock) {
	if mu.tPost == 0 {
		vc[mu.routine] = vc[mu.routine].Inc(mu.routine)
		return
	}

	newRel(mu.id, vc[mu.routine].GetSize())
	vc[mu.routine] = vc[mu.routine].Sync(relW[mu.id])
	vc[mu.routine] = vc[mu.routine].Inc(mu.routine)

	if analysisCases["leak"] {
		addMostRecentAcquireTotal(mu.routine, mu.id, mu.tID, vc[mu.routine], 1)
	}

	if analysisCases["mixedDeadlock"] {
		lockSetAddLock(mu.routine, mu.id, mu.tID, wVc[mu.routine])
	}
}

/*
 * Update and calculate the vector clocks given a runlock operation
 * Args:
 *   mu (*TraceElementMutex): The trace element
 *   vc (map[int]VectorClock): The current vector clocks
 */
func RUnlock(mu *TraceElementMutex, vc map[int]clock.VectorClock) {
	if mu.tPost == 0 {
		vc[mu.routine] = vc[mu.routine].Inc(mu.routine)
		return
	}

	newRel(mu.id, vc[mu.routine].GetSize())
	relR[mu.id] = relR[mu.id].Sync(vc[mu.routine])
	vc[mu.routine] = vc[mu.routine].Inc(mu.routine)

	if analysisCases["mixedDeadlock"] {
		lockSetRemoveLock(mu.routine, mu.id)
	}
}
