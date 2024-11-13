// Copyright (c) 2024 Erik Kassubek
//
// File: analysisUtil_test.go
// Brief: Tests for analysisUtil
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

func TestInfoFromTID(t *testing.T) {
	var tests = []struct {
		name     string
		tID      string
		expFile  string
		expLine  int
		expTPre  int
		expError error
	}{
		{"Correct tid 1", "testfile.go:123@444", "testfile.go", 123, 444, nil},
		{"Correct tid 2", "testfile.go:1@1", "testfile.go", 1, 1, nil},
		{"Missing tPre", "testfile.go:1", "", 0, 0, errors.New("TID not correct: no @: testfile.go:1")},
		{"Missing file", ":1@1", "", 1, 1, nil},
		{"Missing line", "testfile.go@1", "", 0, 0, errors.New("TID not correct: no ':': testfile.go@1")},
		{"tPre not int", "testfile.go:1@a", "", 0, 0, errors.New("strconv.Atoi: parsing \"a\": invalid syntax")},
		{"line not int", "testfile.go:a@1", "", 0, 0, errors.New("strconv.Atoi: parsing \"a\": invalid syntax")},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			file, line, tPre, err := infoFromTID(test.tID)

			if file != test.expFile {
				t.Errorf("Incorrect result for infoFromTID(%s) file. Expected %s. Got %s.",
					test.tID, test.expFile, file)
			}

			if line != test.expLine {
				t.Errorf("Incorrect result for infoFromTID(%s) line. Expected %d. Got %d.",
					test.tID, test.expLine, line)
			}

			if tPre != test.expTPre {
				t.Errorf("Incorrect result for infoFromTID(%s) tPre. Expected %d. Got %d.",
					test.tID, test.expTPre, tPre)
			}

			if err == nil && test.expError != nil {
				t.Errorf("Incorrect result for infoFromTID(%s) error. Expected %s. Got nil.",
					test.tID, test.expError.Error())
			}

			if err != nil && test.expError == nil {
				t.Errorf("Incorrect result for infoFromTID(%s) error. Expected nil. Got %s.",
					test.tID, err.Error())
			}
		})
	}
}

func TestSameRoutine(t *testing.T) {
	tests := []struct {
		name   string
		elems  [][]TraceElement
		expect bool
	}{
		{
			name: "All elements have the same routine",
			elems: [][]TraceElement{
				{&TraceElementAtomic{routine: 1}, &TraceElementAtomic{routine: 1}},
				{&TraceElementAtomic{routine: 1}, &TraceElementAtomic{routine: 1}},
			},
			expect: true,
		},
		{
			name: "Different routines in the same position",
			elems: [][]TraceElement{
				{&TraceElementAtomic{routine: 1}, &TraceElementAtomic{routine: 2}},
				{&TraceElementAtomic{routine: 1}, &TraceElementAtomic{routine: 3}},
			},
			expect: false,
		},
		{
			name: "Different routines in different positions",
			elems: [][]TraceElement{
				{&TraceElementAtomic{routine: 1}, &TraceElementAtomic{routine: 2}},
				{&TraceElementAtomic{routine: 3}, &TraceElementAtomic{routine: 1}},
			},
			expect: false,
		},
		{
			name: "Empty input",
			elems: [][]TraceElement{
				{},
				{},
			},
			expect: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := sameRoutine(test.elems...)
			if result != test.expect {
				t.Errorf("expected %v, got %v", test.expect, result)
			}
		})
	}
}
