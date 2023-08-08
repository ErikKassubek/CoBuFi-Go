package reader

import (
	"analyzer/trace"
	"bufio"
	"os"
	"strings"
)

/*
 * Read a file
 */
func Read(file_path string) {
	file, err := os.Open(file_path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	routine := 0
	for scanner.Scan() {
		routine++
		line := scanner.Text()
		processLine(line, routine)
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}
}

/*
 * Process one line from the log file.
 * Args:
 *   line (string): The line to process
 *   routine (int): The routine id, equal to the line number
 */
func processLine(line string, routine int) {
	// create element for routine in trace
	trace.NewRoutine(routine)

	elements := strings.Split(line, ";")
	for _, element := range elements {
		processElement(element, routine)
	}
}

/*
 * Process one element from the log file.
 * Args:
 *   element (string): The element to process
 *   routine (int): The routine id, equal to the line number
 */
func processElement(element string, routine int) {
	fields := strings.Split(element, ",")
	var err error = nil
	switch fields[0] {
	case "A":
		err = trace.AddTraceElementAtomic(routine, fields[1], fields[2], fields[3])
	case "C":
		err = trace.AddTraceElementChannel(routine, fields[1], fields[2],
			fields[3], fields[4], fields[5], fields[6], fields[7], fields[8],
			fields[9], fields[10])
	case "M":
		err = trace.AddTraceElementMutex(routine, fields[1], fields[2],
			fields[3], fields[4], fields[5], fields[6], fields[7], fields[8])
	case "G":
		err = trace.AddTraceElementRoutine(routine, fields[1], fields[2])
	case "S":
		trace.AddTraceElementSelect(routine, fields[1], fields[2], fields[3],
			fields[4], fields[5], fields[6], fields[7], fields[8])
	case "W":
		trace.AddTraceElementWait(routine, fields[1], fields[2], fields[3],
			fields[4], fields[5], fields[6], fields[7], fields[8])
	default:
		panic("Unknown element type: " + fields[0])
	}

	if err != nil {
		panic(err)
	}

}
