// Copyright (c) 2025 Erik Kassubek
//
// File: fuzzing.go
// Brief: Start the toolchain mode
//
// Author: Erik Kassubek
// Created: 2025-01-06
//
// License: BSD-3-Clause

package modes

import (
	"analyzer/toolchain"
	"fmt"
)

func ModeToolchain(mode, advocate, file, execName, progName, test string,
	timeoutA, timeoutR, numRerecorded int,
	replayAt, meaTime, notExec, stats, keepTraces bool) {
	err := toolchain.Run(mode, advocate, file, execName, progName, test,
		timeoutA, timeoutR, numRerecorded,
		replayAt, meaTime, notExec, stats, keepTraces)
	if err != nil {
		fmt.Println("Failed to run toolchain")
		fmt.Println(err.Error())
	}
}
