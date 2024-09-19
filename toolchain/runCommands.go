// Copyright (c) 2024 Erik Kassubek, Mario Occhinegro
//
// File: runFullWorkflowMain.go
// Brief: Function to run commands
//
// Author: Erik Kassubek, Mario Occhinegro
// Created: 2024-09-18
// Last Changed 2024-09-18
//
// License: BSD-3-Clause

package main

import (
	"os"
	"os/exec"
)

func runCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func runCommandWithOutput(name, outputFile string, args ...string) error {
	cmd := exec.Command(name, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}

	// Write output to the specified file
	return os.WriteFile(outputFile, output, 0644)
}
