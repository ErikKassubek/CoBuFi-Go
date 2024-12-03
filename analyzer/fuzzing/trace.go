// Copyright (c) 2024 Erik Kassubek
//
// File: trace.go
// Brief: Function to parse the trace and get all relevant information
//
// Author: Erik Kassubek
// Created: 2024-11-29
//
// License: BSD-3-Clause

package fuzzing

import "analyzer/analysis"

/*
 * Parse the current trace and record all relevant data
 */
func parseTrace() {
	for _, trace := range *analysis.GetTraces() {
		for _, elem := range trace {
			switch e := elem.(type) {
			case *analysis.TraceElementNew:
				parseNew(e)
			case *analysis.TraceElementChannel:
				parseChannelOp(e)
			case *analysis.TraceElementSelect:
				ch := e.GetChosenCase()
				if ch != nil {
					parseChannelOp(ch)
				}
			}
		}

	}
}

/*
 * Parse a new elem element.
 * For now only channels are considered
 * Add the corresponding info into fuzzingChannel
 */
func parseNew(elem *analysis.TraceElementNew) {
	// only process channels
	if elem.GetObjType() != "NC" {
		return
	}

	fuzzingElem := fuzzingChannel{
		globalID:  elem.GetPos(),
		localID:   elem.GetID(),
		closeInfo: unknown,
		qSize:     elem.GetNum(),
		maxQCount: 0,
	}

	channelInfoTrace[fuzzingElem.localID] = fuzzingElem
}

/*
 * Parse a channel operations.
 * If the operation is a close, update the data in channelInfoTrace
 * If it is an send (with tPost != 0), add it to pairInfoTrace
 * If it is an recv, it is either tPost = 0 (ignore) or will be handled by the send
 */
func parseChannelOp(elem *analysis.TraceElementChannel) {
	op := elem.GetObjType()

	// close -> update channelInfoTrace
	if op == "CC" {
		e := channelInfoTrace[elem.GetID()]
		e.closeInfo = always // before is always unknown
		channelInfoTrace[elem.GetID()] = e
		numberClose++
	} else if op == "CS" {
		if elem.GetTPost() == 0 {
			return
		}

		recv := elem.GetPartner()
		if recv == nil {
			panic("fuzzing parseChannelOp, send without partner: first run find partner")
		}

		sendPos := elem.GetPos()
		recvPos := recv.GetPos()
		chanID := elem.GetID()
		key := sendPos + "-" + recvPos

		if e, ok := pairInfoTrace[key]; ok {
			e.com++
			pairInfoTrace[key] = e
		} else {
			fp := fuzzingPair{
				sendID: sendPos,
				recvID: recvPos,
				chanID: chanID,
				com:    1,
			}
			pairInfoTrace[key] = fp
		}

		channelNew := channelInfoTrace[chanID]
		channelNew.maxQCount = max(channelNew.maxQCount, elem.GetQCount())
	}
}
