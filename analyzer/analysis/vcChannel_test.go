// Copyright (c) 2024 Erik Kassubek
//
// File: vcChannel_test.go
// Brief: Tests for vcChannel.go
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

func TestUnbuffered(t *testing.T) {
	sender := TraceElementChannel{
		routine: 1,
		tPre:    5,
		tPost:   6,
		id:      123,
		opC:     SendOp,
		cl:      false,
		oID:     456,
		qSize:   0,
		pos:     "testfile:999",
		sel:     nil,
		partner: nil,
		vc:      clock.NewVectorClockSet(3, map[int]int{1: 5, 2: 6, 3: 7}),
	}

	recv := TraceElementChannel{
		routine: 2,
		tPre:    4,
		tPost:   7,
		id:      123,
		opC:     RecvOp,
		cl:      false,
		oID:     456,
		qSize:   0,
		pos:     "testfile:888",
		sel:     nil,
		partner: &sender,
		vc:      clock.NewVectorClockSet(3, map[int]int{1: 2, 2: 9, 3: 4}),
	}

	sender.partner = &recv

	vc := map[int]clock.VectorClock{
		1: clock.NewVectorClockSet(3, map[int]int{1: 5, 2: 6, 3: 7}),
		2: clock.NewVectorClockSet(3, map[int]int{1: 2, 2: 9, 3: 4}),
		3: clock.NewVectorClockSet(3, map[int]int{1: 5, 2: 6, 3: 9}),
	}

	expectedVcSend := vc[1].Copy()
	expectedVcRecv := vc[2].Copy()

	expectedVcs := map[int]clock.VectorClock{
		1: clock.NewVectorClockSet(3, map[int]int{1: 6, 2: 9, 3: 7}),
		2: clock.NewVectorClockSet(3, map[int]int{1: 5, 2: 10, 3: 7}),
		3: clock.NewVectorClockSet(3, map[int]int{1: 5, 2: 6, 3: 9}),
	}

	Unbuffered(&sender, &recv, vc)

	if !reflect.DeepEqual(vc, expectedVcs) {
		t.Errorf("Incorrect vc. Expected %v. Got %v.", expectedVcs, vc)
	}

	if !reflect.DeepEqual(sender.vc, expectedVcSend) {
		t.Errorf("Incorrect vc. Expected %v. Got %v.", expectedVcs, vc)
	}

	if !reflect.DeepEqual(recv.vc, expectedVcRecv) {
		t.Errorf("Incorrect vc. Expected %v. Got %v.", expectedVcs, vc)
	}
}

func TestBuffered(t *testing.T) {
	sender := TraceElementChannel{
		routine: 1,
		tPre:    5,
		tPost:   6,
		id:      123,
		opC:     SendOp,
		cl:      false,
		oID:     456,
		qSize:   1,
		pos:     "testfile:999",
		sel:     nil,
		partner: nil,
		vc:      clock.NewVectorClockSet(3, map[int]int{1: 5, 2: 6, 3: 7}),
	}

	recv := TraceElementChannel{
		routine: 2,
		tPre:    4,
		tPost:   7,
		id:      123,
		opC:     RecvOp,
		cl:      false,
		oID:     456,
		qSize:   1,
		pos:     "testfile:888",
		sel:     nil,
		partner: &sender,
		vc:      clock.NewVectorClockSet(3, map[int]int{1: 2, 2: 9, 3: 4}),
	}

	sender.partner = &recv

	vc := map[int]clock.VectorClock{
		1: clock.NewVectorClockSet(3, map[int]int{1: 5, 2: 6, 3: 7}),
		2: clock.NewVectorClockSet(3, map[int]int{1: 2, 2: 9, 3: 4}),
		3: clock.NewVectorClockSet(3, map[int]int{1: 5, 2: 6, 3: 9}),
	}

	Send(&sender, vc, false)
	Recv(&recv, vc, false)

	expectedVcs := map[int]clock.VectorClock{
		1: clock.NewVectorClockSet(3, map[int]int{1: 6, 2: 6, 3: 7}),
		2: clock.NewVectorClockSet(3, map[int]int{1: 6, 2: 10, 3: 7}),
		3: clock.NewVectorClockSet(3, map[int]int{1: 5, 2: 6, 3: 9}),
	}

	if !reflect.DeepEqual(vc, expectedVcs) {
		t.Errorf("Incorrect vc. Expected %v. Got %v.", expectedVcs, vc)
	}
}

func TestStuckChan(t *testing.T) {

	vc := map[int]clock.VectorClock{
		1: clock.NewVectorClockSet(3, map[int]int{1: 5, 2: 6}),
		2: clock.NewVectorClockSet(3, map[int]int{1: 2, 2: 9}),
	}

	StuckChan(1, vc)

	expectedVc := map[int]clock.VectorClock{
		1: clock.NewVectorClockSet(3, map[int]int{1: 6, 2: 6}),
		2: clock.NewVectorClockSet(3, map[int]int{1: 2, 2: 9}),
	}

	if !reflect.DeepEqual(vc, expectedVc) {
		t.Errorf("Incorrect vc. Expected %v. Got %v.", expectedVc, vc)
	}
}

func TestClose(t *testing.T) {
	ch := TraceElementChannel{
		routine: 1,
		tPre:    5,
		tPost:   6,
		id:      123,
		opC:     SendOp,
		cl:      false,
		oID:     456,
		qSize:   1,
		pos:     "testfile:999",
		sel:     nil,
		partner: nil,
		vc:      clock.NewVectorClockSet(3, map[int]int{1: 5, 2: 6, 3: 7}),
	}

	vc := map[int]clock.VectorClock{
		1: clock.NewVectorClockSet(3, map[int]int{1: 5, 2: 6}),
		2: clock.NewVectorClockSet(3, map[int]int{1: 2, 2: 9}),
	}

	Close(&ch, vc)

	expectedVc := map[int]clock.VectorClock{
		1: clock.NewVectorClockSet(3, map[int]int{1: 6, 2: 6}),
		2: clock.NewVectorClockSet(3, map[int]int{1: 2, 2: 9}),
	}

	if !reflect.DeepEqual(vc, expectedVc) {
		t.Errorf("Incorrect vc. Expected %v. Got %v.", expectedVc, vc)
	}

	if !ch.cl {
		t.Errorf("Incorrect cl. Expected true. Got false.")
	}
}
