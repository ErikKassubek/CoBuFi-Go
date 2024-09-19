// Copyrigth (c) 2024 Erik Kassubek
//
// File: main.go
// Brief: Main file and starting point for the analyzer
//
// Author: Erik Kassubek <kassubek.erik@gmail.com>
// Created: 2023-08-08
// LastChange: 2024-09-01
//
// License: BSD-3-Clause

package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	"analyzer/complete"
	"analyzer/explanation"
	"analyzer/io"
	"analyzer/logging"
	"analyzer/rewriter"
	"analyzer/stats"
	"analyzer/trace"
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
	ignoreAtomics := flag.Bool("a", false, "Ignore atomic operations (default false). Use to reduce memory overhead for large traces.")
	explanationIndex := flag.Int("i", 0, "Index of the explanation to create")
	resultFolderTool := flag.String("R", "", "Path where the advocateResult folder created by the pipeline is located")
	programPath := flag.String("P", "", "Path to the program folder")
	preventCopyRewrittenTrace := flag.Bool("n", false, "Do not copy the rewritten trace in the explanation")

	scenarios := flag.String("s", "", "Select which analysis scenario to run, e.g. -s srd for the option s, r and d."+
		"If not set, all scenarios are run.\n"+
		"Options:\n"+
		"\ts: Send on closed channel\n"+
		"\tr: Receive on closed channel\n"+
		"\tw: Done before add on waitGroup\n"+
		"\tn: Close of closed channel\n"+
		"\tb: Concurrent receive on channel\n"+
		"\tl: Leaking routine\n"+
		"\tu: Select case without partner\n",
	)
	// "\tc: Cyclic deadlock\n",
	// "\tm: Mixed deadlock\n"

	// go memorySupervisor() // BUG: does not work properly

	flag.Parse()

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

	outMachine := *resultFolder + "/results_machine.log"
	outReadable := *resultFolder + "/results_readable.log"
	newTrace := *resultFolder + "/rewritten_trace"

	if len(os.Args) >= 2 {

		// ===================== Special cases =====================

		switch os.Args[1] {
		case "stats":
			modeStats(programPath, pathTrace)
		case "explain":
			modeExplain(pathTrace, folderTrace, explanationIndex, preventCopyRewrittenTrace)
		case "check":
			modeCheck(resultFolderTool, programPath)
		case "run":
			modeRun(pathTrace, noPrint, noRewrite, scenarios, level, outReadable,
				outMachine, ignoreAtomics, fifo, ignoreCriticalSection,
				noWarning, folderTrace, newTrace)
		default:
			fmt.Printf("Unknown mode %s", os.Args[1])
			fmt.Printf("Select one mode from 'run', 'stats', 'explain' or 'check'")
			printHelp()
		}
	} else {
		fmt.Printf("No mode selected")
		fmt.Printf("Select one mode from 'run', 'stats', 'explain' or 'check'")
		printHelp()
	}

	// ============== Start the normal program ==============

}

func modeStats(programPath, pathTrace *string) {
	// instead of the normal program, create statistics for the trace
	stats.Create(programPath, pathTrace)
}

func modeExplain(pathTrace *string, folderTrace string, explanationIndex *int,
	preventCopyRewrittenTrace *bool) {
	if *pathTrace == "" || *explanationIndex == 0 {
		fmt.Println("Please provide a path to the trace file and an index (1 based) for the explanation. Set with -t [file] -i [index]")
		return
	}
	err := explanation.CreateOverview(folderTrace, *explanationIndex, *preventCopyRewrittenTrace)
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
	noWarning *bool, folderTrace string, newTrace string) {
	printHeader()

	startTime := time.Now()

	if *pathTrace == "" {
		fmt.Println("Please provide a path to the trace file. Set with -t [file]")
		return
	}

	if *noPrint {
		*noRewrite = true
	}

	analysisCases, err := parseAnalysisCases(*scenarios)
	if err != nil {
		panic(err)
	}

	// run the analysis and, if requested, create a reordered trace file
	// based on the analysis results

	logging.InitLogging(*level, outReadable, outMachine)
	numberOfRoutines, err := io.CreateTraceFromFiles(*pathTrace, *ignoreAtomics)
	if err != nil {
		panic(err)
	}
	trace.SetNumberOfRoutines(numberOfRoutines)

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

	trace.RunAnalysis(*fifo, *ignoreCriticalSection, analysisCases)

	fmt.Print("Analysis finished\n\n")

	numberOfResults := logging.PrintSummary(*noWarning, *noPrint)

	analysisFinishedTime := time.Now()
	err = writeTime(folderTrace, "Analysis", analysisFinishedTime.Sub(startTime).Seconds())
	if err != nil {
		println("Could not write time to file: ", err.Error())
	}

	if !*noRewrite {
		numberRewrittenTrace := 0
		failedRewrites := 0
		notNeededRewrites := 0
		println("\n\nStart rewriting trace files...")
		var rewriteTime time.Duration
		originalTrace := trace.CopyCurrentTrace()
		for resultIndex := 0; resultIndex < numberOfResults; resultIndex++ {
			rewriteStartTime := time.Now()

			needed, err := rewriteTrace(outMachine,
				newTrace+"_"+strconv.Itoa(resultIndex+1)+"/", resultIndex, numberOfRoutines)

			if !needed {
				println("Trace can not be rewritten.")
				notNeededRewrites++
			} else if err != nil {
				println("Failed to rewrite trace: ", err.Error())
				failedRewrites++
				trace.SetTrace(originalTrace)
			} else { // needed && err == nil
				numberRewrittenTrace++
				rewriteTime += time.Now().Sub(rewriteStartTime)
				trace.SetTrace(originalTrace)
			}

			print("\n\n")
		}

		err = writeTime(folderTrace, "AvgRewrite", rewriteTime.Seconds()/float64(numberRewrittenTrace))
		if err != nil {
			println("Could not write time to file: ", err.Error())
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

func memorySupervisor() {
	var stat syscall.Sysinfo_t

	for {
		time.Sleep(2 * time.Second)

		err := syscall.Sysinfo(&stat)
		if err != nil {
			panic(err)
		}

		freeRAM := stat.Freeram * uint64(stat.Unit)

		if freeRAM < 3000000000 {
			println("Not enough free RAM available. Exiting...")
			os.Exit(1)
		}
	}
}

func writeTime(pathTrace string, name string, time float64) error {
	path := pathTrace
	if path[len(path)-1] != os.PathSeparator {
		path += string(os.PathSeparator)
	}
	path += "times.log"

	// Datei lesen
	content, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			// create file
			err = os.WriteFile(path, []byte(""), 0644)
			if err != nil {
				return err
			}
		}
	}

	elems := strings.Split(string(content), "\n")

	found := false
	names := make([]string, 0)
	values := make([]string, 0)
	if len(elems) >= 2 {

		names = strings.Split(elems[0], ",")
		values = strings.Split(elems[1], ",")

		// if name already exists, overwrite the value, if name exists multiple time, delete all and write new
		remove := make([]int, 0)
		for i, n := range names {
			if n == name {
				if !found {
					values[i] = strconv.FormatFloat(time, 'f', 6, 64)
					found = true
				} else {
					remove = append(remove, i)
				}
			}
		}

		// remove all duplicates
		for i := len(remove) - 1; i >= 0; i-- {
			names = append(names[:remove[i]], names[remove[i]+1:]...)
			values = append(values[:remove[i]], values[remove[i]+1:]...)
		}
	}

	// if name not found, append
	if !found {
		names = append(names, name)
		values = append(values, strconv.FormatFloat(time, 'f', 6, 64))
	}

	elem1 := strings.Join(names, ",")
	elem2 := strings.Join(values, ",")

	// Datei schreiben
	err = os.WriteFile(path, []byte(elem1+"\n"+elem2), 0644)
	return err
}

/*
 * Rewrite the trace file based on given analysis results
 * Args:
 *   outMachine (string): The path to the analysis result file
 *   newTrace (string): The path where the new traces folder will be created
 *   resultIndex (int): The index of the result to use for the reordered trace file
 *   numberOfRoutines (int): The number of routines in the trace
 * Returns:
 *   bool: true, if a rewrite was nessesary, false if not (e.g. actual bug, warning)
 *   error: An error if the trace file could not be created
 */
func rewriteTrace(outMachine string, newTrace string, resultIndex int,
	numberOfRoutines int) (bool, error) {

	actual, bug, err := io.ReadAnalysisResults(outMachine, resultIndex)
	if err != nil {
		return false, err
	}

	if actual {
		return false, nil
	}

	rewriteNeeded, code, err := rewriter.RewriteTrace(bug)

	if err != nil {
		return rewriteNeeded, err
	}

	err = io.WriteTrace(newTrace, numberOfRoutines)
	if err != nil {
		return rewriteNeeded, err
	}

	err = io.WriteRewriteInfoFile(newTrace, string(bug.Type), code, resultIndex)
	if err != nil {
		return rewriteNeeded, err
	}

	return rewriteNeeded, nil
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
		case 'u':
			analysisCases["selectWithoutPartner"] = true
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
	println("  -a          Ignore atomic operations (default false). Use to reduce memory overhead for large traces.")
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
	println("  -t [file]   Path to the trace file to create the explanation for (required)")
	println("  -i [index]  Index of the explanation to create (1 based) (required)")
	println("\n\n")
	println("3. Check if all concurrency elements of the program have been executed at least once")
	println("This mode checks if all concurrency elements of the program have been executed at least once.")
	println("It has the following options:")
	println("  -R [folder] Path where the advocateResult folder created by the pipeline is located (required)")
	println("  -P [folder] Path to the program folder (required)")
	println("\n\n")
	println("4. Create statistics about a program")
	println("This creates some statistics about the program and the trace")
	println("Usage: ./analyzer stats [options]")
	println("  -P [folder] Path to the program folder (required)")
	println("  -t [file]   Path to the trace folder (required)")
	println("\n")
}
