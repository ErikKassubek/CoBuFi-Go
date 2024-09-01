// Copyrigth (c) 2024 Erik Kassubek
//
// File: vcOnce.go
// Brief: Update functions of vector clocks for once operations 
// 
// Author: Erik Kassubek <kassubek.erik@gmail.com>
// Created: 2023-07-25
// LastChange: 2024-09-01
//
// License: BSD-3-Clause

package analysis

import "analyzer/clock"

// vector clocks for the successful do
var oSuc map[int]clock.VectorClock = make(map[int]clock.VectorClock)

/*
 * Create a new oSuc if needed
 * Args:
 *   index (int): The id of the atomic variable
 *   nRout (int): The number of routines in the trace
 */
func newOSuc(index int, nRout int) {
	if _, ok := oSuc[index]; !ok {
		oSuc[index] = clock.NewVectorClock(nRout)
	}
}

/*
 * Update and calculate the vector clocks given a successful do operation
 * Args:
 *   routine (int): The routine id
 *   id (int): The id of the atomic variable
 *   vc (map[int]VectorClock): The current vector clocks
 */
func DoSuc(routine int, id int, vc map[int]clock.VectorClock) {
	newOSuc(id, vc[id].GetSize())
	oSuc[id] = vc[routine]
	vc[routine] = vc[routine].Inc(routine)
}

/*
 * Update and calculate the vector clocks given a unsuccessful do operation
 * Args:
 *   routine (int): The routine id
 *   id (int): The id of the atomic variable
 *   vc (map[int]VectorClock): The current vector clocks
 */
func DoFail(routine int, id int, vc map[int]clock.VectorClock) {
	newOSuc(id, vc[id].GetSize())
	vc[routine] = vc[routine].Sync(oSuc[id])
	vc[routine] = vc[routine].Inc(routine)
}
