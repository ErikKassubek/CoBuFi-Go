// Copyright (c) 2025 Erik Kassubek
//
// File: fuzzing.go
// Brief: Start the fuzzing mode
//
// Author: Erik Kassubek
// Created: 2025-01-05
//
// License: BSD-3-Clause

package modes

import (
	"analyzer/fuzzing"
	"fmt"
)

func ModeFuzzing(pathTrace string, progName string, testName string, index int) {
	if pathTrace == "" {
		fmt.Println("Please provide a path to the trace file. Set with -t [folder]")
		return
	}

	if progName == "" {
		fmt.Println("Provide a name for the analyzed program. Set with -N [name]")
		return
	}

	if testName != "" {
		progName = progName + "_" + testName
	}

	fuzzing.Fuzzing(pathTrace, pathTrace, progName, index)
}
