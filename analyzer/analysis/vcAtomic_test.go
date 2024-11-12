// Copyright (c) 2024 Erik Kassubek
//
// File: vcAtomic_test.go
// Brief: Tests for vcAtomic.go
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

func TestNw(t *testing.T) {
	var tests = []struct {
		name       string
		index      []int
		nRout      []int
		expectedLW map[int]clock.VectorClock
	}{
		{"New lw", []int{1, 2, 3}, []int{5, 3, 7}, map[int]clock.VectorClock{1: clock.NewVectorClock(5), 2: clock.NewVectorClock(3), 3: clock.NewVectorClock(7)}},
		{"Existing lw", []int{1, 2, 1}, []int{5, 3, 7}, map[int]clock.VectorClock{1: clock.NewVectorClock(5), 2: clock.NewVectorClock(3)}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			lw = make(map[int]clock.VectorClock)

			for i, index := range test.index {
				newLw(index, test.nRout[i])
			}

			if !reflect.DeepEqual(lw, test.expectedLW) {
				t.Errorf("Incorrect lw. Expected %v. Got %v.", test.expectedLW, lw)
			}
		})
	}
}

func TestWrite(t *testing.T) {
	t.Run("Write", func(t *testing.T) {
		at := TraceElementAtomic{id: 1, routine: 2}
		vc := map[int]clock.VectorClock{2: clock.NewVectorClockSet(3, map[int]int{1: 2, 2: 4, 3: 1})}
		lw = map[int]clock.VectorClock{1: clock.NewVectorClockSet(3, map[int]int{1: 2, 2: 3, 3: 0})}

		expectedLW := map[int]clock.VectorClock{1: vc[2].Copy()}
		expectedVC := map[int]clock.VectorClock{2: clock.NewVectorClockSet(3, map[int]int{1: 2, 2: 5, 3: 1})}

		Write(&at, vc)
		if !reflect.DeepEqual(lw, expectedLW) {
			t.Errorf("Incorrect lw. Expected %v. Got %v.", expectedLW, lw)
		}

		if !reflect.DeepEqual(vc, expectedVC) {
			t.Errorf("Incorrect vc. Expected %v. Got %v.", expectedVC, vc)
		}
	})
}

func TestRead(t *testing.T) {
	var tests = []struct {
		name       string
		sync       bool
		expectedVC map[int]clock.VectorClock
	}{
		{"No sync", false, map[int]clock.VectorClock{1: clock.NewVectorClock(2), 2: clock.NewVectorClockSet(2, map[int]int{1: 0, 2: 1})}},
		{"Sync", true, map[int]clock.VectorClock{1: clock.NewVectorClock(2), 2: clock.NewVectorClockSet(2, map[int]int{1: 2, 2: 4})}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			at := TraceElementAtomic{id: 1, routine: 2}
			vc := map[int]clock.VectorClock{1: clock.NewVectorClock(2), 2: clock.NewVectorClock(2)}
			lw = map[int]clock.VectorClock{1: clock.NewVectorClockSet(2, map[int]int{1: 2, 2: 3})}

			Read(&at, vc, test.sync)

			if !reflect.DeepEqual(vc, test.expectedVC) {
				t.Errorf("Incorrect vc. Expected %v. Got %v.", test.expectedVC, vc)
			}
		})
	}
}

func TestSwap(t *testing.T) {
	var tests = []struct {
		name       string
		sync       bool
		expectedVC map[int]clock.VectorClock
	}{
		{"No sync", false, map[int]clock.VectorClock{1: clock.NewVectorClock(2), 2: clock.NewVectorClockSet(2, map[int]int{1: 0, 2: 2})}},
		{"Sync", true, map[int]clock.VectorClock{1: clock.NewVectorClock(2), 2: clock.NewVectorClockSet(2, map[int]int{1: 2, 2: 5})}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			at := TraceElementAtomic{id: 1, routine: 2}
			vc := map[int]clock.VectorClock{1: clock.NewVectorClock(2), 2: clock.NewVectorClock(2)}
			lw = map[int]clock.VectorClock{1: clock.NewVectorClockSet(2, map[int]int{1: 2, 2: 3})}

			Swap(&at, vc, test.sync)

			if !reflect.DeepEqual(vc, test.expectedVC) {
				t.Errorf("Incorrect vc. Expected %v. Got %v.", test.expectedVC, vc)
			}
		})
	}
}
