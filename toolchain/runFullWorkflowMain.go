// Copyright (c) 2024 Erik Kassubek, Mario Occhinegro
//
// File: runFullWorkflowMain.go
// Brief: Function to run the whole ADVOCATE workflow, including running,
//    analysis and replay on a program with a main function
//
// Author: Erik Kassubek, Mario Occhinegro
// Created: 2024-09-18
// Last Changed 2024-09-19
//
// License: BSD-3-Clause

package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"time"
)

/*
 * Run ADVOCATE on a program with a main function
 * Args:
 *    pathToAdvocate (string): path to the ADVOCATE folder
 *    pathToFile (string): path to the file containing the main function
 *    executableName (string): name of the executable
 * Returns:
 *    error
 */
func runWorkflowMain(pathToAdvocate string, pathToFile string, executableName string) error {
	if _, err := os.Stat(pathToFile); os.IsNotExist(err) {
		return fmt.Errorf("file %s does not exist", pathToFile)
	}

	pathToPatchedGoRuntime := filepath.Join(pathToAdvocate, "go-patch/bin/go")
	pathToGoRoot := filepath.Join(pathToAdvocate, "go-patch")
	pathToAnalyzer := filepath.Join(pathToAdvocate, "analyzer/analyzer")

	// Change to the directory of the main file
	dir := filepath.Dir(pathToFile)
	if err := os.Chdir(dir); err != nil {
		return fmt.Errorf("Failed to change directory: %v", err)
	}
	fmt.Printf("In directory: %s\n", dir)

	os.RemoveAll("advocateResult")
	if err := os.MkdirAll("advocateResult", os.ModePerm); err != nil {
		return fmt.Errorf("Failed to create advocateResult directory: %v", err)
	}

	// Set GOROOT environment variable
	if err := os.Setenv("GOROOT", pathToGoRoot); err != nil {
		return fmt.Errorf("Failed to set GOROOT: %v", err)
	}
	fmt.Println("GOROOT exported")

	// Remove header
	if err := headerRemoverMain(pathToFile); err != nil {
		return fmt.Errorf("Error removing header: %v", err)
	}

	var durationRun time.Duration
	var durationRecord time.Duration
	var durationAnalysis time.Duration
	var durationReplay time.Duration

	// build the program
	if measureTime {
		fmt.Printf("%s build\n", pathToPatchedGoRuntime)
		if err := runCommand(pathToPatchedGoRuntime, "build"); err != nil {
			log.Println("Error in building program, removing header and stopping workflow")
			headerRemoverMain(pathToFile)
			return err
		}

		// run the program
		fmt.Printf("./%s\n", executableName)
		timeStart := time.Now()
		if err := runCommand("./" + executableName); err != nil {
			log.Println("Error in running program, removing header and stopping workflow")
			headerRemoverMain(pathToFile)
			return err
		}
		durationRun = time.Since(timeStart)
	}

	// Add header
	fmt.Printf("Add header to %s\n", pathToFile)
	if err := headerInserterMain(pathToFile, false, "1"); err != nil {
		return fmt.Errorf("Error in adding header: %v", err)
	}

	// build the program
	fmt.Printf("%s build\n", pathToPatchedGoRuntime)
	if err := runCommand(pathToPatchedGoRuntime, "build"); err != nil {
		log.Println("Error in building program, removing header and stopping workflow")
		headerRemoverMain(pathToFile)
		return err
	}

	// run the program
	fmt.Printf("./%s\n", executableName)
	timeStart := time.Now()
	if err := runCommand("./" + executableName); err != nil {
		log.Println("Error in running program, removing header and stopping workflow")
		headerRemoverMain(pathToFile)
		return err
	}
	durationRecord = time.Since(timeStart)

	// Remove header
	if err := headerRemoverMain(pathToFile); err != nil {
		return fmt.Errorf("Error removing header: %v", err)
	}

	// Apply analyzer
	analyzerOutput := filepath.Join(dir, "advocateTrace")
	timeStart = time.Now()
	if err := runCommand(pathToAnalyzer, "run", "-t", analyzerOutput); err != nil {
		return fmt.Errorf("Error applying analyzer: %v", err)
	}
	durationAnalysis = time.Since(timeStart)

	// Find rewritten_trace directories
	rewrittenTraces, err := filepath.Glob(filepath.Join(dir, "rewritten_trace*"))
	if err != nil {
		return fmt.Errorf("Error finding rewritten traces: %v", err)
	}

	// Apply replay header and run tests for each trace
	timeStart = time.Now()
	for _, trace := range rewrittenTraces {
		rtraceNum := extractTraceNum(trace)
		fmt.Printf("Apply replay header for file f %s and trace %s\n", pathToFile, rtraceNum)
		if err := headerInserterMain(pathToFile, true, rtraceNum); err != nil {
			return err
		}

		// build the program
		fmt.Printf("%s build\n", pathToPatchedGoRuntime)
		if err := runCommand(pathToPatchedGoRuntime, "build"); err != nil {
			log.Println("Error in building program, removing header and stopping workflow")
			headerRemoverMain(pathToFile)
			continue
		}

		// run the program
		fmt.Printf("./%s\n", executableName)
		runCommand("./" + executableName)

		fmt.Printf("Remove replay header from %s\n", pathToFile)
		if err := headerRemoverMain(pathToFile); err != nil {
			return err
		}
	}

	durationReplay = time.Since(timeStart)
	if len(rewrittenTraces) > 0 {
		durationReplay = durationReplay / time.Duration(len(rewrittenTraces))
	}

	// Unset GOROOT
	os.Unsetenv("GOROOT")

	resultPath := filepath.Join(dir, "advocateResult")

	moveResults(dir, resultPath)

	// Generate Bug Reports
	fmt.Println("Generate Bug Reports")
	generateBugReports(resultPath, pathToAdvocate)

	if measureTime {
		generateTimeFile(progName, resultPath, durationRun, durationRecord, durationAnalysis, durationReplay)
	}

	if notExecuted {
		fmt.Println("Check for untriggered selects and not executed progs")
		runCommand(pathToAnalyzer, "check", "-R", filepath.Join(dir, "advocateResult"), "-P", dir)
	}

	if stats {
		// create statistics
		fmt.Println("Create statistics")
		runCommand(pathToAnalyzer, "stats", "-R", filepath.Join(dir, "advocateResult"), "-P", dir, "-N", progName)
	}

	return nil
}

func extractTraceNum(tracePath string) string {
	re := regexp.MustCompile(`[0-9]+$`)
	return re.FindString(tracePath)
}
