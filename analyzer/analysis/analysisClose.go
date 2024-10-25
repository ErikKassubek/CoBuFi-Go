// Copyrigth (c) 2024 Erik Kassubek
//
// File: analysisClose.go
// Brief: Trace analysis for send, receive and close on closed channel
//
// Author: Erik Kassubek
// Created: 2024-01-04
//
// License: BSD-3-Clause

package analysis

import (
	"analyzer/clock"
	"analyzer/logging"
	timemeasurement "analyzer/timeMeasurement"
	"strconv"
)

/*
 * Check if a send or receive on a closed channel is possible
 * It it is possible, print a warning or error
 * Args:
 *   ch (*TraceElementChannel): The trace element
 */
func checkForCommunicationOnClosedChannel(ch *TraceElementChannel) {
	// check if there is an earlier send, that could happen concurrently to close
	// println("Check for possible send on closed channel ", analysisCases["sendOnClosed"], hasSend[id])
	if analysisCases["sendOnClosed"] && hasSend[ch.id] {
		timemeasurement.Start("panic")
		defer timemeasurement.End("panic")

		for routine, mrs := range mostRecentSend {
			logging.Debug("Check for possible send on closed channel "+
				strconv.Itoa(ch.id)+" with "+
				mrs[ch.id].Vc.ToString()+" and "+closeData[ch.id].vc.ToString(),
				logging.DEBUG)

			happensBefore := clock.GetHappensBefore(closeData[ch.id].vc, mrs[ch.id].Vc)
			if mrs[ch.id].Elem.GetTID() != "" && happensBefore == clock.Concurrent {

				file1, line1, tPre1, err := infoFromTID(mrs[ch.id].Elem.GetTID()) // send
				if err != nil {
					logging.Debug(err.Error(), logging.ERROR)
					return
				}

				file2, line2, tPre2, err := infoFromTID(ch.tID) // close
				if err != nil {
					logging.Debug(err.Error(), logging.ERROR)
					return
				}

				arg1 := logging.TraceElementResult{ // send
					RoutineID: routine,
					ObjID:     ch.id,
					TPre:      tPre1,
					ObjType:   "CS",
					File:      file1,
					Line:      line1,
				}

				arg2 := logging.TraceElementResult{ // close
					RoutineID: closeData[ch.id].routine,
					ObjID:     ch.id,
					TPre:      tPre2,
					ObjType:   "CC",
					File:      file2,
					Line:      line2,
				}

				logging.Result(logging.CRITICAL, logging.PSendOnClosed,
					"send", []logging.ResultElem{arg1}, "close", []logging.ResultElem{arg2})
			}
		}
	}
	// check if there is an earlier receive, that could happen concurrently to close
	if analysisCases["receiveOnClosed"] && hasReceived[ch.id] {
		timemeasurement.Start("other")
		defer timemeasurement.End("other")

		for routine, mrr := range mostRecentReceive {
			logging.Debug("Check for possible receive on closed channel "+
				strconv.Itoa(ch.id)+" with "+
				mrr[ch.id].Vc.ToString()+" and "+closeData[ch.id].vc.ToString(),
				logging.DEBUG)

			happensBefore := clock.GetHappensBefore(closeData[ch.id].vc, mrr[ch.id].Vc)
			if mrr[ch.id].Elem.GetTID() != "" && (happensBefore == clock.Concurrent || happensBefore == clock.Before) {

				file1, line1, tPre1, err := infoFromTID(mrr[ch.id].Elem.GetTID()) // recv
				if err != nil {
					logging.Debug(err.Error(), logging.ERROR)
					return
				}

				file2, line2, tPre2, err := infoFromTID(ch.tID) // close
				if err != nil {
					logging.Debug(err.Error(), logging.ERROR)
					return
				}

				arg1 := logging.TraceElementResult{ // recv
					RoutineID: routine,
					ObjID:     ch.id,
					TPre:      tPre1,
					ObjType:   "CR",
					File:      file1,
					Line:      line1,
				}

				arg2 := logging.TraceElementResult{ // close
					RoutineID: closeData[ch.id].routine,
					ObjID:     ch.id,
					TPre:      tPre2,
					ObjType:   "CC",
					File:      file2,
					Line:      line2,
				}

				logging.Result(logging.WARNING, logging.PRecvOnClosed,
					"recv", []logging.ResultElem{arg1}, "close", []logging.ResultElem{arg2})
			}
		}

	}

}

/*
 * Lock a found actual send on closed
 * Args:
 *  routineID (int): id of the routine where the send happened
 *  id (int): id of the channel
 *  posSend (string): code location of the send
 */
func foundSendOnClosedChannel(routineID int, id int, posSend string) {
	timemeasurement.Start("panic")
	defer timemeasurement.End("panic")

	if _, ok := closeData[id]; !ok {
		return
	}

	posClose := closeData[id].tID
	if posClose == "" || posSend == "" || posClose == "\n" || posSend == "\n" {
		return
	}

	file1, line1, tPre1, err := infoFromTID(posSend)
	if err != nil {
		logging.Debug(err.Error(), logging.ERROR)
		return
	}

	file2, line2, tPre2, err := infoFromTID(posClose)
	if err != nil {
		logging.Debug(err.Error(), logging.ERROR)
		return
	}

	arg1 := logging.TraceElementResult{ // send
		RoutineID: routineID,
		ObjID:     id,
		TPre:      tPre1,
		ObjType:   "CS",
		File:      file1,
		Line:      line1,
	}

	arg2 := logging.TraceElementResult{ // close
		RoutineID: closeData[id].routine,
		ObjID:     id,
		TPre:      tPre2,
		ObjType:   "CC",
		File:      file2,
		Line:      line2,
	}

	logging.Result(logging.CRITICAL, logging.ASendOnClosed,
		"send", []logging.ResultElem{arg1}, "close", []logging.ResultElem{arg2})

}

/*
 * Log the detection of an actual receive on a closed channel
 * Args:
 *  ch (*TraceElementChannel): The trace element
 */
func foundReceiveOnClosedChannel(ch *TraceElementChannel) {
	timemeasurement.Start("panic")
	defer timemeasurement.End("panic")

	if _, ok := closeData[ch.id]; !ok {
		return
	}

	posClose := closeData[ch.id].tID
	if posClose == "" || ch.tID == "" || posClose == "\n" || ch.tID == "\n" {
		return
	}

	file1, line1, tPre1, err := infoFromTID(ch.tID)
	if err != nil {
		logging.Debug(err.Error(), logging.ERROR)
		return
	}

	file2, line2, tPre2, err := infoFromTID(posClose)
	if err != nil {
		logging.Debug(err.Error(), logging.ERROR)
		return
	}

	arg1 := logging.TraceElementResult{ // recv
		RoutineID: ch.routine,
		ObjID:     ch.id,
		TPre:      tPre1,
		ObjType:   "CR",
		File:      file1,
		Line:      line1,
	}

	arg2 := logging.TraceElementResult{ // close
		RoutineID: closeData[ch.id].routine,
		ObjID:     ch.id,
		TPre:      tPre2,
		ObjType:   "CC",
		File:      file2,
		Line:      line2,
	}

	logging.Result(logging.WARNING, logging.ARecvOnClosed,
		"recv", []logging.ResultElem{arg1}, "close", []logging.ResultElem{arg2})
}

/*
 * Check for a close on a closed channel.
 * Must be called, before the current close operation is added to closePos
 * Args:
 *  ch (*TraceElementChannel): The trace element
 */
func checkForClosedOnClosed(ch *TraceElementChannel) {
	timemeasurement.Start("panic")
	defer timemeasurement.End("panic")

	if oldClose, ok := closeData[ch.id]; ok {
		if oldClose.tID == "" || oldClose.tID == "\n" || ch.tID == "" || ch.tID == "\n" {
			return
		}

		file1, line1, tPre1, err := infoFromTID(oldClose.tID)
		if err != nil {
			logging.Debug(err.Error(), logging.ERROR)
			return
		}

		file2, line2, tPre2, err := infoFromTID(oldClose.tID)
		if err != nil {
			logging.Debug(err.Error(), logging.ERROR)
			return
		}

		arg1 := logging.TraceElementResult{
			RoutineID: ch.routine,
			ObjID:     ch.id,
			TPre:      tPre1,
			ObjType:   "CC",
			File:      file1,
			Line:      line1,
		}

		arg2 := logging.TraceElementResult{
			RoutineID: oldClose.routine,
			ObjID:     ch.id,
			TPre:      tPre2,
			ObjType:   "CC",
			File:      file2,
			Line:      line2,
		}

		logging.Result(logging.CRITICAL, logging.ACloseOnClosed,
			"close", []logging.ResultElem{arg1}, "close", []logging.ResultElem{arg2})
	}
}
