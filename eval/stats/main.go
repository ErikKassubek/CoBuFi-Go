// Copyright (c) 2024 Erik Kassubek
//
// File: main.go
// Brief: Create combined statistics of all progs
//
// Author: Erik Kassubek
// Created: 2024-09-21
//
// License: BSD-3-Clause

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func main() {
	root := "../data/"
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		fileType, progName := getFileInfo(info.Name())

		switch fileType {
		case "statsAll", "statsAnalysis", "statsTrace":
			stats(path, progName, fileType)
		case "times":
			times(path, progName)
		}

		return nil
	})

	if err != nil {
		fmt.Printf("Error walking the path %q: %v\n", root, err)
	}
}

func stats(path string, progName string, progType string) {
	content, err := os.ReadFile(path)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return
	}

	data := string(content)

	lines := strings.Split(data, "\n")

	numberElems := 0
	switch progType {
	case "statsAll":
		numberElems = 85
	case "statsAnalysis":
		numberElems = 5
	case "statsTrace":
		numberElems = 9
	}

	dataSum := make([]int, numberElems)

	for i, line := range lines {
		data := strings.Split(line, ",")

		if i == 0 || len(data) == 0 {
			continue
		}

		// ignore the test name
		data = data[1:]

		for i, s := range data {
			num, err := strconv.Atoi(s)
			if err != nil {
				fmt.Printf("Error converting %s to int: %v\n", s, err)
				return
			}
			dataSum[i] += num
		}
	}

	writeToFile(dataSum, progName, progType, len(lines))
}

func times(path string, progName string) {
	content, err := os.ReadFile(path)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return
	}

	data := string(content)

	lines := strings.Split(data, "\n")

	numberElems := 9

	total := make([]float64, numberElems)

	for i, line := range lines {
		data := strings.Split(line, ",")

		if i == 0 || len(data) == 0 {
			continue
		}

		// ignore the test name
		data = data[1:]

		lastTimeReplay := 0.
		for i, s := range data {
			num, err := strconv.ParseFloat(s, 64)
			if err != nil {
				fmt.Printf("Error converting %s to int: %v\n", s, err)
				return
			}
			if i == len(data)-2 {
				lastTimeReplay = num
				continue
			} else if i == len(data)-1 {
				dur := num * lastTimeReplay
				total[len(data)-2] += dur
				total[i] += num
			} else {
				total[i] += num
			}
		}
	}

	writeToFile[float64](total, progName, "times", 0)
}

func writeToFile[T int | float64](data []T, progName string, progType string, numberOfTests int) {
	filePath := "../data/" + progType + "_total.csv"

	newFile := false
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		newFile = true
	}

	// Open the file for appending, create it if it doesn't exist
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return
	}
	defer file.Close() // Ensure the file is closed after we're done

	if newFile {
		switch progType {
		case "statsAll":
			_, err = file.WriteString("ProgramName,NumberOfTests,NumberOfEvents,NumberOfGoroutines,NumberOfNotEmptyGoroutines,NumberOfSpawnEvents,NumberOfRoutineEndEvents,NumberOfAtomics,NumberOfAtomicEvents,NumberOfChannels,NumberOfBufferedChannels,NumberOfUnbufferedChannels,NumberOfChannelEvents,NumberOfBufferedChannelEvents,NumberOfUnbufferedChannelEvents,NumberOfSelectEvents,NumberOfSelectCases,NumberOfSelectNonDefaultEvents,NumberOfSelectDefaultEvents,NumberOfMutex,NumberOfMutexEvents,NumberOfWaitgroup,NumberOfWaitgroupEvent,NumberOfCondVariables,NumberOfCondVariablesEvents,NumberOfOnce,NumberOfOnceOperations,NumberOfDetectedA01,NumberOfDetectedA02,NumberOfDetectedA03,NumberOfDetectedA04,NumberOfDetectedA05,NumberOfDetectedP01,NumberOfDetectedP02,NumberOfDetectedP03,NumberOfDetectedP04,NumberOfDetectedL00,NumberOfDetectedL01,NumberOfDetectedL02,NumberOfDetectedL03,NumberOfDetectedL04,NumberOfDetectedL05,NumberOfDetectedL06,NumberOfDetectedL07,NumberOfDetectedL08,NumberOfDetectedL09,NumberOfDetectedL10,NumberOfReplayWrittenA01,NumberOfReplayWrittenA02,NumberOfReplayWrittenA03,NumberOfReplayWrittenA04,NumberOfReplayWrittenA05,NumberOfReplayWrittenP01,NumberOfReplayWrittenP02,NumberOfReplayWrittenP03,NumberOfReplayWrittenP04,NumberOfReplayWrittenL00,NumberOfReplayWrittenL01,NumberOfReplayWrittenL02,NumberOfReplayWrittenL03,NumberOfReplayWrittenL04,NumberOfReplayWrittenL05,NumberOfReplayWrittenL06,NumberOfReplayWrittenL07,NumberOfReplayWrittenL08,NumberOfReplayWrittenL09,NumberOfReplayWrittenL10,NumberOfReplaySuccessfulA01,NumberOfReplaySuccessfulA02,NumberOfReplaySuccessfulA03,NumberOfReplaySuccessfulA04,NumberOfReplaySuccessfulA05,NumberOfReplaySuccessfulP01,NumberOfReplaySuccessfulP02,NumberOfReplaySuccessfulP03,NumberOfReplaySuccessfulP04,NumberOfReplaySuccessfulL00,NumberOfReplaySuccessfulL01,NumberOfReplaySuccessfulL02,NumberOfReplaySuccessfulL03,NumberOfReplaySuccessfulL04,NumberOfReplaySuccessfulL05,NumberOfReplaySuccessfulL06,NumberOfReplaySuccessfulL07,NumberOfReplaySuccessfulL08,NumberOfReplaySuccessfulL09,NumberOfReplaySuccessfulL10\n")
		case "statsAnalysis":
			_, err = file.WriteString("ProgramName,NumberOfLeaks,NumberOfLeaksWithRewrite,NumberOfLeaksResolvedViaReplay,NumberOfPanics,NumberOfPanicsVerifiedViaReplay\n")
		case "statsTrace":
			_, err = file.WriteString("ProgramName,NumberOfEvents,NumberOfGoroutines,NumberOfAtomicEvents,NumberOfChannelEvents,NumberOfSelectEvents,NumberOfMutexEvents,NumberOfWaitgroupEvents,NumberOfCondVariablesEvents,NumberOfOnceOperations\n")
		case "times":
			_, err = file.WriteString("TestName,ExecTime,ExecTimeWithTracing,AnalyzerTime,AnalysisTime,HBAnalysisTime,TimeToIdentifyLeaksPlusFindingPoentialPartners,TimeToIdentifyPanicBugs,ReplayTime,NumberReplay\n")
		}
		if err != nil {
			fmt.Printf("Error writing to file: %v\n", err)
			return
		}
	}

	line := progName
	if progType == "statsAll" {

	}

	for _, num := range data {
		line += ","
		switch v := any(num).(type) {
		case int:
			line += strconv.Itoa(v)
		case float64:
			line += strconv.FormatFloat(v, 'f', -1, 64)
		}
	}
	line += "\n"

	_, err = file.WriteString(line)
	if err != nil {
		fmt.Printf("Error writing to file: %v\n", err)
		return
	}

}

// Get the file type and prog name from the file name
func getFileInfo(fileName string) (string, string) {
	name := strings.Split(fileName, ".")[0]

	info := strings.Split(name, "_")
	fileType := info[0]
	progName := strings.Join(info[1:], "_")

	return fileType, progName
}
