// Copyright (c) 2024 Erik Kassubek
//
// File: headerUnitTests.go
// Brief: Util functions for the toolchain
//
// Author: Erik Kassubek
// Created: 2024-10-29
//
// License: BSD-3-Clause

package toolchain

import "strings"

// extractTraceNumber extracts the numeric part from a trace directory name
func extractTraceNumber(trace string) string {
	parts := strings.Split(trace, "rewritten_trace_")
	if len(parts) > 1 {
		return parts[1]
	}
	parts = strings.Split(trace, "advocateTraceReplay_")
	if len(parts) > 1 {
		return parts[1]
	}
	return ""
}
