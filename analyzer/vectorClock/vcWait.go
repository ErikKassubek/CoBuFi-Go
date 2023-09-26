package vectorClock

// vector clock for each wait group
var wg map[int]VectorClock = make(map[int]VectorClock)

/*
 * Create a new wg if needed
 * Args:
 *   index (int): The id of the wait group
 *   nRout (int): The number of routines in the trace
 */
func newWg(index int, nRout int) {
	if _, ok := wg[index]; !ok {
		wg[index] = NewVectorClock(nRout)
	}
}

/*
 * Calculate the new vector clock for a add or done operation and update cv
 * Args:
 *   routine (int): The routine id
 *   id (int): The id of the wait group
 *   numberOfRoutines (int): The number of routines in the trace
 *   vc (map[int]VectorClock): The vector clocks
 * Returns:
 *   (VectorClock): The vector clock for the add or done
 */
func Change(routine int, id int, vc map[int]VectorClock) VectorClock {
	newWg(id, vc[id].size)
	wg[id].Sync(vc[routine])
	return vc[routine].Inc(routine).Copy()
}

/*
 * Calculate the new vector clock for a wait operation and update cv
 * Args:
 *   routine (int): The routine id
 *   id (int): The id of the wait group
 *   numberOfRoutines (int): The number of routines in the trace
 *   vc (*map[int]VectorClock): The vector clocks
 * Returns:
 *   (VectorClock): The vector clock for the wait
 */
func Wait(routine int, id int, vc map[int]VectorClock) VectorClock {
	newWg(id, vc[id].size)
	vc[routine].Sync(wg[id])
	return vc[routine].Inc(routine).Copy()
}
