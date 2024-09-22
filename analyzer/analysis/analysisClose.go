// Copyrigth (c) 2024 Erik Kassubek
//
// File: analysisClose.go
// Brief: Trace analysis for send, receive and close on closed channel
//
// Author: Erik Kassubek <kassubek.erik@gmail.com>
// Created: 2024-01-04
// LastChange: 2024-08-03
//
// License: BSD-3-Clause

package analysis

import (
	"analyzer/clock"
	"analyzer/logging"
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
		for routine, mrs := range mostRecentSend {
			logging.Debug("Check for possible send on closed channel "+
				strconv.Itoa(ch.id)+" with "+
				mrs[ch.id].Vc.ToString()+" and "+closeData[ch.id].Vc.ToString(),
				logging.DEBUG)

			happensBefore := clock.GetHappensBefore(closeData[ch.id].Vc, mrs[ch.id].Vc)
			if mrs[ch.id].TID != "" && happensBefore == clock.Concurrent {

				file1, line1, tPre1, err := infoFromTID(mrs[ch.id].TID) // send
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
					RoutineID: closeData[ch.id].Routine,
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
		for routine, mrr := range mostRecentReceive {
			logging.Debug("Check for possible receive on closed channel "+
				strconv.Itoa(ch.id)+" with "+
				mrr[ch.id].Vc.ToString()+" and "+closeData[ch.id].Vc.ToString(),
				logging.DEBUG)

			happensBefore := clock.GetHappensBefore(closeData[ch.id].Vc, mrr[ch.id].Vc)
			if mrr[ch.id].TID != "" && (happensBefore == clock.Concurrent || happensBefore == clock.Before) {

				file1, line1, tPre1, err := infoFromTID(mrr[ch.id].TID) // recv
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
					RoutineID: closeData[ch.id].Routine,
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
	if _, ok := closeData[id]; !ok {
		return
	}

	posClose := closeData[id].TID
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
		RoutineID: closeData[id].Routine,
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
	if _, ok := closeData[ch.id]; !ok {
		return
	}

	posClose := closeData[ch.id].TID
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
		RoutineID: closeData[ch.id].Routine,
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
	if oldClose, ok := closeData[ch.id]; ok {
		if oldClose.TID == "" || oldClose.TID == "\n" || ch.tID == "" || ch.tID == "\n" {
			return
		}

		file1, line1, tPre1, err := infoFromTID(oldClose.TID)
		if err != nil {
			logging.Debug(err.Error(), logging.ERROR)
			return
		}

		file2, line2, tPre2, err := infoFromTID(oldClose.TID)
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
			RoutineID: oldClose.Routine,
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
