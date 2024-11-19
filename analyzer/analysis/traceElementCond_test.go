// Copyright (c) 2024 Erik Kassubek
//
// File: traceElementCond_test.go
// Brief: Test for traceElementCond.go
//
// Author: Erik Kassubek
// Created: 2024-11-15
//
// License: BSD-3-Clause

package analysis

import (
	"analyzer/utils"
	"errors"
	"testing"
)

func TestTraceElementCondNew(t *testing.T) {
	var tests = []struct {
		name    string
		routine int
		tPre    string
		tPost   string
		id      string
		opN     string
		pos     string

		expRoutine int
		expTPre    int
		expTPost   int
		expID      int
		expOpC     OpCond
		expPos     string
		expErr     error
	}{
		{
			name:       "Valid case with op W",
			routine:    1,
			tPre:       "123",
			tPost:      "234",
			id:         "555",
			opN:        "W",
			pos:        "testfile:999",
			expRoutine: 1,
			expTPre:    123,
			expTPost:   234,
			expID:      555,
			expOpC:     WaitCondOp,
			expPos:     "testfile:999",
			expErr:     nil,
		},
		{
			name:       "Valid case with op S",
			routine:    1,
			tPre:       "123",
			tPost:      "234",
			id:         "555",
			opN:        "S",
			pos:        "testfile:999",
			expRoutine: 1,
			expTPre:    123,
			expTPost:   234,
			expID:      555,
			expOpC:     SignalOp,
			expPos:     "testfile:999",
			expErr:     nil,
		},
		{
			name:       "Valid case with op B",
			routine:    1,
			tPre:       "123",
			tPost:      "234",
			id:         "555",
			opN:        "B",
			pos:        "testfile:999",
			expRoutine: 1,
			expTPre:    123,
			expTPost:   234,
			expID:      555,
			expOpC:     BroadcastOp,
			expPos:     "testfile:999",
			expErr:     nil,
		},
		{
			name:       "Invalid tPre",
			routine:    1,
			tPre:       "12BG",
			tPost:      "234",
			id:         "555",
			opN:        "B",
			pos:        "testfile:999",
			expRoutine: 1,
			expTPre:    123,
			expTPost:   234,
			expID:      555,
			expOpC:     BroadcastOp,
			expPos:     "testfile:999",
			expErr:     errors.New("tpre is not an integer"),
		},
		{
			name:       "Invalid tPost",
			routine:    1,
			tPre:       "123",
			tPost:      "234BG",
			id:         "555",
			opN:        "B",
			pos:        "testfile:999",
			expRoutine: 1,
			expTPre:    123,
			expTPost:   234,
			expID:      555,
			expOpC:     BroadcastOp,
			expPos:     "testfile:999",
			expErr:     errors.New("tpost is not an integer"),
		},
		{
			name:       "Invalid id",
			routine:    1,
			tPre:       "123",
			tPost:      "234",
			id:         "55asd5",
			opN:        "B",
			pos:        "testfile:999",
			expRoutine: 1,
			expTPre:    123,
			expTPost:   234,
			expID:      555,
			expOpC:     BroadcastOp,
			expPos:     "testfile:999",
			expErr:     errors.New("id is not an integer"),
		},
		{
			name:       "Invalid operation",
			routine:    1,
			tPre:       "123",
			tPost:      "234",
			id:         "555",
			opN:        "Y",
			pos:        "testfile:999",
			expRoutine: 1,
			expTPre:    123,
			expTPost:   234,
			expID:      555,
			expOpC:     BroadcastOp,
			expPos:     "testfile:999",
			expErr:     errors.New("op is not a valid operation"),
		},
		{
			name:       "Empty operation",
			routine:    1,
			tPre:       "123",
			tPost:      "234",
			id:         "555",
			opN:        "",
			pos:        "testfile:999",
			expRoutine: 1,
			expTPre:    123,
			expTPost:   234,
			expID:      555,
			expOpC:     BroadcastOp,
			expPos:     "testfile:999",
			expErr:     errors.New("op is not a valid operation"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := AddTraceElementCond(test.routine, test.tPre, test.tPost,
				test.id, test.opN, test.pos)

			if res := utils.GetErrorDiff(test.expErr, err); res != nil {
				t.Errorf(res.Error())
			}

			if err != nil {
				return
			}

			trace := GetTraceFromId(test.routine)
			elem := trace[len(trace)-1].(*TraceElementCond)

			if elem.routine != test.expRoutine {
				t.Errorf("Incorrect routine. Expected %d. Got %d.", test.expRoutine,
					elem.routine)
			}

			if elem.tPre != test.expTPre {
				t.Errorf("Incorrect tPre. Expected %d. Got %d.", test.expTPre,
					elem.tPre)
			}

			if elem.tPost != test.expTPost {
				t.Errorf("Incorrect tPost. Expected %d. Got %d.", test.expTPost,
					elem.tPost)
			}

			if elem.id != test.expID {
				t.Errorf("Incorrect id. Expected %d. Got %d.", test.expID,
					elem.id)
			}

			if elem.opC != test.expOpC {
				t.Errorf("Incorrect opC. Expected %d. Got %d.", test.expOpC,
					elem.opC)
			}

			if elem.pos != test.expPos {
				t.Errorf("Incorrect pos. Expected %s. Got %s.", test.expPos,
					elem.pos)
			}
		})
	}
}

// func TestTraceElementCondGet(t *testing.T) {
// }

// func TestTraceElementCondSet(t *testing.T) {
// }

func TestCondUpdateVectorClockWait(t *testing.T) {

}
