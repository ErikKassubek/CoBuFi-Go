// Copyrigth (c) 2024 Erik Kassubek
//
// File: stats.go
// Brief: Create statistics about programs and traces
//
// Author: Erik Kassubek
// Created: 2023-07-13
//
// License: BSD-3-Clause

package stats

import (
	"fmt"
	"os"
	"path/filepath"
)

/*
 * Create files with the required stats
 * Args:
 *     pathToProgram (string): path to the program
 *     pathToTrace (string): path to the traces
 *     progName (string): name of the analyzed program
 */
func CreateStats(pathToProgram, pathToResults, progName string) error {
	statsProg, err := statsProgram(pathToProgram)
	if err != nil {
		return err
	}

	statsTraces, err := statsTraces(pathToResults)
	if err != nil {
		return err
	}

	statsAnalyzer, err := statsAnalyzer(pathToResults)
	if err != nil {
		return err
	}

	err = writeStatsToFile(pathToResults, progName, statsProg, statsTraces, statsAnalyzer)
	if err != nil {
		return err
	}

	return nil

}

/*
* Write the collected statistics to a file
* The file consists of [............] lines.
* The lines contain the informations always separated by commas
* The first line contains the stats about the program. This consists of the
*   number of files, the number of lines and the number of non empty lines.
* Lines 2 to 4 contain values about the trace.
*   In line 2 those values are
*     number of traces (for a main this should be 1, for tests this should be equal to the number of run tests)
*     total number of routines
*     number of non empty routines
*   Line 3 contains the number of relevant objects, this includes
*     number of atomic variables
*     number of channels
*     number of buffered channels
*     number of unbuffered channels
*     number of selects
*     number of select cases
*     number of mutexes
*     number of wait groups
*     number of cond variables
*     number of once
*   Line 4 contain the number of operations
*     number of total operations in the trace
*     number of spawns
*     number of atomic operations
*     number channel operations
*     number buffered channel operations
*     number unbuffered channel operations
*     number of select operations where a non default case was selected
*     number of select operations where the default case was selected
*     number of mutex operations
*     number of wait group operations
*     number of cond var operations
*     number of once operations
* Line 5 to 7 contain information about the analysis and replay.
*   For each line it contains the number of bugs actual (A), potential (P)
*   and leak (L) and after those three values the more precise for each number
*   Line 5 contains the information about the number of detected bugs
*   Line 6 contains the number of successful rewrites
*   Line 7 contains the number of successful replays
* Args:
*     path (string): path to where the stats file should be created
*     progName (string): name of the program
*     statsProg (map[string]int): statistics about the program
*     statsTraces (map[string]int): statistics about the trace
*     statsAnalyzer (map[string]map[string]int): statistics about the analysis and replay
* Returns:
*     error
 */
func writeStatsToFile(path string, progName string, statsProg, statsTraces map[string]int,
	statsAnalyzer map[string]map[string]int) error {

	f, err := os.Create(filepath.Join(path, "stats_"+progName+".log"))
	if err != nil {
		return err
	}
	defer f.Close()

	statsProgStr := fmt.Sprintf("%d,%d,%d\n", statsProg["numberFiles"],
		statsProg["numberLines"], statsProg["numberNonEmptyLines"])

	statsTraceStr1 := fmt.Sprintf("%d,%d,%d\n", statsTraces["numberTraces"],
		statsTraces["numberRoutines"], statsTraces["numberNonEmptyRoutines"])

	totalNumberOps := statsTraces["numberOfSpawns"] + statsTraces["numberAtomics"] +
		statsTraces["numberChannels"] + statsTraces["numberBufferedChannels"] +
		statsTraces["numberUnbufferedChannels"] + statsTraces["numberSelects"] +
		statsTraces["numberSelectCases"] + statsTraces["numberMutexes"] +
		statsTraces["numberWaitGroups"] + statsTraces["numberCondVars"] +
		statsTraces["numberOnce"]

	statsTraceStr2 := fmt.Sprintf("%d,%d,%d,%d,%d,%d,%d,%d,%d,%d\n",
		statsTraces["numberAtomics"], statsTraces["numberChannels"],
		statsTraces["numberBufferedChannels"], statsTraces["numberUnbufferedChannels"],
		statsTraces["numberSelects"], statsTraces["numberSelectCases"],
		statsTraces["numberMutexes"], statsTraces["numberWaitGroups"],
		statsTraces["numberCondVars"], statsTraces["numberOnce"])

	statsTraceStr3 := fmt.Sprintf("%d,%d,%d,%d,%d,%d,%d,%d,%d,%d,%d,%d\n",
		totalNumberOps, statsTraces["numberOfSpawns"],
		statsTraces["numberAtomicOperations"], statsTraces["numberChannelOperations"],
		statsTraces["numberBufferedOps"], statsTraces["numberUnbufferedOps"],
		statsTraces["numberSelectChanOps"], statsTraces["numberSelectDefaultOps"],
		statsTraces["numberMutexOperations"], statsTraces["numberWaitGroupOperations"],
		statsTraces["numberCondVarOperations"], statsTraces["numberOnceOperations"])

	totalAmount := map[string]map[string]int{
		"detected": {
			"actual":    0,
			"potential": 0,
			"leak":      0,
		},
		"replayWritten": {
			"actual":    0,
			"potential": 0,
			"leak":      0,
		},
		"replaySuccessful": {
			"actual":    0,
			"potential": 0,
			"leak":      0,
		},
	}

	dataString := map[string]string{
		"detected":         "",
		"replayWritten":    "",
		"replaySuccessful": "",
	}

	statType := []string{"detected", "replayWritten", "replaySuccessful"}
	bugCodes := []string{"A01", "A02", "A03", "A04", "A05",
		"P01", "P02", "P03", "P04",
		"L00", "L01", "L02", "L03", "L04", "L05", "L06", "L07", "L08", "L09", "L10"}

	for _, c := range statType {
		m := statsAnalyzer[c]
		for _, key := range bugCodes {
			amount := m[key]
			dataString[c] += fmt.Sprintf(",%d", amount)
			switch string([]rune(key)[0]) {
			case "A":
				totalAmount[c]["actual"] += amount
			case "P":
				totalAmount[c]["potential"] += amount
			case "L":
				totalAmount[c]["leak"] += amount
			default:
				return fmt.Errorf("Unknown bug code %s", key)
			}
		}
	}

	statsAnalyzerStr1 := fmt.Sprintf("%d,%d,%d%s\n",
		totalAmount["detected"]["actual"], totalAmount["detected"]["potential"],
		totalAmount["detected"]["leak"], dataString["detected"])
	statsAnalyzerStr2 := fmt.Sprintf("%d,%d,%d%s\n",
		totalAmount["replayWritten"]["actual"], totalAmount["replayWritten"]["potential"],
		totalAmount["replayWritten"]["leak"], dataString["replayWritten"])
	statsAnalyzerStr3 := fmt.Sprintf("%d,%d,%d%s\n",
		totalAmount["replaySuccessful"]["actual"], totalAmount["replaySuccessful"]["potential"],
		totalAmount["replaySuccessful"]["leak"], dataString["replaySuccessful"])

	f.WriteString(statsProgStr)
	f.WriteString(statsTraceStr1)
	f.WriteString(statsTraceStr2)
	f.WriteString(statsTraceStr3)
	f.WriteString(statsAnalyzerStr1)
	f.WriteString(statsAnalyzerStr2)
	f.WriteString(statsAnalyzerStr3)

	return nil
}
