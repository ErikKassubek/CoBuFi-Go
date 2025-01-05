// Copyright (c) 2025 Erik Kassubek
//
// File: fuzzing.go
// Brief: Start the check mode
//
// Author: Erik Kassubek
// Created: 2025-01-05
//
// License: BSD-3-Clause

package modes

import (
	"analyzer/complete"
	"fmt"
)

func ModeCheck(resultFolderTool, programPath *string) {
	if *resultFolderTool == "" {
		fmt.Println("Please provide the path to the advocateResult folder created by the pipeline. Set with -R [folder]")
		return
	}

	if *programPath == "" {
		fmt.Println("Please provide the path to the program folder. Set with -P [folder]")
		return
	}

	err := complete.Check(*resultFolderTool, *programPath)

	if err != nil {
		panic(err.Error())
	}
}
