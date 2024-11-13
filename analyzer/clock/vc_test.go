// Copyright (c) 2024 Erik Kassubek
//
// File: vc_test.go
// Brief: Tests for vc.go
//
// Author: Erik Kassubek
// Created: 2024-11-11
//
// License: BSD-3-Clause

package clock

import (
	"reflect"
	"testing"
)

func TestVcNew(t *testing.T) {
	var tests = []struct {
		name    string
		size    int
		sizeExp int
		clock   map[int]int
	}{
		{"Valid clock", 5, 5, map[int]int{1: 0, 2: 0, 3: 0, 4: 0, 5: 0}},
		{"Empty clock", 0, 0, map[int]int{}},
		{"Negative length", -1, 0, map[int]int{}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			v := NewVectorClock(test.size)

			if v.GetSize() != test.sizeExp {
				t.Errorf("Incorrect size. Expected %d. Got %d.", test.sizeExp, v.size)
			}

			if !reflect.DeepEqual(v.clock, test.clock) {
				t.Errorf("Incorrect VC. Expected %v. Got %v.", test.clock, v.clock)
			}
		})
	}
}

func TestVectorClockSet(t *testing.T) {
	var tests = []struct {
		name     string
		size     int
		clock    map[int]int
		expClock map[int]int
	}{
		{"Valid clock", 5, map[int]int{1: 2, 2: 6, 3: 1, 4: 0, 5: 99}, map[int]int{1: 2, 2: 6, 3: 1, 4: 0, 5: 99}},
		{"Empty clock", 0, map[int]int{}, map[int]int{}},
		{"Size to small", 1, map[int]int{1: 2, 2: 3}, map[int]int{1: 2}},
		{"Missing value", 3, map[int]int{1: 2, 3: 3}, map[int]int{1: 2, 2: 0, 3: 3}},
		{"Negative length", -1, map[int]int{}, map[int]int{}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			v := NewVectorClockSet(test.size, test.clock)

			if v.size != test.size {
				if !reflect.DeepEqual(v.clock, test.expClock) {
					t.Errorf("Incorrect VC. Expected %v. Got %v.", test.expClock, v.clock)
				}
			}
		})
	}
}

func TestGet(t *testing.T) {
	v := NewVectorClock(3)
	v.clock = map[int]int{1: 3, 2: 1, 3: 6}

	t.Run("Get string", func(t *testing.T) {
		expectString := "[3, 1, 6]"
		s := v.ToString()
		if s != expectString {
			t.Errorf("Incorrect vc string. Expected %s. Got %s", expectString, s)
		}
	})

	t.Run("Get clock", func(t *testing.T) {
		expectedClock := map[int]int{
			1: 3, 2: 1, 3: 6,
		}
		c := v.GetClock()
		if !reflect.DeepEqual(c, expectedClock) {
			t.Errorf("Incorrect vc clock. Expected %v. Got %v.", expectedClock, c)
		}
	})
}

func TestInc(t *testing.T) {
	v := NewVectorClock(3)
	v.clock = map[int]int{1: 3, 2: 2, 3: 1}
	expectString := "[3, 4, 2]"

	t.Run("Increment", func(t *testing.T) {
		v.Inc(2)
		v.Inc(2)
		v.Inc(3)

		res := v.ToString()
		if res != expectString {
			t.Errorf("Incorrect vc clock. Expected %s. Got %s.", expectString, res)
		}
	})

	t.Run("Increment invalid", func(t *testing.T) {
		v.Inc(6) // 4 > len(clock) -> do nothing
		res := v.ToString()
		if res != expectString {
			t.Errorf("Incorrect vc clock. Expected %s. Got %s.", expectString, res)
		}
	})
}

func TestSync(t *testing.T) {
	v1 := NewVectorClock(3)
	v1.clock = map[int]int{1: 3, 2: 2, 3: 1}

	v2 := NewVectorClock(3)
	v2.clock = map[int]int{1: 1, 2: 2, 3: 4}

	t.Run("Sync", func(t *testing.T) {
		v := v1.Sync(v2)

		expectedV := "[3, 2, 4]"
		expectedV1 := "[3, 2, 1]"
		expectedV2 := "[1, 2, 4]"

		if v.ToString() != expectedV {
			t.Errorf("Incorrect value for sync v. Expected %s. Got %s.", expectedV, v.ToString())
		}

		if v1.ToString() != expectedV1 {
			t.Errorf("Incorrect value for sync v1. Expected %s. Got %s.", expectedV1, v1.ToString())
		}

		if v2.ToString() != expectedV2 {
			t.Errorf("Incorrect value for sync v2. Expected %s. Got %s.", expectedV, v2.ToString())
		}
	})
}

func TestCopy(t *testing.T) {
	v := NewVectorClock(3)
	v.clock = map[int]int{1: 1, 2: 2, 3: 3}

	t.Run("Copy", func(t *testing.T) {

		c := v.Copy()

		v.Inc(1)
		v.Inc(2)
		v.Inc(3)

		expectedV := "[2, 3, 4]"
		expectedC := "[1, 2, 3]"

		if expectedV != v.ToString() {
			t.Errorf("Incorrect value in copy v. Expected %s. Got %s.", expectedV, v.ToString())
		}

		if expectedC != c.ToString() {
			t.Errorf("Incorrect value in copy c. Expected %s. Got %s.", expectedV, c.ToString())
		}
	})
}

func TestIsEqual(t *testing.T) {
	tests := []struct {
		name     string
		vc1      VectorClock
		vc2      VectorClock
		expected bool
	}{
		{
			name: "Equal vector clocks",
			vc1: VectorClock{
				size:  3,
				clock: map[int]int{1: 1, 2: 2, 3: 3},
			},
			vc2: VectorClock{
				size:  3,
				clock: map[int]int{1: 1, 2: 2, 3: 3},
			},
			expected: true,
		},
		{
			name: "Different sizes",
			vc1: VectorClock{
				size:  3,
				clock: map[int]int{1: 1, 2: 2, 3: 3},
			},
			vc2: VectorClock{
				size:  2,
				clock: map[int]int{1: 1, 2: 2},
			},
			expected: false,
		},
		{
			name: "Different values",
			vc1: VectorClock{
				size:  3,
				clock: map[int]int{1: 1, 2: 2, 3: 3},
			},
			vc2: VectorClock{
				size:  3,
				clock: map[int]int{1: 1, 2: 2, 3: 4},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.vc1.IsEqual(tt.vc2)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestIsMapVcEqual(t *testing.T) {
	tests := []struct {
		name     string
		v1       map[int]VectorClock
		v2       map[int]VectorClock
		expected bool
	}{
		{
			name: "Equal maps",
			v1: map[int]VectorClock{
				1: {size: 3, clock: map[int]int{1: 1, 2: 2, 3: 3}},
				2: {size: 3, clock: map[int]int{1: 4, 2: 5, 3: 6}},
			},
			v2: map[int]VectorClock{
				1: {size: 3, clock: map[int]int{1: 1, 2: 2, 3: 3}},
				2: {size: 3, clock: map[int]int{1: 4, 2: 5, 3: 6}},
			},
			expected: true,
		},
		{
			name: "Different sizes",
			v1: map[int]VectorClock{
				1: {size: 3, clock: map[int]int{1: 1, 2: 2, 3: 3}},
			},
			v2: map[int]VectorClock{
				1: {size: 3, clock: map[int]int{1: 1, 2: 2, 3: 3}},
				2: {size: 3, clock: map[int]int{1: 4, 2: 5, 3: 6}},
			},
			expected: false,
		},
		{
			name: "Different values",
			v1: map[int]VectorClock{
				1: {size: 3, clock: map[int]int{1: 1, 2: 2, 3: 3}},
				2: {size: 3, clock: map[int]int{1: 4, 2: 5, 3: 6}},
			},
			v2: map[int]VectorClock{
				1: {size: 3, clock: map[int]int{1: 1, 2: 2, 3: 3}},
				2: {size: 3, clock: map[int]int{1: 4, 2: 5, 3: 7}},
			},
			expected: false,
		},
		{
			name: "Different keys",
			v1: map[int]VectorClock{
				1: {size: 3, clock: map[int]int{1: 1, 2: 2, 3: 3}},
			},
			v2: map[int]VectorClock{
				2: {size: 3, clock: map[int]int{1: 1, 2: 2, 3: 3}},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsMapVcEqual(tt.v1, tt.v2)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}
