// Copyrigth (c) 2024 Erik Kassubek
//
// File: traceElements.go
// Brief: Interface for all trace element types
//
// Author: Erik Kassubek <kassubek.erik@gmail.com>
// Created: 2023-08-08
// LastChange: 2024-09-01
//
// License: BSD-3-Clause

package analysis

import "analyzer/clock"

// Interface for trace elements
type TraceElement interface {
	GetID() int
	GetTPre() int
	SetTPre(tPre int)
	getTpost() int
	GetTSort() int
	SetTSort(tSort int)
	SetT(time int)
	SetTWithoutNotExecuted(tSort int)
	GetRoutine() int
	GetPos() string
	GetTID() string
	ToString() string
	updateVectorClock()
	GetVC() clock.VectorClock
	Copy() TraceElement
}
