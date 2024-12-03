// Copyright (c) 2024 Erik Kassubek
//
// File: traceElementNew.go
// Brief: Trace element to store the creation (new) of relevant operations. For now this is only creates the new for channel. This may be expanded later.
//
// Author: Erik Kassubek
// Created: 2024-11-29
//
// License: BSD-3-Clause

package analysis

import (
	"analyzer/clock"
	"errors"
	"fmt"
	"strconv"
)

type newOpType string

const (
	atomic      newOpType = "A"
	channel     newOpType = "C"
	conditional newOpType = "D"
	mutex       newOpType = "M"
	once        newOpType = "O"
	wait        newOpType = "W"
)

/*
 * TraceElementNew is a trace element for the creation of an object / new
 * Fields:
 *   routine (int): The routine id
 *   tPost (int): The timestamp of the new
 *   id (int): The id of the underlying operation
 *   elemType (newOpType): The type of the created object
 *   num (int): Variable field for additional information
 *   pos (string): The position of the new
 * For now this is only creates the new for channel. This may be expanded later.
 */
type TraceElementNew struct {
	routine  int
	tPost    int
	id       int
	elemType newOpType
	num      int
	pos      string
	vc       clock.VectorClock
}

func AddTraceElementNew(routine int, tPost string, id string, elemType string, num string, pos string) error {
	tPostInt, err := strconv.Atoi(tPost)
	if err != nil {
		return errors.New("tpost is not an integer")
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		return errors.New("id is not an integer")
	}

	numInt, err := strconv.Atoi(num)
	if err != nil {
		return errors.New("num is not an integer")
	}

	elem := TraceElementNew{
		routine:  routine,
		tPost:    tPostInt,
		id:       idInt,
		elemType: newOpType(elemType),
		num:      numInt,
		pos:      pos,
	}

	return AddElementToTrace(&elem)
}

func (n *TraceElementNew) GetID() int {
	return n.id
}

func (n *TraceElementNew) GetTPre() int {
	return n.tPost
}

func (n *TraceElementNew) GetTPost() int {
	return n.tPost
}

func (n *TraceElementNew) GetTSort() int {
	return n.tPost
}

func (n *TraceElementNew) GetRoutine() int {
	return n.routine
}

func (n *TraceElementNew) GetPos() string {
	return n.pos
}

func (n *TraceElementNew) GetTID() string {
	return n.pos + "@" + strconv.Itoa(n.tPost)
}

func (n *TraceElementNew) GetObjType() string {
	switch n.elemType {
	case atomic:
		return "NA"
	case channel:
		return "NC"
	case conditional:
		return "ND"
	case mutex:
		return "NM"
	case once:
		return "NO"
	case wait:
		return "NW"
	default:
		return "N"
	}
}

func (n *TraceElementNew) GetVC() clock.VectorClock {
	return n.vc
}

func (n *TraceElementNew) GetNum() int {
	return n.num
}

func (n *TraceElementNew) ToString() string {
	return fmt.Sprintf("N,%d,%d,%s,%d,%s", n.tPost, n.id, string(n.elemType), n.num, n.pos)
}

func (n *TraceElementNew) SetTPre(tSort int) {
	n.tPost = tSort
}

func (n *TraceElementNew) SetT(tSort int) {
	n.tPost = tSort
}

func (n *TraceElementNew) SetTSort(tSort int) {
	n.tPost = tSort
}

func (n *TraceElementNew) SetTWithoutNotExecuted(tSort int) {
	if n.tPost == 0 {
		return
	}
	n.tPost = tSort
}

func (n *TraceElementNew) updateVectorClock() {
	n.vc = currentVCHb[n.routine].Copy()

	currentVCHb[n.routine].Inc(n.routine)
}

func (n *TraceElementNew) Copy() TraceElement {
	return &TraceElementNew{
		routine:  n.routine,
		tPost:    n.tPost,
		id:       n.id,
		elemType: n.elemType,
		pos:      n.pos,
		vc:       n.vc.Copy(),
	}
}
