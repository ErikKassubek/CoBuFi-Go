// Copyright (c) 2024 Erik Kassubek
//
// File: trace.go
// Brief: Function to parse the trace and get all relevant information
//
// Author: Erik Kassubek
// Created: 2024-11-29
//
// License: BSD-3-Clause

package fuzzing

import "analyzer/analysis"

/*
 * Parse the current trace and record all relevant data
 */
func parseTrace() {
	for _, trace := range *analysis.GetTraces() {
		for _, elem := range trace {
			switch elem.(type) {
			case *analysis.TraceElementNew:

			case *analysis.TraceElementChannel:

			}
		}

	}

}
