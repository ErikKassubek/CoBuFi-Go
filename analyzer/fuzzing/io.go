// Copyright (c) 2024 Erik Kassubek
//
// File: io.go
// Brief: Functions to read ans write the fuzzing files
//
// Author: Erik Kassubek
// Created: 2024-11-28
//
// License: BSD-3-Clause

package fuzzing

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

const fileBlockSep = "###########"

/*
 * read the fuzzing file and store in data
 * Args:
 *   path: path to the fuzzing file
 */
func readFile(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		println("Error opening file: " + filePath)
		panic(err)
	}

	scanner := bufio.NewScanner(file)

	state := 0
	for scanner.Scan() {
		line := scanner.Text()

		// # if the block separator was found, change to the next state
		if line == fileBlockSep {
			state++
			continue
		}

		switch state {
		case 0:
			numberOfPreviousRuns, err = strconv.Atoi(line)
		case 1:
			err = readChannelInfo(line)
		case 2:
			err = readPairInfo(line)
		}
		if err != nil {
			return fmt.Errorf("Fuzzing File invalid format: %s", err.Error())
		}
	}

	if state != 2 {
		return fmt.Errorf("Fuzzing File invalid format: Incorrect number of blocks")
	}

	return nil
}

/*
 * Read the info for a channel. The line should have the following format:
 * fileCreate:lineCreate;closeInfo;qSize;maxQSize
 * closeInfo can be
 *   n: never been closed
 *   c: always been closed
 *   b: in some runs it was closed, but not in all
 * qSize is the buffer size
 * maxQSize is the maximum fullness of the buffer over all operations
 */
func readChannelInfo(line string) error {
	// TODO: implement
	return nil
}

/*
 * Read the info for a send/receive pair
 * It must have the following form:
 * fileSend:lineSend:selCaseSend;fileRecv:lineRecv:selCaseRecv;avgNumberCom
 * selCaseSend and selCaseRecv identify the cases in a select. If the send/recv is in a select, the value is set to the number of the case. If it is not part of a select, it is set to 0.
 */
func readPairInfo(line string) error {
	// TODO implement
	return nil
}

/*
 * Write the current info to a file
 */
func writeFile(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// write the number of previous runs
	file.WriteString(strconv.Itoa(numberOfPreviousRuns))

	file.WriteString(fileBlockSep)

	// write the channel block
	for _, line := range channelInfo {
		file.WriteString(line.toString())
	}

	file.WriteString(fileBlockSep)

	// write the operation pair block
	for _, line := range pairInfo {
		file.WriteString(line.toString())
	}

	return nil
}
