// Copyrigth (c) 2024 Erik Kassubek
//
// File: stats.go
// Brief: Create statistics about programs and traces
//
// Author: Erik Kassubek <kassubek.erik@gmail.com>
// Created: 2023-07-13
// LastChange: 2024-09-01
//
// License: BSD-3-Clause

package stats

import "fmt"

/*
 * Create files with the required stats
 * Args:
 *     pathToProgram (string): path to the program
 *     pathToTrace (string): path to the traces
 */
func CreateStats(pathToProgram *string, pathToResults *string) error {
	if pathToProgram == nil && pathToResults == nil {
		panic("Please provide at least one of the following flags: -t [file] or -P [file]")
	}

	statsProg, err := statsProgram(*pathToProgram)
	if err != nil {
		return err
	}

	statsTraces, err := statsTraces(*pathToResults)
	if err != nil {
		return err
	}

	statsAnalyzer, err := statsAnalyzer(*pathToResults)
	if err != nil {
		return err
	}

	fmt.Println(statsProg)
	fmt.Println(statsTraces)
	fmt.Println(statsAnalyzer["detected"])
	fmt.Println(statsAnalyzer["replayWritten"])
	fmt.Println(statsAnalyzer["replaySuccessful"])

	return nil

}

// 	pathToStats := ""
// 	if pathToTrace != nil {
// 		pathToStats = filepath.Dir(*pathToTrace)
// 	} else {
// 		pathToStats = *pathToProgram
// 	}

// 	pathToCSV := pathToStats + "/stats.csv"
// 	err := createFile(pathToCSV)
// 	if err != nil {
// 		panic(err)
// 	}

// 	if pathToProgram != nil {
// 		parseProgramToCSV(*pathToProgram, pathToCSV)
// 	}

// 	if pathToTrace != nil {
// 		parseTraceToCSV(*pathToTrace, pathToCSV)
// 	}
// }

// func createFile(path string) error {
// 	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
// 		return err
// 	}

// 	file, err := os.Create(path)
// 	if err != nil {
// 		return err
// 	}

// 	defer file.Close()

// 	return nil
// }

// func contains(arr []string, str string) bool {
// 	for _, a := range arr {
// 		if a == str {
// 			return true
// 		}
// 	}
// 	return false
// }
