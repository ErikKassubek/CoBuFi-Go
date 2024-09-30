package main

import (
	"fmt"
	"strings"
)

var bugCodes = []string{
	"A", "P", "L",
	"A01", "A02", "A03", "A04", "A05",
	"P01", "P02", "P03", "P04",
	"L00", "L01", "L02", "L03", "L04", "L05", "L06", "L07", "L08", "L09", "L10",
}

func createCsv(fileName string) {
	table := getCsvTopLine()

	timeRun := 0.
	timeRecord := 0.
	timeAnalysis := 0.
	timeReplay := 0.

	for _, prog := range progs {
		line, tInfo := getCsvLine(prog)
		table += line
		timeRun += tInfo.timeRun
		timeRecord += tInfo.timeRecord
		timeAnalysis += tInfo.timeAnalysis
		timeReplay += tInfo.timeReplay
	}

	table += getCsvAvgTime(timeRun, timeRecord, timeAnalysis, timeReplay, float64(len(progs)))

	writeToFile(fileName, table)
}

func getCsvAvgTime(timeRun, timeRecord, timeAnalysis, timeReplay, numProgs float64) string {
	timeRunAvg := timeRun / numProgs
	timeRecordingAvg := timeRecord / numProgs
	timeAnalysisAvg := timeAnalysis / numProgs
	timeReplayAvg := timeReplay / numProgs

	overheadRecording := -1.
	overheadReplay := -1.
	rationTimeRunAnalysis := -1.

	if timeRunAvg != 0 {
		overheadRecording = (timeRecordingAvg - timeRunAvg) / timeRunAvg * 100.
	}

	if timeRunAvg != 0 {
		overheadReplay = (timeReplayAvg - timeRunAvg) / timeRunAvg * 100.
	}

	if timeRunAvg != 0 {
		rationTimeRunAnalysis = timeAnalysisAvg / timeRunAvg * 100
	}

	line := "Avg.,"
	line += "-,-,-,-,-,-,-,-,-,-,-,-,-,-,-,-,-,-,-,-,-,-,-,-,-,-,-,-,-,-"
	line += "-,-,-,-,-,-,-,-,-,-,-,-,-,-,-,-,-,-,-,-,-,-,-,"
	line += "-,-,-,-,-,-,-,-,-,-,-,-,-,-,-,-,-,-,-,-,-,-,-,-,-,-,"
	line += "-,-,-,-,-,-,-,-,-,-,-,-,-,-,-,-,-,-,-,-"
	line += fmt.Sprintf("%f,%f,%f,%f,%f,%f,%f", timeRunAvg, timeRecordingAvg, timeAnalysisAvg, timeReplayAvg, rationTimeRunAnalysis, overheadRecording, overheadReplay)
	line += "\n"

	return line
}

func getCsvTopLine() string {

	line := "name,"
	line += "numberTests,numberFiles,numberLines,numberNonEmptyLines,"
	line += "numberTraces,numberRoutines,numberNonEmptyRoutines,"
	line += "numberAtomics,numberChannels,numberBufferedChannels,numberUnbufferedChannels,numberSelects,numberSelectCases,numberMutexes,numberWaitGroups,numberCondVar,numberOnce,"
	line += "numberOpsTotal,numberOpsSpawn,numberOpsRoutineTerm,numberOpsAtomic,numberOpsChan,numberOpsChanBuf,numberOpsChanUnbuf,numberOpsSelectCase,numberOpsSelectDefault,numberOpsMutex,numberOpsWait,numberOpsCondVar,numberOpsWait,"
	line += "numberDetectedA,numberDetectedP,numberDetectedL,"
	line += "numberDetectedA01,numberDetectedA02,numberDetectedA03,numberDetectedA04,numberDetectedA05,"
	line += "numberDetectedP01,numberDetectedP02,numberDetectedP03,numberDetectedP04,"
	line += "numberDetectedL00,numberDetectedL01,numberDetectedL02,numberDetectedL03,numberDetectedL04,numberDetectedL05,numberDetectedL06,numberDetectedL07,numberDetectedL08,numberDetectedL09,numberDetectedL10,"
	line += "numberRewrittenA,numberRewrittenP,numberRewrittenL,"
	line += "numberRewrittenA01,numberRewrittenA02,numberRewrittenA03,numberRewrittenA04,numberRewrittenA05,"
	line += "numberRewrittenP01,numberRewrittenP02,numberRewrittenP03,numberRewrittenP04,"
	line += "numberRewrittenL00,numberRewrittenL01,numberRewrittenL02,numberRewrittenL03,numberRewrittenL04,numberRewrittenL05,numberRewrittenL06,numberRewrittenL07,numberRewrittenL08,numberRewrittenL09,numberRewrittenL10,"
	line += "numberReplaySucA,numberReplaySucP,numberReplaySucL,"
	line += "numberReplaySucA01,numberReplaySucA02,numberReplaySucA03,numberReplaySucA04,numberReplaySucA05,"
	line += "numberReplaySucP01,numberReplaySucP02,numberReplaySucP03,numberReplaySucP04,"
	line += "numberReplaySucL00,numberReplaySucL01,numberReplaySucL02,numberReplaySucL03,numberReplaySucL04,numberReplaySucL05,numberReplaySucL06,numberReplaySucL07,numberReplaySucL08,numberReplaySucL09,numberReplaySucL10,"
	line += "timeRun,timeRecording,timeAnalysis,avgTimeReplay,rationTimeRunAnalysis,overheadRecording,overheadReplay"
	line += "\n"
	return line
}

type timeInfo struct {
	timeRun      float64
	timeRecord   float64
	timeAnalysis float64
	timeReplay   float64
}

func getCsvLine(data progData) (string, timeInfo) {
	values := []string{
		data.name,
		data.numberTests, data.numberFiles, data.numberLines, data.numberNonEmptyLines,
		data.numberTraces, data.numberRoutines, data.numberNonEmptyRoutines,
		data.numberAtomics, data.numberChannels, data.numberBuffereChannels, data.numberUnbufferedChannels,
		data.numberSelects, data.numberSelectCases, data.numberMutexes, data.numberWaitGroups, data.numberCondVariables, data.numberOnce,
		data.numberOperations, data.numberSpawnOps, data.numberRoutineTermOps, data.numberAtomicOps, data.numberChannelOps, data.numberBuffereChannelOps,
		data.numberUnbufferedChannelOps, data.numberSelectCaseOps, data.numberSelectDefaultOps, data.numberMutexOps, data.numberWaitOps, data.numberCondVarOps, data.numberOnceOps}

	for _, op := range []map[string]string{data.numberDetected, data.numberRewritten, data.numberReplayed} {
		for _, code := range bugCodes {
			values = append(values, op[code])
		}
	}

	values = append(values, fmt.Sprintf("%f", data.timeRun))
	values = append(values, fmt.Sprintf("%f", data.timeRecord))
	values = append(values, fmt.Sprintf("%f", data.timeAnalysis))
	values = append(values, fmt.Sprintf("%f", data.timeReplay))

	overheadRecord := -1.
	overheadReplay := -1.
	rationTimeRunAnalysis := -1.

	if data.timeRun != 0 {
		overheadRecord = (data.timeRecord - data.timeRun) / data.timeRun * 100.
	}

	if data.timeRun != 0 {
		overheadReplay = (data.timeReplay - data.timeRun) / data.timeRun * 100.
	}

	if data.timeRun != 0 {
		rationTimeRunAnalysis = data.timeAnalysis / data.timeRun * 100
	}

	values = append(values, fmt.Sprintf("%.4f", rationTimeRunAnalysis))
	values = append(values, fmt.Sprintf("%.4f", overheadRecord))
	values = append(values, fmt.Sprintf("%.4f", overheadReplay))

	res := strings.Join(values, ",")
	res += "\n"

	tInfo := timeInfo{
		timeRun:      data.timeRun,
		timeRecord:   data.timeRecord,
		timeAnalysis: data.timeAnalysis,
		timeReplay:   data.timeReplay,
	}

	return res, tInfo
}

func getCsvAvg() {}
