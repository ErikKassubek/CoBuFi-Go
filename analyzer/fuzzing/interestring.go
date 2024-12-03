// Copyright (c) 2024 Erik Kassubek
//
// File: interestring.go
// Brief: Functions to determine whether a run was interesting
//
// Author: Erik Kassubek
// Created: 2024-11-28
//
// License: BSD-3-Clause

package fuzzing

/*
 * A run is considered interesting, if at least one of the following conditions is met
 * The run contains a new pair of channel operations (new meaning it has not been seen in any of the previous runs)
 * An operation pair's execution counter changes significantly from previous order.
 * A new channel operation is triggered, such as creating, closing or not closing a channel for the first time
 * A buffered channel gets a larger maximum fullness than in all previous executions (MaxChBufFull)
 */
func isInteresting() bool {
	// TODO: implement
	return false
}
