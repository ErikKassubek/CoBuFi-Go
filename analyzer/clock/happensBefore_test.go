// Copyright (c) 2024 Erik Kassubek
//
// File: happensBefore_test.go
// Brief: Test for happensBefore.go
//
// Author: Erik Kassubek
// Created: 2024-11-11
//
// License: BSD-3-Clause

package clock

import "testing"

func TestIsCause(t *testing.T) {
	v1 := NewVectorClock(3)
	v1.clock = map[int]int{1: 1, 2: 2, 3: 3}

	v2 := NewVectorClock(3)
	v2.clock = map[int]int{1: 2, 2: 3, 3: 4}

	v3 := NewVectorClock(3)
	v3.clock = map[int]int{1: 3, 2: 2, 3: 1}

	var tests = []struct {
		name     string
		v1       VectorClock
		v2       VectorClock
		expected bool
	}{
		{"True", v1, v2, true},
		{"Reversed", v2, v1, false},
		{"Concurrent 1", v1, v3, false},
		{"Concurrent 2", v3, v1, false},
		{"Equal", v1, v1, false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res := isCause(test.v1, test.v2)
			if res != test.expected {
				t.Errorf("Incorrect result in isCause(%s, %s). Expected %t. Got %t.", test.v1.ToString(), test.v2.ToString(), test.expected, res)
			}
		})
	}
}

func TestGetHappensBefore(t *testing.T) {
	v1 := NewVectorClock(3)
	v1.clock = map[int]int{1: 1, 2: 2, 3: 3}

	v2 := NewVectorClock(3)
	v2.clock = map[int]int{1: 2, 2: 3, 3: 4}

	v3 := NewVectorClock(3)
	v3.clock = map[int]int{1: 3, 2: 2, 3: 1}

	v4 := NewVectorClock(2)

	var tests = []struct {
		name     string
		v1       VectorClock
		v2       VectorClock
		expected HappensBefore
	}{
		{"Before", v1, v2, Before},
		{"After", v2, v1, After},
		{"Concurrent 1", v1, v3, Concurrent},
		{"Concurrent 2", v3, v1, Concurrent},
		{"Equal", v1, v1, Concurrent},
		{"Different length", v1, v4, None}, // None if size is not equal
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res := GetHappensBefore(test.v1, test.v2)
			if res != test.expected {
				t.Errorf("Incorrect result in GetHappensBefore(%s, %s). Expected %v. Got %v.", test.v1.ToString(), test.v2.ToString(), test.expected, res)
			}
		})
	}
}
