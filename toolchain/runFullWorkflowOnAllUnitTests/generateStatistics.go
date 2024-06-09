package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
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
	possibleCodes := make(map[string]int)
	for _, code := range codes {
		possibleCodes[code] = 0
	}
	fmt.Println("Starting Program")
	bugCodes := getBugCodes("./results_machine.log")
	for _, code := range bugCodes {
		_, ok := possibleCodes[code]
		if ok {
			possibleCodes[code]++
		}
	}
	fmt.Println(possibleCodes)
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
