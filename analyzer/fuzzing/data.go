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
)

var (
	numberOfPreviousRuns = 0
	// Info for the current trace
	channelInfoTrace = make([]fuzzingChannel, 0)
	pairInfoTrace    = make([]fuzzingPair, 0)
	// Info from the file/the previous runs
	channelInfoFile = make([]fuzzingChannel, 0)
	pairInfoFile    = make([]fuzzingPair, 0)
)

/*
 * For each channel that has ever been created, store the
 * following information:
 *   id: file:line of creation with new
 *   qSize: buffer size of the channel
 *   maxQSize: maximum buffer fullness over all runs
 *   whether the channel has always/never/sometimes been closed
 */
type fuzzingChannel struct {
	id        string
	closeInfo closeInfo
	qSize     int
	maxQsize  int
}

/*
 * For each pair of channel operations, that have communicated, store the following information:
 *    idSend: file:line:caseSend of the send
 *      caseSend: If the send is in a select, the case ID, otherwise 0
 *    idRecv: file:line:Recv of the recv
 *      caseRecv: If the recv is in a select, the case ID, otherwise 0
 *    com: avg number of communication from all the run
 */
type fuzzingPair struct {
	idSend string
	idRecv string
	com    float64
}

func (f fuzzingChannel) toString() string {
	return fmt.Sprintf("%s;%s;%d;%d", f.id, f.closeInfo, f.qSize, f.maxQsize)
}

func (f fuzzingPair) toString() string {
	return fmt.Sprintf("%s;%s;%f", f.idSend, f.idRecv, f.com)
}

func addFuzzingChannel(id string, closeInfo closeInfo, qSize int, maxQSize int) {
	fc := fuzzingChannel{id: id, closeInfo: closeInfo, qSize: qSize, maxQsize: maxQSize}
	channelInfoFile = append(channelInfoFile, fc)
}

func addFuzzingPair(idSend string, idRecv string, com float64) {
	fp := fuzzingPair{idSend: idSend, idRecv: idRecv, com: com}
	pairInfoFile = append(pairInfoFile, fp)
}
