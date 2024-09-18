// Copyright (c) 2024 Erik Kassubek, Mario Occhinegro
//
// File: runFullWorkflowMain.go
// Brief: Function to run the whole ADVOCATE workflow, including running,
//    analysis and replay on all unit tests of a program
//
// Author: Erik Kassubek, Mario Occhinegro
// Created: 2024-09-18
// Last Changed 2024-09-18
//
// License: BSD-3-Clause

package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

/*
 * Run ADVOCATE for all given unit tests
 * Args:
 *    pathToAdvocate (string): pathToAdvocate
 *    dir (string): path to the folder containing the unit tests
 * Returns:
 *    error
 */
func runWorkflowUnit(pathToAdvocate string, dir string) error {
	// Validate required inputs
	if pathToAdvocate == "" {
		return errors.New("Path to advocate is empty")
	}
	if dir == "" {
		return errors.New("Directory is empty")
	}

	// Set paths
	pathToAnalyzer := filepath.Join(pathToAdvocate, "analyzer/analyzer")
	// pathToFullWorkflowExecutor := filepath.Join(pathToAdvocate, "toolchain/unitTestFullWorkflow/unitTestFullWorkflow.bash")

	// Change to the directory
	if err := os.Chdir(dir); err != nil {
		return fmt.Errorf("Failed to change directory: %v", err)
	}
	fmt.Printf("In directory: %s\n", dir)

	// Create the advocateResult directory
	if err := os.MkdirAll("advocateResult", os.ModePerm); err != nil {
		return fmt.Errorf("Failed to create advocateResult directory: %v", err)
	}

	// Find all _test.go files in the directory
	testFiles, err := findTestFiles(dir)
	if err != nil {
		return fmt.Errorf("Failed to find test files: %v", err)
	}

	totalFiles := len(testFiles)
	attemptedTests, skippedTests, currentFile := 0, 0, 1

	// Process each test file
	for _, file := range testFiles {
		fmt.Printf("Progress: %d/%d\n", currentFile, totalFiles)
		fmt.Printf("Processing file: %s\n", file)

		packagePath := filepath.Dir(file)
		testFunctions, err := findTestFunctions(file)
		if err != nil {
			log.Printf("Failed to find test functions in %s: %v", file, err)
			continue
		}

		for _, testFunc := range testFunctions {
			attemptedTests++
			packageName := filepath.Base(packagePath)
			fileName := filepath.Base(file)
			fmt.Printf("Running full workflow for test: %s in package: %s in file: %s\n", testFunc, packageName, file)

			adjustedPackagePath := strings.TrimPrefix(packagePath, dir)
			directoryName := fmt.Sprintf("advocateResult/file(%d)-test(%d)-%s-%s", currentFile, attemptedTests, fileName, testFunc)
			if err := os.MkdirAll(directoryName, os.ModePerm); err != nil {
				log.Printf("Failed to create directory %s: %v", directoryName, err)
				continue
			}

			// Execute full workflow
			err = unitTestFullWorkflow(pathToAdvocate, dir, testFunc,
				adjustedPackagePath, file, filepath.Join(directoryName, "output.txt"))

			if err != nil {
				fmt.Printf("File %d with Test %d failed, check output.txt for more information. Skipping...\n", currentFile, attemptedTests)
				skippedTests++
				continue
			}

			// Move logs and results to the appropriate directory
			moveResults(packagePath, directoryName)
		}
		currentFile++
	}

	// Generate Bug Reports
	fmt.Println("Generate Bug Reports")
	generateBugReports(pathToAdvocate, filepath.Join(dir, "advocateResult"))

	// Check for untriggered selects
	fmt.Println("Check for untriggered selects")
	runCommand(pathToAnalyzer, "-o", "-R", filepath.Join(dir, "advocateResult"), "-P", dir)

	// Output test summary
	fmt.Println("Finished full workflow for all tests")
	fmt.Printf("Attempted tests: %d\n", attemptedTests)
	fmt.Printf("Skipped tests: %d\n", skippedTests)

	return nil
}

/*
 * Function to find all _test.go files in the specified directory
 * Args:
 *    dir (string): folder to search in
 * Returns:
 *    []string: found files
 *    error
 */
func findTestFiles(dir string) ([]string, error) {
	var testFiles []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if strings.HasSuffix(info.Name(), "_test.go") {
			testFiles = append(testFiles, path)
		}
		return nil
	})
	return testFiles, err
}

/*
 * Function to find all test function in the specified file
 * Args:
 *    file (string): file to search in
 * Returns:
 *    []string: functions
 *    error
 */func findTestFunctions(file string) ([]string, error) {
	content, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	var testFunctions []string
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		if strings.Contains(line, "Test") && strings.Contains(line, "func") && strings.Contains(line, "*testing.T") {
			testFunc := strings.TrimSpace(strings.Split(line, "(")[0])
			testFunc = strings.TrimPrefix(testFunc, "func ")
			testFunctions = append(testFunctions, testFunc)
		}
	}
	return testFunctions, nil
}

/*
 * Function to move results files from the package directory to the destination directory
 * Args:
 *    packagePath (string): path to the package directory
 *    destination (string): path to the destination directory
 */
func moveResults(packagePath, destination string) {
	filesToMove := []string{
		"advocateTrace",
		"results_machine.log",
		"results_readable.log",
		"times.log",
	}

	for _, file := range filesToMove {
		src := filepath.Join(packagePath, file)
		dest := filepath.Join(destination, file)
		if err := os.Rename(src, dest); err != nil {
			log.Printf("Failed to move %s to %s: %v", src, dest, err)
		}
	}

	// Move any rewritten_trace directories
	rewrittenTraces, _ := filepath.Glob(filepath.Join(packagePath, "rewritten_trace*"))
	for _, trace := range rewrittenTraces {
		dest := filepath.Join(destination, filepath.Base(trace))
		if err := os.Rename(trace, dest); err != nil {
			log.Printf("Failed to move %s to %s: %v", trace, dest, err)
		}
	}

	// Move advocateCommand.log if present
	if _, err := os.Stat("advocateCommand.log"); err == nil {
		os.Rename("advocateCommand.log", filepath.Join(destination, "advocateCommand.log"))
	}
}

/*
 * Run the full workflow for a given unit test
 * Args:
 *    pathToAdvocate (string): path to advocate
 *    dir (string): path to the package to test
 *    testName (string): name of the test
 *    pkg (string): adjusted package path
 *    file (string): file with the test
 *    output (string): write all outputs to this file
 * Returns:
 *    error
 */
func unitTestFullWorkflow(pathToAdvocate string, dir string,
	testName string, pkg string, file string, output string) error {

	outFile, err := os.OpenFile(output, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("Failed to open log file: %v", err)
	}
	defer outFile.Close()

	// Redirect stdout and stderr to the file
	os.Stdout = outFile
	os.Stderr = outFile

	// Validate required inputs
	if pathToAdvocate == "" {
		return errors.New("Path to advocate is empty")
	}
	if dir == "" {
		return errors.New("Directory is empty")
	}
	if testName == "" {
		return errors.New("Test name is empty")
	}
	if pkg == "" {
		return errors.New("Package is empty")
	}
	if file == "" {
		return errors.New("Test file is empty")
	}

	pathToPatchedGoRuntime := filepath.Join(pathToAdvocate, "go-patch/bin/go")
	pathToGoRoot := filepath.Join(pathToAdvocate, "go-patch")
	pathToOverheadInserter := filepath.Join(pathToAdvocate, "toolchain/unitTestOverheadInserter/unitTestOverheadInserter")
	pathToOverheadRemover := filepath.Join(pathToAdvocate, "toolchain/unitTestOverheadRemover/unitTestOverheadRemover")
	pathToAnalyzer := filepath.Join(pathToAdvocate, "analyzer/analyzer")

	// Change to the directory
	if err := os.Chdir(dir); err != nil {
		return fmt.Errorf("Failed to change directory: %v", err)
	}
	fmt.Printf("In directory: %s\n", dir)

	// Set GOROOT
	os.Setenv("GOROOT", pathToGoRoot)
	fmt.Println("GOROOT exported")

	// Create advocateCommand.log
	logFile, err := os.OpenFile("advocateCommand.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("Failed to create advocateCommand.log: %v", err)
	}
	defer logFile.Close()

	writeLog := func(data string) {
		if _, err := logFile.WriteString(data + "\n"); err != nil {
			log.Printf("Failed to write to log: %v", err)
		}
	}

	writeLog(file)
	writeLog(testName)

	// Remove overhead just in case
	writeLog(fmt.Sprintf("%s -f %s -t %s", pathToOverheadRemover, file, testName))
	if err := runCommand(pathToOverheadRemover, "-f", file, "-t", testName); err != nil {
		return fmt.Errorf("Failed to remove overhead: %v", err)
	}

	// Add overhead
	writeLog(fmt.Sprintf("%s -f %s -t %s", pathToOverheadInserter, file, testName))
	if err := runCommand(pathToOverheadInserter, "-f", file, "-t", testName); err != nil {
		return fmt.Errorf("Error in adding overhead: %v", err)
	}

	// Run the test
	writeLog(fmt.Sprintf("%s test -count=1 -run=%s ./%s", pathToPatchedGoRuntime, testName, pkg))
	err = runCommand(pathToPatchedGoRuntime, "test", "-count=1", "-run="+testName, "./"+pkg)

	if err != nil {
		log.Println("Test failed, removing overhead and exiting")
		runCommand(pathToOverheadRemover, "-f", file, "-t", testName)
		return errors.New("Error running test")
	}

	// Remove overhead after the test
	writeLog(fmt.Sprintf("%s -f %s -t %s", pathToOverheadRemover, file, testName))
	runCommand(pathToOverheadRemover, "-f", file, "-t", testName)

	// Apply analyzer
	writeLog(fmt.Sprintf("%s -t %s/%s/advocateTrace", pathToAnalyzer, dir, pkg))
	runCommand(pathToAnalyzer, "-t", filepath.Join(dir, pkg, "advocateTrace"))

	// Apply reorder overhead for rewritten traces
	rewrittenTraces, _ := filepath.Glob(filepath.Join(pkg, "rewritten_trace*"))
	for _, trace := range rewrittenTraces {
		traceNum := extractTraceNumber(trace)
		writeLog(fmt.Sprintf("%s -f %s -t %s -r true -n %s", pathToOverheadInserter, file, testName, traceNum))
		runCommand(pathToOverheadInserter, "-f", file, "-t", testName, "-r", "true", "-n", traceNum)

		writeLog(fmt.Sprintf("%s test -count=1 -run=%s ./%s", pathToPatchedGoRuntime, testName, pkg))
		runCommandWithTee(pathToPatchedGoRuntime, filepath.Join(trace, "reorder_output.txt"), "test", "-count=1", "-run="+testName, "./"+pkg)

		// Remove reorder overhead
		writeLog(fmt.Sprintf("%s -f %s -t %s", pathToOverheadRemover, file, testName))
		runCommand(pathToOverheadRemover, "-f", file, "-t", testName)
	}

	// Unset GOROOT
	os.Unsetenv("GOROOT")

	return nil
}

// runCommandWithTee runs a command and writes output to a file
func runCommandWithTee(name, outputFile string, args ...string) error {
	cmd := exec.Command(name, args...)
	outfile, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer outfile.Close()
	cmd.Stdout = outfile
	cmd.Stderr = outfile
	return cmd.Run()
}

// extractTraceNumber extracts the numeric part from a trace directory name
func extractTraceNumber(trace string) string {
	parts := strings.Split(trace, "rewritten_trace")
	if len(parts) > 1 {
		return parts[1]
	}
	return ""
}
