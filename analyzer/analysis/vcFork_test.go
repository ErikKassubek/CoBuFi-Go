// Copyright (c) 2024 Erik Kassubek
//
// File: vcFork_test.go
// Brief: Tests for vcFork.go
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

func TestFork(t *testing.T) {
	fo := TraceElementFork{
		routine: 1,
		tPost:   5,
		id:      2, // new routine
		pos:     "testfile:999",
		vc:      clock.NewVectorClockSet(2, map[int]int{1: 5, 2: 6, 3: 7}),
	}

	vc := map[int]clock.VectorClock{
		1: clock.NewVectorClockSet(2, map[int]int{1: 5, 2: 6}),
		2: clock.NewVectorClockSet(2, map[int]int{1: 0, 2: 0}),
	}

	vc2 := map[int]clock.VectorClock{
		1: clock.NewVectorClockSet(2, map[int]int{1: 5, 2: 6}),
		2: clock.NewVectorClockSet(2, map[int]int{1: 0, 2: 0}),
	}

	Fork(&fo, vc, vc2)

	expextedVcs := map[int]clock.VectorClock{
		1: clock.NewVectorClockSet(2, map[int]int{1: 6, 2: 6}),
		2: clock.NewVectorClockSet(2, map[int]int{1: 5, 2: 7}),
	}

	if !reflect.DeepEqual(vc, expextedVcs) {
		t.Errorf("Expected %v, got %v", expextedVcs, vc)
	}

}
