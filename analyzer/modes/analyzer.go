// Copyright (c) 2025 Erik Kassubek
//
// File: analyzer.go
// Brief: Start the analyze mode
//
// Author: Erik Kassubek
// Created: 2025-01-05
//
// License: BSD-3-Clause

package modes

import (
	"analyzer/analysis"
	"analyzer/bugs"
	"analyzer/io"
	"analyzer/results"
	"analyzer/rewriter"
	timemeasurement "analyzer/timeMeasurement"
	"analyzer/utils"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func ModeAnalyzer(pathTrace *string, noPrint *bool, noRewrite *bool,
	scenarios *string, outReadable string, outMachine string,
	ignoreAtomics *bool, fifo *bool, ignoreCriticalSection *bool,
	noWarning *bool, rewriteAll *bool, folderTrace string, newTrace string, timeout *int, ignoreRewrite *string) {
	// printHeader()

	if *pathTrace == "" {
		fmt.Println("Please provide a path to the trace files. Set with -trace [folder]")
		return
	}

	if *noPrint {
		*noRewrite = true
	}

	// set timeout
	if timeout != nil && *timeout > 0 {
		go func() {
			<-time.After(time.Duration(*timeout) * time.Second)
			os.Exit(1)
		}()
	}

	analysisCases, err := parseAnalysisCases(*scenarios)
	if err != nil {
		panic(err)
	}

	// clean data in case of fuzzing
	if analysis.DataUsed {
		analysis.ClearData()
		analysis.ClearTrace()
		analysis.DataUsed = true
	}

	// run the analysis and, if requested, create a reordered trace file
	// based on the analysis results

	results.InitResults(outReadable, outMachine)

	// done and separate routine to implement timeout
	done := make(chan bool)
	numberOfRoutines := 0
	containsElems := false
	go func() {
		defer func() { done <- true }()

		numberOfRoutines, containsElems, err = io.CreateTraceFromFiles(*pathTrace, *ignoreAtomics)
		if err != nil {
			panic(err)
		}

		if !containsElems {
			fmt.Println("Trace does not contain any elem")
			fmt.Println("Skip analysis")
			return
		}

		analysis.SetNumberOfRoutines(numberOfRoutines)

		if analysisCases["all"] {
			fmt.Println("Start Analysis for all scenarios")
		} else {
			fmt.Println("Start Analysis for the following scenarios:")
			for key, value := range analysisCases {
				if value {
					fmt.Println("\t", key)
				}
			}
		}

		timemeasurement.Start("analysis")
		analysis.RunAnalysis(*fifo, *ignoreCriticalSection, analysisCases)
		timemeasurement.End("analysis")

		timemeasurement.Print()
	}()

	if timeout != nil && *timeout > 0 {
		select {
		case <-done:
			fmt.Print("Analysis finished\n\n")
		case <-time.After(time.Duration(*timeout) * time.Second):
			fmt.Printf("Analysis ended by timeout after %d seconds\n\n", *timeout)
		}
	} else {
		<-done
	}

	numberOfResults := results.PrintSummary(*noWarning, *noPrint)

	if !*noRewrite {
		numberRewrittenTrace := 0
		failedRewrites := 0
		notNeededRewrites := 0
		println("\n\nStart rewriting trace file ", *pathTrace)
		originalTrace := analysis.CopyCurrentTrace()

		analysis.ClearData()

		rewrittenBugs := make(map[bugs.ResultType][]string) // bugtype -> paths string

		addAlreadyProcessed(rewrittenBugs, *ignoreRewrite)

		file := filepath.Base(*pathTrace)
		rewriteNr := "0"
		spl := strings.Split(file, "_")
		if len(spl) > 1 {
			rewriteNr = spl[len(spl)-1]
		}

		for resultIndex := 0; resultIndex < numberOfResults; resultIndex++ {
			needed, double, err := rewriteTrace(outMachine,
				newTrace+"_"+strconv.Itoa(resultIndex+1)+"/", resultIndex, numberOfRoutines, &rewrittenBugs, !*rewriteAll)

			if !needed {
				println("Trace can not be rewritten.")
				notNeededRewrites++
				if double {
					fmt.Printf("Bugreport info: %s_%d,double", rewriteNr, resultIndex+1)
				} else {
					fmt.Printf("Bugreport info: %s_%d,fail", rewriteNr, resultIndex+1)
				}
			} else if err != nil {
				println("Failed to rewrite trace: ", err.Error())
				failedRewrites++
				analysis.SetTrace(originalTrace)
				fmt.Printf("Bugreport info: %s_%d,fail", rewriteNr, resultIndex+1)
			} else { // needed && err == nil
				numberRewrittenTrace++
				analysis.SetTrace(originalTrace)
				fmt.Printf("Bugreport info: %s_%d,suc", rewriteNr, resultIndex+1)
			}

			print("\n\n")
		}

		println("Finished Rewrite")
		println("\n\n\tNumber Results: ", numberOfResults)
		println("\tSuccessfully rewrites: ", numberRewrittenTrace)
		println("\tNo need/not possible to rewrite: ", notNeededRewrites)
		if failedRewrites > 0 {
			println("\tFailed rewrites: ", failedRewrites)
		} else {
			println("\tFailed rewrites: ", failedRewrites)
		}
	}

	print("\n\n\n")
}

/*
 * Parse the given analysis cases
 * Args:
 *   cases (string): The string of analysis cases to parse
 * Returns:
 *   map[string]bool: A map of the analysis cases and if they are set
 *   error: An error if the cases could not be parsed
 */
func parseAnalysisCases(cases string) (map[string]bool, error) {
	analysisCases := map[string]bool{
		"all":                  false, // all cases enabled
		"sendOnClosed":         false,
		"receiveOnClosed":      false,
		"doneBeforeAdd":        false,
		"closeOnClosed":        false,
		"concurrentRecv":       false,
		"leak":                 false,
		"selectWithoutPartner": false,
		"cyclicDeadlock":       false,
		"mixedDeadlock":        false,
	}

	if cases == "" {
		analysisCases["all"] = true
		analysisCases["sendOnClosed"] = true
		analysisCases["receiveOnClosed"] = true
		analysisCases["doneBeforeAdd"] = true
		analysisCases["closeOnClosed"] = true
		analysisCases["concurrentRecv"] = true
		analysisCases["leak"] = true
		analysisCases["selectWithoutPartner"] = true
		analysisCases["unlockBeforeLock"] = true
		// analysisCases["cyclicDeadlock"] = true
		// analysisCases["mixedDeadlock"] = true

		return analysisCases, nil
	}

	for _, c := range cases {
		switch c {
		case 's':
			analysisCases["sendOnClosed"] = true
		case 'r':
			analysisCases["receiveOnClosed"] = true
		case 'w':
			analysisCases["doneBeforeAdd"] = true
		case 'n':
			analysisCases["closeOnClosed"] = true
		case 'b':
			analysisCases["concurrentRecv"] = true
		case 'l':
			analysisCases["leak"] = true
		case 'p':
			analysisCases["selectWithoutPartner"] = true
		case 'u':
			analysisCases["unlockBeforeLock"] = true
		// case 'c':
		// 	analysisCases["cyclicDeadlock"] = true
		// case 'm':
		// analysisCases["mixedDeadlock"] = true
		default:
			return nil, fmt.Errorf("Invalid analysis case: %c", c)
		}
	}
	return analysisCases, nil
}

func addAlreadyProcessed(alreadyProcessed map[bugs.ResultType][]string, ignoreRewrite string) {
	if ignoreRewrite == "" {
		return
	}

	data, err := os.ReadFile(ignoreRewrite)
	if err != nil {
		return
	}
	for _, bugStr := range strings.Split(string(data), "\n") {
		_, bug, err := bugs.ProcessBug(bugStr)
		if err != nil {
			continue
		}

		if _, ok := alreadyProcessed[bug.Type]; !ok {
			alreadyProcessed[bug.Type] = make([]string, 0)
		} else {
			if utils.ContainsString(alreadyProcessed[bug.Type], bugStr) {
				continue
			}
		}
		alreadyProcessed[bug.Type] = append(alreadyProcessed[bug.Type], bugStr)
	}
}

/*
 * Rewrite the trace file based on given analysis results
 * Args:
 *   outMachine (string): The path to the analysis result file
 *   newTrace (string): The path where the new traces folder will be created
 *   resultIndex (int): The index of the result to use for the reordered trace file
 *   numberOfRoutines (int): The number of routines in the trace
 *   rewrittenTrace (*map[string][]string): set of bugs that have been already rewritten
 * Returns:
 *   bool: true, if a rewrite was nessesary, false if not (e.g. actual bug, warning)
 *   bool: true if rewrite was skipped because of double
 *   error: An error if the trace file could not be created
 */
func rewriteTrace(outMachine string, newTrace string, resultIndex int,
	numberOfRoutines int, rewrittenTrace *map[bugs.ResultType][]string, rewriteOnce bool) (bool, bool, error) {

	actual, bug, err := io.ReadAnalysisResults(outMachine, resultIndex)
	if err != nil {
		return false, false, err
	}

	if rewriteOnce {
		bugString := bug.GetBugString()
		println(resultIndex, bugString)
		if _, ok := (*rewrittenTrace)[bug.Type]; !ok {
			(*rewrittenTrace)[bug.Type] = make([]string, 0)
		} else {
			if utils.ContainsString((*rewrittenTrace)[bug.Type], bugString) {
				fmt.Println("Bug was already rewritten before")
				fmt.Println("Skip rewrite")
				return false, true, nil
			}
		}
		(*rewrittenTrace)[bug.Type] = append((*rewrittenTrace)[bug.Type], bugString)
	}

	if actual {
		return false, false, nil
	}

	rewriteNeeded, code, err := rewriter.RewriteTrace(bug, 0)

	if err != nil {
		return rewriteNeeded, false, err
	}

	err = io.WriteTrace(newTrace, numberOfRoutines)
	if err != nil {
		return rewriteNeeded, false, err
	}

	err = io.WriteRewriteInfoFile(newTrace, string(bug.Type), code, resultIndex)
	if err != nil {
		return rewriteNeeded, false, err
	}

	return rewriteNeeded, false, nil
}
