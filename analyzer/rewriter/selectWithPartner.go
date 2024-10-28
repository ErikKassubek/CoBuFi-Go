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

	// // remove everything that is concurrent or after to the select
	hb := clock.GetHappensBefore(sel.GetVC(), partner.GetVC())

	analysis.RemoveConcurrentOrAfter(sel, 0)

	if hb != clock.Before {
		analysis.AddElementToTrace(partner)
	}

	if partner.GetTSort() == 0 {
		if ca.ObjType == "CS" {
			partner.SetTSort(sel.GetTSort() + 1)
		} else {
			sel.(*analysis.TraceElementSelect).SetCase(ca.ID, analysis.RecvOp)
			partner.SetTSort(max(1, sel.GetTSort()-1))
		}
	}

	if ca.ObjType == "CS" {
		sel.(*analysis.TraceElementSelect).SetCase(ca.ID, analysis.SendOp)
	} else {
		sel.(*analysis.TraceElementSelect).SetCase(ca.ID, analysis.RecvOp)
	}

	analysis.AddTraceElementReplay(max(partner.GetTSort(), sel.GetTSort())+1, exitCodeNone)

	return nil
}
