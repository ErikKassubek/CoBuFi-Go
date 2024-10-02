// Copyrigth (c) 2024 Erik Kassubek
//
// File: replay.go
// Brief: Read the info about the rewrite and replay of the bug
//
// Author: Erik Kassubek
// Created: 2024-06-18
//
// License: BSD-3-Clause

package explanation

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func getRewriteInfo(bugType string, codes map[string]string, index int) map[string]string {
	res := make(map[string]string)

	rewPos := rewriteType[bugType]

	res["description"] = ""
	res["exitCode"] = ""
	res["exitCodeExplanation"] = ""
	res["replaySuc"] = "was not possible"

	var err error

	if rewPos == "Actual" {
		res["description"] += "The bug is an actual bug. Therefore no rewrite is possibel."
	} else if rewPos == "Possible" {
		res["description"] += "The bug is a potential bug.\n"
		res["description"] += "The analyzer has tries to rewrite the trace in such a way, "
		res["description"] += "that the bug will be triggered when replaying the trace."
	} else if rewPos == "LeakPos" {
		res["description"] += "The analyzer found a leak in the recorded trace.\n"
		res["description"] += "The analyzer found a way to resolve the leak, meaning the "
		res["description"] += "leak should not reappear in the rewritten trace."
	} else if rewPos == "Leak" {
		res["description"] += "The analyzer found a leak in the recorded trace.\n"
		res["description"] += "The analyzer could not find a way to resolve the leak. "
		res["description"] += "No rewritten trace was created. This does not need to mean, "
		res["description"] += "that the leak can not be resolved, especially because the "
		res["description"] += "analyzer is only aware of executed operations."
	}
	res["exitCode"], res["exitCodeExplanation"], res["replaySuc"], err = getReplayInfo(codes, index)

	if err != nil {
		fmt.Println("Error getting replay info: ", err)
	}

	return res

}

func getOutputCodes(path string) map[string]string {
	output := filepath.Join(path, "output.log")
	if _, err := os.Stat(output); os.IsNotExist(err) {
		res := "No replay info available. Output.log does not exist."
		return map[string]string{"AdvocateFailExplanationInfo": res, "AdvocateFailResplaySucInfo": "information not available"}
	}

	// read the output file
	content, err := os.ReadFile(output)
	if err != nil {
		res := "No replay info available. Could not read output.log file"
		return map[string]string{"AdvocateFailExplanationInfo": res, "AdvocateFailResplaySucInfo": "information not available"}
	}

	lines := strings.Split(string(content), "\n")

	replayPos := make(map[string]bool)
	replayCode := make(map[string]string)
	bugrepPrefix := "Bugreport info: "
	replayReadPrefix := "Reading trace from rewritten_trace_"
	exitCodePrefix := "Exit Replay with code"

	lastReplayIndex := ""
	lastReplayIndexInfoFound := true

	for _, line := range lines {
		if strings.HasPrefix(line, bugrepPrefix) {
			line = strings.TrimPrefix(line, bugrepPrefix)
			lineSplit := strings.Split(line, ",")
			if len(lineSplit) == 2 {
				index := lineSplit[0]
				if lineSplit[1] == "suc" {
					replayPos[index] = true
				} else {
					replayPos[index] = false
				}
				if lineSplit[1] == "double" {
					replayCode[index] = "double"
				}
				if lineSplit[1] == "fail" {
					replayCode[index] = "fail"
				}
			}
		} else if strings.HasPrefix(line, replayReadPrefix) {
			if !lastReplayIndexInfoFound {
				replayCode[lastReplayIndex] = "panic"
			}
			lastReplayIndex = strings.TrimPrefix(line, replayReadPrefix)
			lastReplayIndexInfoFound = false
		} else if strings.HasPrefix(line, exitCodePrefix) {
			line = strings.TrimPrefix(line, exitCodePrefix)
			line = strings.TrimSpace(line)
			replayCode[lastReplayIndex] = strings.Split(line, " ")[0]
			lastReplayIndexInfoFound = true
		}
	}

	return replayCode
}

func getReplayInfo(replayCode map[string]string, index int) (string, string, string, error) {
	if _, ok := replayCode["AdvocateFailExplanationInfo"]; ok {
		fmt.Println("Could not read")
		return "", replayCode["AdvocateFailExplanationInfo"], replayCode["AdvocateFailResplaySucInfo"], fmt.Errorf("Could not read output file")
	}

	exitCode := replayCode[fmt.Sprint(index)]
	replaySuc := "failed"
	if exitCode == "double" {
		replaySuc = "was already performed for this bug in another test"
		return "double", "", replaySuc, nil
	}
	if exitCode == "fail" {
		return "fail", exitCodeExplanation["fail"], "was not run", nil
	}
	if exitCode == "panic" {
		return "panic", exitCodeExplanation["panic"], "was terminated unexpectedly", nil
	}

	exitCodeInt, err := strconv.Atoi(exitCode)
	if err != nil {
		res := fmt.Sprintf("Invalid format in output.log. Could not convert exit code %s to int for index %d", exitCode, index)
		return "", res, "failed", errors.New(res)
	}
	if exitCodeInt >= 20 || exitCodeInt == 0 {
		replaySuc = "was successful"
	}

	return exitCode, exitCodeExplanation[exitCode], replaySuc, nil
}
