// Copyrigth (c) 2024 Erik Kassubek
//
// File: vcWait.go
// Brief: Update functions of vector groups for wait group operations
//        Some function start analysis functions
//
// Author: Erik Kassubek <kassubek.erik@gmail.com>
// Created: 2023-07-25
// LastChange: 2024-09-01
//
// License: BSD-3-Clause

package analysis

import "analyzer/clock"

// vector clock for each wait group
var wg map[int]clock.VectorClock = make(map[int]clock.VectorClock)

/*
 * Create a new wg if needed
 * Args:
 *   index (int): The id of the wait group
 *   nRout (int): The number of routines in the trace
 */
func newWg(index int, nRout int) {
	if _, ok := wg[index]; !ok {
		wg[index] = clock.NewVectorClock(nRout)
	}
}

/*
 * Calculate the new vector clock for a add or done operation and update cv
 * Args:
 *   wa (*TraceElementWait): The trace element
 *   vc (map[int]VectorClock): The vector clocks
 */
func Change(wa *TraceElementWait, vc map[int]clock.VectorClock) {
	newWg(wa.id, vc[wa.id].GetSize())
	wg[wa.id] = wg[wa.id].Sync(vc[wa.routine])
	vc[wa.routine] = vc[wa.routine].Inc(wa.routine)

	if analysisCases["doneBeforeAdd"] {
		checkForDoneBeforeAddChange(wa)
	}
}

/*
 * Calculate the new vector clock for a wait operation and update cv
 * Args:
 *   wa (*TraceElementWait): The trace element
 *   vc (*map[int]VectorClock): The vector clocks
 */
func Wait(wa *TraceElementWait, vc map[int]clock.VectorClock) {
	newWg(wa.id, vc[wa.id].GetSize())
	if wa.tPost != 0 {
		vc[wa.routine] = vc[wa.routine].Sync(wg[wa.id])
		vc[wa.routine] = vc[wa.routine].Inc(wa.routine)
	}
}
