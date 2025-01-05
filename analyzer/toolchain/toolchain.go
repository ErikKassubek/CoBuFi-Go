// Copyright (c) 2024 Erik Kassubek, Mario Occhinegro
//
// File: headerUnitTests.go
// Brief: Functions to add and remove the ADVOCATE header into file containing
//    unit tests
//
// Author: Erik Kassubek, Mario Occhinegro
// Created: 2024-09-18
//
// License: BSD-3-Clause

package toolchain

import (
	"fmt"
	"os"
	"strings"
)

/*
 * Main function to run the toolchain for main functions
 * TODO: add stats and measure time for main
 * Args:
 * 	pathToAdvocate (string): path to the ADVOCATE folder
 * 	pathToFile (string): path to the main file
 * 	executableName (string): name of the program executable
 * 	progName (string): name of the programs, used for stats
 * 	stats (bool): whether to create stats
 * 	measureTime (bool): whether to measure the runtime
 * 	timeoutAna (int): timeout of analysis
 * 	timeoutReplay (int): timeout for replay
 */
func runMain(pathToAdvocate string, pathToFile string, executableName string,
	progName string, stats, measureTime bool, timeoutAna, timeoutReplay int) error {
	// replace ~ in path with home
	home, _ := os.UserHomeDir()
	pathToAdvocate = strings.Replace(pathToAdvocate, "~", home, -1)
	pathToFile = strings.Replace(pathToFile, "~", home, -1)

	if pathToAdvocate == "" {
		return fmt.Errorf("Path to advocate required for mode main")
	}
	if pathToFile == "" {
		return fmt.Errorf("Path to file required")
	}
	if executableName == "" {
		return fmt.Errorf("Name of the executable required")
	}
	if (stats || measureTime) && progName == "" {
		return fmt.Errorf("If stats or measureTime is set, progName cannot be empty")
	}
	return runWorkflowMain(pathToAdvocate, pathToFile, executableName, timeoutAna, timeoutReplay)
}

/*
 * Main function to run the toolchain for tests
 * Args:
 * 	pathToAdvocate (string): path to the ADVOCATE folder
 * 	pathToTests (string): path to the test folder
 * 	progName (string): name of the programs, used for stats
 * 	stats (bool): whether to create stats
 * 	measureTime (bool): whether to measure the runtime
 * 	notExecuted (bool): check for never executed operations
 * 	timeoutAna (int): timeout of analysis
 * 	timeoutReplay (int): timeout for replay
 */
func runTest(pathToAdvocate, pathToTests, progName string, stats,
	measureTime bool, notExecuted bool, timeoutAna, timeoutReplay int) error {
	// replace ~ in path with home
	home, _ := os.UserHomeDir()
	pathToAdvocate = strings.Replace(pathToAdvocate, "~", home, -1)
	pathToTests = strings.Replace(pathToTests, "~", home, -1)

	if pathToAdvocate == "" {
		return fmt.Errorf("Path to advocate required")
	}
	if pathToTests == "" {
		return fmt.Errorf("Path to test folder required for mode main")
	}
	if (stats || measureTime) && progName == "" {
		return fmt.Errorf("If -s or -t is set, -N [name] must be set as well")
	}
	return runWorkflowUnit(pathToAdvocate, pathToTests, progName, measureTime, notExecuted, stats, timeoutAna, timeoutReplay)
}

/*
 * Generate bug reports
 * Args:
 * 	pathToFolder (string): path to the result folder
 */
func runReport(pathToFolder string) error {
	// replace ~ in path with home
	home, _ := os.UserHomeDir()
	pathToFolder = strings.Replace(pathToFolder, "~", home, -1)

	if pathToFolder == "" {
		return fmt.Errorf("Path to test folder required for mode main")
	}
	generateBugReports(pathToFolder)
	return nil
}
