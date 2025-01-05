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
	pathTrace := flag.String("t", "", "Path to the trace folder to analyze or rewrite")
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
	lastIndexFuzzing := flag.Int("i", 0, "Index of last fuzzing run for this program")

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
	case "stats":
		modes.ModeStats(*pathTrace, *progName, *testName)
	case "explain":
		modes.ModeExplain(*pathTrace, !*rewriteAll)
	case "check":
		modes.ModeCheck(resultFolderTool, programPath)
	case "run":
		modes.ModeAnalyzer(pathTrace, noPrint, noRewrite, scenarios, outReadable,
			outMachine, ignoreAtomics, fifo, ignoreCriticalSection,
			noWarning, rewriteAll, folderTrace, newTrace, timeout, ignoreRewrite)
	case "fuzzing":
		modes.ModeFuzzing(*pathTrace, *progName, *testName, *lastIndexFuzzing)
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
	println("4. Create statistics about a program")
	println("5. Create new runs for fuzzing\n\n")
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
	println("\n\n")
	println("5. Create runs for fuzzing")
	println("This creates and updates the information required for the fuzzing runs")
	println("Usage: ./analyzer fuzzing [options]")
	println("  -t [file]   Path to the folder containing the results_machine file (required)")
	println("  -N [name]   Name of the program")
	println("  -M [name]   Name of the test (only if used on tests)")
	println("  -i [index]  Index of the last fuzzing")

}
