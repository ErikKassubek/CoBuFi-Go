package vectorClock

// vector clocks for last release times
var relW map[int]VectorClock = make(map[int]VectorClock)
var relR map[int]VectorClock = make(map[int]VectorClock)

/*
 * Create a new relW and relR if needed
 * Args:
 *   index (int): The id of the atomic variable
 *   nRout (int): The number of routines in the trace
 */
func newRel(index int, nRout int) {
	if _, ok := relW[index]; !ok {
		relW[index] = NewVectorClock(nRout)
		relR[index] = NewVectorClock(nRout)
	}
}

/*
 * Update and calculate the vector clocks given a lock operation
 * Args:
 *   vc (vectorClock): The current vector clocks
 * Returns:
 *   (vectorClock): The new vector clock
 */
func Lock(routine int, id int, nRout int, vc *[]VectorClock) VectorClock {
	newRel(id, nRout)
	(*vc)[routine] = (*vc)[routine].Sync(relW[id])
	(*vc)[routine] = (*vc)[routine].Sync(relR[id])
	return (*vc)[routine].Inc(routine)
}

/*
 * Update and calculate the vector clocks given a unlock operation
 * Args:
 *   vc (vectorClock): The current vector clocks
 * Returns:
 *   (vectorClock): The new vector clock
 */
func Unlock(routine int, id int, nRout int, vc *[]VectorClock) VectorClock {
	newRel(id, nRout)
	relW[id] = (*vc)[routine]
	relR[id] = (*vc)[routine]
	return (*vc)[routine].Inc(routine)
}

/*
 * Update and calculate the vector clocks given a rlock operation
 * Args:
 *   vc (vectorClock): The current vector clocks
 * Returns:
 *   (vectorClock): The new vector clock
 */
func RLock(routine int, id int, nRout int, vc *[]VectorClock) VectorClock {
	newRel(id, nRout)
	(*vc)[routine] = (*vc)[routine].Sync(relW[id])
	return (*vc)[routine].Inc(routine)
}

/*
 * Update and calculate the vector clocks given a runlock operation
 * Args:
 *   vc (vectorClock): The current vector clocks
 * Returns:
 *   (vectorClock): The new vector clock
 */
func RUnlock(routine int, id int, nRout int, vc *[]VectorClock) VectorClock {
	newRel(id, nRout)
	relR[id] = (*vc)[routine].Sync(relR[id])
	return (*vc)[routine].Inc(routine)
}
