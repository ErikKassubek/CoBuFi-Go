// Copyrigth (c) 2024 Erik Kassubek
//
// File: analysisWaitGroup.go
// Brief: Trace analysis for possible negative wait group counter
//
// Author: Erik Kassubek <kassubek.erik@gmail.com>
// Created: 2023-11-24
// LastChange: 2024-09-01
//
// License: BSD-3-Clause

package analysis

import (
	"analyzer/clock"
	"analyzer/logging"
	"analyzer/utils"
	"errors"
	"strconv"
)

func checkForDoneBeforeAddChange(wa *TraceElementWait) {
	if wa.delta > 0 {
		checkForDoneBeforeAddAdd(wa)
	} else if wa.delta < 0 {
		checkForDoneBeforeAddDone(wa)
	} else {
		// checkForImpossibleWait(routine, id, pos, vc)
	}
}

func checkForDoneBeforeAddAdd(wa *TraceElementWait) {
	// if necessary, create maps and lists
	if _, ok := wgAdd[wa.id]; !ok {
		wgAdd[wa.id] = make(map[int][]*TraceElementWait)
	}
	if _, ok := wgAdd[wa.id][wa.routine]; !ok {
		wgAdd[wa.id][wa.routine] = make([]*TraceElementWait, 0)
	}

	// add the vector clock and position to the list
	for i := 0; i < wa.delta; i++ {
		if wa.delta > 1 {
			wa.tID = wa.tID + "+" + strconv.Itoa(i) // add a unique identifier to the position
		}
		wgAdd[wa.id][wa.routine] = append(wgAdd[wa.id][wa.routine], wa)
	}
}

func checkForDoneBeforeAddDone(wa *TraceElementWait) {
	// if necessary, create maps and lists
	if _, ok := wgDone[wa.id]; !ok {
		wgDone[wa.id] = make(map[int][]*TraceElementWait)

	}
	if _, ok := wgDone[wa.id][wa.routine]; !ok {
		wgDone[wa.id][wa.routine] = make([]*TraceElementWait, 0)
	}

	// add the vector clock and position to the list
	wgDone[wa.id][wa.routine] = append(wgDone[wa.id][wa.routine], wa)
}

/*
 * Build a st graph for a wait group.
 * The graph has the following structure:
 * - a start node s
 * - a end node t
 * - edges from s to all done operations
 * - edges from all add operations to t
 * - edges from done to add if the add happens before the done
 * Args:
 *   adds (map[int][]|*TraceElementWait): The add operations
 *   dones (map[int][]|*TraceElementWait): The done operations
 * Returns:
 *   []Edge: The graph
 */
func buildResidualGraph(adds map[int][]*TraceElementWait, dones map[int][]*TraceElementWait) map[string][]string {
	graph := make(map[string][]string, 0)
	graph["s"] = []string{}
	graph["t"] = []string{}

	// add edges from s to all done operations
	for _, done := range dones {
		for _, vc := range done {
			graph[vc.tID] = []string{}
			graph["s"] = append(graph["s"], vc.tID)
		}
	}

	// add edges from all add operations to t
	for _, add := range adds {
		for _, vc := range add {
			graph[vc.tID] = []string{"t"}
		}
	}

	// add edge from done to add if the add happens before the done
	for _, done := range dones {
		for _, vcDone := range done {
			for _, add := range adds {
				for _, vcAdd := range add {
					if clock.GetHappensBefore(vcAdd.vc, vcDone.vc) == clock.Before {
						graph[vcDone.tID] = append(graph[vcDone.tID], vcAdd.tID)

					}
				}
			}
		}
	}

	return graph
}

/*
 * Calculate the maximum flow of a graph using the ford fulkerson algorithm
 * Args:
 *   graph ([]Edge): The graph
 * Returns:
 *   int: The maximum flow
 */
func calculateMaxFlow(graph map[string][]string) (int, map[string][]string) {
	maxFlow := 0
	for {
		path, flow := findPath(graph)
		if flow == 0 {
			break
		}

		maxFlow += flow
		for i := 0; i < len(path)-1; i++ {
			graph[path[i]] = append(graph[path[i]], path[i+1])
			graph[path[i+1]] = remove(graph[path[i+1]], path[i])
		}
	}

	return maxFlow, graph
}

/*
 * Find a path in a graph using a breadth-first search
 * Args:
 *   graph ([]Edge): The graph
 * Returns:
 *   []string: The path
 *   int: The flow
 */
func findPath(graph map[string][]string) ([]string, int) {
	visited := make(map[string]bool, 0)
	queue := []string{"s"}
	visited["s"] = true
	parents := make(map[string]string, 0)

	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]

		if node == "t" {
			path := []string{}
			for node != "s" {
				path = append(path, node)
				node = parents[node]
			}
			path = append(path, "s")

			return path, 1
		}

		for _, neighbor := range graph[node] {
			if !visited[neighbor] {
				visited[neighbor] = true
				queue = append(queue, neighbor)
				parents[neighbor] = node
			}
		}
	}

	return []string{}, 0
}

/*
 * Remove an element from a list
 * Args:
 *   list ([]string): The list
 *   element (string): The element to remove
 * Returns:
 *   []string: The list without the element
 */
func remove(list []string, element string) []string {
	for i, e := range list {
		if e == element {
			list = append(list[:i], list[i+1:]...)
			return list
		}
	}
	return list
}

func numberDone(id int) int {
	res := 0
	for _, dones := range wgDone[id] {
		res += len(dones)
	}
	return res
}

/*
- Check if a wait group counter could become negative
- For each done operation, build a bipartite st graph.
- Use the Ford-Fulkerson algorithm to find the maximum flow.
- If the maximum flow is smaller than the number of done operations, a negative wait group counter is possible.
*/
func CheckForDoneBeforeAdd() {
	for id := range wgAdd { // for all waitgroups
		graph := buildResidualGraph(wgAdd[id], wgDone[id])

		maxFlow, graph := calculateMaxFlow(graph)
		nrDone := numberDone(id)

		addsVcTIDs := []*TraceElementWait{}
		donesVcTIDs := []*TraceElementWait{}

		if maxFlow < nrDone {
			// sort the adds and dones, that do not have a partner is such a way,
			// that the i-th add in the result message is concurrent with the
			// i-th done in the result message

			for _, adds := range wgAdd[id] {
				for _, add := range adds {
					if !utils.Contains(graph["t"], add.tID) {
						addsVcTIDs = append(addsVcTIDs, add)
					}
				}
			}
			for _, dones := range graph["s"] {
				doneVcTID, err := getDoneElemFromTID(id, dones)
				if err != nil {
					logging.Debug(err.Error(), logging.ERROR)
				} else {
					donesVcTIDs = append(donesVcTIDs, doneVcTID)
				}
			}

			addsVcTIDSorted := make([]*TraceElementWait, 0)
			donesVcTIDSorted := make([]*TraceElementWait, 0)

			for i := 0; i < len(addsVcTIDs); i++ {
				for j := 0; j < len(donesVcTIDs); j++ {
					if clock.GetHappensBefore(addsVcTIDs[i].vc, addsVcTIDs[j].vc) == clock.Concurrent {
						addsVcTIDSorted = append(addsVcTIDSorted, addsVcTIDs[i])
						donesVcTIDSorted = append(donesVcTIDSorted, donesVcTIDs[j])
						// remove the element from the list
						addsVcTIDs = append(addsVcTIDs[:i], addsVcTIDs[i+1:]...)
						donesVcTIDs = append(donesVcTIDs[:j], donesVcTIDs[j+1:]...)
						// fix the index
						i--
						j = 0
					}
				}
			}

			args1 := []logging.ResultElem{} // adds
			args2 := []logging.ResultElem{} // dones

			for _, add := range addsVcTIDSorted {
				if add.tID == "\n" {
					continue
				}
				file, line, tPre, err := infoFromTID(add.tID)
				if err != nil {
					logging.Debug(err.Error(), logging.ERROR)
					return
				}

				args1 = append(args1, logging.TraceElementResult{
					RoutineID: add.routine,
					ObjID:     id,
					TPre:      tPre,
					ObjType:   "WA",
					File:      file,
					Line:      line,
				})

			}

			for _, done := range donesVcTIDSorted {
				if done.tID == "\n" {
					continue
				}
				file, line, tPre, err := infoFromTID(done.tID)
				if err != nil {
					logging.Debug(err.Error(), logging.ERROR)
					return
				}

				args2 = append(args2, logging.TraceElementResult{
					RoutineID: done.routine,
					ObjID:     id,
					TPre:      tPre,
					ObjType:   "WD",
					File:      file,
					Line:      line,
				})
			}

			logging.Result(logging.CRITICAL, logging.PNegWG,
				"add", args1, "done", args2)
		}
	}
}

func getDoneElemFromTID(id int, tID string) (*TraceElementWait, error) {
	for _, dones := range wgDone[id] {
		for _, done := range dones {
			if done.tID == tID {
				return done, nil
			}
		}
	}
	return nil, errors.New("Could not find done operation with tID " + tID)
}
