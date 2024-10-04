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
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	timeout = "10m"
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
 *    timeout (int): Set a timeout in seconds for the analysis
 * Returns:
 *    error
 */
func runWorkflowUnit(pathToAdvocate, dir, progName string,
	measureTime, notExecuted, stats bool, timeoutAna int) error {
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
		return fmt.Errorf("Failed to change directory: %v", dir)
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

	resultPath := filepath.Join(dir, "advocateResult")

	// Process each test file
	for _, file := range testFiles {
		fmt.Printf("\n\nProgress: %d/%d\n", currentFile, totalFiles)
		fmt.Printf("\nProcessing file: %s\n", file)

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
			directoryPath := filepath.Join(dir, directoryName)
			if err := os.MkdirAll(directoryName, os.ModePerm); err != nil {
				log.Printf("Failed to create directory %s: %v", directoryName, err)
				continue
			}

			// Execute full workflow
			times, nrReplay, err := unitTestFullWorkflow(pathToAdvocate, dir, testFunc, adjustedPackagePath, file, directoryName, measureTime, timeoutAna)

			if measureTime {
				updateTimeFiles(progName, testFunc, resultPath, times, nrReplay)
			}

			// Move logs and results to the appropriate directory
			moveResults(packagePath, directoryName)

			if err != nil {
				fmt.Printf("File %d with Test %d failed, check output.log for more information.\n", currentFile, attemptedTests)
				skippedTests++
			}

			generateBugReports(directoryPath, pathToAdvocate)

			if stats {
				updateStatsFiles(pathToAnalyzer, progName, directoryPath)
				// create statistics
			}
		}

		currentFile++
	}

	// Check for untriggered selects
	if notExecuted {
		fmt.Println("Check for untriggered selects and not executed progs")
		err := runCommand(pathToAnalyzer, "check", "-R", filepath.Join(dir, "advocateResult"), "-P", dir)
		if err != nil {
			fmt.Println("Could not run check for untriggered select and not executed progs")
		}
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
 *     progName (string): name of the program
 *     folderName (string): path to the destination of the file,
 *     time (map[string]time.Durations): runtimes
 *     numberReplay (int): number of replay
 */
func updateTimeFiles(progName string, testName string, folderName string, times map[string]time.Duration, numberReplay int) {
	// fmt.Println("Generate time file")

	timeFilePath := filepath.Join(folderName, "times_"+progName+".csv")

	newFile := false
	_, err := os.Stat(timeFilePath)
	if os.IsNotExist(err) {
		newFile = true
	}

	file, err := os.OpenFile(timeFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening or creating file:", err)
		return
	}
	defer file.Close()

	if newFile {
		csvTitels := "TestName,ExecTime,ExecTimeWithTracing,AnalyzerTime,AnalysisTime,HBAnalysisTime,TimeToIdentifyLeaksPlusFindingPoentialPartners,TimeToIdentifyPanicBugs,ReplayTime,NumberReplay\n"
		if _, err := file.WriteString(csvTitels); err != nil {
			fmt.Println("Could not write time: ", err)
		}
	}

	timeInfo := fmt.Sprintf(
		"%s,%.5f#%.5f#%.5f#%.5f#%.5f#%.5f#%.5f#%.5f#%d\n", testName,
		times["run"].Seconds(), times["record"].Seconds(),
		times["analyzer"].Seconds(), times["analysis"].Seconds(),
		times["hb"].Seconds(), times["leak"].Seconds(), times["panic"].Seconds(), times["replay"].Seconds(),
		numberReplay)

	// Write to the file
	if _, err := file.WriteString(timeInfo); err != nil {
		fmt.Println("Could not write time: ", err)
	}
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
		if strings.HasPrefix(line, "func Test") && strings.Contains(line, "*testing.T") {
			testFunc := strings.TrimSpace(strings.Split(line, "(")[0])
			testFunc = strings.TrimPrefix(testFunc, "func ")
			testFunctions = append(testFunctions, testFunc)
		}
	}
	return testFunctions, nil
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
 *    measureTime (bool): if true, run the test once without recording/replay to measure time
 *    timeout (int): Set a timeout in seconds for the analysis
 * Returns:
 *    map[string]time.Duration
 *    int: number of run replays
 *    error
 */
func unitTestFullWorkflow(pathToAdvocate string, dir string,
	testName string, pkg string, file string, outputDir string,
	measureTime bool, timeoutAna int) (map[string]time.Duration, int, error) {

	resTimes := make(map[string]time.Duration)

	output := filepath.Join(outputDir, "output.log")

	outFile, err := os.OpenFile(output, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return resTimes, 0, fmt.Errorf("Failed to open log file: %v", err)
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
		return resTimes, 0, errors.New("Path to advocate is empty")
	}
	if dir == "" {
		return resTimes, 0, errors.New("Directory is empty")
	}
	if testName == "" {
		return resTimes, 0, errors.New("Test name is empty")
	}
	// if pkg == "" {
	// 	return 0, 0, 0, 0, errors.New("Package is empty")
	// }
	if file == "" {
		return resTimes, 0, errors.New("Test file is empty")
	}

	pathToPatchedGoRuntime := filepath.Join(pathToAdvocate, "go-patch/bin/go")
	pathToGoRoot := filepath.Join(pathToAdvocate, "go-patch")
	pathToAnalyzer := filepath.Join(pathToAdvocate, "analyzer/analyzer")

	// Change to the directory
	if err := os.Chdir(dir); err != nil {
		return resTimes, 0, fmt.Errorf("Failed to change directory: %v", err)
	}
	fmt.Printf("In directory: %s\n", dir)

	// run the tests without recording/replay
	resTimes["run"] = time.Duration(0)
	if measureTime {
		// Remove header just in case
		if err := headerRemoverUnit(file); err != nil {
			fmt.Println(err)
		}

		os.Unsetenv("GOROOT")
		fmt.Println("GOROOT removed")

		timeStart := time.Now()
		fmt.Println("Run T0")
		err = runCommand("go", "test", "-timeout", timeout, "-count=1", "-run="+testName, "./"+pkg)
		if err != nil {
			log.Println("Test failed: ", err)
		}
		resTimes["run"] = time.Since(timeStart)
	}

	fmt.Println("FileName: ", file)
	fmt.Println("TestName: ", testName)

	// Remove header just in case
	fmt.Println(fmt.Sprintf("Remove header for %s", file))
	if err := headerRemoverUnit(file); err != nil {
		fmt.Printf("Error in removing header: %v\n", err)
		return resTimes, 0, fmt.Errorf("Failed to remove header: %v", err)
	}

	// Add header
	fmt.Println(fmt.Sprintf("Add header for %s: %s", file, testName))
	if err := headerInserterUnit(file, testName, false, "0"); err != nil {
		fmt.Printf("Error in adding header: %v\n", err)
		return resTimes, 0, fmt.Errorf("Error in adding header: %v", err)
	}

	// Run the test
	fmt.Println("\nRun Recording")

	// Set GOROOT
	os.Setenv("GOROOT", pathToGoRoot)
	fmt.Println("GOROOT = " + pathToGoRoot + " exported")

	timeStart := time.Now()
	err = runCommand(pathToPatchedGoRuntime, "test", "-timeout", timeout, "-count=1", "-run="+testName, "./"+pkg)
	if err != nil {
		log.Println(err)
		// log.Println("Test failed, removing header and exiting")
		// headerRemoverUnit(file)
		// return 0, 0, 0, 0, errors.New("Error running test")
	}
	resTimes["record"] = time.Since(timeStart)

	os.Unsetenv("GOROOT")
	fmt.Println("GOROOT removed")

	// Remove header after the test
	fmt.Println(fmt.Sprintf("Remove header for %s", file))
	headerRemoverUnit(file)

	// Apply analyzer
	fmt.Println(fmt.Sprintf("Run the analyzer for %s/%s/advocateTrace", dir, pkg))
	startTime := time.Now()
	err = runCommand(pathToAnalyzer, "run", "-t", filepath.Join(dir, pkg, "advocateTrace"), "-T", strconv.Itoa(timeoutAna))
	if err != nil {
		log.Println("Analyzer failed", err)
	}
	resTimes["analyzer"] = time.Since(startTime)

	fileOuputRead, err := os.OpenFile(output, os.O_RDONLY, 0644)
	if err != nil {
		fmt.Println("Could not open file: ", err)
	}
	outFileContent, err := io.ReadAll(fileOuputRead)
	if err != nil {
		fmt.Println("Could not read file: ", err)
	}
	lines := strings.Split(string(outFileContent), "\n")

	durationOther := time.Duration(0)
	for _, line := range lines {
		if strings.HasPrefix(line, "AdvocateAnalysisTimes:") {
			line = strings.TrimPrefix(line, "AdvocateAnalysisTimes:")
			elems := strings.Split(line, "#")

			timeAnaFloat, _ := strconv.ParseFloat(elems[0], 64)
			timeLeakFloat, _ := strconv.ParseFloat(elems[1], 64)
			timePanicFloat, _ := strconv.ParseFloat(elems[2], 64)
			timeOtherFloat, _ := strconv.ParseFloat(elems[3], 64)

			resTimes["analysis"] = time.Duration(timeAnaFloat * float64(time.Second))
			resTimes["leak"] = time.Duration(timeLeakFloat * float64(time.Second))
			resTimes["panic"] = time.Duration(timePanicFloat * float64(time.Second))
			durationOther = time.Duration(timeOtherFloat * float64(time.Second))
		}
	}
	fileOuputRead.Close()
	resTimes["hb"] = resTimes["analysis"] - resTimes["leak"] - resTimes["panic"] - durationOther
	if resTimes["hb"] < 0 {
		resTimes["hb"] = 0
	}

	fmt.Println("Finished Analyzer")

	pathPkg := filepath.Join(dir, pkg)
	rewrittenTraces, _ := filepath.Glob(filepath.Join(pathPkg, "rewritten_trace_*"))
	fmt.Printf("Found %d rewritten traces\n", len(rewrittenTraces))

	resTimes["replay"] = time.Duration(0)
	for i, trace := range rewrittenTraces {
		traceNum := extractTraceNumber(trace)
		fmt.Printf("Insert replay header or %s: %s for trace %s\n", file, testName, traceNum)
		headerInserterUnit(file, testName, true, traceNum)

		os.Setenv("GOROOT", pathToGoRoot)
		fmt.Println("GOROOT = " + pathToGoRoot + " exported")

		fmt.Printf("\nRun replay %d/%d\n", i+1, len(rewrittenTraces))
		startTime := time.Now()
		runCommand(pathToPatchedGoRuntime, "test", "-timeout", timeout, "-count=1", "-run="+testName, "./"+pkg)
		resTimes["replay"] += time.Since(startTime)

		os.Unsetenv("GOROOT")
		fmt.Println("GOROOT removed")

		// Remove reorder header
		fmt.Printf("Remove header for %s\n", file)
		headerRemoverUnit(file)
	}

	return resTimes, len(rewrittenTraces), nil
}

// extractTraceNumber extracts the numeric part from a trace directory name
func extractTraceNumber(trace string) string {
	parts := strings.Split(trace, "rewritten_trace_")
	if len(parts) > 1 {
		return parts[1]
	}
	return ""
}
