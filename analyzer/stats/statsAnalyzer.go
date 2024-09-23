// Copyright (c) 2024 Erik Kassubek
//
// File: statsAnalyzer.go
// Brief: Collect stats about the analysis and the replay
//
// Author: Erik Kassubek
// Created: 2024-09-20
//
// License: BSD-3-Clause

package stats

import (
	"analyzer/explanation"
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

/*
 * Parse the analyzer and replay output to collect the corresponding information
 * Args:
 *     pathToResults (string): path to the advocateResult folder
 * Returns:
 *     map[string]int: map with information
 *     error
 */
func statsAnalyzer(pathToResults string) (map[string]map[string]int, error) {
	detected := map[string]int{
		"A01": 0, "A02": 0, "A03": 0, "A04": 0, "A05": 0, "P01": 0, "P02": 0,
		"P03": 0, "P04": 0, "L01": 0, "L02": 0, "L03": 0, "L04": 0, "L05": 0,
		"L06": 0, "L07": 0, "L08": 0, "L09": 0, "L10": 0}
	replayWriten := map[string]int{
		"A01": 0, "A02": 0, "A03": 0, "A04": 0, "A05": 0, "P01": 0, "P02": 0,
		"P03": 0, "P04": 0, "L01": 0, "L02": 0, "L03": 0, "L04": 0, "L05": 0,
		"L06": 0, "L07": 0, "L08": 0, "L09": 0, "L10": 0}
	replaySuccessful := map[string]int{
		"A01": 0, "A02": 0, "A03": 0, "A04": 0, "A05": 0, "P01": 0, "P02": 0,
		"P03": 0, "P04": 0, "L01": 0, "L02": 0, "L03": 0, "L04": 0, "L05": 0,
		"L06": 0, "L07": 0, "L08": 0, "L09": 0, "L10": 0}

	res := map[string]map[string]int{
		"detected":         detected,
		"replayWritten":    replayWriten,
		"replaySuccessful": replaySuccessful,
	}

	alreadyProcessed := make(map[string][]processedBug)

	err := filepath.Walk(pathToResults, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && info.Name() == "bug.md" {
			err := processBugFile(path, &res, &alreadyProcessed)
			if err != nil {
				fmt.Println(err)
			}
		}

		return nil
	})

	return res, err
}

type processedBug struct {
	paths         string
	detectedBug   string
	replayWritten bool
	replaySuc     bool
}

/*
 * Parse a bug file to get the information
 * Args:
 *     filePath (string): path to the bug file
 *     info (*map[string]map[string]int): map to store the info in
 * Returns:
 *     error
 */
func processBugFile(filePath string, info *map[string]map[string]int, alreadyProcessed *map[string][]processedBug) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	bugType := ""

	// paths of bug elems
	elemPaths := make([]string, 0)
	bug := processedBug{}

	// read the file
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// get detected bug
		if strings.HasPrefix(line, "# ") {
			textSplit := strings.Split(line, ": ")
			if len(textSplit) != 2 {
				continue
			}

			line = textSplit[1]

			bugType = explanation.GetCodeFromDescription(line)
			if bugType == "" {
				return fmt.Errorf("unknown error type %s", line)
			}
			bug.detectedBug = bugType
		}

		// get bug id based on the path of the elements
		if strings.HasPrefix(line, "->") {
			path := strings.TrimPrefix(line, "-> ")
			elemPaths = append(elemPaths, path)
		}

		// replay written
		if line == "The rewritten trace can be found in the `rewritten_trace` folder." {
			bug.replayWritten = true
		}

		// replay result
		if strings.HasPrefix(line, "It exited with the following code: ") {
			code := strings.TrimPrefix(line, "It exited with the following code: ")

			num, err := strconv.Atoi(code)
			if err != nil {
				return err
			}

			if num < 10 || num >= 20 {
				bug.replaySuc = true
			}
		}
	}

	// do not count the same bug twice
	sort.Strings(elemPaths)
	bug.paths = strings.Join(elemPaths, ">")

	if val, ok := (*alreadyProcessed)[bug.paths]; ok { // path already counted -> only count if bug type is not equal
		for _, b := range val {
			if bug.detectedBug == b.detectedBug {
				continue
			}

			(*info)["detected"][bugType]++

			if bug.replayWritten {
				(*info)["replayWritten"][bugType]++
			}

			if bug.replaySuc {
				(*info)["replaySuccessful"][bugType]++
			}

			(*alreadyProcessed)[bug.paths] = append((*alreadyProcessed)[bug.paths], bug)
		}
	} else { // path not in already processed
		(*info)["detected"][bugType]++

		if bug.replayWritten {
			(*info)["replayWritten"][bugType]++
		}

		if bug.replaySuc {
			(*info)["replaySuccessful"][bugType]++
		}

		(*alreadyProcessed)[bug.paths] = []processedBug{bug}
	}

	return nil
}
