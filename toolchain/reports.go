// Copyright (c) 2024 Erik Kassubek, Mario Occhinegro
//
// File: headerUnitTests.go
// Brief: Functions to generate bug reports
//
// Author: Erik Kassubek, Mario Occhinegro
// Created: 2024-09-18
// Last Changed 2024-09-21
//
// License: BSD-3-Clause

package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
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

	for _, file := range files {
		folder := filepath.Dir(file)
		// advocateTraceFolder := filepath.Join(folder, "advocateTrace")
		analyzerPath := filepath.Join(advocateRoot, "analyzer", "analyzer")
		cmd := exec.Command(analyzerPath, "explain", "-t", folder)
		err := cmd.Run()
		if err != nil {
			fmt.Println(err)
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

/*
 * Function to move results files from the package directory to the destination directory
 * Args:
 *    packagePath (string): path to the package directory
 *    destination (string): path to the destination directory
 */
func moveResults(packagePath, destination string) {
	filesToMove := []string{
		"advocateTrace",
		"results_machine.log",
		"results_readable.log",
		"output.log",
	}

	for _, file := range filesToMove {
		src := filepath.Join(packagePath, file)
		dest := filepath.Join(destination, file)
		if err := os.Rename(src, dest); err != nil && file != "output.log" {
			log.Printf("Failed to move %s to %s: %v", src, dest, err)
		}
	}

	// Move any rewritten_trace directories
	rewrittenTraces, _ := filepath.Glob(filepath.Join(packagePath, "rewritten_trace*"))
	for _, trace := range rewrittenTraces {
		dest := filepath.Join(destination, filepath.Base(trace))
		if err := os.Rename(trace, dest); err != nil {
			log.Printf("Failed to move %s to %s: %v", trace, dest, err)
		}
	}
}
