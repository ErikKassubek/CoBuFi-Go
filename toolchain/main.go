// Copyright (c) 2024 Erik Kassubek
//
// File: headerUnitTests.go
// Brief: Functions to add and remove the ADVOCATE header into file containing
//    unit tests
//
// Author: Erik Kassubek
// Created: 2024-09-18
// Last Changed 2024-09-19
//
// License: BSD-3-Clause

package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

var (
	pathToAdvocate string
	pathToFile     string
	progName       string
	executableName string
	help           bool
	measureTime    bool
	notExecuted    bool
	stats          bool
	timeoutAna     int
)

func init() {
	flag.BoolVar(&help, "h", false, "Help")
	flag.StringVar(&pathToAdvocate, "a", "", "path to the ADVOCATE folder")
	flag.StringVar(&pathToFile, "f", "", "main: path to the main program file, tests: path to the folder with the program and the tests")
	flag.StringVar(&progName, "N", "", "name of the analyzed program. Only required if -s or -t is set")
	flag.StringVar(&executableName, "E", "", "name of the executable. Only required for main")
	flag.BoolVar(&measureTime, "t", false, "set to measure the duration of the"+
		"different steps. This will also run the program/tests once without any recording"+
		"to get a base value")
	flag.BoolVar(&notExecuted, "m", false, "check for not executed operations")
	flag.BoolVar(&stats, "s", false, "create statistic files")
	flag.IntVar(&timeoutAna, "T", -1, "Set a timeout in seconds for each run of the analyzer")
}

func main() {
	flag.Parse()

	var mode string
	if len(os.Args) > 2 {
		mode = os.Args[1]
		flag.CommandLine.Parse(os.Args[2:])
	}

	if help {
		switch mode {
		case "main":
			printHelpMain()
		case "test", "tests":
			printHelpUnit()
		default:
			printHelp()
		}
		return
	}

	// replace ~ in path with home
	home, _ := os.UserHomeDir()
	pathToAdvocate = strings.Replace(pathToAdvocate, "~", home, -1)
	pathToFile = strings.Replace(pathToFile, "~", home, -1)

	var err error
	switch mode {
	case "main":
		if pathToAdvocate == "" {
			fmt.Println("Path to advocate required for mode main")
			printHelpMain()
			return
		}
		if pathToFile == "" {
			fmt.Println("Path to file required")
			printHelpMain()
			return
		}
		if executableName == "" {
			fmt.Println("Name of the executable required")
			printHelpMain()
		}
		if (stats || measureTime) && progName == "" {
			fmt.Println("If -s or -t is set, -N [name] must be set as well")
			printHelpMain()
			return
		}
		err = runWorkflowMain(pathToAdvocate, pathToFile, executableName, timeoutAna)
	case "test", "tests":
		if pathToAdvocate == "" {
			fmt.Println("Path to advocate required")
			printHelpUnit()
			return
		}
		if pathToFile == "" {
			fmt.Println("Path to test folder required for mode main")
			printHelpUnit()
			return
		}
		if (stats || measureTime) && progName == "" {
			fmt.Println("If -s or -t is set, -N [name] must be set as well")
			printHelpUnit()
			return
		}
		err = runWorkflowUnit(pathToAdvocate, pathToFile, progName, measureTime, notExecuted, stats, timeoutAna)
	case "explain":
		if pathToAdvocate == "" {
			fmt.Println("Path to advocate required")
			printHelpUnit()
			return
		}
		if pathToFile == "" {
			fmt.Println("Path to test folder required for mode main")
			printHelpUnit()
			return
		}
		generateBugReports(pathToFile, pathToAdvocate)
	default:
		fmt.Println("Choose one mode from 'main' or 'test'")
		printHelp()
	}

	if err != nil {
		fmt.Println(err)
	}
}

func printHelp() {
	fmt.Println("Usage: ./toolchain <mode> [options]")
	fmt.Println("Modes:")
	fmt.Println("  main:   Run the workflow for a program with a main function")
	fmt.Println("  test:   Run the workflow for unit tests")
	fmt.Println("Use ./toolchain <mode> -h for more help")
}

func printHelpMain() {
	fmt.Println("Usage: ./toolchain main [options]")
	fmt.Println("Required Flags:")
	fmt.Println("  -a [path]: path to the ADCOVATE folder")
	fmt.Println("  -f [path]: path to the file containing the main function")
	fmt.Println("  -E [name]: name of the program executable")
	fmt.Println("  -t       : measure the runtimes")
	fmt.Println("  -m       : check for never executed operations")
	fmt.Println("  -s       : create statistics about the analyzed program")
	fmt.Println("  -N [name]: give a name for the analyzed program. Only required if -s or -t is set")
}

func printHelpUnit() {
	fmt.Println("Usage: ./toolchain test [options]")
	fmt.Println("Required Flags:")
	fmt.Println("  -a [path]: path to the ADCOVATE folder")
	fmt.Println("  -f [path]: path to the folder containing the tests")
	fmt.Println("  -t       : measure the runtimes")
	fmt.Println("  -m       : check for never executed operations")
	fmt.Println("  -s       : create statistics about the analyzed program")
	fmt.Println("  -N [name]: give a name for the analyzed program. Only required if -s or -t is set")
}
