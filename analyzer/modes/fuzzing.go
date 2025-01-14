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

// TODO (FUZZING): make fuzzing work
func ModeFuzzing(advocate, testPath, progName, testName string) {
	if progName == "" {
		fmt.Println("Provide a name for the analyzed program. Set with -prog [name]")
		return
	}

	if testName == "" {
		fmt.Println("Provide a name for the analyzed test. Set with -test [name]")
		return
	}

	fuzzing.Fuzzing(advocate, testPath, progName, testName)
}
