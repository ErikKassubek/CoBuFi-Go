// Copyright (c) 2024 Erik Kassubek
//
// File: statsProgram.go
// Brief: Collect statistics about the program
//
// Author: Erik Kassubek
// Created: 2024-09-20
//
// License: BSD-3-Clause

package stats

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

/*
 * Parse a program to measure the number of files, and number of lines
 * Args:
 *     programPath (string): path to the folder containing the program
 * Returns:
 *     map[string]int: map with numberFiles, numberLines, numberNonEmptyLines
 *     error
 */
func statsProgram(programPath string) (map[string]int, error) {
	res := make(map[string]int)
	res["numberFiles"] = 0
	res["numberLines"] = 0
	res["numberNonEmptyLines"] = 0

	err := filepath.Walk(programPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if filepath.Ext(path) == ".go" {
			resFile, err := parseProgramFile(path)
			if err != nil {
				return err
			}

			res["numberFiles"]++
			res["numberLines"] += resFile["numberLines"]
			res["numberNonEmptyLines"] += resFile["numberNonEmptyLines"]
		}

		return nil
	})
	return res, err
}

/*
 * Parse one program file to measure the number of lines
 * Args:
 *     programPath (string): path to the file
 * Returns:
 *     map[string]int: map with numberLines, numberNonEmptyLines
 *     error
 */
func parseProgramFile(filePath string) (map[string]int, error) {
	res := make(map[string]int)
	res["numberLines"] = 0
	res["numberNonEmptyLines"] = 0

	// open the file
	file, err := os.Open(filePath)
	if err != nil {
		return res, err
	}
	defer file.Close()

	// read the file
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text := strings.TrimSpace(scanner.Text())

		res["numberLines"]++
		if text != "" && text != "\n" && !strings.HasPrefix(text, "//") {
			res["numberNonEmptyLines"]++
		}
	}

	if err := scanner.Err(); err != nil {
		return res, err
	}

	return res, nil
}
