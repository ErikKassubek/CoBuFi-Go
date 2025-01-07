// Copyrigth (c) 2024 Erik Kassubek
//
// File: analysisMixedDeadlock.go
// Brief: Trace analysis for mixed deadlocks. Currently not used.
//
// Author: Erik Kassubek
// Created: 2024-01-27
//
// License: BSD-3-Clause

package analysis

import (
	"analyzer/clock"
	"log"
	"strconv"
)

/*
 * Add a lock to the lockSet of a routine. Also save the vector clock of the acquire
 * Args:
 *   routine (int): The routine id
 *   lock (int): The id of the mutex
 *   tId (string): The trace id of the mutex operation
 *   vc (VectorClock): The current vector clock
 */
func lockSetAddLock(routine int, lock int, tID string, vc clock.VectorClock) {
	if _, ok := lockSet[routine]; !ok {
		lockSet[routine] = make(map[int]string)
	}
	if _, ok := mostRecentAcquire[routine]; !ok {
		mostRecentAcquire[routine] = make(map[int]VectorClockTID)
	}

	if _, ok := lockSet[routine][lock]; ok {
		// TODO: TODO: add a result. Deadlock detection is currently disabled
		// errorMsg := "Lock " + strconv.Itoa(lock) +
		// 	" already in lockSet for routine " + strconv.Itoa(routine)
		// results.Debug(errorMsg, results.ERROR)

		// // this is a double locking
		// found := "Double locking:\n"
		// found += "\tlock1: " + posOld + "\n"
		// found += "\tlock2: " + tID
		// results.Result(found, results.CRITICAL)
	}

	lockSet[routine][lock] = tID
	mostRecentAcquire[routine][lock] = VectorClockTID{vc, tID, routine}
}

/*
 * Remove a lock from the lockSet of a routine
 * Args:
 *   routine (int): The routine id
 *   lock (int): The id of the mutex
 */
func lockSetRemoveLock(routine int, lock int) {
	if _, ok := lockSet[routine][lock]; !ok {
		errorMsg := "Lock " + strconv.Itoa(lock) +
			" not in lockSet for routine " + strconv.Itoa(routine)
		log.Print(errorMsg)
		// TODO: we can give up on deadlock detection when this happens - at least for now
		return
	}
	delete(lockSet[routine], lock)
}

/*
 * Check for mixed deadlocks
 * Args:
 *   routineSend (int): The routine id of the send operation
 *   routineRevc (int): The routine id of the receive operation
 *   tIDSend (string): The trace id of the channel send
 *   tIDSend (string): The trace id of the channel recv
 */
func checkForMixedDeadlock(routineSend int, routineRevc int, tIDSend string, tIDRecv string) {
	for m := range lockSet[routineSend] {
		_, ok1 := mostRecentAcquire[routineRevc][m]
		_, ok2 := mostRecentAcquire[routineSend][m]
		if ok1 && ok2 && mostRecentAcquire[routineSend][m].TID != mostRecentAcquire[routineRevc][m].TID {
			// found possible mixed deadlock
			// TODO: add a result. Deadlock detection is currently disabled
			// found := "Possible mixed deadlock:\n"
			// found += "\tlocks: \t\t" + mostRecentAcquire[routineSend][m].TID + "\t\t" + mostRecentAcquire[routineRevc][m].TID + "\n"
			// found += "\tsend/close-recv: \t\t" + tIDSend + "\t\t" + tIDRecv

			// results.Result(found, results.CRITICAL)
		}
	}

	for m := range lockSet[routineRevc] {
		_, ok1 := mostRecentAcquire[routineRevc][m]
		_, ok2 := mostRecentAcquire[routineSend][m]
		if ok1 && ok2 && mostRecentAcquire[routineSend][m].TID != mostRecentAcquire[routineRevc][m].TID {
			// found possible mixed deadlock
			// TODO: add a result. Deadlock detection is currently disabled
			// found := "Possible mixed deadlock:\n"
			// found += "\tlocks: \t\t" + mostRecentAcquire[routineSend][m].TID + "\t\t" + mostRecentAcquire[routineRevc][m].TID + "\n"
			// found += "\tsend/close-recv: \t\t" + tIDSend + "\t\t" + tIDRecv

			// results.Result(found, results.CRITICAL)
		}
	}
}

/*
func checkForMixedDeadlock2(routine int) {
	for m := range lockSet[routine] {
		// if the lock was not acquired by the routine, continue. Should not happen
		vc1, okS := mostRecentAcquire[routine][m]
		if !okS {
			continue
		}

		for routine2, acquire := range mostRecentAcquire {
			if routine == routine2 {
				continue
			}

			if vc2, ok := acquire[m]; ok {
				weakHappensBefore := clock.GetHappensBefore(vc1, vc2)
				if weakHappensBefore != Concurrent {
					continue
				}

				// found possible mixed deadlock
				found := "Possible mixed deadlock:\n"
				found += "\tlock1: " + lockSet[routine][m] + "\n"
				found += "\tlock2: " + lockSet[routine2][m]

				results.Result(found, results.CRITICAL)
			}

		}
	}
}
*/
