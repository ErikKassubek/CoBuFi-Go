// Copyright (c) 2024 Erik Kassubek
//
// File: vcWait_test.go
// Brief: Tests for vcWait
//
// Author: Erik Kassubek
// Created: 2024-11-12
//
// License: BSD-3-Clause

package analysis

import (
	"analyzer/clock"
	"reflect"
	"testing"
)

func TestChange(t *testing.T) {
	wa := TraceElementWait{
		routine: 1,
		tPre:    12,
		tPost:   13,
		id:      123,
		opW:     ChangeOp,
		delta:   1,
		val:     1,
		pos:     "testfile.go:999",
		vc:      clock.NewVectorClock(2),
	}

	vc := map[int]clock.VectorClock{
		1: clock.NewVectorClockSet(2, map[int]int{1: 5, 2: 6}),
		2: clock.NewVectorClockSet(2, map[int]int{1: 2, 2: 9}),
	}

	lastChangeWG[wa.id] = clock.NewVectorClockSet(2, map[int]int{1: 3, 2: 8})

	expectedVc := map[int]clock.VectorClock{
		1: clock.NewVectorClockSet(2, map[int]int{1: 6, 2: 6}),
		2: clock.NewVectorClockSet(2, map[int]int{1: 2, 2: 9}),
	}
	expectedWg := clock.NewVectorClockSet(2, map[int]int{1: 5, 2: 8})

	Change(&wa, vc)

	t.Run("VC", func(t *testing.T) {
		if !reflect.DeepEqual(vc, expectedVc) {
			t.Errorf("Incorrect result. Expected %v. Got %v.", expectedVc, vc)
		}
	})

	t.Run("Wg", func(t *testing.T) {
		if !reflect.DeepEqual(lastChangeWG[wa.id], expectedWg) {
			t.Errorf("Incorrect result. Expected %v. Got %v.", expectedWg, lastChangeWG[wa.id])
		}
	})
}

func TestWgWait(t *testing.T) {
	wa := TraceElementWait{
		routine: 1,
		tPre:    12,
		tPost:   13,
		id:      123,
		opW:     WaitOp,
		delta:   1,
		val:     1,
		pos:     "testfile.go:999",
		vc:      clock.NewVectorClock(2),
	}

	vc := map[int]clock.VectorClock{
		1: clock.NewVectorClockSet(2, map[int]int{1: 5, 2: 6}),
		2: clock.NewVectorClockSet(2, map[int]int{1: 2, 2: 9}),
	}

	lastChangeWG[wa.id] = clock.NewVectorClockSet(2, map[int]int{1: 3, 2: 8})

	expectedVc := map[int]clock.VectorClock{
		1: clock.NewVectorClockSet(2, map[int]int{1: 6, 2: 8}),
		2: clock.NewVectorClockSet(2, map[int]int{1: 2, 2: 9}),
	}
	expectedWg := lastChangeWG[wa.id].Copy()

	Wait(&wa, vc)

	t.Run("VC", func(t *testing.T) {
		if !reflect.DeepEqual(vc, expectedVc) {
			t.Errorf("Incorrect result. Expected %v. Got %v.", expectedVc, vc)
		}
	})

	t.Run("Wg", func(t *testing.T) {
		if !reflect.DeepEqual(lastChangeWG[wa.id], expectedWg) {
			t.Errorf("Incorrect result. Expected %v. Got %v.", expectedWg, lastChangeWG[wa.id])
		}
	})
}
