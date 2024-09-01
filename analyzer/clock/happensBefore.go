// Copyrigth (c) 2024 Erik Kassubek
//
// File: happensBefore.go
// Brief: Type for happens before   
// 
// Author: Erik Kassubek <kassubek.erik@gmail.com>
// Created: 2023-11-30
// LastChange: 2024-09-01
//
// License: BSD-3-Clause

package clock

type HappensBefore int

const (
	Before HappensBefore = iota
	After
	Concurrent
	None
)
