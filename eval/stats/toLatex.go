// Copyright (c) 2024 Erik Kassubek
//
// File: toLatex.go
// Brief: Create latex tables from the data
//
// Author: Erik Kassubek
// Created: 2024-09-30
//
// License: BSD-3-Clause

package main

import (
	"fmt"
)

func createLatex(fileName string) {
	writeBoilerplate("start", fileName)

	createLatexTable("dataMin", fileName)
	createLatexTable("dataMinSize", fileName)
	createLatexTable("dataSize", fileName)
	createLatexTable("dataActual", fileName)
	createLatexTable("dataPotential", fileName)
	createLatexTable("dataLeak", fileName)

	createLatexTable("time", fileName)

	writeExplanation(fileName)

	writeBoilerplate("end", fileName)
}

func createLatexTable(name string, fileName string) {
	table := getTableTopLine(name)

	for _, prog := range progs {
		table += getTableRows(prog, name)
	}

	if name == "time" {
		table += getAvgTime()
	}

	table += "\\end{tabular}\n\\caption{"
	table += name
	table += "}\n\\label{Tab:"
	table += name
	table += "}\n\\end{table}"

	writeToFile(fileName, table)
}

func getTableRows(data progData, size string) string {
	switch size {
	case "dataMin":
		return fmt.Sprintf("%s & %s & %s & %s & %s & %s & %s & %s  \\\\ \\hline\n",
			data.name,
			gv(data.numberDetected["A"]),
			gv(data.numberDetected["P"]),
			gv(data.numberDetected["L"]),
			gv(data.numberRewritten["P"]),
			gv(data.numberRewritten["L"]),
			gv(data.numberReplayed["P"]),
			gv(data.numberReplayed["L"]),
		)
	case "dataMinSize":
		return fmt.Sprintf("%s & %s & %s & %s & %s & %s & %s & %s & %s & %s & %s  \\\\ \\hline\n",
			data.name,
			gv(data.numberLines),
			gv(data.numberTests),
			gv(data.numberOperations),
			gv(data.numberDetected["A"]),
			gv(data.numberDetected["P"]),
			gv(data.numberDetected["L"]),
			gv(data.numberRewritten["P"]),
			gv(data.numberRewritten["L"]),
			gv(data.numberReplayed["P"]),
			gv(data.numberReplayed["L"]),
		)
	case "dataSize":
		return fmt.Sprintf("%s & %s & %s & %s   \\\\ \\hline\n",
			data.name,
			gv(data.numberLines),
			gv(data.numberTests),
			gv(data.numberOperations),
		)
	case "dataActual":
		return fmt.Sprintf("%s & %s & %s & %s & %s & %s \\\\ \\hline\n",
			data.name,
			gv(data.numberDetected["A01"]),
			gv(data.numberDetected["A02"]),
			gv(data.numberDetected["A03"]),
			gv(data.numberDetected["A04"]),
			gv(data.numberDetected["A05"]),
		)
	case "dataPotential":
		return fmt.Sprintf("%s & %s & %s & %s & %s & %s & %s & %s & %s & %s & %s & %s & %s \\\\ \\hline\n",
			data.name,
			gv(data.numberDetected["P01"]),
			gv(data.numberDetected["P02"]),
			gv(data.numberDetected["P03"]),
			gv(data.numberDetected["P04"]),
			gv(data.numberRewritten["P01"]),
			gv(data.numberRewritten["P02"]),
			gv(data.numberRewritten["P03"]),
			gv(data.numberRewritten["P04"]),
			gv(data.numberReplayed["P01"]),
			gv(data.numberReplayed["P02"]),
			gv(data.numberReplayed["P03"]),
			gv(data.numberReplayed["P04"]),
		)
	case "dataLeak":
		return fmt.Sprintf("%s & %s & %s & %s & %s & %s & %s & %s & %s & %s & %s & %s & %s & %s & %s & %s & %s & %s & %s& %s & %s & %s & %s & %s & %s & %s & %s & %s & %s& %s & %s  \\\\ \\hline\n",
			data.name,
			gv(data.numberDetected["P01"]),
			gv(data.numberDetected["P02"]),
			gv(data.numberDetected["P03"]),
			gv(data.numberDetected["P04"]),
			gv(data.numberDetected["P05"]),
			gv(data.numberDetected["P06"]),
			gv(data.numberDetected["P07"]),
			gv(data.numberDetected["P08"]),
			gv(data.numberDetected["P09"]),
			gv(data.numberDetected["P10"]),
			gv(data.numberRewritten["P01"]),
			gv(data.numberRewritten["P02"]),
			gv(data.numberRewritten["P03"]),
			gv(data.numberRewritten["P04"]),
			gv(data.numberRewritten["P05"]),
			gv(data.numberRewritten["P06"]),
			gv(data.numberRewritten["P07"]),
			gv(data.numberRewritten["P08"]),
			gv(data.numberRewritten["P09"]),
			gv(data.numberRewritten["P10"]),
			gv(data.numberReplayed["P01"]),
			gv(data.numberReplayed["P02"]),
			gv(data.numberReplayed["P03"]),
			gv(data.numberReplayed["P04"]),
			gv(data.numberReplayed["P05"]),
			gv(data.numberReplayed["P06"]),
			gv(data.numberReplayed["P07"]),
			gv(data.numberReplayed["P08"]),
			gv(data.numberReplayed["P09"]),
			gv(data.numberReplayed["P10"]),
		)
	case "time":
		overheadRecord := -1.
		overheadReplay := -1.

		if data.timeRun != 0 {
			overheadRecord = (data.timeRecord - data.timeRun) / data.timeRun * 100.
		}
		if data.timeRun != 0 {
			overheadReplay = (data.timeReplay - data.timeRun) / data.timeRun * 100.
		}

		return fmt.Sprintf("%s & %.2f & %.2f & %.2f & %.2f & %.2f & %.2f \\\\ \\hline\n",
			data.name,
			data.timeRun,
			data.timeRecord,
			data.timeAnalysis,
			data.timeReplay,
			max(0, overheadRecord),
			max(0, overheadReplay),
		)
	}
	return ""
}

func writeExplanation(fileName string) {
	text := "\\begin{itemize}\n"
	text += "\\item $\\mathcal{S}_L$: number lines of program\n"
	text += "\\item $\\mathcal{S}_T$: number traces / number of run tests\n"
	text += "\\item $\\mathcal{S}_O$: size of traces (total number of operations in all traces)\n"
	text += "\\item $\\mathcal{D}_A$: total number of detection of unique actual detected bugs\n"
	text += "\\item $\\mathcal{D}_P$: total number of detection of unique potential bugs\n"
	text += "\\item $\\mathcal{D}_L$: total number of detection of unique leaks\n"
	text += "\\item $\\mathcal{R}_P$: total number of rewrites of  unique potential bugs\n"
	text += "\\item $\\mathcal{R}_L$: total number of rewrites of  unique leaks\n"
	text += "\\item $\\mathcal{P}_P$: total number of successful replays of unique potential bugs\n"
	text += "\\item $\\mathcal{P}_L$: total number of successful replays of unique leaks\n"
	text += "\\item $\\mathcal{D}_Ax$: total number of detection of unique actual detected bug of type x\n"
	text += "\\item $\\mathcal{D}_Px$: total number of detection of unique potential bug of type x\n"
	text += "\\item $\\mathcal{D}_Lx$: total number of detection of unique leak of type x\n"
	text += "\\item $\\mathcal{R}_Px$: total number of rewrites of  unique potential bug of type x\n"
	text += "\\item $\\mathcal{R}_Lx$: total number of rewrites of  unique leak of type x\n"
	text += "\\item $\\mathcal{P}_Px$: total number of successful replays of unique potential bug of type x\n"
	text += "\\item $\\mathcal{P}_Lx$: total number of successful replays of unique leak of type x\n"
	text += "\\item $\\mathcal{T}_0$: runtime without recording/replay\n"
	text += "\\item $\\mathcal{T}_R$: runtime of recording in s\n"
	text += "\\item $\\mathcal{T}_A$: runtime of analysis in s\n"
	text += "\\item $\\mathcal{T}_P$: avg. runtime of replay in s\n"
	text += "\\item $\\Delta_R$: overhead of recording compared to $\\mathcal{T}_0$\n"
	text += "\\item $\\Delta_P$: avg overhead of replay compared to $\\mathcal{T}_0$\n"
	text += "\\end{itemize}\n"

	writeToFile(fileName, text)
}

func writeBoilerplate(part string, fileName string) {
	switch part {
	case "start":
		text := "\\documentclass{article}\n\\usepackage[english]{babel}\n" +
			"\\usepackage[a4paper,top=2cm,bottom=2cm,left=3cm,right=3cm,marginparwidth=1.75cm]{geometry}\n\n" +
			"\\begin{document}"
		writeToFile(fileName, text)
	case "end":
		text := "\\end{document}"
		writeToFile(fileName, text)
	default:
		println("unknown")
	}
}

// dataMin -> A,P,L without detail, with replay, without sizes
// dataMinSize -> A,P,L without detail, with replay, with sizes
// dataSize -> only sizes
// dataActual -> P with detail, with replay, without sizes
// dataPotential -> P with detail, with replay, without sizes
// dataLeak -> L with detail, with replay, without sizes
func getTableTopLine(size string) string {
	tableStarter := "\\begin{table}[ht]\n\\centering\n\\begin{tabular}"
	switch size {
	case "dataMin":
		return tableStarter + "{|l|c|c|c|c|c|c|c|}\n" +
			"\\hline\nname " +
			"& $\\mathcal{D}_A$ & $\\mathcal{D}_P$ & $\\mathcal{D}_L$ " +
			"& $\\mathcal{R}_P$ & $\\mathcal{R}_L$ & $\\mathcal{P}_P$ & $\\mathcal{P}_L$ " +
			"\\\\ \\hline\n"
	case "dataMinSize":
		return tableStarter + "{|l|c|c|c|c|c|c|c|c|c|c|}\n" +
			"\\hline\nname  & $\\mathcal{S}_L$ & $\\mathcal{S}_T$ & $\\mathcal{S}_O$ " +
			"& $\\mathcal{D}_A$ & $\\mathcal{D}_P$ & $\\mathcal{D}_L$ " +
			"& $\\mathcal{R}_P$ & $\\mathcal{R}_L$ & $\\mathcal{P}_P$ & $\\mathcal{P}_L$ " +
			"\\\\ \\hline\n"
	case "dataSize":
		return tableStarter + "{|l|c|c|c|}\n" +
			"\\hline\nname  & $\\mathcal{S}_L$ & $\\mathcal{S}_T$ & $\\mathcal{S}_O$ " +
			"\\\\ \\hline\n"
	case "dataActual":
		return tableStarter + "{|l|c|c|c|c|c|}\n" +
			"\\hline\nname " +
			"& $\\mathcal{D}_{A1}$ & $\\mathcal{D}_{A2}$ & $\\mathcal{D}_{A3}$ " +
			"& $\\mathcal{D}_{A4}$ & $\\mathcal{D}_{A5}$ " +
			"\\\\ \\hline\n"
	case "dataPotential":
		return tableStarter + "{|l|c|c|c|c|c|c|c|c|c|c|c|c|}\n" +
			"\\hline\nname " +
			"& $\\mathcal{D}_{P1}$ & $\\mathcal{D}_{P2}$ & $\\mathcal{D}_{P3}$ " +
			"& $\\mathcal{D}_{P4}$ & $\\mathcal{R}_{P1}$ & $\\mathcal{R}_{P2}$ " +
			"& $\\mathcal{R}_{P3}$ & $\\mathcal{R}_{P4}$ & $\\mathcal{P}_{P1}$" +
			"& $\\mathcal{P}_{P2}$ & $\\mathcal{P}_{P3}$ & $\\mathcal{P}_{P4}$" +
			"\\\\ \\hline\n"
	case "dataLeak":
		return tableStarter + "{|l|c|c|c|c|c|c|c|c|c|c|c|c|c|c|c|c|c|c|c|c|c|c|c|c|c|c|c|c|c|c|}\n" +
			"\\hline\nname " +
			"& $\\mathcal{D}_{L1}$ & $\\mathcal{D}_{L2}$ & $\\mathcal{D}_{L3}$ " +
			"& $\\mathcal{D}_{L4}$ & $\\mathcal{D}_{L5}$ & $\\mathcal{D}_{L6}$ " +
			"& $\\mathcal{D}_{L7}$ & $\\mathcal{D}_{L8}$ & $\\mathcal{D}_{L9}$ " +
			"& $\\mathcal{D}_{L10}$ " +
			"& $\\mathcal{R}_{L1}$ & $\\mathcal{R}_{L2}$ & $\\mathcal{R}_{L3}$ " +
			"& $\\mathcal{R}_{L4}$ & $\\mathcal{R}_{L5}$ & $\\mathcal{R}_{L6}$ " +
			"& $\\mathcal{R}_{L7}$ & $\\mathcal{R}_{L8}$ & $\\mathcal{R}_{L9}$ " +
			"& $\\mathcal{R}_{L10}$ " +
			"& $\\mathcal{P}_{L1}$ & $\\mathcal{P}_{L2}$ & $\\mathcal{P}_{L3}$ " +
			"& $\\mathcal{P}_{L4}$ & $\\mathcal{P}_{L5}$ & $\\mathcal{P}_{L6}$ " +
			"& $\\mathcal{P}_{L7}$ & $\\mathcal{P}_{L8}$ & $\\mathcal{P}_{L9}$ " +
			"& $\\mathcal{P}_{L10}$" +
			"\\\\ \\hline\n"
	case "time":
		return tableStarter + "{|l|c|c|c|c|c|c|}\n" +
			"\\hline\nname " +
			"& $\\mathcal{T}_0 [s]$ & $\\mathcal{T}_R [s]$ & $\\mathcal{T}_A [s]$ " +
			"& $\\mathcal{T}_P [s]$ & $\\Delta_R [\\%] & $\\Delta_P [\\%]" +
			"\\\\ \\hline\n"
	}

	return ""
}

func getAvgTime() string {
	totalTimeRun := 0.
	totalTimeRecord := 0.
	totalTimeAnalysis := 0.
	totalTimeReplay := 0.

	for _, prog := range progs {
		totalTimeRun += prog.timeRun
		totalTimeRecord += prog.timeRecord
		totalTimeAnalysis += prog.timeAnalysis
		totalTimeReplay += prog.timeReplay
	}

	count := float64(len(progs))
	avgTimeRun := totalTimeRun / count
	avgTimeRecord := totalTimeRecord / count
	avgTimeReplay := totalTimeReplay / count

	overheadRecord := (avgTimeRecord - avgTimeRun) / avgTimeRun * 100.
	overheadReplay := (avgTimeReplay - avgTimeRun) / avgTimeRun * 100.

	return fmt.Sprintf("Average & - & - & - & - & %.2f & %.2f \\\\ \\hline\n",
		max(0, overheadRecord),
		max(0, overheadReplay),
	)
}
