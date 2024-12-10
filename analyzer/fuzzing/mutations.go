// Copyright (c) 2024 Erik Kassubek
//
// File: mutations.go
// Brief: Create the mutations
//
// Author: Erik Kassubek
// Created: 2024-12-03
//
// License: BSD-3-Clause

package fuzzing

func createMutations(numberMutations int) []map[string][]fuzzingSelect {
	res := make([]map[string][]fuzzingSelect, 0)

	for i := 0; i < numberMutations; i++ {
		res = append(res, getMutations())
	}

	return res
}

func getMutations() map[string][]fuzzingSelect {
	res := make(map[string][]fuzzingSelect)

	for key, listSel := range allSelects {
		res[key] = make([]fuzzingSelect, 0)
		for _, sel := range listSel {
			res[key] = append(res[key], sel.getCopyRandom(false))
		}
	}

	return res
}
