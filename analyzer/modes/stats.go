// Copyright (c) 2025 Erik Kassubek
//
// File: fuzzing.go
// Brief: Start the stats mode
//
// Author: Erik Kassubek
// Created: 2025-01-05
//
// License: BSD-3-Clause

package modes

import (
	"analyzer/stats"
	"fmt"
)

func ModeStats(pathFolder string, progName string, testName string) {
	// instead of the normal program, create statistics for the trace
	if pathFolder == "" {
		fmt.Println("Provide the path to the folder containing the results_machine file. Set with -t [path]")
		return
	}

	if progName == "" {
		fmt.Println("Provide a name for the analyzed program. Set with -N [name]")
		return
	}

	if testName == "" {
		testName = progName
	}

	stats.CreateStats(pathFolder, progName, testName)
}
