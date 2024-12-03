// Copyright (c) 2024 Erik Kassubek
//
// File: data.go
// Brief: File to define and contain the fuzzing data
//
// Author: Erik Kassubek
// Created: 2024-11-28
//
// License: BSD-3-Clause

package fuzzing

import "fmt"

type closeInfo string

const (
	always    closeInfo = "a"
	never     closeInfo = "n"
	sometimes closeInfo = "s"
	unknown   closeInfo = "u"
)

var (
	numberOfPreviousRuns = 0
	maxScore             = 0.0
	// Info for the current trace
	channelInfoTrace = make(map[int]fuzzingChannel) // localID -> fuzzingChannel
	pairInfoTrace    = make(map[string]fuzzingPair) // posSend-posRecv -> fuzzing pair  // TODO: do not store paths twice
	numberClose      = 0
	// Info from the file/the previous runs
	channelInfoFile = make([]fuzzingChannel, 0)
	pairInfoFile    = make([]fuzzingPair, 0)
)

/*
 * For each channel that has ever been created, store the
 * following information:
 *   globalId: file:line of creation with new
 *   localId: id in this run
 *   qSize: buffer size of the channel
 *   maxQSize: maximum buffer fullness over all runs
 *   whether the channel has always/never/sometimes been closed
 */
type fuzzingChannel struct {
	globalID  string
	localID   int
	closeInfo closeInfo
	qSize     int
	maxQCount int
}

/*
 * For each pair of channel operations, that have communicated, store the following information:
 *    sendID: file:line:caseSend of the send
 *      caseSend: If the send is in a select, the case ID, otherwise 0
 *    recvID: file:line:Recv of the recv
 *      caseRecv: If the recv is in a select, the case ID, otherwise 0
 *    chanID: local ID of the channel
 *    com: avg number of communication from all the run
 */
type fuzzingPair struct {
	sendID string
	recvID string
	chanID int
	com    float64
}

func (f fuzzingChannel) toString() string {
	return fmt.Sprintf("%s;%s;%d;%d", f.globalID, f.closeInfo, f.qSize, f.maxQCount)
}

func (f fuzzingPair) toString() string {
	return fmt.Sprintf("%s;%s;%f", f.sendID, f.recvID, f.com)
}

func addFuzzingChannel(id string, closeInfo closeInfo, qSize int, maxQSize int) {
	fc := fuzzingChannel{globalID: id, closeInfo: closeInfo, qSize: qSize, maxQCount: maxQSize}
	channelInfoFile = append(channelInfoFile, fc)
}

func addFuzzingPair(idSend string, idRecv string, com float64) {
	fp := fuzzingPair{sendID: idSend, recvID: idRecv, com: com}
	pairInfoFile = append(pairInfoFile, fp)
}
