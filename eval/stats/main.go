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
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type progData struct {
	name string

	numberTests         string
	numberFiles         string
	numberLines         string
	numberNonEmptyLines string

	numberRoutines         string
	numberNonEmptyRoutines string

	numberAtomics            string
	numberChannels           string
	numberBuffereChannels    string
	numberUnbufferedChannels string
	numberSelects            string
	numberSelectCases        string
	numberMutexes            string
	numberWaitGroups         string
	numberCondVariables      string
	numberOnce               string

	numberOperations           string
	numberSpawnOps             string
	numberRoutineTermOps       string
	numberAtomicOps            string
	numberChannelOps           string
	numberBuffereChannelOps    string
	numberUnbufferedChannelOps string
	numberSelectCaseOps        string
	numberSelectDefaultOps     string
	numberMutexOps             string
	numberWaitOps              string
	numberCondVarOps           string
	numberOnceOps              string

	numberDetected  map[string]string
	numberRewritten map[string]string
	numberReplayed  map[string]string

	timeRun      float64
	timeRecord   float64
	timeAnalysis float64
	timeReplay   float64
	numberReplay int
}

var progs = make(map[string]progData)

func main() {
	statsPath := flag.String("f", "", "Path to the stat and time files")
	flag.Parse()

	if *statsPath == "" {
		fmt.Println("Please set the path to the folder containing the stats and time files")
		return
	}

	err := filepath.Walk(*statsPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			if strings.Contains(info.Name(), "stats_") {
				err := readStats(path)
				if err != nil {
					fmt.Println("Failed to read stats for ", info.Name(), err)
				}
			} else if strings.Contains(info.Name(), "times_") {
				err := readTime(path)
				if err != nil {
					fmt.Println("Failed to read times for ", info.Name(), err)
				}
			}
		}

		return nil
	})

	if err != nil {
		log.Fatalf("Error walking the path %v: %v\n", *statsPath, err)
	}

	fileTime := time.Now().Format("2006-01-02_15:04:05")
	fileNameCvs := "tables_" + fileTime + ".csv"
	fileNameLatex := "tables_" + fileTime + ".tex"

	createCsv(fileNameCvs)
	createLatex(fileNameLatex)

}

func readStats(path string) error {
	name := getProgNameFromFile(path)

	data, err := readFile(path)
	if err != nil {
		return err
	}

	dataSplit := strings.Split(data, "\n")

	if len(dataSplit) < 7 {
		return fmt.Errorf("The stats file expected at least 7 lines but got %d", len(dataSplit))
	}

	infoProg := strings.Split(dataSplit[0], ",")
	infoTrace := strings.Split(dataSplit[1], ",")
	infoObjects := strings.Split(dataSplit[2], ",")
	infoOperations := strings.Split(dataSplit[3], ",")
	infoDetected := strings.Split(dataSplit[4], ",")
	infoRewritten := strings.Split(dataSplit[5], ",")
	infoReplay := strings.Split(dataSplit[6], ",")

	if val, ok := progs[name]; ok {
		val.numberFiles = infoProg[0]
		val.numberLines = infoProg[1]
		val.numberNonEmptyLines = infoProg[2]

		val.numberTests = infoTrace[0]
		val.numberRoutines = infoTrace[1]
		val.numberNonEmptyRoutines = infoTrace[2]

		val.numberAtomics = infoObjects[0]
		val.numberChannels = infoObjects[1]
		val.numberBuffereChannels = infoObjects[2]
		val.numberUnbufferedChannels = infoObjects[3]
		val.numberSelects = infoObjects[4]
		val.numberSelectCases = infoObjects[5]
		val.numberMutexes = infoObjects[6]
		val.numberWaitGroups = infoObjects[7]
		val.numberCondVariables = infoObjects[8]
		val.numberOnce = infoObjects[9]

		val.numberOperations = infoOperations[0]
		val.numberSpawnOps = infoOperations[1]
		val.numberRoutineTermOps = infoOperations[2]
		val.numberAtomicOps = infoOperations[3]
		val.numberChannelOps = infoOperations[4]
		val.numberBuffereChannelOps = infoOperations[5]
		val.numberUnbufferedChannelOps = infoOperations[6]
		val.numberSelectCaseOps = infoOperations[7]
		val.numberSelectDefaultOps = infoOperations[8]
		val.numberMutexOps = infoOperations[9]
		val.numberWaitOps = infoOperations[10]
		val.numberCondVarOps = infoOperations[11]
		val.numberOnceOps = infoOperations[12]

		val.numberDetected = mapCodes(infoDetected)
		val.numberRewritten = mapCodes(infoRewritten)
		val.numberReplayed = mapCodes(infoReplay)

		progs[name] = val
	} else {
		progs[name] = progData{
			name: name,

			numberFiles:         infoProg[0],
			numberLines:         infoProg[1],
			numberNonEmptyLines: infoProg[2],

			numberTests:            infoTrace[0],
			numberRoutines:         infoTrace[1],
			numberNonEmptyRoutines: infoTrace[2],

			numberAtomics:            infoObjects[0],
			numberChannels:           infoObjects[1],
			numberBuffereChannels:    infoObjects[2],
			numberUnbufferedChannels: infoObjects[3],
			numberSelects:            infoObjects[4],
			numberSelectCases:        infoObjects[5],
			numberMutexes:            infoObjects[6],
			numberWaitGroups:         infoObjects[7],
			numberCondVariables:      infoObjects[8],
			numberOnce:               infoObjects[9],

			numberOperations:           infoOperations[0],
			numberSpawnOps:             infoOperations[1],
			numberAtomicOps:            infoOperations[2],
			numberChannelOps:           infoOperations[3],
			numberBuffereChannelOps:    infoOperations[4],
			numberUnbufferedChannelOps: infoOperations[5],
			numberSelectCaseOps:        infoOperations[6],
			numberSelectDefaultOps:     infoOperations[7],
			numberMutexOps:             infoOperations[8],
			numberWaitOps:              infoOperations[9],
			numberCondVarOps:           infoOperations[10],
			numberOnceOps:              infoOperations[11],

			numberDetected:  mapCodes(infoDetected),
			numberRewritten: mapCodes(infoRewritten),
			numberReplayed:  mapCodes(infoReplay),
		}
	}

	return nil
}

func readTime(path string) error {
	name := getProgNameFromFile(path)

	data, err := readFile(path)
	if err != nil {
		return err
	}

	dataSplit := strings.Split(data, "\n")

	for _, test := range dataSplit {
		if test == "" {
			continue
		}
		lineSplit := strings.Split(test, "#")
		timeRun, _ := strconv.ParseFloat(lineSplit[0], 64)
		timeRecord, _ := strconv.ParseFloat(lineSplit[1], 64)
		timeAnalysis, _ := strconv.ParseFloat(lineSplit[2], 64)
		timeReplay, _ := strconv.ParseFloat(lineSplit[3], 64)
		numberReplay, _ := strconv.Atoi(lineSplit[4])

		if val, ok := progs[name]; ok {
			val.timeRun += timeRun
			val.timeRecord += timeRecord
			val.timeAnalysis += timeAnalysis
			val.timeReplay += timeReplay
			val.numberReplay += numberReplay
			progs[name] = val
		} else {
			progs[name] = progData{
				name:         name,
				timeRun:      timeRun,
				timeRecord:   timeRecord,
				timeAnalysis: timeAnalysis,
				timeReplay:   timeReplay,
				numberReplay: numberReplay,
			}
		}
	}

	return nil

}

func mapCodes(data []string) map[string]string {
	return map[string]string{
		"A":   data[0],
		"P":   data[1],
		"L":   data[2],
		"A01": data[3],
		"A02": data[4],
		"A03": data[5],
		"A04": data[6],
		"A05": data[7],
		"P01": data[8],
		"P02": data[9],
		"P03": data[10],
		"P04": data[11],
		"L00": data[12],
		"L01": data[13],
		"L02": data[14],
		"L03": data[15],
		"L04": data[16],
		"L05": data[17],
		"L06": data[18],
		"L07": data[19],
		"L08": data[20],
		"L09": data[21],
		"L10": data[22],
	}
}

func readFile(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

func getProgNameFromFile(path string) string {
	path = strings.Replace(path, "stats_", "", -1)
	path = strings.Replace(path, "times_", "", -1)
	path = strings.Replace(path, ".log", "", -1)
	return filepath.Base(path)
}

func gv(val string) string {
	if val == "" {
		return "0"
	}
	return val
}

func writeToFile(fileName, content string) {
	// Open the file in append mode, or create it if it doesn't exist
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()
	_, err = file.WriteString(content + "\n\n\n")
	if err != nil {
		fmt.Println("Error writing to file:", err)
	}
}
