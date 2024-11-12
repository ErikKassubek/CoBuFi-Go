// Copyrigth (c) 2024 Erik Kassubek
//
// File: analysisUtil.go
// Brief: Collection of utility functiond for trace analysis
//
// Author: Erik Kassubek
// Created: 2024-05-29
//
// License: BSD-3-Clause

package analysis

import (
	"analyzer/utils"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

/*
 * Get the info from a TID
 * Args:
 *   tID (string): The TID
 * Return:
 *   string: the file
 *   int: the line
 *   int: the tPre
 *   error: the error
 */
func infoFromTID(tID string) (string, int, int, error) {
	spilt1 := utils.SplitAtLast(tID, "@")

	if len(spilt1) != 2 {
		return "", 0, 0, errors.New(fmt.Sprint("TID not correct: no @: ", tID))
	}

	split2 := strings.Split(spilt1[0], ":")
	if len(split2) != 2 {
		return "", 0, 0, errors.New(fmt.Sprint("TID not correct: no ':': ", tID))
	}

	tPre, err := strconv.Atoi(spilt1[1])
	if err != nil {
		return "", 0, 0, err
	}

	line, err := strconv.Atoi(split2[1])
	if err != nil {
		return "", 0, 0, err
	}

	return split2[0], line, tPre, nil
}

func sameRoutine(elems ...[]TraceElement) bool {
	ids := make(map[int]int)
	for _, elem := range elems {
		for i, e := range elem {
			if _, ok := ids[i]; !ok {
				ids[i] = e.GetRoutine()
			} else if ids[i] != e.GetRoutine() {
				return false
			}
		}
	}

	return true
}
