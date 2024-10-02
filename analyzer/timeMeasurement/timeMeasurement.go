// Copyright (c) 2024 Erik Kassubek
//
// File: timeMeasurement.go
// Brief: Measure times
//
// Author: Erik Kassubek
// Created: 2024-10-02
//
// License: BSD-3-Clause

package timemeasurement

import (
	"fmt"
	"time"
)

var duration = make(map[string]time.Duration)
var start = make(map[string]time.Time)

// counter : total,leak,panic, io, rewrite, other (other beeing untriggered select, recv on closed usw)
// time for HBAnalysis: total - everythingElse

func Start(counter string) {
	start[counter] = time.Now()
}

func End(counter string) {
	if _, ok := duration[counter]; !ok {
		duration[counter] = time.Since(start[counter])
	} else {
		duration[counter] += time.Since(start[counter])
	}
	start[counter] = time.Now()
}

func Print() {
	fmt.Printf("AdvocateAnalysisTimes:%.5f#%.5f#%.5f#%.5f\n",
		duration["analysis"].Seconds(),
		duration["leak"].Seconds(), duration["panic"].Seconds(), duration["other"].Seconds())
}
