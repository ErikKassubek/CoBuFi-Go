package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	folderName := flag.String("f", "", "path to the file")
	flag.Parse()
	if *folderName == "" {
		fmt.Fprintln(os.Stderr, "Usage generateStatistics -f <folder>")
		os.Exit(1)
	}
	predictedCodes, err := getPredictedBugCounts(*folderName)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(predictedCodes)
}

func getBugCodes(filePath string) []string {
	bugCodes := make([]string, 0)
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		idx := strings.Index(line, ",")
		if idx != -1 {
			bugcode := line[:idx]
			bugCodes = append(bugCodes, bugcode)
		} else {
			fmt.Println("no comma found in line -> invalid format")
		}
	}
	return bugCodes
}

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

func getPredictedBugCounts(folderPath string) (map[string]int, error) {
	codes := []string{
		"A1",
		"A2",
		"A3",
		"A4",
		"A5",
		"P1",
		"P2",
		"P3",
		"L1",
		"L2",
		"L3",
		"L4",
		"L5",
		"L6",
		"L7",
		"L8",
		"L9",
		"L0",
	}
	predictedCodes := make(map[string]int)
	for _, code := range codes {
		predictedCodes[code] = 0
	}

	files, err := getFiles(folderPath, "results_machine.log")
	if err != nil {
		fmt.Println(err)
	}
	for _, file := range files {
		bugCodes := getBugCodes(file)
		for _, code := range bugCodes {
			_, ok := predictedCodes[code]
			if ok {
				predictedCodes[code]++
			}
		}
	}

	return predictedCodes, nil
}
