// Copyrigth (c) 2024 Erik Kassubek
//
// File: analysisSelectPartner.go
// Brief: Trace analysis for detection of select cases without any possible partners
//
// Author: Erik Kassubek
// Created: 2024-03-04
//
// License: BSD-3-Clause

package analysis

import (
	"analyzer/clock"
	"analyzer/logging"
	"strconv"
	"strings"
)

/*
* CheckForSelectCaseWithoutPartner checks for select cases without a valid
* partner. Call when all elements have been processed.
 */
func CheckForSelectCaseWithoutPartner() {
	// check if not selected cases could be partners
	for i, c1 := range selectCases {
		for j := i + 1; j < len(selectCases); j++ {
			c2 := selectCases[j]

			// if c1.partnerFound && c2.partnerFound {
			// 	continue
			// }

			if c1.chanID != c2.chanID || c1.vcTID.TID == c2.vcTID.TID || c1.send == c2.send {
				continue
			}

			if c2.send { // c1 should be send, c2 should be recv
				c1, c2 = c2, c1
			}

			hb := clock.GetHappensBefore(c1.vcTID.Vc, c2.vcTID.Vc)
			found := false
			if c1.buffered && (hb == clock.Concurrent || hb == clock.After) {
				found = true
			} else if !c1.buffered && hb == clock.Concurrent {
				found = true
			}

			if found {
				selectCases[i].partnerFound = true
				selectCases[j].partnerFound = true
				selectCases[i].partner = append(selectCases[i].partner, VectorClockTID3{selectCases[j].sel, selectCases[j].sel.GetVC(), 0})
				selectCases[j].partner = append(selectCases[j].partner, VectorClockTID3{selectCases[i].sel, selectCases[i].sel.GetVC(), 0})
			}
		}
	}

	if len(selectCases) == 0 {
		return
	}

	// collect all cases with no partner and all not triggered cases with partner

	casesWithoutPartner := make(map[string][]logging.ResultElem) // tID -> cases
	casesWithoutPartnerInfo := make(map[string][]int)            // tID -> [routine, selectID]

	for cIndex, c := range selectCases {
		opjType := "C"
		if c.send {
			opjType += "S"
		} else {
			opjType += "R"
		}

		partnerResult := make([]logging.ResultElem, 0)

		if c.partnerFound {
			if c.exec {
				continue
			}

			file, line, tPre, err := infoFromTID(c.vcTID.TID)
			if err != nil {
				continue
			}

			sel := logging.TraceElementResult{
				RoutineID: c.vcTID.Routine,
				ObjID:     c.sel.GetID(),
				TPre:      tPre,
				ObjType:   "SS",
				File:      file,
				Line:      line,
			}

			ca := logging.SelectCaseResult{
				SelID:   c.sel.GetID(),
				ObjID:   c.chanID,
				ObjType: opjType,
				Routine: c.vcTID.Routine,
				Index:   cIndex,
			}

			for _, p := range c.partner {
				pos := strings.Split(p.Elem.GetPos(), ":")
				if len(pos) < 2 {
					continue
				}

				line, err := strconv.Atoi(pos[1])
				if err != nil {
					continue
				}

				partner := logging.TraceElementResult{
					RoutineID: p.Elem.GetRoutine(),
					ObjID:     p.Elem.GetID(),
					TPre:      p.Elem.GetTPre(),
					ObjType:   "SS",
					File:      pos[0],
					Line:      line,
				}

				partnerResult = append(partnerResult, partner)
			}

			if len(partnerResult) == 0 {
				continue
			}

			logging.Result(logging.INFORMATION, logging.SNotExecutedWithPartner,
				"select", []logging.ResultElem{sel, ca}, "partner", partnerResult)
			continue
		}

		arg2 := logging.SelectCaseResult{
			SelID:   c.sel.GetID(),
			ObjID:   c.chanID,
			ObjType: opjType,
			Routine: c.vcTID.Routine,
			Index:   cIndex,
		}

		if _, ok := casesWithoutPartner[c.vcTID.TID]; !ok {
			casesWithoutPartner[c.vcTID.TID] = make([]logging.ResultElem, 0)
			casesWithoutPartnerInfo[c.vcTID.TID] = []int{c.vcTID.Routine, c.sel.GetID()}
		}

		casesWithoutPartner[c.vcTID.TID] = append(casesWithoutPartner[c.vcTID.TID], arg2)
	}

	for tID, cases := range casesWithoutPartner {
		if len(cases) == 0 {
			continue
		}

		info := casesWithoutPartnerInfo[tID]
		if len(info) != 2 {
			logging.Debug("info should have 2 elements", logging.ERROR)
			continue
		}

		file, line, tPre, err := infoFromTID(tID)
		if err != nil {
			logging.Debug(err.Error(), logging.ERROR)
			continue
		}

		arg1 := logging.TraceElementResult{
			RoutineID: info[0],
			ObjID:     info[1],
			TPre:      tPre,
			ObjType:   "SS",
			File:      file,
			Line:      line,
		}

		logging.Result(logging.WARNING, logging.ASelCaseWithoutPartner,
			"select", []logging.ResultElem{arg1}, "case", cases)
	}
}

/*
* CheckForSelectCaseWithoutPartnerSelect checks for select cases without a valid
* partner. Call whenever a select is processed.
* Args:
*   se (*TraceElementSelect): The trace element
*   ids ([]int): The ids of the channels
*   bufferedInfo ([]bool): The buffer status of the channels
*   sendInfo ([]bool): The send status of the channels
*   vc (VectorClock): The vector clock
 */
//  func CheckForSelectCaseWithoutPartnerSelect(routine int, selectID int, caseChanIds []int, bufferedInfo []bool,
func CheckForSelectCaseWithoutPartnerSelect(se *TraceElementSelect, caseChanIds []int, bufferedInfo []bool,
	sendInfo []bool, vc clock.VectorClock) {
	for i, id := range caseChanIds {
		buffered := bufferedInfo[i]
		send := sendInfo[i]

		found := false
		executed := false
		var partner = make([]VectorClockTID3, 0)

		if i == se.chosenIndex && se.tPost != 0 {
			// no need to check if the channel is the chosen case
			executed = true
			p := se.GetPartner()
			if p != nil {
				found = true
				vcTID := VectorClockTID3{
					p, p.vc.Copy(), 0,
				}
				partner = append(partner, vcTID)
			}
		} else {
			// not select cases
			if send {
				for _, mrr := range mostRecentReceive {
					if possiblePartner, ok := mrr[id]; ok {
						hb := clock.GetHappensBefore(vc, possiblePartner.Vc)
						if buffered && (hb == clock.Concurrent || hb == clock.Before) {
							found = true
							partner = append(partner, possiblePartner)
						} else if !buffered && hb == clock.Concurrent {
							found = true
							partner = append(partner, possiblePartner)
						}
					}
				}
			} else { // recv
				for _, mrs := range mostRecentSend {
					if possiblePartner, ok := mrs[id]; ok {
						hb := clock.GetHappensBefore(vc, possiblePartner.Vc)
						if buffered && (hb == clock.Concurrent || hb == clock.After) {
							found = true
							partner = append(partner, possiblePartner)
						} else if !buffered && hb == clock.Concurrent {
							found = true
							partner = append(partner, possiblePartner)
						}
					}
				}
			}
		}

		selectCases = append(selectCases,
			allSelectCase{se, id, VectorClockTID{vc, se.tID, se.routine}, send, buffered, found, partner, executed})

	}
}

/*
* CheckForSelectCaseWithoutPartnerChannel checks for select cases without a valid
* partner. Call whenever a channel operation is processed.
* Args:
*   id (int): The id of the channel
*   vc (VectorClock): The vector clock
*   tID (string): The position of the channel operation in the program
*   send (bool): True if the operation is a send
*   buffered (bool): True if the channel is buffered
*   sel (bool): True if the operation is part of a select statement
 */
func CheckForSelectCaseWithoutPartnerChannel(ch TraceElement, vc clock.VectorClock,
	send bool, buffered bool) {

	for i, c := range selectCases {
		if c.partnerFound || c.chanID != ch.GetID() || c.send == send || c.vcTID.TID == ch.GetTID() {
			continue
		}

		hb := clock.GetHappensBefore(vc, c.vcTID.Vc)
		found := false
		if send {
			if buffered && (hb == clock.Concurrent || hb == clock.Before) {
				found = true
			} else if !buffered && hb == clock.Concurrent {
				found = true
			}
		} else {
			if buffered && (hb == clock.Concurrent || hb == clock.After) {
				found = true
			} else if !buffered && hb == clock.Concurrent {
				found = true
			}
		}

		if found {
			selectCases[i].partnerFound = true
			selectCases[i].partner = append(selectCases[i].partner, VectorClockTID3{ch, vc, 0})
		}
	}
}

/*
* CheckForSelectCaseWithoutPartnerClose checks for select cases without a valid
* partner. Call whenever a close operation is processed.
* Args:
*   id (int): The id of the channel
*   vc (VectorClock): The vector clock
 */
func CheckForSelectCaseWithoutPartnerClose(cl *TraceElementChannel, vc clock.VectorClock) {
	for i, c := range selectCases {
		if c.partnerFound || c.chanID != cl.id || c.send {
			continue
		}

		hb := clock.GetHappensBefore(vc, c.vcTID.Vc)
		found := false
		if c.buffered && (hb == clock.Concurrent || hb == clock.After) {
			found = true
		} else if !c.buffered && hb == clock.Concurrent {
			found = true
		}

		if found {
			selectCases[i].partnerFound = true
			selectCases[i].partner = append(selectCases[i].partner, VectorClockTID3{cl, vc, 0})
		}
	}
}
