// Copyrigth (c) 2024 Erik Kassubek
//
// File: vcCond.go
// Brief: Update functions for vector clocks from conditional variables operations
//
// Author: Erik Kassubek <kassubek.erik@gmail.com>
// Created: 2024-01-09
// LastChange: 2024-09-01
//
// License: BSD-3-Clause

package analysis

import "analyzer/clock"

var lastCondRelease = make(map[int]int) // -> id -> routine

/*
 * Update and calculate the vector clocks given a wait operation
 * Args:
 *   co (*TraceElementCond): The trace element
 *   vc (map[int]VectorClock): The current vector clocks
 */
func CondWait(co *TraceElementCond, vc map[int]clock.VectorClock) {
	if co.tPost != 0 { // not leak
		vc[co.routine].Sync(vc[lastCondRelease[co.id]])
	}
	vc[co.routine].Inc(co.routine)
}

/*
 * Update and calculate the vector clocks given a signal operation
 * Args:
 *   co (*TraceElementCond): The trace element
 *   vc (map[int]VectorClock): The current vector clocks
 */
func CondSignal(co *TraceElementCond, vc map[int]clock.VectorClock) {
	vc[co.routine].Inc(co.routine)

	lastCondRelease[co.id] = co.routine
}

/*
 * Update and calculate the vector clocks given a broadcast operation
 * Args:
 *   co (*TraceElementCond): The trace element
 *   vc (map[int]VectorClock): The current vector clocks
 */
func CondBroadcast(co *TraceElementCond, vc map[int]clock.VectorClock) {
	vc[co.routine].Inc(co.routine)
	lastCondRelease[co.id] = co.routine
}
