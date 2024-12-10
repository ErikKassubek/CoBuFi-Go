// ADVOCATE-FILE_START

// Copyright (c) 2024 Erik Kassubek
//
// File: advocate_time.go
// Brief: Set of functions using time
//
// Author: Erik Kassubek
// Created: 2024-12-04
//
// License: BSD-3-Clause

package runtime

func sleep(seconds float64) {
	start := nanotime()
	durationNano := int64(seconds * 1e9)
	for nanotime()-start < durationNano {
	}
}

func sToNs(seconds int64) int64 {
	return seconds * 1e9
}

func hasTimePast(startNs int64, durationS int64) bool {
	durationNano := durationS * 1e9
	return nanotime()-startNs > durationNano
}

func currentTime() int64 {
	return nanotime()
}

// ADVOCATE-FILE-END
