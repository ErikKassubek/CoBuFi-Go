package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"analyzer/io"
	"analyzer/logging"
	"analyzer/rewriter"
	"analyzer/trace"
)

func main() {
	pathTrace := flag.String("t", "", "Path to the trace folder to analyze or rewrite")
	level := flag.Int("d", 1, "Debug Level, 0 = silent, 1 = errors, 2 = info, 3 = debug (default 1)")
	fifo := flag.Bool("f", false, "Assume a FIFO ordering for buffered channels (default false)")
	ignoreCriticalSection := flag.Bool("c", false, "Ignore happens before relations of critical sections (default false)")
	noRewrite := flag.Bool("x", false, "Do not ask to create a reordered trace file after the analysis (default false)")
	noWarning := flag.Bool("w", false, "Do not print warnings (default false)")
	noPrint := flag.Bool("p", false, "Do not print the results to the terminal (default false). Automatically set -x to true")
	resultFolder := flag.String("r", "", "Path to where the result file should be saved.")
	ignoreAtomics := flag.Bool("a", false, "Ignore atomic operations (default false). Use to reduce memory overhead for large traces.")

	scenarios := flag.String("s", "", "Select which analysis scenario to run, e.g. -s srd for the option s, r and d. Options:\n"+
		"\ts: Send on closed channel\n"+
		"\tr: Receive on closed channel\n"+
		"\tw: Done before add on waitGroup\n"+
		"\tn: Close of closed channel\n"+
		"\tb: Concurrent receive on channel\n"+
		"\tl: Leaking routine\n"+
		"\tu: Select case without partner\n"+
		"\tc: Cyclic deadlock\n",
	// "\tm: Mixed deadlock\n"
	)

	startTime := time.Now()

	flag.Parse()

	printHeader()

	if *pathTrace == "" {
		fmt.Println("Please provide a path to the trace file. Set with -t [file]")
		return
	}

	if *noPrint {
		*noRewrite = true
	}

	folder := filepath.Dir(*pathTrace) + string(os.PathSeparator)
	if *resultFolder != "" {
		folder = *resultFolder
		if folder[len(folder)-1] != os.PathSeparator {
			folder += string(os.PathSeparator)
		}
	}

	analysisCases, err := parseAnalysisCases(*scenarios)
	if err != nil {
		panic(err)
	}

	outMachine := folder + "results_machine.log"
	outReadable := folder + "results_readable.log"
	newTrace := folder + "rewritten_trace/"

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
	fmt.Println("Analysis finished\n")

	numberOfResults := logging.PrintSummary(*noWarning, *noPrint)

	analysisFinishedTime := time.Now()
	writeTime(*pathTrace, "analysis", analysisFinishedTime.Sub(startTime).Seconds())

	if numberOfResults != 0 && !*noRewrite {
		fmt.Println("\n\n\n")

		for {
			fmt.Print("Enter the index of the result to use for the reordered trace file (0 to exit): ")

			resultIndex := -1

			_, err := fmt.Scanf("%d", &resultIndex)
			if err != nil {
				if numberOfResults == 1 {
					fmt.Println("Please enter 1 to create the reordered trace file or 0 to exit\n")
				} else {
					fmt.Print("Please enter a valid number between 1 and ", numberOfResults, " or 0 to exit\n")
				}
				continue
			}
			if resultIndex == 0 { // no rewrite
				break
			}

			if resultIndex < 0 || resultIndex > numberOfResults {
				if numberOfResults == 1 {
					fmt.Println("Please enter 1 to create the reordered trace file or 0 to exit\n")
				} else {
					fmt.Print("Please enter a valid number between 1 and ", numberOfResults, " or 0 to exit\n")
				}
				continue
			}

			rewriteStartTime := time.Now()

			print("\n")
			if err := rewriteTrace(outMachine, newTrace, resultIndex, numberOfRoutines); err != nil {
				println("Could not rewrite trace file: ", err.Error())
			}

			writeTime(*pathTrace, "rewrite", time.Now().Sub(rewriteStartTime).Seconds())

			break // exit for loop
		}
	}
}

func writeTime(pathTrace string, name string, time float64) error {
	path := pathTrace
	if path[len(path)-1] != os.PathSeparator {
		path += string(os.PathSeparator)
	}
	path += "time.log"

	// Datei lesen
	content, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	elems := strings.Split(string(content), "\n")
	if len(elems) > 2 {
		return fmt.Errorf("time.log contains more than 2 lines")
	}

	if elems[0] != "" {
		elems[0] += ","
		elems[1] += ","
	}

	elems[0] += name
	elems[1] += strconv.FormatFloat(time, 'f', 6, 64)

	// Datei schreiben
	err = os.WriteFile(path, []byte(strings.Join(elems, "\n")), 0644)
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
 *   error: An error if the trace file could not be created
 */
func rewriteTrace(outMachine string, newTrace string, resultIndex int,
	numberOfRoutines int) error {

	actual, bug, err := io.ReadAnalysisResults(outMachine, resultIndex)
	if err != nil {
		return err
	}

	if actual {
		// copy the file of the tracePath to the outputPath
		io.CopyFile(*&outMachine, newTrace)
		println("Trace created")
		return nil
	}

	err = rewriter.RewriteTrace(bug)
	if err != nil {
		return err
	}

	err = io.WriteTrace(newTrace, numberOfRoutines)
	if err != nil {
		return err
	}

	return nil
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
		"concurrentReceive":    false,
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
		analysisCases["concurrentReceive"] = true
		analysisCases["leak"] = true
		analysisCases["selectWithoutPartner"] = true
		analysisCases["cyclicDeadlock"] = true
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
			analysisCases["concurrentReceive"] = true
		case 'l':
			analysisCases["leak"] = true
		case 'u':
			analysisCases["selectWithoutPartner"] = true
		case 'c':
			analysisCases["cyclicDeadlock"] = true
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
