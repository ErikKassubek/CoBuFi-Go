// Copyright (c) 2024 Erik Kassubek
//
// File: fuzzing.go
// Brief: Main file for fuzzing
//
// Author: Erik Kassubek
// Created: 2024-12-03
//
// License: BSD-3-Clause

package fuzzing

import (
	"analyzer/io"
	"fmt"
	"path/filepath"
)

/*
 * Create the fuzzing data
 * Args:
 * 	pathFuzzing (string): path to the fuzzing data
 * 	pathTrace (string): path to the trace
 * 	progName (string): unique identifier for the program or test
 * 	lastID (int): last ID of fuzzing select traces
 * Returns:
 * 	int: last ID of created fuzzing traces
 */
func Fuzzing(pathFuzzing, pathTrace string, progName string, lastID int) int {
	fuzzingFilePath := filepath.Join(pathFuzzing, fmt.Sprintf("fuzzingFile_%s.info", progName))
	readFile(fuzzingFilePath)

	io.CreateTraceFromFiles(pathTrace, true)
	parseTrace()

	// if the run was not interesting, there is nothing else to do
	if !isInteresting() {
		return lastID
	}

	numMut := numberMutations()
	muts := createMutations(numMut)

	updateFileData()
	writeFileInfo(fuzzingFilePath)

	return writeMutationsToFile(pathFuzzing, lastID, muts, progName)
}
