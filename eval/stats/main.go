package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type progData struct {
	name string

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

	createLatexTable(*statsPath)

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

func getTableLineFromData(data progData) string {
	return fmt.Sprintf("%s & %s & %s & %s & %s & %s & %s & %s & %s & %s & %s & %s & %s & %s \\\\ \\hline\n",
		data.name,
		data.numberLines,
		data.numberOperations,
		data.numberDetected["A"],
		data.numberDetected["P"],
		data.numberDetected["L"],
		data.numberRewritten["P"],
		data.numberDetected["L"],
		data.numberReplayed["P"],
		data.numberDetected["L"],
		data.timeRun,
		data.timeRecord,
		data.timeAnalysis,
		data.timeReplay,
	)
}

func createLatexTable(statsPath string) {
	table := "\\begin{table}[h]\n\\begin{tabular}{|l|c|c|c|c|c|c|c|c|c|c|c|c|c|}\n" +
		"\\hline\nname  & $\\mathcal{L}$ & $\\mathcal{O}$ " +
		"& $\\mathcal{D}_A$ & $\\mathcal{D}_P$ & $\\mathcal{D}_L$ " +
		"& $\\mathcal{R}_P$ & $\\mathcal{D}_L$ & $\\mathcal{P}_P$ & $\\mathcal{P}_L$ " +
		"& $\\mathcal{T}_0$ & $\\mathcal{T}_R$ & $\\mathcal{T}_A$ & $\\mathcal{T}_P$ \\\\ \\hline\n"

	for _, prog := range progs {
		table += getTableLineFromData(prog)
	}

	table += "\\end{tabular}\n"
	table += "\\caption{$\\mathcal{L}$: number of line, " +
		"$\\mathcal{O}$: number of recorded operations, $\\mathcal{D}_A$: number of detected actual bugs, " +
		"$\\mathcal{D}_P$: number of detected potential bugs, $\\mathcal{D}_A$: number of detected leaks, " +
		"$\\mathcal{R}_P$: number of rewritten potential bugs, $\\mathcal{R}_L$: number of rewritten leaks, " +
		"$\\mathcal{P}_P$: number of successfully replayed potential bugs, " +
		"$\\mathcal{R}_L$: number of successfully replayed leaks, " +
		"$\\mathcal{T}_0$: base runtime [s], $\\mathcal{T}_R$: time for recording [s], " +
		"$\\mathcal{T}_A$: time for analysis [s], $\\mathcal{T}_P$: time for replay [s]}\n"
	table += "\\label{}\n"
	table += "\\end{table}"

	file, err := os.Create(filepath.Join(statsPath, "table.tex"))
	if err != nil {
		log.Fatalf("Error creating file: %v", err)
	}
	defer file.Close()

	// Write the string to the file
	_, err = file.WriteString(table)
	if err != nil {
		log.Fatalf("Error writing to file: %v", err)
	}
}

// F: number files
// L: size in lines
// O: trace size (number ops)

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
