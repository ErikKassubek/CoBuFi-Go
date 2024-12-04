// Copyright (c) 2024 Erik Kassubek
//
// File: select.go
// Brief: File for the selects for fuzzing
//
// Author: Erik Kassubek
// Created: 2024-12-04
//
// License: BSD-3-Clause

package fuzzing

import (
	"fmt"
	"sort"

	"math/rand"
)

var (
	allSelects = make(map[string][]fuzzingSelect) // id -> []fuzzingSelects
)

/*
 * Struct to handle the selects for fuzzing
 *   id (string): file:line
 *   t: tpost of the select execution, used for order
 *   chosenCase (int): id of the chosen case, -1 for default
 *   numberCases (int): number of cases not including default
 *   containsDefault (bool): true if contains default case, otherwise false
 */
type fuzzingSelect struct {
	id              string
	t               int
	chosenCase      int
	numberCases     int
	containsDefault bool
}

func addFuzzingSelect(id string, t int, chosenCase int, numberCases int, containsDefault bool) {
	fs := fuzzingSelect{
		id:              id,
		t:               t,
		chosenCase:      chosenCase,
		numberCases:     numberCases,
		containsDefault: containsDefault,
	}

	allSelects[id] = append(allSelects[id], fs)
}

func sortSelects() {
	for key := range allSelects {
		sort.Slice(allSelects[key], func(i, j int) bool {
			return allSelects[key][i].t < allSelects[key][j].t
		})
	}
}

func (fs fuzzingSelect) toString() string {
	return fmt.Sprintf("%s;%d;%d", fs.id, fs.chosenCase, fs.numberCases)
}

/*
 * Get a copy of fs with a randomly selected case id
 * Args:
 *   def (bool): if true, default is a possible value, if false it is not
 * Return:
 *   (int): the chosen case ID
 */
func (fs fuzzingSelect) getCopyRandom(def bool) fuzzingSelect {
	// if def == false -> rand between 0 and fs.numberCases - 1
	// otherwise rand between -1 and fs.numberCases - 1
	shift := 0
	if def && fs.containsDefault {
		shift = 1
	}

	chosenCase := rand.Intn(fs.numberCases+shift) - shift

	return fuzzingSelect{id: fs.id, t: fs.t, chosenCase: chosenCase, numberCases: fs.numberCases, containsDefault: fs.containsDefault}
}
