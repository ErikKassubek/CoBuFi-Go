// Copyright (c) 2024 Erik Kassubek
//
// File: vcMutex_test.go
// Brief: Tests for vcMutex.go
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

func TestLock(t *testing.T) {
	mu := TraceElementMutex{
		routine: 1,
		tPre:    4,
		tPost:   5,
		id:      123,
		rw:      false,
		opM:     LockOp,
		suc:     true,
		pos:     "testfile.go:999",
		vc:      clock.NewVectorClock(2),
	}

	vc := map[int]clock.VectorClock{
		1: clock.NewVectorClockSet(3, map[int]int{1: 1, 2: 5, 3: 9}),
		2: clock.NewVectorClockSet(3, map[int]int{1: 1, 2: 5, 3: 9}),
		3: clock.NewVectorClockSet(3, map[int]int{1: 1, 2: 5, 3: 9}),
	}
	vcw := map[int]clock.VectorClock{
		1: clock.NewVectorClockSet(3, map[int]int{1: 1, 2: 5, 3: 9}),
		2: clock.NewVectorClockSet(3, map[int]int{1: 1, 2: 5, 3: 9}),
		3: clock.NewVectorClockSet(3, map[int]int{1: 1, 2: 5, 3: 9}),
	} // not tested

	relR = map[int]clock.VectorClock{
		123: clock.NewVectorClockSet(3, map[int]int{1: 3, 2: 4, 3: 7}),
	}
	relW = map[int]clock.VectorClock{
		123: clock.NewVectorClockSet(3, map[int]int{1: 2, 2: 7, 3: 8}),
	}

	expectedRelR := relR[123].Copy()
	expectedRelW := relW[123].Copy()

	Lock(&mu, vc, vcw)

	expectedVc := map[int]clock.VectorClock{
		1: clock.NewVectorClockSet(3, map[int]int{1: 4, 2: 7, 3: 9}),
		2: clock.NewVectorClockSet(3, map[int]int{1: 1, 2: 5, 3: 9}),
		3: clock.NewVectorClockSet(3, map[int]int{1: 1, 2: 5, 3: 9}),
	}

	if !reflect.DeepEqual(vc, expectedVc) {
		t.Errorf("Expected %v, got %v", expectedVc, vc[1])
	}

	if !reflect.DeepEqual(relR[123], expectedRelR) {
		t.Errorf("Expected %v, got %v", expectedRelR, relR[123])
	}

	if !reflect.DeepEqual(relW[123], expectedRelW) {
		t.Errorf("Expected %v, got %v", expectedRelW, relW[123])
	}
}

func TestUnlock(t *testing.T) {
	mu := TraceElementMutex{
		routine: 1,
		tPre:    4,
		tPost:   5,
		id:      123,
		rw:      false,
		opM:     UnlockOp,
		suc:     true,
		pos:     "testfile.go:999",
		vc:      clock.NewVectorClock(2),
	}

	vc := map[int]clock.VectorClock{
		1: clock.NewVectorClockSet(3, map[int]int{1: 1, 2: 5, 3: 9}),
		2: clock.NewVectorClockSet(3, map[int]int{1: 1, 2: 5, 3: 9}),
		3: clock.NewVectorClockSet(3, map[int]int{1: 1, 2: 5, 3: 9}),
	}

	relR = map[int]clock.VectorClock{
		123: clock.NewVectorClockSet(3, map[int]int{1: 3, 2: 4, 3: 7}),
	}
	relW = map[int]clock.VectorClock{
		123: clock.NewVectorClockSet(3, map[int]int{1: 2, 2: 7, 3: 8}),
	}

	expectedVc := map[int]clock.VectorClock{
		1: clock.NewVectorClockSet(3, map[int]int{1: 2, 2: 5, 3: 9}),
		2: clock.NewVectorClockSet(3, map[int]int{1: 1, 2: 5, 3: 9}),
		3: clock.NewVectorClockSet(3, map[int]int{1: 1, 2: 5, 3: 9}),
	}

	expectedRelR := vc[1].Copy()
	expectedRelW := vc[1].Copy()

	Unlock(&mu, vc)

	if !reflect.DeepEqual(vc, expectedVc) {
		t.Errorf("Expected %v, got %v", expectedVc, vc[1])
	}

	if !reflect.DeepEqual(relR[123], expectedRelR) {
		t.Errorf("Expected %v, got %v", expectedRelR, relR[123])
	}

	if !reflect.DeepEqual(relW[123], expectedRelW) {
		t.Errorf("Expected %v, got %v", expectedRelW, relW[123])
	}
}

func TestRLock(t *testing.T) {
	mu := TraceElementMutex{
		routine: 1,
		tPre:    4,
		tPost:   5,
		id:      123,
		rw:      true,
		opM:     RLockOp,
		suc:     true,
		pos:     "testfile.go:999",
		vc:      clock.NewVectorClock(2),
	}

	vc := map[int]clock.VectorClock{
		1: clock.NewVectorClockSet(3, map[int]int{1: 1, 2: 5, 3: 9}),
		2: clock.NewVectorClockSet(3, map[int]int{1: 1, 2: 5, 3: 9}),
		3: clock.NewVectorClockSet(3, map[int]int{1: 1, 2: 5, 3: 9}),
	}
	vcw := map[int]clock.VectorClock{
		1: clock.NewVectorClockSet(3, map[int]int{1: 1, 2: 5, 3: 9}),
		2: clock.NewVectorClockSet(3, map[int]int{1: 1, 2: 5, 3: 9}),
		3: clock.NewVectorClockSet(3, map[int]int{1: 1, 2: 5, 3: 9}),
	} // not tested

	relR = map[int]clock.VectorClock{
		123: clock.NewVectorClockSet(3, map[int]int{1: 3, 2: 4, 3: 7}),
	}
	relW = map[int]clock.VectorClock{
		123: clock.NewVectorClockSet(3, map[int]int{1: 2, 2: 7, 3: 8}),
	}

	expectedRelR := relR[123].Copy()
	expectedRelW := relW[123].Copy()

	RLock(&mu, vc, vcw)

	expectedVc := map[int]clock.VectorClock{
		1: clock.NewVectorClockSet(3, map[int]int{1: 3, 2: 7, 3: 9}),
		2: clock.NewVectorClockSet(3, map[int]int{1: 1, 2: 5, 3: 9}),
		3: clock.NewVectorClockSet(3, map[int]int{1: 1, 2: 5, 3: 9}),
	}

	if !reflect.DeepEqual(vc, expectedVc) {
		t.Errorf("Expected %v, got %v", expectedVc, vc[1])
	}

	if !reflect.DeepEqual(relR[123], expectedRelR) {
		t.Errorf("Expected %v, got %v", expectedRelR, relR[123])
	}

	if !reflect.DeepEqual(relW[123], expectedRelW) {
		t.Errorf("Expected %v, got %v", expectedRelW, relW[123])
	}
}

func TestRUnlock(t *testing.T) {
	mu := TraceElementMutex{
		routine: 1,
		tPre:    4,
		tPost:   5,
		id:      123,
		rw:      true,
		opM:     RUnlockOp,
		suc:     true,
		pos:     "testfile.go:999",
		vc:      clock.NewVectorClock(2),
	}

	vc := map[int]clock.VectorClock{
		1: clock.NewVectorClockSet(3, map[int]int{1: 1, 2: 5, 3: 9}),
		2: clock.NewVectorClockSet(3, map[int]int{1: 1, 2: 5, 3: 9}),
		3: clock.NewVectorClockSet(3, map[int]int{1: 1, 2: 5, 3: 9}),
	}

	relR = map[int]clock.VectorClock{
		123: clock.NewVectorClockSet(3, map[int]int{1: 3, 2: 4, 3: 7}),
	}
	relW = map[int]clock.VectorClock{
		123: clock.NewVectorClockSet(3, map[int]int{1: 2, 2: 7, 3: 8}),
	}

	expectedVc := map[int]clock.VectorClock{
		1: clock.NewVectorClockSet(3, map[int]int{1: 2, 2: 5, 3: 9}),
		2: clock.NewVectorClockSet(3, map[int]int{1: 1, 2: 5, 3: 9}),
		3: clock.NewVectorClockSet(3, map[int]int{1: 1, 2: 5, 3: 9}),
	}

	expectedRelR := clock.NewVectorClockSet(3, map[int]int{1: 3, 2: 5, 3: 9})
	expectedRelW := relW[123].Copy()

	RUnlock(&mu, vc)

	t.Run("VC", func(t *testing.T) {
		if !reflect.DeepEqual(vc, expectedVc) {
			t.Errorf("Expected %v, got %v", expectedVc, vc[1])
		}
	})

	t.Run("relR", func(t *testing.T) {
		if !reflect.DeepEqual(relR[123], expectedRelR) {
			t.Errorf("Expected %v, got %v", expectedRelR, relR[123])
		}
	})

	t.Run("relW", func(t *testing.T) {
		if !reflect.DeepEqual(relW[123], expectedRelW) {
			t.Errorf("Expected %v, got %v", expectedRelW, relW[123])
		}
	})
}
