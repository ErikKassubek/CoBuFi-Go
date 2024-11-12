// Copyright (c) 2024 Erik Kassubek
//
// File: utils_test.go
// Brief: Test for utils.go
//
// Author: Erik Kassubek
// Created: 2024-11-11
//
// License: BSD-3-Clause

package utils

import (
	"reflect"
	"testing"
)

func TestContainsString(t *testing.T) {
	var tests = []struct {
		name     string
		list     []string
		char     string
		expected bool
	}{
		{"Contains true", []string{"a", "b", "c"}, "a", true},
		{"Contains false", []string{"a", "b", "c"}, "d", false},
		{"Contains empty", []string{}, "", false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res := ContainsString(test.list, test.char)
			if res != test.expected {
				t.Errorf("Incorrect result for contains(%v, %s). Expected %t. Got %t", test.list, test.char, test.expected, res)
			}
		})
	}
}

func TestContainsInt(t *testing.T) {
	var tests = []struct {
		name     string
		list     []int
		number   int
		expected bool
	}{
		{"Contains true", []int{1, 2, 3}, 1, true},
		{"Contains false", []int{1, 2, 3}, 4, false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res := ContainsInt(test.list, test.number)
			if res != test.expected {
				t.Errorf("Incorrect result. Expected %t. Got %t", test.expected, res)
			}
		})
	}
}

func TestSplitAtLast(t *testing.T) {
	var tests = []struct {
		name     string
		str      string
		sep      string
		expected []string
	}{
		{"Split middle short", "abc", "b", []string{"a", "c"}},
		{"Split middle long", "abcdefg", "d", []string{"abc", "efg"}},
		{"Split at beginning", "abcdefg", "a", []string{"", "bcdefg"}},
		{"Split at end", "abcdefg", "g", []string{"abcdef", ""}},
		{"Split at not existing", "abcdefg", "x", []string{"abcdefg"}},
		{"Split at empty separator", "abcdefg", "", []string{"abcdefg"}},
		{"Split empty sting", "", "d", []string{""}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res := SplitAtLast(test.str, test.sep)
			if !reflect.DeepEqual(res, test.expected) {
				t.Errorf("Incorrect result for TestSplitAtLast(%s, %s). Expected %v. Got %v", test.str, test.sep, test.expected, res)
			}
		})
	}
}
