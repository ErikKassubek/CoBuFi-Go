// Copyrigth (c) 2024 Erik Kassubek
//
// File: progInfo.go
// Brief: Read the info required for running the program
//
// Author: Erik Kassubek
// Created: 2024-06-18
//
// License: BSD-3-Clause

package explanation

import (
	"errors"
	"os"
	"strings"
)

func readProgInfo(path string, index int) (map[string]string, error) {
	res := make(map[string]string)

	file, err := os.ReadFile(path + "output.log")
	if err != nil {
		return res, err
	}

	lines := strings.Split(string(file), "\n")

	if len(lines) < 3 {
		return res, errors.New("output file is too short")
	}

	for i := 0; i < len(lines); i++ {
		if lines[i] == "" {
			continue
		}

		println(lines[i])

		// if strings.Contains(lines[i], "unitTestheaderInserter") {
		// 	if strings.Contains(lines[i], "-r true") {
		// 		line := lines[i][:strings.LastIndex(lines[i], " ")]
		// 		res["inserterReplay"] = line + " " + strconv.Itoa(index)
		// 	} else {
		// 		res["inserterRecord"] = lines[i]
		// 	}
		// } else if strings.Contains(lines[i], "unitTestheaderRemover") {
		// 	res["remover"] = lines[i]
		// } else if strings.Contains(lines[i], "-run") {
		// 	res["run"] = lines[i]
		if strings.Contains(lines[i], "FileName: ") {
			println("FOUND")
			res["file"] = strings.TrimPrefix(lines[i], "FileName: ")
			println(res["file"])
		} else if strings.Contains(lines[i], "TestName: ") {
			res["name"] = strings.TrimPrefix(lines[i], "TestName: ")
		} else if strings.Contains(lines[i], "Import added at line: ") {
			res["importLine"] = strings.TrimPrefix(lines[i], "Import added at line: ")
		} else if strings.Contains(lines[i], "header added at line: ") {
			res["headerLine"] = strings.TrimPrefix(lines[i], "header added at line: ")
		}
	}

	return res, nil
}

func getProgInfo(info map[string]string, key string) string {
	if _, ok := info[key]; !ok {
		return "Failed to read command for " + key
	}

	if info[key] == "" {
		return "Failed to read command for " + key
	}

	return info[key]
}
