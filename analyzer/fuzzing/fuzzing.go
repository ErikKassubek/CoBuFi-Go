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

func Fuzzing(fuzzingFilePath string) {
	readFile(fuzzingFilePath)

	// if the run was not interesting, there is nothing else to do
	if !isInteresting() {
		return
	}

	numMut := numberMutations()
	mutate(numMut)
}
