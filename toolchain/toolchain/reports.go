// Copyright (c) 2024 Erik Kassubek, Mario Occhinegro
//
// File: headerUnitTests.go
// Brief: Functions to add and remove the ADVOCATE header into file containing
//    unit tests
//
// Author: Erik Kassubek, Mario Occhinegro
// Created: 2024-09-18
// Last Changed 2024-09-18
//
// License: BSD-3-Clause

package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

/*
 * Generate the bug reports
 * Args:
 *    folderName string: path to folder containing the results
 *    advocateRoot string: path to ADVOCATE
 */
func generateBugReports(folderName string, advocateRoot string) {
	files, err := getFiles(folderName, "results_machine.log")
	if err != nil {
		fmt.Println(err)
	}

	analyzerPath := filepath.Join(advocateRoot, "analyzer", "analyzer")
	for _, file := range files {
		folder := filepath.Dir(file)
		advocateTraceFolder := filepath.Join(folder, "advocateTrace")
		cmd := exec.Command("wc", "-l", file)
		out, err := cmd.Output()
		if err != nil {
			fmt.Println(err)
		}

		lineCount, err := strconv.Atoi(strings.Fields(string(out))[0])
		if err != nil {
			fmt.Println(err)
		}

		for i := 1; i <= lineCount; i++ {
			cmd := exec.Command(analyzerPath, "explain", "-t", advocateTraceFolder, "-i", strconv.Itoa(i))
			err := cmd.Run()
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}

/*
 * Get all files in folder path with name fileName
 * Args:
 *    folderPath (string): path to the folder to search in
 *    fileName (string): name of the files to search for
 * Returns:
 *    []string: list of the paths of the files
 *    error
 */
func getFiles(folderPath string, fileName string) ([]string, error) {
	var files []string
	err := filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if filepath.Base(path) == fileName {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return files, nil
}
