// Copyright (c) 2024 Erik Kassubek
//
// File: headerUnitTests.go
// Brief: Functions to add and remove the ADVOCATE header into file containing
//    unit tests
//
// Author: Erik Kassubek
// Created: 2024-09-18
// Last Changed 2024-09-18
//
// License: BSD-3-Clause

package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		printHelp()
		return
	}

	mode := os.Args[1]

	help := flag.Bool("h", false, "Help")
	pathToAdvocate := flag.String("a", "", "path to the ADVOCATE folder")
	pathToFile := flag.String("f", "", "main: path to the main program file, tests: path to the folder with the program and the tests")
	flag.Parse()

	if *help {
		switch mode {
		case "main":
			printHelpMain()
		case "tests":
			printHelpUnit()
		default:
			printHelp()
		}
		return
	}

	switch mode {
	case "main":
		if *pathToAdvocate == "" {
			fmt.Println("Path to advocate required for mode main")
			printHelpMain()
			return
		}
		if *pathToFile == "" {
			fmt.Println("Path to file required for mode main")
			printHelpMain()
		}
		runWorkflowMain(*pathToAdvocate, *pathToFile)
	case "tests":
		runWorkflowUnit(*pathToAdvocate, *pathToFile)
	default:
		printHelp()
	}
}

func printHelp() {
	fmt.Println("Usage: ./toolchain <mode> [options]")
	fmt.Println("Modes:")
	fmt.Println("  main:   Run the workflow for a program with a main function")
	fmt.Println("  tests:  Run the workflow for unit tests")
	fmt.Println("Use ./toolchain <mode> -h for mor help")
}

func printHelpMain() {
	fmt.Println("Usage: ./toolchain main [options]")
	fmt.Println("Required Flags:")
	fmt.Println("  -a [path]: path to the ADCOVATE folder")
	fmt.Println("  -f [path]: path to the file containing the main function")
}

func printHelpUnit() {
	fmt.Println("Usage: ./toolchain tests [options]")
	fmt.Println("Required Flags:")
	fmt.Println("  -a [path]: path to the ADCOVATE folder")
	fmt.Println("  -f [path]: path to the folder containing the ")
}
