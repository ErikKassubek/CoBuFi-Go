// Copyright (c) 2024 Erik Kassubek
//
// File: traceElementAtomic_test.go
// Brief: Tests for traceElementAtomic
//
// Author: Erik Kassubek
// Created: 2024-11-12
//
// License: BSD-3-Clause

package analysis

import (
	"analyzer/clock"
	"errors"
	"reflect"
	"testing"
)

func TestTraceElementAtomicNew(t *testing.T) {
	var tests = []struct {
		name       string
		routine    int
		tPost      string
		id         string
		operation  string
		position   string
		expRoutine int
		expTPost   int
		expID      int
		expOp      opAtomic
		expPos     string
		expError   error
	}{
		{"Valid Load", 1, "213", "123", "L", "testfile.go:99", 1, 213, 123, LoadOp, "testfile.go:99", nil},
		{"Valid Store", 1, "213", "123", "S", "testfile.go:99", 1, 213, 123, StoreOp, "testfile.go:99", nil},
		{"Valid Add", 1, "213", "123", "A", "testfile.go:99", 1, 213, 123, AddOp, "testfile.go:99", nil},
		{"Valid Swap", 1, "213", "123", "W", "testfile.go:99", 1, 213, 123, SwapOp, "testfile.go:99", nil},
		{"Valid CompSwap", 1, "213", "123", "C", "testfile.go:99", 1, 213, 123, CompSwapOp, "testfile.go:99", nil},
		{"Invalid ID", 1, "ABC", "123", "L", "testfile.go:99", 0, 0, 0, LoadOp, "", errors.New("tpost is not an integer")},
		{"Invalid tPost", 1, "213", "ABC", "L", "testfile.go:99", 0, 0, 0, LoadOp, "", errors.New("id is not an integer")},
		{"Invalid operation", 1, "213", "321", "Q", "testfile.go:99", 0, 0, 0, LoadOp, "", errors.New("operation is not a valid operation")},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := AddTraceElementAtomic(test.routine, test.tPost, test.id, test.operation, test.position)

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

			elem := trace[len(trace)-1].(*TraceElementAtomic)

			if elem.routine != test.expRoutine {
				t.Errorf("Incorrect routine. Expected %d. Got %d.", test.expRoutine, elem.routine)
			}

			if elem.tPost != test.expTPost {
				t.Errorf("Incorrect tPost. Expected %d. Got %d.", test.expTPost, elem.tPost)
			}

			if elem.id != test.expID {
				t.Errorf("Incorrect ID. Expected %d. Got %d.", test.expID, elem.id)
			}

			if elem.opA != test.expOp {
				t.Errorf("Incorrect op. Expected %d. Got %d.", test.expOp, elem.opA)
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

func TestTraceElementAtomicSet(t *testing.T) {
}

func TestAtomicUpdateVectorClock(t *testing.T) {
	t.Run("LoadOp", func(t *testing.T) {
		at := TraceElementAtomic{id: 1, routine: 2, opA: LoadOp}
		currentVCHb = map[int]clock.VectorClock{2: clock.NewVectorClockSet(3, map[int]int{1: 2, 2: 4, 3: 1})}
		lw = map[int]clock.VectorClock{1: clock.NewVectorClockSet(3, map[int]int{1: 2, 2: 3, 3: 0})}

		expectedVC := map[int]clock.VectorClock{2: clock.NewVectorClockSet(3, map[int]int{1: 2, 2: 5, 3: 1})}
		expectedAtVC := currentVCHb[at.routine].Copy()

		at.updateVectorClock()

		if !reflect.DeepEqual(currentVCHb, expectedVC) {
			t.Errorf("Incorrect vc. Expected %v. Got %v.", expectedVC, currentVCHb)
		}

		if !reflect.DeepEqual(at.vc, expectedAtVC) {
			t.Errorf("Incorrect at vc. Expected %v. Got %v.", expectedAtVC, at.vc)
		}
	})

	t.Run("Store", func(t *testing.T) {
		at := TraceElementAtomic{id: 1, routine: 2, opA: StoreOp}
		currentVCHb = map[int]clock.VectorClock{2: clock.NewVectorClockSet(3, map[int]int{1: 2, 2: 4, 3: 1})}
		lw = map[int]clock.VectorClock{1: clock.NewVectorClockSet(3, map[int]int{1: 2, 2: 3, 3: 0})}

		expectedLW := map[int]clock.VectorClock{1: currentVCHb[2].Copy()}
		expectedVC := map[int]clock.VectorClock{2: clock.NewVectorClockSet(3, map[int]int{1: 2, 2: 5, 3: 1})}
		expectedAtVC := currentVCHb[at.routine].Copy()

		at.updateVectorClock()

		if !reflect.DeepEqual(lw, expectedLW) {
			t.Errorf("Incorrect lw. Expected %v. Got %v.", expectedLW, lw)
		}

		if !reflect.DeepEqual(currentVCHb, expectedVC) {
			t.Errorf("Incorrect vc. Expected %v. Got %v.", expectedVC, currentVCHb)
		}

		if !reflect.DeepEqual(at.vc, expectedAtVC) {
			t.Errorf("Incorrect at vc. Expected %v. Got %v.", expectedAtVC, at.vc)
		}
	})

	t.Run("Add", func(t *testing.T) {
		at := TraceElementAtomic{id: 1, routine: 2, opA: AddOp}
		currentVCHb = map[int]clock.VectorClock{2: clock.NewVectorClockSet(3, map[int]int{1: 2, 2: 4, 3: 1})}
		lw = map[int]clock.VectorClock{1: clock.NewVectorClockSet(3, map[int]int{1: 2, 2: 3, 3: 0})}

		expectedLW := map[int]clock.VectorClock{1: currentVCHb[2].Copy()}
		expectedVC := map[int]clock.VectorClock{2: clock.NewVectorClockSet(3, map[int]int{1: 2, 2: 5, 3: 1})}
		expectedAtVC := currentVCHb[at.routine].Copy()

		at.updateVectorClock()

		if !reflect.DeepEqual(lw, expectedLW) {
			t.Errorf("Incorrect lw. Expected %v. Got %v.", expectedLW, lw)
		}

		if !reflect.DeepEqual(currentVCHb, expectedVC) {
			t.Errorf("Incorrect vc. Expected %v. Got %v.", expectedVC, currentVCHb)
		}

		if !reflect.DeepEqual(at.vc, expectedAtVC) {
			t.Errorf("Incorrect at vc. Expected %v. Got %v.", expectedAtVC, at.vc)
		}
	})

	t.Run("Swap", func(t *testing.T) {
		at := TraceElementAtomic{id: 1, routine: 2, opA: SwapOp}
		currentVCHb = map[int]clock.VectorClock{2: clock.NewVectorClockSet(3, map[int]int{1: 2, 2: 4, 3: 1})}
		lw = map[int]clock.VectorClock{1: clock.NewVectorClockSet(3, map[int]int{1: 2, 2: 3, 3: 0})}

		expectedVC := map[int]clock.VectorClock{2: clock.NewVectorClockSet(3, map[int]int{1: 2, 2: 6, 3: 1})}
		expectedAtVC := currentVCHb[at.routine].Copy()

		at.updateVectorClock()

		if !reflect.DeepEqual(currentVCHb, expectedVC) {
			t.Errorf("Incorrect vc. Expected %v. Got %v.", expectedVC, currentVCHb)
		}

		if !reflect.DeepEqual(at.vc, expectedAtVC) {
			t.Errorf("Incorrect at vc. Expected %v. Got %v.", expectedAtVC, at.vc)
		}
	})

	t.Run("CompSwap", func(t *testing.T) {
		at := TraceElementAtomic{id: 1, routine: 2, opA: CompSwapOp}
		currentVCHb = map[int]clock.VectorClock{2: clock.NewVectorClockSet(3, map[int]int{1: 2, 2: 4, 3: 1})}
		lw = map[int]clock.VectorClock{1: clock.NewVectorClockSet(3, map[int]int{1: 2, 2: 3, 3: 0})}

		expectedVC := map[int]clock.VectorClock{2: clock.NewVectorClockSet(3, map[int]int{1: 2, 2: 6, 3: 1})}
		expectedAtVC := currentVCHb[at.routine].Copy()

		at.updateVectorClock()

		if !reflect.DeepEqual(currentVCHb, expectedVC) {
			t.Errorf("Incorrect vc. Expected %v. Got %v.", expectedVC, currentVCHb)
		}

		if !reflect.DeepEqual(at.vc, expectedAtVC) {
			t.Errorf("Incorrect at vc. Expected %v. Got %v.", expectedAtVC, at.vc)
		}
	})
}
