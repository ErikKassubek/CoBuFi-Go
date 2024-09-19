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
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
)

/*
 * Run ADVOCATE on a program with a main function
 * Args:
 *    pathToAdvocate (string): path to the ADVOCATE folder
 *    pathToFile (string): path to the file containing the main function
 * Returns:
 *    error
 */
func runWorkflowMain(pathToAdvocate string, pathToFile string) error {
	if pathToAdvocate == "" {
		return errors.New("path to ADVOCATE not set")
	}

	if pathToFile == "" {
		return errors.New("No file to analyze given")
	}

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

	// Set GOROOT environment variable
	if err := os.Setenv("GOROOT", pathToGoRoot); err != nil {
		return fmt.Errorf("Failed to set GOROOT: %v", err)
	}
	fmt.Println("GOROOT exported")

	// Remove header
	if err := headerRemoverMain(pathToFile); err != nil {
		return fmt.Errorf("Error removing header: %v", err)
	}

	// Add header
	fmt.Printf("Add header to %s", pathToFile)
	if err := headerInserterMain(pathToFile, false, "1"); err != nil {
		return fmt.Errorf("Error in adding header: %v", err)
	}

	// Run the program
	fmt.Printf("%s run %s\n", pathToPatchedGoRuntime, pathToFile)
	if err := runCommand(pathToPatchedGoRuntime, "run", pathToFile); err != nil {
		log.Println("Error in running program, removing header and stopping workflow")
		headerRemoverMain(pathToFile)
		return err
	}

	// Remove header
	if err := headerRemoverMain(pathToFile); err != nil {
		return fmt.Errorf("Error removing header: %v", err)
	}

	// Apply analyzer
	analyzerOutput := filepath.Join(dir, "advocateTrace")
	if err := runCommand(pathToAnalyzer, "run", "-t", analyzerOutput); err != nil {
		return fmt.Errorf("Error applying analyzer: %v", err)
	}

	// Find rewritten_trace directories
	rewrittenTraces, err := filepath.Glob(filepath.Join(dir, "rewritten_trace*"))
	if err != nil {
		return fmt.Errorf("Error finding rewritten traces: %v", err)
	}

	// Apply replay header and run tests for each trace
	for _, trace := range rewrittenTraces {
		rtraceNum := extractTraceNum(trace)
		fmt.Printf("Apply replay header for file f %s and trace %s\n", pathToFile, rtraceNum)
		if err := headerInserterMain(pathToFile, true, rtraceNum); err != nil {
			return err
		}

		outputFile := filepath.Join(trace, "replay_output.log")
		runCommandWithOutput(pathToPatchedGoRuntime, outputFile, "run", pathToFile)

		fmt.Printf("Remove replay header from %s\n", pathToFile)
		if err := headerRemoverMain(pathToFile); err != nil {
			return err
		}
	}

	// Unset GOROOT
	os.Unsetenv("GOROOT")

	return nil
}

func extractTraceNum(tracePath string) string {
	re := regexp.MustCompile(`[0-9]+$`)
	return re.FindString(tracePath)
}
