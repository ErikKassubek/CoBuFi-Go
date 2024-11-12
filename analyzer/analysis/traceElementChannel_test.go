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

// func TestTraceElementAtomicGet(t *testing.T) {
// 	var tests = []struct {
// 		name          string
// 		routine       int
// 		tPost         string
// 		id            string
// 		operation     string
// 		position      string
// 		expRoutine    int
// 		expTPre       int
// 		expTSort      int
// 		expTPost      int
// 		expID         int
// 		expPos        string
// 		expTID        string
// 		expObjectType string
// 		expVC         clock.VectorClock
// 		expString     string
// 	}{
// 		{"Valid Load", 1, "213", "123", "L", "testfile.go:99", 1, 213, 213, 213, 123, "testfile.go:99", "testfile.go:99@213", "AL", clock.NewVectorClock(0), "A,213,123,L,testfile.go:99"},
// 		{"Valid Store", 1, "213", "123", "S", "testfile.go:99", 1, 213, 213, 213, 123, "testfile.go:99", "testfile.go:99@213", "AS", clock.NewVectorClock(0), "A,213,123,S,testfile.go:99"},
// 		{"Valid Add", 1, "213", "123", "A", "testfile.go:99", 1, 213, 213, 213, 123, "testfile.go:99", "testfile.go:99@213", "AA", clock.NewVectorClock(0), "A,213,123,A,testfile.go:99"},
// 		{"Valid Swap", 1, "213", "123", "W", "testfile.go:99", 1, 213, 213, 213, 123, "testfile.go:99", "testfile.go:99@213", "AW", clock.NewVectorClock(0), "A,213,123,W,testfile.go:99"},
// 		{"Valid CompSwap", 1, "213", "123", "C", "testfile.go:99", 1, 213, 213, 213, 123, "testfile.go:99", "testfile.go:99@213", "AC", clock.NewVectorClock(0), "A,213,123,C,testfile.go:99"},
// 	}

// 	for _, test := range tests {
// 		AddTraceElementAtomic(test.routine, test.tPost, test.id, test.operation, test.position)

// 		trace := GetTraceFromId(test.routine)

// 		elem := trace[len(trace)-1].(*TraceElementAtomic)

// 		if elem.GetRoutine() != test.expRoutine {
// 			t.Errorf("Incorrect routine. Expected %d. Got %d.", test.expRoutine, elem.routine)
// 		}

// 		if elem.GetTPre() != test.expTPre {
// 			t.Errorf("Incorrect tPre. Expected %d. Got %d.", test.expTPre, elem.GetTPre())
// 		}

// 		if elem.GetTSort() != test.expTSort {
// 			t.Errorf("Incorrect tSort. Expected %d. Got %d.", test.expTSort, elem.GetTSort())
// 		}

// 		if elem.getTpost() != test.expTPost {
// 			t.Errorf("Incorrect tPost. Expected %d. Got %d.", test.expTPost, elem.tPost)
// 		}

// 		if elem.GetID() != test.expID {
// 			t.Errorf("Incorrect ID. Expected %d. Got %d.", test.expID, elem.GetID())
// 		}

// 		if elem.GetPos() != test.expPos {
// 			t.Errorf("Incorrect position. Expected %s. Got %s.", test.expPos, elem.GetPos())
// 		}

// 		if elem.GetTID() != test.expTID {
// 			t.Errorf("Incorrect tID. Expected %s. Got %s.", test.expTID, elem.GetTID())
// 		}

// 		if elem.GetObjType() != test.expObjectType {
// 			t.Errorf("Incorrect object type. Expected %s. Got %s.", test.expObjectType, elem.GetObjType())
// 		}

// 		if !elem.GetVC().IsEqual(test.expVC) {
// 			t.Errorf("Incorrect VC. Expected %v. Got %v.", test.expVC, elem.GetVC())
// 		}

// 		if elem.ToString() != test.expString {
// 			t.Errorf("Incorrect string. Expected %s. Got %s.", test.expString, elem.ToString())
// 		}
// 	}
// }

// func TestUpdateVectorClock(t *testing.T) {
// 	t.Run("LoadOp", func(t *testing.T) {
// 		at := TraceElementAtomic{id: 1, routine: 2, opA: LoadOp}
// 		currentVCHb = map[int]clock.VectorClock{2: clock.NewVectorClockSet(3, map[int]int{1: 2, 2: 4, 3: 1})}
// 		lw = map[int]clock.VectorClock{1: clock.NewVectorClockSet(3, map[int]int{1: 2, 2: 3, 3: 0})}

// 		expectedVC := map[int]clock.VectorClock{2: clock.NewVectorClockSet(3, map[int]int{1: 2, 2: 5, 3: 1})}
// 		expectedAtVC := currentVCHb[at.routine].Copy()

// 		at.updateVectorClock()

// 		if !reflect.DeepEqual(currentVCHb, expectedVC) {
// 			t.Errorf("Incorrect vc. Expected %v. Got %v.", expectedVC, currentVCHb)
// 		}

// 		if !reflect.DeepEqual(at.vc, expectedAtVC) {
// 			t.Errorf("Incorrect at vc. Expected %v. Got %v.", expectedAtVC, at.vc)
// 		}
// 	})

// 	t.Run("Store", func(t *testing.T) {
// 		at := TraceElementAtomic{id: 1, routine: 2, opA: StoreOp}
// 		currentVCHb = map[int]clock.VectorClock{2: clock.NewVectorClockSet(3, map[int]int{1: 2, 2: 4, 3: 1})}
// 		lw = map[int]clock.VectorClock{1: clock.NewVectorClockSet(3, map[int]int{1: 2, 2: 3, 3: 0})}

// 		expectedLW := map[int]clock.VectorClock{1: currentVCHb[2].Copy()}
// 		expectedVC := map[int]clock.VectorClock{2: clock.NewVectorClockSet(3, map[int]int{1: 2, 2: 5, 3: 1})}
// 		expectedAtVC := currentVCHb[at.routine].Copy()

// 		at.updateVectorClock()

// 		if !reflect.DeepEqual(lw, expectedLW) {
// 			t.Errorf("Incorrect lw. Expected %v. Got %v.", expectedLW, lw)
// 		}

// 		if !reflect.DeepEqual(currentVCHb, expectedVC) {
// 			t.Errorf("Incorrect vc. Expected %v. Got %v.", expectedVC, currentVCHb)
// 		}

// 		if !reflect.DeepEqual(at.vc, expectedAtVC) {
// 			t.Errorf("Incorrect at vc. Expected %v. Got %v.", expectedAtVC, at.vc)
// 		}
// 	})

// 	t.Run("Add", func(t *testing.T) {
// 		at := TraceElementAtomic{id: 1, routine: 2, opA: AddOp}
// 		currentVCHb = map[int]clock.VectorClock{2: clock.NewVectorClockSet(3, map[int]int{1: 2, 2: 4, 3: 1})}
// 		lw = map[int]clock.VectorClock{1: clock.NewVectorClockSet(3, map[int]int{1: 2, 2: 3, 3: 0})}

// 		expectedLW := map[int]clock.VectorClock{1: currentVCHb[2].Copy()}
// 		expectedVC := map[int]clock.VectorClock{2: clock.NewVectorClockSet(3, map[int]int{1: 2, 2: 5, 3: 1})}
// 		expectedAtVC := currentVCHb[at.routine].Copy()

// 		at.updateVectorClock()

// 		if !reflect.DeepEqual(lw, expectedLW) {
// 			t.Errorf("Incorrect lw. Expected %v. Got %v.", expectedLW, lw)
// 		}

// 		if !reflect.DeepEqual(currentVCHb, expectedVC) {
// 			t.Errorf("Incorrect vc. Expected %v. Got %v.", expectedVC, currentVCHb)
// 		}

// 		if !reflect.DeepEqual(at.vc, expectedAtVC) {
// 			t.Errorf("Incorrect at vc. Expected %v. Got %v.", expectedAtVC, at.vc)
// 		}
// 	})

// 	t.Run("Swap", func(t *testing.T) {
// 		at := TraceElementAtomic{id: 1, routine: 2, opA: SwapOp}
// 		currentVCHb = map[int]clock.VectorClock{2: clock.NewVectorClockSet(3, map[int]int{1: 2, 2: 4, 3: 1})}
// 		lw = map[int]clock.VectorClock{1: clock.NewVectorClockSet(3, map[int]int{1: 2, 2: 3, 3: 0})}

// 		expectedVC := map[int]clock.VectorClock{2: clock.NewVectorClockSet(3, map[int]int{1: 2, 2: 6, 3: 1})}
// 		expectedAtVC := currentVCHb[at.routine].Copy()

// 		at.updateVectorClock()

// 		if !reflect.DeepEqual(currentVCHb, expectedVC) {
// 			t.Errorf("Incorrect vc. Expected %v. Got %v.", expectedVC, currentVCHb)
// 		}

// 		if !reflect.DeepEqual(at.vc, expectedAtVC) {
// 			t.Errorf("Incorrect at vc. Expected %v. Got %v.", expectedAtVC, at.vc)
// 		}
// 	})

// 	t.Run("CompSwap", func(t *testing.T) {
// 		at := TraceElementAtomic{id: 1, routine: 2, opA: CompSwapOp}
// 		currentVCHb = map[int]clock.VectorClock{2: clock.NewVectorClockSet(3, map[int]int{1: 2, 2: 4, 3: 1})}
// 		lw = map[int]clock.VectorClock{1: clock.NewVectorClockSet(3, map[int]int{1: 2, 2: 3, 3: 0})}

// 		expectedVC := map[int]clock.VectorClock{2: clock.NewVectorClockSet(3, map[int]int{1: 2, 2: 6, 3: 1})}
// 		expectedAtVC := currentVCHb[at.routine].Copy()

// 		at.updateVectorClock()

// 		if !reflect.DeepEqual(currentVCHb, expectedVC) {
// 			t.Errorf("Incorrect vc. Expected %v. Got %v.", expectedVC, currentVCHb)
// 		}

// 		if !reflect.DeepEqual(at.vc, expectedAtVC) {
// 			t.Errorf("Incorrect at vc. Expected %v. Got %v.", expectedAtVC, at.vc)
// 		}
// 	})
// }
