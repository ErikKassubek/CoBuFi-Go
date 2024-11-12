// Copyright (c) 2024 Erik Kassubek
//
// File: vcCond_test.go
// Brief: Tests for vcCond.go
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

func TestWait(t *testing.T) {
	co1 := TraceElementCond{
		routine: 1,
		tPre:    5,
		tPost:   6,
		id:      123,
		opC:     WaitCondOp,
		pos:     "testfile.go:999",
		vc:      clock.NewVectorClock(2),
	}

	co2 := TraceElementCond{
		routine: 2,
		tPre:    7,
		tPost:   8,
		id:      123,
		opC:     WaitCondOp,
		pos:     "testfile.go:999",
		vc:      clock.NewVectorClock(2),
	}

	vc := map[int]clock.VectorClock{
		1: clock.NewVectorClockSet(3, map[int]int{1: 5, 2: 6, 3: 2}),
		2: clock.NewVectorClockSet(3, map[int]int{1: 2, 2: 9, 3: 2}),
		3: clock.NewVectorClockSet(3, map[int]int{1: 2, 2: 9, 3: 2}),
	}

	expectedVc := map[int]clock.VectorClock{
		1: clock.NewVectorClockSet(3, map[int]int{1: 6, 2: 6, 3: 2}),
		2: clock.NewVectorClockSet(3, map[int]int{1: 2, 2: 10, 3: 2}),
		3: clock.NewVectorClockSet(3, map[int]int{1: 2, 2: 9, 3: 2}),
	}
	expectedCurrentlyWaiting := map[int][]int{
		123: {1, 2},
	}

	CondWait(&co1, vc)
	CondWait(&co2, vc)

	t.Run("VC", func(t *testing.T) {
		if !reflect.DeepEqual(vc, expectedVc) {
			t.Errorf("Incorrect result. Expected %v. Got %v.", expectedVc, vc)
		}
	})

	t.Run("Currently Waiting", func(t *testing.T) {
		if !reflect.DeepEqual(currentlyWaiting, expectedCurrentlyWaiting) {
			t.Errorf("Incorrect result. Expected %v. Got %v.", expectedCurrentlyWaiting, currentlyWaiting)
		}
	})
}

func TestSignal(t *testing.T) {
	co1 := TraceElementCond{
		routine: 1,
		tPre:    5,
		tPost:   6,
		id:      123,
		opC:     WaitCondOp,
		pos:     "testfile.go:999",
		vc:      clock.NewVectorClock(2),
	}

	co2 := TraceElementCond{
		routine: 2,
		tPre:    7,
		tPost:   8,
		id:      123,
		opC:     WaitCondOp,
		pos:     "testfile.go:999",
		vc:      clock.NewVectorClock(2),
	}

	co3 := TraceElementCond{
		routine: 3,
		tPre:    9,
		tPost:   10,
		id:      123,
		opC:     SignalOp,
		pos:     "testfile.go:999",
		vc:      clock.NewVectorClock(2),
	}

	currentlyWaiting = map[int][]int{}

	vc := map[int]clock.VectorClock{
		1: clock.NewVectorClockSet(3, map[int]int{1: 5, 2: 6, 3: 2}),
		2: clock.NewVectorClockSet(3, map[int]int{1: 2, 2: 9, 3: 2}),
		3: clock.NewVectorClockSet(3, map[int]int{1: 2, 2: 9, 3: 2}),
	}

	expectedVc := map[int]clock.VectorClock{
		1: clock.NewVectorClockSet(3, map[int]int{1: 6, 2: 9, 3: 2}),
		2: clock.NewVectorClockSet(3, map[int]int{1: 2, 2: 10, 3: 2}),
		3: clock.NewVectorClockSet(3, map[int]int{1: 2, 2: 9, 3: 3}),
	}
	expectedCurrentlyWaiting := map[int][]int{
		123: {2},
	}

	CondWait(&co1, vc)
	CondWait(&co2, vc)
	CondSignal(&co3, vc)

	t.Run("VC", func(t *testing.T) {
		if !reflect.DeepEqual(vc, expectedVc) {
			t.Errorf("Incorrect result. Expected %v. Got %v.", expectedVc, vc)
		}
	})

	t.Run("Currently Waiting", func(t *testing.T) {
		if !reflect.DeepEqual(currentlyWaiting, expectedCurrentlyWaiting) {
			t.Errorf("Incorrect result. Expected %v. Got %v.", expectedCurrentlyWaiting, currentlyWaiting)
		}
	})
}

func TestBroadcast(t *testing.T) {
	co1 := TraceElementCond{
		routine: 1,
		tPre:    5,
		tPost:   6,
		id:      123,
		opC:     WaitCondOp,
		pos:     "testfile.go:999",
		vc:      clock.NewVectorClock(2),
	}

	co2 := TraceElementCond{
		routine: 2,
		tPre:    7,
		tPost:   8,
		id:      123,
		opC:     WaitCondOp,
		pos:     "testfile.go:999",
		vc:      clock.NewVectorClock(2),
	}

	co3 := TraceElementCond{
		routine: 3,
		tPre:    9,
		tPost:   10,
		id:      123,
		opC:     BroadcastOp,
		pos:     "testfile.go:999",
		vc:      clock.NewVectorClock(2),
	}

	currentlyWaiting = map[int][]int{}

	vc := map[int]clock.VectorClock{
		1: clock.NewVectorClockSet(3, map[int]int{1: 5, 2: 6, 3: 2}),
		2: clock.NewVectorClockSet(3, map[int]int{1: 2, 2: 9, 3: 2}),
		3: clock.NewVectorClockSet(3, map[int]int{1: 2, 2: 9, 3: 3}),
	}

	expectedVc := map[int]clock.VectorClock{
		1: clock.NewVectorClockSet(3, map[int]int{1: 6, 2: 9, 3: 3}),
		2: clock.NewVectorClockSet(3, map[int]int{1: 2, 2: 10, 3: 3}),
		3: clock.NewVectorClockSet(3, map[int]int{1: 2, 2: 9, 3: 4}),
	}
	expectedCurrentlyWaiting := map[int][]int{
		123: {},
	}

	CondWait(&co1, vc)
	CondWait(&co2, vc)
	CondBroadcast(&co3, vc)

	t.Run("VC", func(t *testing.T) {
		if !reflect.DeepEqual(vc, expectedVc) {
			t.Errorf("Incorrect result. Expected %v. Got %v.", expectedVc, vc)
		}
	})

	t.Run("Currently Waiting", func(t *testing.T) {
		if !reflect.DeepEqual(currentlyWaiting, expectedCurrentlyWaiting) {
			t.Errorf("Incorrect result. Expected %v. Got %v.", expectedCurrentlyWaiting, currentlyWaiting)
		}
	})
}
