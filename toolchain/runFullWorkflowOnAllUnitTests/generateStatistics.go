package main

import (
	"bufio"
	"time"
	"fmt"
	"os"
	"strings"
)

func main() {
	fmt.Println("Starting Program")
	bugCodes := getBugCodes("./results_machine.log")
	fmt.Print(bugCodes)
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
		time.Sleep(5000)
	}
	return bugCodes
}
