package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type progData struct {
	name string

	numberTests         string
	numberFiles         string
	numberLines         string
	numberNonEmptyLines string

	numberTraces           string
	numberRoutines         string
	numberNonEmptyRoutines string

	numberAtomics            string
	numberChannels           string
	numberBuffereChannels    string
	numberUnbufferedChannels string
	numberSelects            string
	numberSelectCases        string
	numberMutexes            string
	numberWaitGroups         string
	numberCondVariables      string
	numberOnce               string

	numberOperations           string
	numberSpawnOps             string
	numberAtomicOps            string
	numberChannelOps           string
	numberBuffereChannelOps    string
	numberUnbufferedChannelOps string
	numberSelectCaseOps        string
	numberSelectDefaultOps     string
	numberMutexOps             string
	numberWaitOps              string
	numberCondVarOps           string
	numberOnceOps              string

	numberDetected  map[string]string
	numberRewritten map[string]string
	numberReplayed  map[string]string

	timeRun      string
	timeRecord   string
	timeAnalysis string
	timeReplay   string
}

var progs = make(map[string]progData)

func main() {
	statsPath := flag.String("f", "", "Path to the stat and time files")
	flag.Parse()

	if *statsPath == "" {
		fmt.Println("Please set the path to the folder containing the stats and time files")
		return
	}

	err := filepath.Walk(*statsPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			if strings.Contains(info.Name(), "stats_") {
				err := readStats(path)
				if err != nil {
					fmt.Println("Failed to read stats for ", info.Name(), err)
				}
			} else if strings.Contains(info.Name(), "times_") {
				err := readTime(path)
				if err != nil {
					fmt.Println("Failed to read times for ", info.Name(), err)
				}
			}
		}

		return nil
	})

	if err != nil {
		log.Fatalf("Error walking the path %v: %v\n", *statsPath, err)
	}

	fileName := "tables_" + time.Now().Format("2006-01-02_15:04:05") + ".tex"

	writeBoilerplate("start", fileName)

	createLatexTable("dataMin", fileName)
	createLatexTable("dataMinSize", fileName)
	createLatexTable("dataSize", fileName)
	createLatexTable("dataActual", fileName)
	createLatexTable("dataPotential", fileName)
	createLatexTable("dataLeak", fileName)

	writeExplanation(fileName)

	writeBoilerplate("end", fileName)
}

func readStats(path string) error {
	name := getProgNameFromFile(path)

	data, err := readFile(path)
	if err != nil {
		return err
	}

	dataSplit := strings.Split(data, "\n")

	if len(dataSplit) < 7 {
		return fmt.Errorf("The stats file expected at least 7 lines but got %d", len(dataSplit))
	}

	infoProg := strings.Split(dataSplit[0], ",")
	infoTrace := strings.Split(dataSplit[1], ",")
	infoObjects := strings.Split(dataSplit[2], ",")
	infoOperations := strings.Split(dataSplit[3], ",")
	infoDetected := strings.Split(dataSplit[4], ",")
	infoRewritten := strings.Split(dataSplit[5], ",")
	infoReplay := strings.Split(dataSplit[6], ",")

	if val, ok := progs[name]; ok {
		val.numberFiles = infoProg[0]
		val.numberLines = infoProg[1]
		val.numberNonEmptyLines = infoProg[2]

		val.numberTraces = infoTrace[0]
		val.numberRoutines = infoTrace[1]
		val.numberNonEmptyRoutines = infoTrace[1]

		val.numberAtomics = infoObjects[0]
		val.numberChannels = infoObjects[1]
		val.numberBuffereChannels = infoObjects[2]
		val.numberUnbufferedChannels = infoObjects[3]
		val.numberSelects = infoObjects[4]
		val.numberSelectCases = infoObjects[5]
		val.numberMutexes = infoObjects[6]
		val.numberWaitGroups = infoObjects[7]
		val.numberCondVariables = infoObjects[8]
		val.numberOnce = infoObjects[9]

		val.numberOperations = infoOperations[0]
		val.numberSpawnOps = infoOperations[1]
		val.numberAtomicOps = infoOperations[2]
		val.numberChannelOps = infoOperations[3]
		val.numberBuffereChannelOps = infoOperations[4]
		val.numberUnbufferedChannelOps = infoOperations[5]
		val.numberSelectCaseOps = infoOperations[6]
		val.numberSelectDefaultOps = infoOperations[7]
		val.numberMutexOps = infoOperations[8]
		val.numberWaitOps = infoOperations[9]
		val.numberCondVarOps = infoOperations[10]
		val.numberOnceOps = infoOperations[11]

		val.numberDetected = mapCodes(infoDetected)
		val.numberRewritten = mapCodes(infoRewritten)
		val.numberReplayed = mapCodes(infoReplay)

		progs[name] = val
	} else {
		progs[name] = progData{
			name: name,

			numberFiles:         infoProg[0],
			numberLines:         infoProg[1],
			numberNonEmptyLines: infoProg[2],

			numberTraces:           infoTrace[0],
			numberRoutines:         infoTrace[1],
			numberNonEmptyRoutines: infoTrace[1],

			numberAtomics:            infoObjects[0],
			numberChannels:           infoObjects[1],
			numberBuffereChannels:    infoObjects[2],
			numberUnbufferedChannels: infoObjects[3],
			numberSelects:            infoObjects[4],
			numberSelectCases:        infoObjects[5],
			numberMutexes:            infoObjects[6],
			numberWaitGroups:         infoObjects[7],
			numberCondVariables:      infoObjects[8],
			numberOnce:               infoObjects[9],

			numberOperations:           infoOperations[0],
			numberSpawnOps:             infoOperations[1],
			numberAtomicOps:            infoOperations[2],
			numberChannelOps:           infoOperations[3],
			numberBuffereChannelOps:    infoOperations[4],
			numberUnbufferedChannelOps: infoOperations[5],
			numberSelectCaseOps:        infoOperations[6],
			numberSelectDefaultOps:     infoOperations[7],
			numberMutexOps:             infoOperations[8],
			numberWaitOps:              infoOperations[9],
			numberCondVarOps:           infoOperations[10],
			numberOnceOps:              infoOperations[11],

			numberDetected:  mapCodes(infoDetected),
			numberRewritten: mapCodes(infoRewritten),
			numberReplayed:  mapCodes(infoReplay),
		}
	}

	return nil
}

func readTime(path string) error {
	name := getProgNameFromFile(path)

	data, err := readFile(path)
	if err != nil {
		return err
	}

	dataSplit := strings.Split(data, "\n")

	if len(dataSplit) < 4 {
		return fmt.Errorf("The time file expected at least 4 lines but got %d", len(dataSplit))
	}

	timeRun := strings.Split(dataSplit[0], ": ")[1]
	timeRecord := strings.Split(dataSplit[1], ": ")[1]
	timeAnalysis := strings.Split(dataSplit[2], ": ")[1]
	timeReplay := strings.Split(dataSplit[3], ": ")[1]

	if val, ok := progs[name]; ok {
		val.timeRun = timeRun
		val.timeRecord = timeRecord
		val.timeAnalysis = timeAnalysis
		val.timeReplay = timeReplay
		progs[name] = val
	} else {
		progs[name] = progData{
			name:         name,
			timeRun:      timeRun,
			timeRecord:   timeRecord,
			timeAnalysis: timeAnalysis,
			timeReplay:   timeReplay,
		}
	}

	return nil

}

func mapCodes(data []string) map[string]string {
	return map[string]string{
		"A":   data[0],
		"P":   data[1],
		"L":   data[2],
		"A01": data[3],
		"A02": data[4],
		"A03": data[5],
		"A04": data[6],
		"A05": data[7],
		"P01": data[8],
		"P02": data[9],
		"P03": data[10],
		"L01": data[11],
		"L02": data[12],
		"L03": data[13],
		"L04": data[14],
		"L05": data[15],
		"L06": data[16],
		"L07": data[17],
		"L08": data[18],
		"L09": data[19],
		"L10": data[20],
	}
}

func readFile(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

func getProgNameFromFile(path string) string {
	path = strings.Replace(path, "stats_", "", -1)
	path = strings.Replace(path, "times_", "", -1)
	path = strings.Replace(path, ".log", "", -1)
	return filepath.Base(path)
}

func getTableRows1(data progData, size string) string {
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
			gv(data.numberTraces),
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
			gv(data.numberTraces),
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
	}
	return ""
}

func gv(val string) string {
	if val == "" {
		return "0"
	}
	return val
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
	}

	return ""
}

func createLatexTable(name string, fileName string) {
	table := getTableTopLine(name)

	for _, prog := range progs {
		table += getTableRows1(prog, name)
	}

	table += "\\end{tabular}\n\\caption{"
	table += name
	table += "}\n\\label{Tab:"
	table += name
	table += "}\n\\end{table}"

	writeToFile(fileName, table)
}

func writeToFile(fileName, content string) {
	// Open the file in append mode, or create it if it doesn't exist
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()
	_, err = file.WriteString(content + "\n\n\n")
	if err != nil {
		fmt.Println("Error writing to file:", err)
	}
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

// S_T: number traces
// S_L: size in lines
// S_O: trace size (number ops)

// D_A: detected Actual
// D_P: detected Potentail
// D_L: detected Leal
// R_P: rewritten potential
// R_L: rewritten leak
// P_P: replayed Potentail
// P_L: replatyed leak

// T_0: time run
// T_R: time recording
// T_A: time analysis
// T_P: time replay

// "$\\mathcal{T}_0$: base runtime [s], $\\mathcal{T}_R$: time for recording [s], " +
// 		"$\\mathcal{T}_A$: time for analysis [s], $\\mathcal{T}_P$: time for replay [s]}
