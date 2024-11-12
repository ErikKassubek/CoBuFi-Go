// Copyrigth (c) 2024 Erik Kassubek
//
// File: analysisConcurrentCommunication.go
// Brief: Trace analysis of concurrent reveice on the same channel
//
// Author: Erik Kassubek
// Created: 2024-01-27
//
// License: BSD-3-Clause

package analysis

import (
	"analyzer/clock"
	"analyzer/logging"
	timemeasurement "analyzer/timeMeasurement"
)

/*
 * Check if there are multiple concurrent receive operations on the same channel.
 * Such concurrent recv can lead to nondeterministic behaviour.
 * If such a situation is detected, it is logged.
 * Call this function on a recv.
 * Args:
 *  ch (*TraceElementChannel): The trace element
 *  routine (int): routine of the recv
 *  tID (int): tID of the recv operation
 *  vc (int): vector clock of the recv operation
 */
func checkForConcurrentRecv(ch *TraceElementChannel, vc map[int]clock.VectorClock) {
	timemeasurement.Start("other")
	defer timemeasurement.End("other")

	for r, elem := range lastRecvRoutine {
		if r == ch.routine {
			continue
		}

		if elem[ch.id].Vc.GetClock() == nil {
			continue
		}

		happensBefore := clock.GetHappensBefore(elem[ch.id].Vc, vc[ch.routine])
		if happensBefore == clock.Concurrent {

			file1, line1, tPre1, err := infoFromTID(ch.GetTID())
			if err != nil {
				logging.Debug(err.Error(), logging.ERROR)
				return
			}

			file2, line2, tPre2, err := infoFromTID(lastRecvRoutine[r][ch.id].TID)

			arg1 := logging.TraceElementResult{
				RoutineID: ch.routine,
				ObjID:     ch.id,
				TPre:      tPre1,
				ObjType:   "CR",
				File:      file1,
				Line:      line1,
			}

			arg2 := logging.TraceElementResult{
				RoutineID: r,
				ObjID:     ch.id,
				TPre:      tPre2,
				ObjType:   "CR",
				File:      file2,
				Line:      line2,
			}

			logging.Result(logging.WARNING, logging.AConcurrentRecv,
				"recv", []logging.ResultElem{arg1}, "recv", []logging.ResultElem{arg2})
		}
	}

	if ch.tPost != 0 {
		if _, ok := lastRecvRoutine[ch.routine]; !ok {
			lastRecvRoutine[ch.routine] = make(map[int]VectorClockTID)
		}

		lastRecvRoutine[ch.routine][ch.id] = VectorClockTID{vc[ch.routine].Copy(), ch.GetTID(), ch.routine}
	}
}
