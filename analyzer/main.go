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
	"strings"
	"time"

	"analyzer/modes"

	"github.com/shirou/gopsutil/mem"
)

func main() {
	help := flag.Bool("h", false, "Print this help")

	pathToAdvocate := flag.String("advocate", "", "Path to advocate")

	pathTrace := flag.String("trace", "", "Path to the trace folder to analyze or rewrite")
	programPath := flag.String("dir", "", "Path to the program folder, for toolMain: path to main file, for toolTest: path to test folder")

	progName := flag.String("prog", "", "Name of the program")
	testName := flag.String("test", "", "Name of the test")
	execName := flag.String("exec", "", "Name of the executable")

	timeoutAnalysis := flag.Int("timeout", -1, "Set a timeout in seconds for the analysis")
	timeoutReplay := flag.Int("timeoutReplay", -1, "Set a timeout in seconds for the replay")
	recordTime := flag.Bool("time", true, "measure the runtime")

	resultFolder := flag.String("out", "", "Path to where the result file should be saved.")
	resultFolderTool := flag.String("resultTool", "", "Path where the advocateResult folder created by the pipeline is located")
	outM := flag.String("outM", "results_machine", "Name for the result machine file")
	outR := flag.String("outR", "results_readable", "Name for the result readable file")
	outT := flag.String("outT", "rewritten_trace", "Name for the rewritten traces")

	fifo := flag.Bool("fifo", false, "Assume a FIFO ordering for buffered channels (default false)")
	ignoreCriticalSection := flag.Bool("ignCritSec", false, "Ignore happens before relations of critical sections (default false)")
	ignoreAtomics := flag.Bool("ignoreAtomics", false, "Ignore atomic operations (default false). Use to reduce memory header for large traces.")
	ignoreRewrite := flag.String("ignoreRew", "", "Path to a result machine file. If a found bug is already in this file, it will not be rewritten")

	rewriteAll := flag.Bool("rewriteAll", false, "If a the same position is flagged multiple times, run the replay for each of them. "+
		"If not set, only the first occurence is rewritten")

	noRewrite := flag.Bool("noRewrite", false, "Do not rewrite the trace file (default false)")
	noWarning := flag.Bool("noWarning", false, "Do not print warnings (default false)")
	noPrint := flag.Bool("noPrint", false, "Do not print the results to the terminal (default false). Automatically set -noRewrite to true")
	keepTraces := flag.Bool("keepTrace", false, "If set, the traces are not deleted after analysis. Can result in very large output folders")

	notExec := flag.Bool("notExec", false, "Find never executed operations, *notExec, *stats")
	statistics := flag.Bool("stats", false, "Create statistics")

	scenarios := flag.String("scen", "", "Select which analysis scenario to run, e.g. -scen srd for the option s, r and d."+
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

	if *help {
		printHelp()
		return
	}

	var mode string
	if len(os.Args) >= 2 {
		mode = os.Args[1]
		flag.CommandLine.Parse(os.Args[2:])
	} else {
		fmt.Println("No mode selected")
		fmt.Println("Select one mode from 'run', 'stats', 'explain' or 'check'")
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
	case "toolMain":
		modes.ModeToolchain("main", *pathToAdvocate, *programPath, *execName, *progName, *testName, *timeoutAnalysis, *timeoutReplay, 0, *ignoreAtomics, *recordTime, *notExec, *statistics, *keepTraces)
	case "toolTest":
		modes.ModeToolchain("test", *pathToAdvocate, *programPath, "", *progName, *testName, *timeoutAnalysis, *timeoutReplay, 0, *ignoreAtomics, *recordTime, *notExec, *statistics, *keepTraces)
	case "stats":
		modes.ModeStats(*pathTrace, *progName, *testName)
	case "explain":
		modes.ModeExplain(*pathTrace, !*rewriteAll)
	case "check":
		modes.ModeCheck(resultFolderTool, programPath)
	case "run":
		modes.ModeAnalyzer(pathTrace, noPrint, noRewrite, scenarios, outReadable,
			outMachine, ignoreAtomics, fifo, ignoreCriticalSection,
			noWarning, rewriteAll, folderTrace, newTrace, timeoutAnalysis, ignoreRewrite)
	case "fuzzing":
		modes.ModeFuzzing(*pathToAdvocate, *programPath, *progName, *testName)
	default:
		fmt.Printf("Unknown mode %s\n", os.Args[1])
		fmt.Println("Select one mode from 'run', 'stats', 'explain' or 'check'")
		printHelp()
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
		"analyzer with the -ignCritSecflag to ignore the happens before relations of critical sections (mutex lock/unlock operations).\n" +
		"For the first analysis this is not recommended, because it increases the likelihood of false positives." +
		"\n\n\n"

	fmt.Print(headerInfo)
}

func printHelp() {
	println("Usage: ./analyzer [mode] [options]\n")
	println("There are different modes of operation:")
	println("1. Analyze a trace file and create a reordered trace file based on the analysis results (Default)")
	println("2. Create an explanation for a found bug")
	println("3. Check if all concurrency elements of the program have been executed at least once")
	println("4. Create statistics about a program")
	println("5. Run the toolchain on tests")
	println("6. Run the toolchain on a main function")
	println("7. Create new runs for fuzzing\n\n")
	println("1. Analyze a trace file and create a reordered trace file based on the analysis results (Default)")
	println("This mode is the default mode and analyzes a trace file and creates a reordered trace file based on the analysis results.")
	println("Usage: ./analyzer run [options]")
	println("It has the following options:")
	println("  -trace [file]          Path to the trace folder to analyze or rewrite (required)")
	println("  -fifo                  Assume a FIFO ordering for buffered channels (default false)")
	println("  -ignCritSec            Ignore happens before relations of critical sections (default false)")
	println("  -noRewrite             Do not rewrite the trace file (default false)")
	println("  -noWarning             Do not print warnings (default false)")
	println("  -noPrint               Do not print the results to the terminal (default false). Automatically set -noRewrite to true")
	println("  -keepTrace             Do not delete the trace files after analysis finished")
	println("  -out [folder]          Path to where the result file should be saved. (default parallel to -t)")
	println("  -ignoreAtomics         Ignore atomic operations (default false). Use to reduce memory header for large traces.")
	println("  -rewriteAll            If the same bug is detected multiple times, run the replay for each of them. If not set, only the first occurence is rewritten")
	println("  -timeout [second]      Set a timeout in seconds for the analysis")
	println("  -scen [cases]          Select which analysis scenario to run, e.g. -scen srd for the option s, r and d.")
	println("                         If it is not set, all scenarios are run")
	println("                         Options:")
	println("                             s: Send on closed channel")
	println("                             r: Receive on closed channel")
	println("                             w: Done before add on waitGroup")
	println("                             n: Close of closed channel")
	println("                             b: Concurrent receive on channel")
	println("                             l: Leaking routine")
	println("                             u: Select case without partner")
	// println("                             c: Cyclic deadlock")
	// println("                             m: Mixed deadlock")
	println("\n\n")
	println("2. Create an explanation for a found bug")
	println("Usage: ./analyzer explain [options]")
	println("This mode creates an explanation for a found bug in the trace file.")
	println("It has the following options:")
	println("  -trace [file]          Path to the folder containing the machine readable result file (required)")
	println("\n\n")
	println("3. Check if all concurrency elements of the program have been executed at least once")
	println("Usage: ./analyzer check [options]")
	println("This mode checks if all concurrency elements of the program have been executed at least once.")
	println("It has the following options:")
	println("  -resultTool [folder]   Path where the advocateResult folder created by the pipeline is located (required)")
	println("  -dir [folder]          Path to the program folder (required)")
	println("\n\n")
	println("4. Create statistics about a program")
	println("This creates some statistics about the program and the trace")
	println("Usage: ./analyzer stats [options]")
	// println("  -dir [folder] Path to the program folder (required)")
	println("  -trace [file]          Path to the folder containing the results_machine file (required)")
	println("  -prog [name]           Name of the program")
	println("  -test [name]           Name of the test")
	println("\n\n")
	println("5. Run the toolchain on tests")
	println("This runs the toolchain on a given main function")
	println("Usage: ./analyzer toolMain [options]")
	println("  -advocate [path]       Path to advocate")
	println("  -dir [path]            Path to the folder containing the program and tests")
	println("  -test [name]           Name of the test to run. If not set, all tests are run")
	println("  -prog [name]           Name of the program (used for statistics)")
	println("  -timeout [sec]         Timeout for the analysis")
	println("  -timeoutRelay [sec]    Timeout for the replay")
	println("  -ignoreAtomics         Set to ignore atomics in replay")
	println("  -recordTime            Set to record runtimes")
	println("  -notExec               Set to determine never executed operations")
	println("  -stats                 Set to create statistics")
	println("  -keepTrace             Do not delete the trace files after analysis finished")
	println("\n\n")
	println("6. Run the toolchain on a main function")
	println("This runs the toolchain on a given main function")
	println("Usage: ./analyzer toolMain [options]")
	println("  -advocate [path]       Path to advocate")
	println("  -dir [path]            Path to the file containing the main function")
	println("  -exec [name]           Name of the executable")
	println("  -prog [name]           Name of the program (used for statistics)")
	println("  -timeout [sec]         Timeout for the analysis")
	println("  -timeoutRelay [sec]    Timeout for the replay")
	println("  -ignoreAtomics         Set to ignore atomics in replay")
	println("  -recordTime            Set to record runtimes")
	println("  -notExec               Set to determine never executed operations")
	println("  -stats                 Set to create statistics")
	println("  -keepTrace             Do not delete the trace files after analysis finished")
	println("\n\n")
	println("7. Create runs for fuzzing")
	println("This creates and updates the information required for the fuzzing runs")
	println("Usage: ./analyzer fuzzing [options]")
	println("  -prog [name]           Name of the program")
	println("  -test [name]           Name of the test (only if used on tests)")
}
