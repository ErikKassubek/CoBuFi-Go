// Copyrigth (c) 2024 Erik Kassubek
//
// File: main.go
// Brief: Main file and starting point for the analyzer
//
// Author: Erik Kassubek
// Created: 2023-08-08
//
// License: BSD-3-Clause

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"analyzer/analysis"
	"analyzer/bugs"
	"analyzer/complete"
	"analyzer/explanation"
	"analyzer/io"
	"analyzer/logging"
	"analyzer/rewriter"
	"analyzer/stats"
	timemeasurement "analyzer/timeMeasurement"
	"analyzer/utils"

	"github.com/shirou/gopsutil/mem"
)

func main() {
	help := flag.Bool("h", false, "Print this help")
	pathTrace := flag.String("t", "", "Path to the trace folder to analyze or rewrite")
	level := flag.Int("d", 1, "Debug Level, 0 = silent, 1 = errors, 2 = info, 3 = debug (default 1)")
	fifo := flag.Bool("f", false, "Assume a FIFO ordering for buffered channels (default false)")
	ignoreCriticalSection := flag.Bool("c", false, "Ignore happens before relations of critical sections (default false)")
	noRewrite := flag.Bool("x", false, "Do not rewrite the trace file (default false)")
	noWarning := flag.Bool("w", false, "Do not print warnings (default false)")
	noPrint := flag.Bool("p", false, "Do not print the results to the terminal (default false). Automatically set -x to true")
	resultFolder := flag.String("r", "", "Path to where the result file should be saved.")
	ignoreAtomics := flag.Bool("a", false, "Ignore atomic operations (default false). Use to reduce memory header for large traces.")
	resultFolderTool := flag.String("R", "", "Path where the advocateResult folder created by the pipeline is located")
	programPath := flag.String("P", "", "Path to the program folder")
	progName := flag.String("N", "", "Name of the program")
	testName := flag.String("M", "", "Name of the test")
	rewriteAll := flag.Bool("S", false, "If a the same position is flagged multiple times, run the replay for each of them. "+
		"If not set, only the first occurence is rewritten")
	timeout := flag.Int("T", -1, "Set a timeout in seconds for the analysis")
	outM := flag.String("outM", "results_machine", "Name for the result machine file")
	outR := flag.String("outR", "results_readable", "Name for the result readable file")
	outT := flag.String("outT", "rewritten_trace", "Name for the rewritten traces")
	ignoreRewrite := flag.String("ignoreRew", "", "Path to a result machine file. If a found bug is already in this file, it will not be rewritten")

	scenarios := flag.String("s", "", "Select which analysis scenario to run, e.g. -s srd for the option s, r and d."+
		"If not set, all scenarios are run.\n"+
		"Options:\n"+
		"\ts: Send on closed channel\n"+
		"\tr: Receive on closed channel\n"+
		"\tw: Done before add on waitGroup\n"+
		"\tn: Close of closed channel\n"+
		"\tb: Concurrent receive on channel\n"+
		"\tl: Leaking routine\n"+
		"\tp: Select case without partner\n"+
		"\tu: Unlock of unlocked mutex\n",
	)
	// "\tc: Cyclic deadlock\n",
	// "\tm: Mixed deadlock\n"

	go memorySupervisor() // panic if not enough ram

	flag.Parse()

	var mode string
	if len(os.Args) > 2 {
		mode = os.Args[1]
		flag.CommandLine.Parse(os.Args[2:])
	} else {
		fmt.Printf("No mode selected")
		fmt.Printf("Select one mode from 'run', 'stats', 'explain' or 'check'")
		printHelp()
	}

	if *help {
		printHelp()
		return
	}

	folderTrace, err := filepath.Abs(*pathTrace)
	if err != nil {
		panic(err)
	}

	// remove last folder from path
	folderTrace = folderTrace[:strings.LastIndex(folderTrace, string(os.PathSeparator))+1]

	if *resultFolder == "" {
		*resultFolder = folderTrace
		if (*resultFolder)[len(*resultFolder)-1] != os.PathSeparator {
			*resultFolder += string(os.PathSeparator)
		}
	}

	outMachine := filepath.Join(*resultFolder, *outM) + ".log"
	outReadable := filepath.Join(*resultFolder, *outR) + ".log"
	newTrace := filepath.Join(*resultFolder, *outT)
	if *ignoreRewrite != "" {
		*ignoreRewrite = filepath.Join(*resultFolder, *ignoreRewrite)
	}

	switch mode {
	case "stats":
		modeStats(*pathTrace, *progName, *testName)
	case "explain":
		modeExplain(pathTrace, !*rewriteAll)
	case "check":
		modeCheck(resultFolderTool, programPath)
	case "run":
		modeRun(pathTrace, noPrint, noRewrite, scenarios, level, outReadable,
			outMachine, ignoreAtomics, fifo, ignoreCriticalSection,
			noWarning, rewriteAll, folderTrace, newTrace, timeout, ignoreRewrite)
	default:
		fmt.Printf("Unknown mode %s", os.Args[1])
		fmt.Printf("Select one mode from 'run', 'stats', 'explain' or 'check'")
		printHelp()
	}
}

func modeStats(pathFolder string, progName string, testName string) {
	// instead of the normal program, create statistics for the trace
	if pathFolder == "" {
		fmt.Println("Provide the path to the folder containing the results_machine file. Set with -t [path]")
		return
	}

	if progName == "" {
		fmt.Println("Provide a name for the analyzed program. Set with -N [name]")
		return
	}

	if testName == "" {
		testName = progName
	}

	stats.CreateStats(pathFolder, progName, testName)
}

func modeExplain(pathTrace *string, ignoreDouble bool) {
	if *pathTrace == "" {
		fmt.Println("Please provide a path to the trace files for the explanation. Set with -t [file]")
		return
	}

	err := explanation.CreateOverview(*pathTrace, ignoreDouble)
	if err != nil {
		fmt.Println("Error creating explanation: ", err.Error())
	}
}

func modeCheck(resultFolderTool, programPath *string) {
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

func modeRun(pathTrace *string, noPrint *bool, noRewrite *bool,
	scenarios *string, level *int, outReadable string, outMachine string,
	ignoreAtomics *bool, fifo *bool, ignoreCriticalSection *bool,
	noWarning *bool, rewriteAll *bool, folderTrace string, newTrace string, timeout *int, ignoreRewrite *string) {
	// printHeader()

	if *pathTrace == "" {
		fmt.Println("Please provide a path to the trace file. Set with -t [file]")
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

	// run the analysis and, if requested, create a reordered trace file
	// based on the analysis results

	logging.InitLogging(*level, outReadable, outMachine)

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

	numberOfResults := logging.PrintSummary(*noWarning, *noPrint)

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
			if utils.Contains(alreadyProcessed[bug.Type], bugStr) {
				continue
			}
		}
		alreadyProcessed[bug.Type] = append(alreadyProcessed[bug.Type], bugStr)
	}
}

func memorySupervisor() {
	thresholdRAM := uint64(1 * 1024 * 1024 * 1024) // 1GB
	thresholdSwap := uint64(200 * 1024 * 1024)     // 200mb
	for {
		// Get the memory stats
		v, err := mem.VirtualMemory()
		if err != nil {
			log.Fatalf("Error getting memory info: %v", err)
		}

		// Get the swap stats
		s, err := mem.SwapMemory()
		if err != nil {
			log.Fatalf("Error getting swap info: %v", err)
		}

		// fmt.Printf("Available RAM: %v MB, Available Swap: %v MB\n", v.Available/1024/1024, s.Free/1024/1024)

		// Panic if available RAM or swap is below the threshold
		if v.Available < thresholdRAM {
			log.Panicf("Available RAM is below threshold! Available: %v MB, Threshold: %v MB", v.Available/1024/1024, thresholdRAM/1024/1024)
		}

		if s.Free < thresholdSwap {
			log.Panicf("Available Swap is below threshold! Available: %v MB, Threshold: %v MB", s.Free/1024/1024, thresholdSwap/1024/1024)
		}

		// Sleep for a while before checking again
		time.Sleep(5 * time.Second)
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
			if utils.Contains((*rewrittenTrace)[bug.Type], bugString) {
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

func printHeader() {
	fmt.Print("\n")
	fmt.Println(" $$$$$$\\  $$$$$$$\\  $$\\    $$\\  $$$$$$\\   $$$$$$\\   $$$$$$\\ $$$$$$$$\\ $$$$$$$$\\ ")
	fmt.Println("$$  __$$\\ $$  __$$\\ $$ |   $$ |$$  __$$\\ $$  __$$\\ $$  __$$\\\\__$$  __|$$  _____|")
	fmt.Println("$$ /  $$ |$$ |  $$ |$$ |   $$ |$$ /  $$ |$$ /  \\__|$$ /  $$ |  $$ |   $$ |      ")
	fmt.Println("$$$$$$$$ |$$ |  $$ |\\$$\\  $$  |$$ |  $$ |$$ |      $$$$$$$$ |  $$ |   $$$$$\\    ")
	fmt.Println("$$  __$$ |$$ |  $$ | \\$$\\$$  / $$ |  $$ |$$ |      $$  __$$ |  $$ |   $$  __|   ")
	fmt.Println("$$ |  $$ |$$ |  $$ |  \\$$$  /  $$ |  $$ |$$ |  $$\\ $$ |  $$ |  $$ |   $$ |      ")
	fmt.Println("$$ |  $$ |$$$$$$$  |   \\$  /    $$$$$$  |\\$$$$$$  |$$ |  $$ |  $$ |   $$$$$$$$\\ ")
	fmt.Println("\\__|  \\__|\\_______/     \\_/     \\______/  \\______/ \\__|  \\__|  \\__|   \\________|")

	headerInfo := "\n\n\n" +
		"Welcome to the trace analyzer and rewriter.\n" +
		"This program analyzes a trace file and detects common concurrency bugs in Go programs.\n" +
		"It can also create a reordered trace file based on the analysis results.\n" +
		"Be aware, that the analysis is based on the trace file and may not be complete.\n" +
		"Be aware, that the analysis may contain false positives and false negatives.\n" +
		"\n" +
		"If the rewrite of a trace file does not create the expected result, it can help to run the\n" +
		"analyzer with the -c flag to ignore the happens before relations of critical sections (mutex lock/unlock operations).\n" +
		"For the first analysis this is not recommended, because it increases the likelihood of false positives." +
		"\n\n\n"

	fmt.Print(headerInfo)
}

func printHelp() {
	println("Usage: ./analyzer [mode] [options]\n")
	println("There are four modes of operation:")
	println("1. Analyze a trace file and create a reordered trace file based on the analysis results (Default)")
	println("2. Create an explanation for a found bug")
	println("3. Check if all concurrency elements of the program have been executed at least once")
	println("4. Create statistics about a program\n\n")
	println("1. Analyze a trace file and create a reordered trace file based on the analysis results (Default)")
	println("This mode is the default mode and analyzes a trace file and creates a reordered trace file based on the analysis results.")
	println("Usage: ./analyzer run [options]")
	println("It has the following options:")
	println("  -t [file]   Path to the trace folder to analyze or rewrite (required)")
	println("  -d [level]  Debug Level, 0 = silent, 1 = errors, 2 = info, 3 = debug (default 1)")
	println("  -f          Assume a FIFO ordering for buffered channels (default false)")
	println("  -c          Ignore happens before relations of critical sections (default false)")
	println("  -x          Do not rewrite the trace file (default false)")
	println("  -w          Do not print warnings (default false)")
	println("  -p          Do not print the results to the terminal (default false). Automatically set -x to true")
	println("  -r [folder] Path to where the result file should be saved. (default parallel to -t)")
	println("  -a          Ignore atomic operations (default false). Use to reduce memory header for large traces.")
	println("  -S          If the same bug is detected multiple times, run the replay for each of them. If not set, only the first occurence is rewritten")
	println("  -T [second] Set a timeout in seconds for the analysis")
	println("  -s [cases]  Select which analysis scenario to run, e.g. -s srd for the option s, r and d.")
	println("              If it is not set, all scenarios are run")
	println("              Options:")
	println("                  s: Send on closed channel")
	println("                  r: Receive on closed channel")
	println("                  w: Done before add on waitGroup")
	println("                  n: Close of closed channel")
	println("                  b: Concurrent receive on channel")
	println("                  l: Leaking routine")
	println("                  u: Select case without partner")
	// println("                  c: Cyclic deadlock")
	// println("                  m: Mixed deadlock")
	println("\n\n")
	println("2. Create an explanation for a found bug")
	println("Usage: ./analyzer explain [options]")
	println("This mode creates an explanation for a found bug in the trace file.")
	println("It has the following options:")
	println("  -t [file]   Path to the folder containing the machine readable result file (required)")
	println("\n\n")
	println("3. Check if all concurrency elements of the program have been executed at least once")
	println("Usage: ./analyzer check [options]")
	println("This mode checks if all concurrency elements of the program have been executed at least once.")
	println("It has the following options:")
	println("  -R [folder] Path where the advocateResult folder created by the pipeline is located (required)")
	println("  -P [folder] Path to the program folder (required)")
	println("\n\n")
	println("4. Create statistics about a program")
	println("This creates some statistics about the program and the trace")
	println("Usage: ./analyzer stats [options]")
	// println("  -P [folder] Path to the program folder (required)")
	println("  -t [file]   Path to the folder containing the results_machine file (required)")
	println("  -N [name]   Name of the program")
	println("  -M [name]   Name of the test")
	println("\n")
}
