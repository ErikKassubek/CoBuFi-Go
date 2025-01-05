// Copyright (c) 2025 Erik Kassubek
//
// File: fuzzing.go
// Brief: Start the explain mode
//
// Author: Erik Kassubek
// Created: 2025-01-05
//
// License: BSD-3-Clause

package modes

import (
	"analyzer/explanation"
	"fmt"
	"log"
)

func ModeExplain(pathTrace string, ignoreDouble bool) {
	if pathTrace == "" {
		fmt.Println("Please provide a path to the trace files for the explanation. Set with -t [file]")
		return
	}

	err := explanation.CreateOverview(pathTrace, ignoreDouble)
	if err != nil {
		log.Println("Error creating explanation: ", err.Error())
	}
}
