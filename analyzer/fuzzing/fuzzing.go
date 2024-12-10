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

import "path/filepath"

/*
 * Create the fuzzing data
 * Args:
 * 	path (string): path to the fuzzing data
 * 	lastID (int): last ID of fuzzing select traces
 * Returns:
 * 	int: last ID of created fuzzing traces
 */
func Fuzzing(path string, lastID int) int {
	fuzzingFilePath := filepath.Join(path, "fuzzingFile.info")
	readFile(fuzzingFilePath)

	// if the run was not interesting, there is nothing else to do
	if !isInteresting() {
		return lastID
	}

	numMut := numberMutations()
	muts := createMutations(numMut)

	updateFileData()
	writeFileInfo(fuzzingFilePath)

	return writeMutationsToFile(path, lastID, muts)
}
