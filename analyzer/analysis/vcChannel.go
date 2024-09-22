// Copyrigth (c) 2024 Erik Kassubek
//
// File: vcChannel.go
// Brief: Update functions for vector clocks from channel operations
//        Some of the update function also start analysis functions
//
// Author: Erik Kassubek <kassubek.erik@gmail.com>
// Created: 2023-07-27
// LastChange: 2024-09-01
//
// License: BSD-3-Clause

package analysis

import (
	"analyzer/clock"
	"analyzer/logging"
	"strconv"
)

// elements for buffered channel internal vector clock
type bufferedVC struct {
	occupied    bool
	oID         int
	vc          clock.VectorClock
	routineSend int
	tID         string
}

/*
 * Update and calculate the vector clocks given a send/receive pair on a unbuffered
 * channel.
 * Args:
 * 	ch (*TraceElementChannel): The trace element
 * 	routSend (int): the route of the sender
 * 	routRecv (int): the route of the receiver
 * 	tID_send (string): the position of the send in the program
 * 	tID_recv (string): the position of the receive in the program
 * 	vc (map[int]VectorClock): the current vector clocks
 */
func Unbuffered(ch *TraceElementChannel, routSend int, routRecv int, tIDSend string,
	tIDRecv string, vc map[int]clock.VectorClock) {
	if analysisCases["concurrentRecv"] {
		checkForConcurrentRecv(routRecv, ch.id, tIDRecv, vc, ch.tPost)
	}

	if ch.tPost != 0 {

		if mostRecentReceive[routRecv] == nil {
			mostRecentReceive[routRecv] = make(map[int]VectorClockTID3)
		}
		if mostRecentSend[routSend] == nil {
			mostRecentSend[routSend] = make(map[int]VectorClockTID3)
		}

		vc[routRecv] = vc[routRecv].Sync(vc[routSend])
		vc[routSend] = vc[routRecv].Copy()
		vc[routSend] = vc[routSend].Inc(routSend)
		vc[routRecv] = vc[routRecv].Inc(routRecv)

		// for detection of send on closed
		hasSend[ch.id] = true
		mostRecentSend[routSend][ch.id] = VectorClockTID3{routSend, tIDSend, mostRecentSend[routSend][ch.id].Vc.Sync(vc[routSend]).Copy(), ch.id}

		// for detection of receive on closed
		hasReceived[ch.id] = true
		mostRecentReceive[routRecv][ch.id] = VectorClockTID3{routRecv, tIDRecv, mostRecentReceive[routRecv][ch.id].Vc.Sync(vc[routRecv]).Copy(), ch.id}

		logging.Debug("Set most recent send of "+strconv.Itoa(ch.id)+" to "+mostRecentSend[routSend][ch.id].Vc.ToString(), logging.DEBUG)
		logging.Debug("Set most recent recv of "+strconv.Itoa(ch.id)+" to "+mostRecentReceive[routRecv][ch.id].Vc.ToString(), logging.DEBUG)

	} else {
		vc[routSend] = vc[routSend].Inc(routSend)
	}

	if analysisCases["sendOnClosed"] {
		if _, ok := closeData[ch.id]; ok {
			foundSendOnClosedChannel(routSend, ch.id, tIDSend)
		}
	}

	if analysisCases["mixedDeadlock"] {
		CheckForSelectCaseWithoutPartnerChannel(ch.id, vc[routSend], tIDSend, true, false)
		CheckForSelectCaseWithoutPartnerChannel(ch.id, vc[routRecv], tIDRecv, false, false)
		checkForMixedDeadlock(routSend, routRecv, tIDSend, tIDRecv)
	}

	if analysisCases["selectWithoutPartner"] {
		CheckForSelectCaseWithoutPartnerChannel(ch.id, vc[routSend], tIDSend, true, false)
		CheckForSelectCaseWithoutPartnerChannel(ch.id, vc[routRecv], tIDRecv, false, false)
	}

	if analysisCases["leak"] {
		CheckForLeakChannelRun(routSend, ch.id, VectorClockTID{vc[routSend].Copy(), tIDSend, routSend}, 0, false)
		CheckForLeakChannelRun(routRecv, ch.id, VectorClockTID{vc[routRecv].Copy(), tIDRecv, routRecv}, 1, false)
	}

}

type holdObj struct {
	ch   *TraceElementChannel
	vc   map[int]clock.VectorClock
	fifo bool
}

var holdSend = make([]holdObj, 0)
var holdRecv = make([]holdObj, 0)

/*
 * Update and calculate the vector clocks given a send on a buffered channel.
 * Args:
 * 	ch (*TraceElementChannel): The trace element
 * 	vc (map[int]VectorClock): the current vector clocks
 *  fifo (bool): true if the channel buffer is assumed to be fifo
 */
func Send(ch *TraceElementChannel, vc map[int]clock.VectorClock, fifo bool) {

	if ch.tPost == 0 {
		vc[ch.routine] = vc[ch.routine].Inc(ch.routine)
		return
	}

	if mostRecentSend[ch.routine] == nil {
		mostRecentSend[ch.routine] = make(map[int]VectorClockTID3)
	}

	newBufferedVCs(ch.id, ch.qSize, vc[ch.routine].GetSize())

	count := bufferedVCsCount[ch.id]

	if len(bufferedVCs[ch.id]) <= count {
		holdSend = append(holdSend, holdObj{ch, vc, fifo})
		return
		// panic("BufferedVCsCount is bigger than the buffer qSize for chan " + strconv.Itoa(id) + " with count " + strconv.Itoa(count) + " and qSize " + strconv.Itoa(qSize) + "\n\tand tID " + tID)
	}

	if count > ch.qSize || bufferedVCs[ch.id][count].occupied {
		logging.Debug("Write to occupied buffer position or to big count", logging.ERROR)
	}

	v := bufferedVCs[ch.id][count].vc
	vc[ch.routine] = vc[ch.routine].Sync(v)

	if fifo {
		vc[ch.routine] = vc[ch.routine].Sync(mostRecentSend[ch.routine][ch.id].Vc)
	}

	bufferedVCs[ch.id][count] = bufferedVC{true, ch.oID, vc[ch.routine].Copy(), ch.routine, ch.tID}

	bufferedVCsCount[ch.id]++

	// for detection of send on closed
	hasSend[ch.id] = true
	mostRecentSend[ch.routine][ch.id] = VectorClockTID3{ch.routine, ch.tID, mostRecentSend[ch.routine][ch.id].Vc.Sync(vc[ch.routine]), ch.id}

	vc[ch.routine] = vc[ch.routine].Inc(ch.routine)

	if analysisCases["sendOnClosed"] {
		if _, ok := closeData[ch.id]; ok {
			foundSendOnClosedChannel(ch.routine, ch.id, ch.tID)
		}
	}

	if analysisCases["selectWithoutPartner"] {
		CheckForSelectCaseWithoutPartnerChannel(ch.id, vc[ch.routine], ch.tID, true, true)
	}

	if analysisCases["leak"] {
		CheckForLeakChannelRun(ch.routine, ch.id, VectorClockTID{vc[ch.routine].Copy(), ch.tID, ch.routine}, 0, true)
	}

	for i, hold := range holdRecv {
		if hold.ch.id == ch.id {
			Recv(hold.ch, hold.vc, hold.fifo)
			holdRecv = append(holdRecv[:i], holdRecv[i+1:]...)
			break
		}
	}

}

/*
 * Update and calculate the vector clocks given a receive on a buffered channel.
 * Args:
 * 	ch (*TraceElementChannel): The trace element
 * 	vc (map[int]VectorClock): the current vector clocks
 *  fifo (bool): true if the channel buffer is assumed to be fifo
 */
func Recv(ch *TraceElementChannel, vc map[int]clock.VectorClock, fifo bool) {

	if analysisCases["concurrentRecv"] {
		checkForConcurrentRecv(ch.routine, ch.id, ch.tID, vc, ch.tPost)
	}

	if ch.tPost == 0 {
		vc[ch.routine] = vc[ch.routine].Inc(ch.routine)
		return
	}

	if mostRecentReceive[ch.routine] == nil {
		mostRecentReceive[ch.routine] = make(map[int]VectorClockTID3)
	}

	newBufferedVCs(ch.id, ch.qSize, vc[ch.routine].GetSize())

	if bufferedVCsCount[ch.id] == 0 {
		holdSend = append(holdSend, holdObj{ch, vc, fifo})
		return
		// logging.Debug("Read operation on empty buffer position", logging.ERROR)
	}
	bufferedVCsCount[ch.id]--

	if bufferedVCs[ch.id][0].oID != ch.oID {
		found := false
		for i := 1; i < ch.qSize; i++ {
			if bufferedVCs[ch.id][i].oID == ch.oID {
				found = true
				bufferedVCs[ch.id][0] = bufferedVCs[ch.id][i]
				bufferedVCs[ch.id][i] = bufferedVC{false, 0, vc[ch.routine].Copy(), 0, ""}
				break
			}
		}
		if !found {
			err := "Read operation on wrong buffer position - ID: " + strconv.Itoa(ch.id) + ", OID: " + strconv.Itoa(ch.oID) + ", SIZE: " + strconv.Itoa(ch.qSize)
			logging.Debug(err, logging.INFO)
		}
	}
	v := bufferedVCs[ch.id][0].vc
	routSend := bufferedVCs[ch.id][0].routineSend
	tIDSend := bufferedVCs[ch.id][0].tID

	vc[ch.routine] = vc[ch.routine].Sync(v)

	if fifo {
		vc[ch.routine] = vc[ch.routine].Sync(mostRecentReceive[ch.routine][ch.id].Vc)
	}

	bufferedVCs[ch.id] = bufferedVCs[ch.id][1:]
	bufferedVCs[ch.id] = append(bufferedVCs[ch.id], bufferedVC{false, 0, vc[ch.routine].Copy(), 0, ""})

	// for detection of receive on closed
	hasReceived[ch.id] = true
	mostRecentReceive[ch.routine][ch.id] = VectorClockTID3{ch.routine, ch.tID, mostRecentReceive[ch.routine][ch.id].Vc.Sync(vc[ch.routine]), ch.id}

	vc[ch.routine] = vc[ch.routine].Inc(ch.routine)

	if analysisCases["selectWithoutPartner"] {
		CheckForSelectCaseWithoutPartnerChannel(ch.id, vc[ch.routine], ch.tID, true, true)
	}

	if analysisCases["mixedDeadlock"] {
		checkForMixedDeadlock(routSend, ch.routine, tIDSend, ch.tID)
	}
	if analysisCases["leak"] {
		CheckForLeakChannelRun(ch.routine, ch.id, VectorClockTID{vc[ch.routine].Copy(), ch.tID, ch.routine}, 1, true)
	}

	for i, hold := range holdSend {
		if hold.ch.id == ch.id {
			Send(hold.ch, hold.vc, hold.fifo)
			holdSend = append(holdSend[:i], holdSend[i+1:]...)
			break
		}
	}
}

/*
 * Update and calculate the vector clocks for a stuck channel element
 * Args:
 *  routint (int): the route of the operation
 *  vc (map[int]VectorClock): the current vector clocks
 */
func StuckChan(routine int, vc map[int]clock.VectorClock) {
	vc[routine] = vc[routine].Inc(routine)
}

/*
 * Update and calculate the vector clocks given a close on a channel.
 * Args:
 * 	ch (*TraceElementChannel): The trace element
 * 	vc (map[int]VectorClock): the current vector clocks
 */
func Close(ch *TraceElementChannel, vc map[int]clock.VectorClock) {
	if ch.tPost == 0 {
		return
	}

	if analysisCases["closeOnClosed"] {
		checkForClosedOnClosed(ch.routine, ch.id, ch.tID) // must be called before closePos is updated
	}

	vc[ch.routine] = vc[ch.routine].Inc(ch.routine)

	closeData[ch.id] = VectorClockTID3{Routine: ch.routine, TID: ch.tID, Vc: vc[ch.routine].Copy(), Val: ch.id}

	if analysisCases["sendOnClosed"] || analysisCases["receiveOnClosed"] {
		checkForCommunicationOnClosedChannel(ch.id, ch.tID)
	}

	if analysisCases["selectWithoutPartner"] {
		CheckForSelectCaseWithoutPartnerClose(ch.id, vc[ch.routine])
	}

	if analysisCases["leak"] {
		CheckForLeakChannelRun(ch.routine, ch.id, VectorClockTID{vc[ch.routine].Copy(), ch.tID, ch.routine}, 2, true)
	}
}

func SendC(ch *TraceElementChannel) {
	if analysisCases["sendOnClosed"] {
		foundSendOnClosedChannel(ch.routine, ch.id, ch.tID)
	}
}

/*
 * Update and calculate the vector clocks given a receive on a closed channel.
 * Args:
 * 	ch (*TraceElementChannel): The trace element
 * 	vc (map[int]VectorClock): the current vector clocks
 *  buffered (bool): true if the channel is buffered
 */
func RecvC(ch *TraceElementChannel, vc map[int]clock.VectorClock, buffered bool) {
	if ch.tPost == 0 {
		return
	}

	if analysisCases["receiveOnClosed"] {
		foundReceiveOnClosedChannel(ch.routine, ch.id, ch.tID)
	}

	vc[ch.routine] = vc[ch.routine].Sync(closeData[ch.id].Vc)
	vc[ch.routine] = vc[ch.routine].Inc(ch.routine)

	if analysisCases["selectWithoutPartner"] {
		CheckForSelectCaseWithoutPartnerChannel(ch.id, vc[ch.routine], ch.tID, false, buffered)
	}

	if analysisCases["mixedDeadlock"] {
		checkForMixedDeadlock(closeData[ch.id].Routine, ch.routine, closeData[ch.id].TID, ch.tID)
	}
	if analysisCases["leak"] {
		CheckForLeakChannelRun(ch.routine, ch.id, VectorClockTID{vc[ch.routine].Copy(), ch.tID, ch.routine}, 1, buffered)
	}
}

/*
 * Create a new map of buffered vector clocks for a channel if not already in
 * bufferedVCs.
 * Args:
 * 	id (int): the id of the channel
 * 	qSize (int): the buffer qSize of the channel
 * 	numRout (int): the number of routines
 */
func newBufferedVCs(id int, qSize int, numRout int) {
	if _, ok := bufferedVCs[id]; !ok {
		bufferedVCs[id] = make([]bufferedVC, qSize)
		for i := 0; i < qSize; i++ {
			bufferedVCsCount[id] = 0
			bufferedVCs[id][i] = bufferedVC{false, 0, clock.NewVectorClock(numRout), 0, ""}
		}
	}
}

/*
 * Set the channel as the last send operation.
 * Used for not executed select send
 * Args:
 * 	id (int): the id of the channel
 * 	routine (int): the route of the operation
 *  vc (VectorClock): the vector clock of the operation
 *  tID (string): the position of the send in the program
 */
func SetChannelAsLastSend(id int, routine int, vc clock.VectorClock, tID string) {
	if mostRecentSend[routine] == nil {
		mostRecentSend[routine] = make(map[int]VectorClockTID3)
	}
	mostRecentSend[routine][id] = VectorClockTID3{routine, tID, vc, id}
	hasSend[id] = true
}

/*
 * Set the channel as the last recv operation.
 * Used for not executed select recv
 * Args:
 * 	id (int): the id of the channel
 * 	rout (int): the route of the operation
 *  vc (VectorClock): the vector clock of the operation
 *  tID (string): the position of the recv in the program
 */
func SetChannelAsLastReceive(id int, routine int, vc clock.VectorClock, tID string) {
	if mostRecentReceive[routine] == nil {
		mostRecentReceive[routine] = make(map[int]VectorClockTID3)
	}
	mostRecentReceive[routine][id] = VectorClockTID3{routine, tID, vc, id}
	hasReceived[id] = true
}
