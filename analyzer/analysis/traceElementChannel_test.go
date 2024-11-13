// Copyright (c) 2024 Erik Kassubek
//
// File: traceElementChannel_test.go
// Brief: Tests for traceElementChannel
//
// Author: Erik Kassubek
// Created: 2024-11-12
//
// License: BSD-3-Clause

package analysis

import (
	"analyzer/clock"
	"errors"
	"testing"
)

func TestTraceElementChannelNew(t *testing.T) {
	var tests = []struct {
		name       string
		routine    int
		tPre       string
		tPost      string
		id         string
		op         string
		cl         string
		oID        string
		qSize      string
		pos        string
		expRoutine int
		expTPre    int
		expTPost   int
		expID      int
		expOp      OpChannel
		expCl      bool
		expOID     int
		expQSize   int
		expPos     string
		expError   error
	}{
		// Valid cases
		{
			name:       "Valid case with op S and cl t",
			routine:    1,
			tPre:       "10",
			tPost:      "20",
			id:         "100",
			op:         "S",
			cl:         "t",
			oID:        "200",
			qSize:      "10",
			pos:        "testfile.go:10",
			expRoutine: 1,
			expTPre:    10,
			expTPost:   20,
			expID:      100,
			expOp:      SendOp,
			expCl:      true,
			expOID:     200,
			expQSize:   10,
			expPos:     "testfile.go:10",
			expError:   nil,
		},
		{
			name:       "Valid case with op R and cl f",
			routine:    2,
			tPre:       "15",
			tPost:      "25",
			id:         "101",
			op:         "R",
			cl:         "f",
			oID:        "201",
			qSize:      "11",
			pos:        "testfile.go:11",
			expRoutine: 2,
			expTPre:    15,
			expTPost:   25,
			expID:      101,
			expOp:      RecvOp,
			expCl:      false,
			expOID:     201,
			expQSize:   11,
			expPos:     "testfile.go:11",
			expError:   nil,
		},
		{
			name:       "Valid case with op C and cl t",
			routine:    3,
			tPre:       "20",
			tPost:      "30",
			id:         "102",
			op:         "C",
			cl:         "t",
			oID:        "202",
			qSize:      "12",
			pos:        "testfile.go:12",
			expRoutine: 3,
			expTPre:    20,
			expTPost:   30,
			expID:      102,
			expOp:      CloseOp,
			expCl:      true,
			expOID:     202,
			expQSize:   12,
			expPos:     "testfile.go:12",
			expError:   nil,
		},
		{
			name:       "Valid case with id *",
			routine:    4,
			tPre:       "25",
			tPost:      "35",
			id:         "*",
			op:         "S",
			cl:         "f",
			oID:        "203",
			qSize:      "13",
			pos:        "testfile.go:13",
			expRoutine: 4,
			expTPre:    25,
			expTPost:   35,
			expID:      -1,
			expOp:      SendOp,
			expCl:      false,
			expOID:     203,
			expQSize:   13,
			expPos:     "testfile.go:13",
			expError:   nil,
		},
		// Invalid cases
		{
			name:     "Invalid tPre",
			routine:  1,
			tPre:     "invalid",
			tPost:    "20",
			id:       "100",
			op:       "S",
			cl:       "t",
			oID:      "200",
			qSize:    "10",
			pos:      "testfile.go:10",
			expError: errors.New("tPre is not an integer"),
		},
		{
			name:     "Invalid tPost",
			routine:  1,
			tPre:     "10",
			tPost:    "invalid",
			id:       "100",
			op:       "S",
			cl:       "t",
			oID:      "200",
			qSize:    "10",
			pos:      "testfile.go:10",
			expError: errors.New("tPost is not an integer"),
		},
		{
			name:     "Invalid id",
			routine:  1,
			tPre:     "10",
			tPost:    "20",
			id:       "invalid",
			op:       "S",
			cl:       "t",
			oID:      "200",
			qSize:    "10",
			pos:      "testfile.go:10",
			expError: errors.New("id is not an integer"),
		},
		{
			name:     "Invalid op",
			routine:  1,
			tPre:     "10",
			tPost:    "20",
			id:       "100",
			op:       "invalid",
			cl:       "t",
			oID:      "200",
			qSize:    "10",
			pos:      "testfile.go:10",
			expError: errors.New("opC is not a valid operation"),
		},
		{
			name:     "Invalid cl",
			routine:  1,
			tPre:     "10",
			tPost:    "20",
			id:       "100",
			op:       "S",
			cl:       "invalid",
			oID:      "200",
			qSize:    "10",
			pos:      "testfile.go:10",
			expError: errors.New("cl is not a boolean"),
		},
		{
			name:     "Invalid oID",
			routine:  1,
			tPre:     "10",
			tPost:    "20",
			id:       "100",
			op:       "S",
			cl:       "t",
			oID:      "invalid",
			qSize:    "10",
			pos:      "testfile.go:10",
			expError: errors.New("oId is not an integer"),
		},
		{
			name:     "Invalid qSize",
			routine:  1,
			tPre:     "10",
			tPost:    "20",
			id:       "100",
			op:       "S",
			cl:       "t",
			oID:      "200",
			qSize:    "invalid",
			pos:      "testfile.go:10",
			expError: errors.New("qSize is not an integer"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := AddTraceElementChannel(test.routine, test.tPre, test.tPost, test.id, test.op, test.cl, test.oID, test.qSize, test.pos)

			if test.expError == nil && err != nil {
				t.Errorf("Received unexpected error %s", err.Error())
				return
			}

			if test.expError != nil && err == nil {
				t.Errorf("Expected error, but no error was triggered")
				return
			}

			if err != nil && err.Error() != test.expError.Error() {
				t.Errorf("Incorrect Error. Expected %s. Got %s.", test.expError, err)
				return
			}

			if err != nil {
				return
			}

			trace := GetTraceFromId(test.routine)

			elem := trace[len(trace)-1].(*TraceElementChannel)

			if elem.routine != test.expRoutine {
				t.Errorf("Incorrect routine. Expected %d. Got %d.", test.expRoutine, elem.routine)
			}

			if elem.tPre != test.expTPre {
				t.Errorf("Incorrect tPre. Expected %d. Got %d.", test.expTPre, elem.tPre)
			}

			if elem.tPost != test.expTPost {
				t.Errorf("Incorrect tPost. Expected %d. Got %d.", test.expTPost, elem.tPost)
			}

			if elem.id != test.expID {
				t.Errorf("Incorrect ID. Expected %d. Got %d.", test.expID, elem.id)
			}

			if elem.opC != test.expOp {
				t.Errorf("Incorrect op. Expected %d. Got %d.", test.expOp, elem.opC)
			}

			if elem.qSize != test.expQSize {
				t.Errorf("Incorrect qSize. Expected %d. Got %d.", test.expQSize, elem.qSize)
			}

			if elem.pos != test.expPos {
				t.Errorf("Incorrect pos. Expected %s. Got %s.", test.expPos, elem.pos)
			}
		})
	}
}

// func TestTraceElementChannelGet(t *testing.T) {
// }

// func TestTraceElementChannelSet(t *testing.T) {
// }

func TestChannelUpdateVectorClockUnbufferedSend(t *testing.T) {
	t.Run("Unbuffered Send", func(t *testing.T) {
		sendUnbuffered := TraceElementChannel{
			routine: 1,
			tPre:    4,
			tPost:   6,
			id:      1,
			opC:     SendOp,
			cl:      false,
			oID:     1,
			qSize:   0,
			pos:     "exampleFile.go:111",
			vc:      clock.NewVectorClock(2),
		}

		recvUnbuffered := TraceElementChannel{
			routine: 2,
			tPre:    5,
			tPost:   7,
			id:      1,
			opC:     RecvOp,
			cl:      false,
			oID:     1,
			qSize:   0,
			pos:     "exampleFile.go:111",
			vc:      clock.NewVectorClock(2),
		}

		ClearTrace()
		AddElementToTrace(&sendUnbuffered)
		AddElementToTrace(&recvUnbuffered)

		sendT, _ := GetTraceElementFromTID(sendUnbuffered.GetTID())
		recvT, _ := GetTraceElementFromTID(recvUnbuffered.GetTID())

		send := (*sendT).(*TraceElementChannel)
		recv := (*recvT).(*TraceElementChannel)

		currentVCHb = map[int]clock.VectorClock{
			1: clock.NewVectorClockSet(2, map[int]int{1: 2, 2: 5}),
			2: clock.NewVectorClockSet(2, map[int]int{1: 7, 2: 3}),
		}

		expChVcSend := currentVCHb[sendUnbuffered.routine].Copy()
		expChVcRecv := currentVCHb[recvUnbuffered.routine].Copy()

		expVc := map[int]clock.VectorClock{
			1: clock.NewVectorClockSet(2, map[int]int{1: 8, 2: 5}),
			2: clock.NewVectorClockSet(2, map[int]int{1: 7, 2: 6}),
		}

		(*send).updateVectorClock()

		if !(*send).vc.IsEqual(expChVcSend) {
			t.Errorf("Incorrect ch vc send. Expected %v. Got %v", expChVcSend, (*send).vc)
		}

		if !(*recv).vc.IsEqual(expChVcRecv) {
			t.Errorf("Incorrect ch vc recv. Expected %v. Got %v", expChVcRecv, (*recv).vc)
		}

		if !clock.IsMapVcEqual(currentVCHb, expVc) {
			t.Errorf("Incorrect currentVCHb send. Expected %v. Got %v", expVc, currentVCHb)
		}
	})
}

func TestChannelUpdateVectorClockUnbufferedRecv(t *testing.T) {
	sendElem := TraceElementChannel{
		routine: 1,
		tPre:    4,
		tPost:   6,
		id:      1,
		opC:     SendOp,
		cl:      false,
		oID:     1,
		qSize:   0,
		pos:     "exampleFile.go:111",
		vc:      clock.NewVectorClock(2),
	}

	recvElem := TraceElementChannel{
		routine: 2,
		tPre:    5,
		tPost:   7,
		id:      1,
		opC:     RecvOp,
		cl:      false,
		oID:     1,
		qSize:   0,
		pos:     "exampleFile.go:111",
		vc:      clock.NewVectorClock(2),
	}

	ClearTrace()
	AddElementToTrace(&sendElem)
	AddElementToTrace(&recvElem)

	sendT, _ := GetTraceElementFromTID(sendElem.GetTID())
	recvT, _ := GetTraceElementFromTID(recvElem.GetTID())

	send := (*sendT).(*TraceElementChannel)
	recv := (*recvT).(*TraceElementChannel)

	currentVCHb = map[int]clock.VectorClock{
		1: clock.NewVectorClockSet(2, map[int]int{1: 2, 2: 5}),
		2: clock.NewVectorClockSet(2, map[int]int{1: 7, 2: 3}),
	}

	expChVcSend := currentVCHb[sendElem.routine].Copy()
	expChVcRecv := currentVCHb[recvElem.routine].Copy()

	expVc := map[int]clock.VectorClock{
		1: clock.NewVectorClockSet(2, map[int]int{1: 8, 2: 5}),
		2: clock.NewVectorClockSet(2, map[int]int{1: 7, 2: 6}),
	}

	(*recv).updateVectorClock()

	if !(*send).vc.IsEqual(expChVcSend) {
		t.Errorf("Incorrect ch vc send. Expected %v. Got %v", expChVcSend, (*send).vc)
	}

	if !(*recv).vc.IsEqual(expChVcRecv) {
		t.Errorf("Incorrect ch vc recv. Expected %v. Got %v", expChVcRecv, (*recv).vc)
	}

	if !clock.IsMapVcEqual(currentVCHb, expVc) {
		t.Errorf("Incorrect currentVCHb recv. Expected %v. Got %v", expVc, currentVCHb)
	}
}

func TestChannelUpdateVectorClockBufferedSend(t *testing.T) {
	sendElem := TraceElementChannel{
		routine: 1,
		tPre:    4,
		tPost:   6,
		id:      1,
		opC:     SendOp,
		cl:      false,
		oID:     1,
		qSize:   1,
		pos:     "exampleFile.go:111",
		vc:      clock.NewVectorClock(2),
	}

	traceElem := TraceElementChannel{
		routine: 2,
		tPre:    5,
		tPost:   7,
		id:      1,
		opC:     RecvOp,
		cl:      false,
		oID:     1,
		qSize:   1,
		pos:     "exampleFile.go:111",
		vc:      clock.NewVectorClock(2),
	}

	ClearTrace()
	AddElementToTrace(&sendElem)
	AddElementToTrace(&traceElem)

	sendT, _ := GetTraceElementFromTID(sendElem.GetTID())

	send := (*sendT).(*TraceElementChannel)

	currentVCHb = map[int]clock.VectorClock{
		1: clock.NewVectorClockSet(2, map[int]int{1: 2, 2: 5}),
		2: clock.NewVectorClockSet(2, map[int]int{1: 7, 2: 3}),
	}

	expChVcSend := currentVCHb[sendElem.routine].Copy()

	expVc := map[int]clock.VectorClock{
		1: clock.NewVectorClockSet(2, map[int]int{1: 3, 2: 5}),
		2: clock.NewVectorClockSet(2, map[int]int{1: 7, 2: 3}),
	}

	(*send).updateVectorClock()

	if !(*send).vc.IsEqual(expChVcSend) {
		t.Errorf("Incorrect ch vc send. Expected %v. Got %v", expChVcSend, (*send).vc)
	}

	if !clock.IsMapVcEqual(currentVCHb, expVc) {
		t.Errorf("Incorrect currentVCHb send. Expected %v. Got %v", expVc, currentVCHb)
	}
}

func TestChannelUpdateVectorClockBufferedRecv(t *testing.T) {
	sendElem := TraceElementChannel{
		routine: 1,
		tPre:    4,
		tPost:   6,
		id:      111,
		opC:     SendOp,
		cl:      false,
		oID:     1,
		qSize:   1,
		pos:     "exampleFile.go:111",
		vc:      clock.NewVectorClock(2),
	}

	recvElem := TraceElementChannel{
		routine: 2,
		tPre:    5,
		tPost:   7,
		id:      111,
		opC:     RecvOp,
		cl:      false,
		oID:     1,
		qSize:   1,
		pos:     "exampleFile.go:222",
		vc:      clock.NewVectorClock(2),
	}

	ClearTrace()
	AddElementToTrace(&sendElem)
	AddElementToTrace(&recvElem)

	sendT, _ := GetTraceElementFromTID(sendElem.GetTID())
	recvT, _ := GetTraceElementFromTID(recvElem.GetTID())

	send := (*sendT).(*TraceElementChannel)
	recv := (*recvT).(*TraceElementChannel)

	currentVCHb = map[int]clock.VectorClock{
		1: clock.NewVectorClockSet(2, map[int]int{1: 9, 2: 2}),
		2: clock.NewVectorClockSet(2, map[int]int{1: 8, 2: 6}),
	}

	expChVcSend := currentVCHb[sendElem.routine].Copy()
	expChVcRecv := currentVCHb[recvElem.routine].Copy()

	expVc := map[int]clock.VectorClock{
		1: clock.NewVectorClockSet(2, map[int]int{1: 10, 2: 2}),
		2: clock.NewVectorClockSet(2, map[int]int{1: 10, 2: 7}),
	}

	(*send).updateVectorClock()
	(*recv).updateVectorClock()

	if !(*send).vc.IsEqual(expChVcSend) {
		t.Errorf("Incorrect ch vc send. Expected %v. Got %v", expChVcSend, (*send).vc)
	}

	if !(*recv).vc.IsEqual(expChVcRecv) {
		t.Errorf("Incorrect ch vc recv. Expected %v. Got %v", expChVcRecv, (*recv).vc)
	}

	if !clock.IsMapVcEqual(currentVCHb, expVc) {
		t.Errorf("Incorrect currentVCHb send. Expected %v. Got %v", expVc, currentVCHb)
	}
}

func TestChannelUpdateVectorClockClose(t *testing.T) {
	closeElem := TraceElementChannel{
		routine: 2,
		tPre:    5,
		tPost:   7,
		id:      1,
		opC:     CloseOp,
		cl:      false,
		oID:     1,
		qSize:   1,
		pos:     "exampleFile.go:111",
		vc:      clock.NewVectorClock(2),
	}

	currentVCHb = map[int]clock.VectorClock{
		1: clock.NewVectorClockSet(2, map[int]int{1: 2, 2: 5}),
		2: clock.NewVectorClockSet(2, map[int]int{1: 7, 2: 3}),
	}
	expClAt := currentVCHb[closeElem.routine].Copy()
	expCurrentVCHb := map[int]clock.VectorClock{
		1: clock.NewVectorClockSet(2, map[int]int{1: 2, 2: 5}),
		2: clock.NewVectorClockSet(2, map[int]int{1: 7, 2: 4}),
	}

	closeElem.updateVectorClock()

	if !closeElem.vc.IsEqual(expClAt) {
		t.Errorf("Incorrect ch vc. Expected %v. Got %v", expClAt, closeElem.vc)
	}

	if !clock.IsMapVcEqual(currentVCHb, expCurrentVCHb) {
		t.Errorf("Incorrect currentVCHb. Expected %v. Got %v", expCurrentVCHb, currentVCHb)
	}

}

func TestChannelFindPartner(t *testing.T) {
	sendElem := TraceElementChannel{
		routine: 1,
		tPre:    4,
		tPost:   6,
		id:      1,
		opC:     SendOp,
		cl:      false,
		oID:     1,
		qSize:   0,
		pos:     "exampleFile.go:111",
		vc:      clock.NewVectorClock(2),
	}

	recvElem := TraceElementChannel{
		routine: 2,
		tPre:    5,
		tPost:   7,
		id:      1,
		opC:     RecvOp,
		cl:      false,
		oID:     1,
		qSize:   0,
		pos:     "exampleFile.go:111",
		vc:      clock.NewVectorClock(2),
	}

	ClearTrace()
	AddElementToTrace(&sendElem)
	AddElementToTrace(&recvElem)

	sendT, _ := GetTraceElementFromTID(sendElem.GetTID())
	recvT, _ := GetTraceElementFromTID(recvElem.GetTID())

	send := (*sendT).(*TraceElementChannel)
	recv := (*recvT).(*TraceElementChannel)

	res := send.findPartner()

	if res != recv.routine {
		t.Errorf("Incorrect result for send.findPartner. Expected %d. Got %d.", recv.GetRoutine(), res)
	}

	if send.partner.ToString() != recv.ToString() {
		t.Errorf("Incorrect value for partner send. Expected %s. Got %s.", recv.ToString(), send.partner.ToString())
	}

	if recv.partner.ToString() != send.ToString() {
		t.Errorf("Incorrect value for partner recv. Expected %s. Got %s.", send.ToString(), recv.partner.ToString())
	}
}
