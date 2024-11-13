// Copyright (c) 2024 Erik Kassubek
//
// File: vcOnce_test.go
// Brief: Tests for vcOnce
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

func TestOnceDo(t *testing.T) {
	on1 := TraceElementOnce{
		routine: 1,
		tPre:    10,
		tPost:   11,
		id:      123,
		suc:     true,
		pos:     "testfile.go:123",
		vc:      clock.NewVectorClock(2),
	}

	on2 := TraceElementOnce{
		routine: 2,
		tPre:    13,
		tPost:   14,
		id:      123,
		suc:     false,
		pos:     "testfile.go:123",
		vc:      clock.NewVectorClock(2),
	}

	vc := map[int]clock.VectorClock{
		1: clock.NewVectorClockSet(2, map[int]int{1: 5, 2: 6}),
		2: clock.NewVectorClockSet(2, map[int]int{1: 2, 2: 9}),
	}

	vcExp1 := map[int]clock.VectorClock{
		1: clock.NewVectorClockSet(2, map[int]int{1: 6, 2: 6}),
		2: clock.NewVectorClockSet(2, map[int]int{1: 2, 2: 9}),
	}

	vcExp2 := map[int]clock.VectorClock{
		1: clock.NewVectorClockSet(2, map[int]int{1: 6, 2: 6}),
		2: clock.NewVectorClockSet(2, map[int]int{1: 5, 2: 10}),
	}

	vcOSuc := clock.NewVectorClockSet(2, map[int]int{1: 5, 2: 6})

	DoSuc(&on1, vc)

	t.Run("Suc", func(t *testing.T) {
		if !reflect.DeepEqual(vc, vcExp1) {
			t.Errorf("Incorrect result. Expected %v. Got %v.", vcExp1, vc)
		}
	})

	t.Run("vcOSuc", func(t *testing.T) {
		if !reflect.DeepEqual(oSuc[on1.id], vcOSuc) {
			t.Errorf("Incorrect result. Expected %v. Got %v.", vcOSuc, oSuc[on1.id])
		}
	})

	DoFail(&on2, vc)

	t.Run("Fail", func(t *testing.T) {
		if !reflect.DeepEqual(vc, vcExp2) {
			t.Errorf("Incorrect result. Expected %v. Got %v.", vcExp2, vc)
		}
	})

	t.Run("vcOFail", func(t *testing.T) {
		if !reflect.DeepEqual(oSuc[on2.id], vcOSuc) {
			t.Errorf("Incorrect result. Expected %v. Got %v.", vcOSuc, oSuc[on1.id])
		}
	})

}
