// Copyright (c) 2024 Erik Kassubek, Mario Occhinegro
//
// File: runFullWorkflowMain.go
// Brief: Function to run the whole ADVOCATE workflow, including running,
//    analysis and replay on all unit tests of a program
//
// Author: Erik Kassubek, Mario Occhinegro
// Created: 2024-09-18
// Last Changed 2024-09-20
//
// License: BSD-3-Clause

package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

/*
 * Run ADVOCATE for all given unit tests
 * Args:
 *    pathToAdvocate (string): pathToAdvocate
 *    dir (string): path to the folder containing the unit tests
 *    progName (string): name of the analyzed program
 *    measureTime (bool): if true, measure the time for all steps. This
 *      also runs the tests once without any recoding/replay to get a base value
 *    notExecuted (bool): if true, check for never executed operations
 *    stats (bool): create a stats file
 * Returns:
 *    error
 */
func runWorkflowUnit(pathToAdvocate, dir, progName string,
	measureTime, notExecuted, stats bool) error {
	// Validate required inputs
	if pathToAdvocate == "" {
		return errors.New("Path to advocate is empty")
	}
	if dir == "" {
		return errors.New("Directory is empty")
	}

	pathToAnalyzer := filepath.Join(pathToAdvocate, "analyzer/analyzer")

	// Change to the directory
	if err := os.Chdir(dir); err != nil {
		return fmt.Errorf("Failed to change directory: %v", err)
	}
	fmt.Printf("In directory: %s\n", dir)

	os.RemoveAll("advocateResult")
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

	var durationRun time.Duration
	var durationRecord time.Duration
	var durationAnalysis time.Duration
	var durationReplay time.Duration

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
			fileNameWithoutEnding := strings.TrimSuffix(fileName, ".go")
			directoryName := fmt.Sprintf("advocateResult/file(%d)-test(%d)-%s-%s", currentFile, attemptedTests, fileNameWithoutEnding, testFunc)
			if err := os.MkdirAll(directoryName, os.ModePerm); err != nil {
				log.Printf("Failed to create directory %s: %v", directoryName, err)
				continue
			}

			// run the tests without recording/replay
			if measureTime {
				if err := os.Chdir(dir); err != nil {
					fmt.Printf("Failed to change directory: %v\n", err)
				}

				// Remove header just in case
				if err := headerRemoverUnit(file); err != nil {
					fmt.Println(err)
				}

				timeStart := time.Now()
				err = runCommand("go", "test", "-count=1", "-run="+testFunc, "./"+adjustedPackagePath)
				if err != nil {
					log.Println(err)
					log.Println("Test failed, removing header and exiting")
				}
				duration := time.Since(timeStart)
				durationRun += duration
			}

			// Execute full workflow
			timeRecord, timeAnalysis, timeReplay, err := unitTestFullWorkflow(pathToAdvocate, dir, testFunc, adjustedPackagePath, file, directoryName)
			durationRecord += timeRecord
			durationAnalysis += timeAnalysis
			durationReplay += timeReplay
			if err != nil {
				fmt.Printf("File %d with Test %d failed, check output.log for more information. Skipping...\n", currentFile, attemptedTests)
				skippedTests++
				continue
			}

			// Move logs and results to the appropriate directory
			moveResults(packagePath, directoryName)
		}
		currentFile++
	}

	resultPath := filepath.Join(dir, "advocateResult")

	// Generate Bug Reports
	fmt.Println("Generate Bug Reports")
	generateBugReports(resultPath, pathToAdvocate)

	if measureTime {
		generateTimeFile(resultPath, durationRun, durationRecord, durationAnalysis, durationReplay)
	}

	// Check for untriggered selects
	if notExecuted {
		fmt.Println("Check for untriggered selects and not executed progs")
		runCommand(pathToAnalyzer, "check", "-R", filepath.Join(dir, "advocateResult"), "-P", dir)
	}

	if stats {
		// create statistics
		fmt.Println("Create statistics")
		runCommand(pathToAnalyzer, "stats", "-R", filepath.Join(dir, "advocateResult"), "-P", dir, "-N", progName)
	}

	// Output test summary
	fmt.Println("Finished full workflow for all tests")
	fmt.Printf("Attempted tests: %d\n", attemptedTests)
	fmt.Printf("Skipped tests: %d\n", skippedTests)

	return nil
}

/*
 * Function to write the time information to a file
 * Args:
 *     folderName (string): path to the destination of the file,
 *     durationRun (time.Duration): time to run all tests
 *     durationRecord (time.Duration): time to record all tests
 *     durationAnalysis (time.Duration): time to analyze all traces
 *     durationReplay (time.Duration): time to run all replays. For each test the avg time is used
 */
func generateTimeFile(folderName string, durationRun, durationRecord, durationAnalysis,
	durationReplay time.Duration) {
	fmt.Println("Generate time file")
	file, err := os.Create(filepath.Join(folderName, "times.log"))
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	timeInfo := fmt.Sprintf(
		"Run: %f\nRecord: %f\nAnalysis: %f\nReplay: %f",
		durationRun.Seconds(), durationRecord.Seconds(),
		durationAnalysis.Seconds(), durationReplay.Seconds())

	fmt.Println(timeInfo)

	_, err = writer.WriteString(timeInfo)
	if err != nil {
		fmt.Println("Failed to write time file: ", err)
	}
	writer.Flush()
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
 *    time.Duration: time of trace recording
 *    time.Duration: time for analysis
 *    time.Duration: avg. time of trace replay
 *    error
 */
func unitTestFullWorkflow(pathToAdvocate string, dir string,
	testName string, pkg string, file string, outputDir string) (time.Duration, time.Duration, time.Duration, error) {

	output := filepath.Join(outputDir, "output.log")

	outFile, err := os.OpenFile(output, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("Failed to open log file: %v", err)
	}
	defer outFile.Close()

	// Redirect stdout and stderr to the file
	origStdout := os.Stdout
	origStderr := os.Stderr

	os.Stdout = outFile
	os.Stderr = outFile

	defer func() {
		os.Stdout = origStdout
		os.Stderr = origStderr
	}()

	// Validate required inputs
	if pathToAdvocate == "" {
		return 0, 0, 0, errors.New("Path to advocate is empty")
	}
	if dir == "" {
		return 0, 0, 0, errors.New("Directory is empty")
	}
	if testName == "" {
		return 0, 0, 0, errors.New("Test name is empty")
	}
	if pkg == "" {
		return 0, 0, 0, errors.New("Package is empty")
	}
	if file == "" {
		return 0, 0, 0, errors.New("Test file is empty")
	}

	pathToPatchedGoRuntime := filepath.Join(pathToAdvocate, "go-patch/bin/go")
	pathToGoRoot := filepath.Join(pathToAdvocate, "go-patch")
	pathToAnalyzer := filepath.Join(pathToAdvocate, "analyzer/analyzer")

	// Change to the directory
	if err := os.Chdir(dir); err != nil {
		return 0, 0, 0, fmt.Errorf("Failed to change directory: %v", err)
	}
	fmt.Printf("In directory: %s\n", dir)

	// Set GOROOT
	os.Setenv("GOROOT", pathToGoRoot)
	fmt.Println("GOROOT exported")

	fmt.Println("FileName: ", file)
	fmt.Println("TestName: ", testName)

	// Remove header just in case
	fmt.Println(fmt.Sprintf("Remove header for %s", file))
	if err := headerRemoverUnit(file); err != nil {
		return 0, 0, 0, fmt.Errorf("Failed to remove header: %v", err)
	}

	// Add header
	fmt.Println(fmt.Sprintf("Add header for %s: %s", file, testName))
	if err := headerInserterUnit(file, testName, false, "0"); err != nil {
		return 0, 0, 0, fmt.Errorf("Error in adding header: %v", err)
	}

	// Run the test
	fmt.Println(fmt.Sprintf("%s test -count=1 -run=%s ./%s", pathToPatchedGoRuntime, testName, pkg))

	timeStart := time.Now()
	err = runCommand(pathToPatchedGoRuntime, "test", "-count=1", "-run="+testName, "./"+pkg)
	if err != nil {
		log.Println(err)
		log.Println("Test failed, removing header and exiting")
		headerRemoverUnit(file)
		return 0, 0, 0, errors.New("Error running test")
	}
	durationRecord := time.Since(timeStart)

	// Remove header after the test
	fmt.Println(fmt.Sprintf("Remove header for %s", file))
	headerRemoverUnit(file)

	// Apply analyzer
	fmt.Println(fmt.Sprintf("Run the analyzer for %s/%s/advocateTrace", dir, pkg))
	startTime := time.Now()
	err = runCommand(pathToAnalyzer, "run", "-t", filepath.Join(dir, pkg, "advocateTrace"))
	if err != nil {
		log.Println("Analyzer failed", err)
	}
	timeAnalysis := time.Since(startTime)
	fmt.Println("Finished Analyzer")

	pathPkg := filepath.Join(dir, pkg)
	rewrittenTraces, _ := filepath.Glob(filepath.Join(pathPkg, "rewritten_trace_*"))
	fmt.Printf("Found %d rewritten traces\n", len(rewrittenTraces))

	var timeReplay time.Duration
	for _, trace := range rewrittenTraces {
		traceNum := extractTraceNumber(trace)
		fmt.Printf("Insert replay header or %s: %s for trace %s\n", file, testName, traceNum)
		headerInserterUnit(file, testName, true, traceNum)

		fmt.Printf("%s test -count=1 -run=%s ./%s\n", pathToPatchedGoRuntime, testName, pkg)
		startTime := time.Now()
		runCommand(pathToPatchedGoRuntime, "test", "-count=1", "-run="+testName, "./"+pkg)
		timeReplay += time.Since(startTime)

		// Remove reorder header
		fmt.Printf("Remove header for %s\n", file)
		headerRemoverUnit(file)
	}

	if len(rewrittenTraces) > 0 {
		timeReplay = timeReplay / time.Duration(len(rewrittenTraces))
	}

	// Unset GOROOT
	os.Unsetenv("GOROOT")
	fmt.Println("GOROOT removed")

	return durationRecord, timeAnalysis, timeReplay, nil
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
	parts := strings.Split(trace, "rewritten_trace_")
	if len(parts) > 1 {
		return parts[1]
	}
	return ""
}
