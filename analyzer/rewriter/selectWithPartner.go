// Copyrigth (c) 2024 Erik Kassubek
//
// File: rewriter.go
// Brief: Main functions to rewrite the trace
//
// Author: Erik Kassubek
// Created: 2024-10-25
//
// License: BSD-3-Clause

// Package rewriter provides functions for rewriting traces.

package rewriter

import (
	"analyzer/analysis"
	"analyzer/bugs"
	"analyzer/clock"
	"fmt"
)

// ========= Not executed select with partner =========================
func rewriteNotExecutedSelect(bug bugs.Bug, index int) error {
	sel := bug.TraceElement1[0]
	ca := bug.TraceElement1Sel[0]
	partners := bug.TraceElement2
	if index >= len(partners) {
		return fmt.Errorf("Index %d not in possible partner (len %d)", index, len(partners))
	}
	partner := bug.TraceElement2[index]

	// remove everything that is concurrent or after to the select
	hb := clock.GetHappensBefore(sel.GetVC(), partner.GetVC())

	if hb == clock.Before {
		analysis.RemoveConcurrentOrAfter(sel, 0)
		if partner.GetObjType() == "CC" || partner.GetObjType() == "CS" {
			partner.SetTSort(sel.GetTSort() - 1)
		} else {
			partner.SetTSort(sel.GetTSort() + 1)
		}

		analysis.AddElementToTrace(partner)

	} else {
		analysis.RemoveConcurrentOrAfter(partner, 0)
		if partner.GetObjType() == "CC" || partner.GetObjType() == "CS" {
			sel.SetTSort(partner.GetTSort() + 1)
		} else {
			sel.SetTSort(partner.GetTSort() - 1)
		}

		analysis.AddElementToTrace(sel)
	}

	sel.(*analysis.TraceElementSelect).SetChosenCase(ca.Index)

	analysis.AddTraceElementReplay(partner.GetTSort()+3, exitCodeNone)

	return nil
}
