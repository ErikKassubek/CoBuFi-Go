// ADVOCATE-FILE_START

package runtime

var AdvocateRoutines map[uint64]*AdvocateRoutine
var AdvocateRoutinesLock = mutex{}

var projectPath string

/*
 * AdvocateRoutine is a struct to store the trace of a routine
 * id: the id of the routine
 * G: the g struct of the routine
 * Trace: the trace of the routine
 * ignoreDead: if true, ignore for checkdead
 */
type AdvocateRoutine struct {
	id          uint64
	maxObjectId uint64
	G           *g
	Trace       []string
	ignoreDead  bool
	file string
	line int32
	// lock    *mutex
}

/*
 * Create a new advocate routine
 * Params:
 * 	g: the g struct of the routine
 * 	ignoreDead: if true, the routine will be ignored fore checkdead
 * 	file: file where the routine was created
 * 	line: line where the routine was created
 * Return:
 * 	the new advocate routine
 */
func newAdvocateRoutine(g *g, ignoreDead bool, file string, line int32) *AdvocateRoutine {
	routine := &AdvocateRoutine{id: GetAdvocateRoutineID(), maxObjectId: 0,
		G:          g,
		Trace:      make([]string, 0),
		ignoreDead: ignoreDead,
		file: file,
		line: line,
	}

	lock(&AdvocateRoutinesLock)
	defer unlock(&AdvocateRoutinesLock)

	if AdvocateRoutines == nil {
		AdvocateRoutines = make(map[uint64]*AdvocateRoutine)
	}

	AdvocateRoutines[routine.id] = routine

	return routine
}

/*
 * Add an element to the trace of the current routine
 * Params:
 * 	elem: the element to add
 * Return:
 * 	the index of the element in the trace
 */
func (gi *AdvocateRoutine) addToTrace(elem string) int {
	// do nothing if tracer disabled
	if advocateTracingDisabled {
		return -1
	}

	// do nothing while trace writing disabled
	// this is used to avoid writing to the trace, while the trace is written
	// to the file in case of a too high memory usage
	// for advocateTraceWritingDisabled {
	// 	slowExecution()
	// }

	// never needed in actual code, without it the compiler tests fail
	if gi == nil {
		return -1
	}

	gi.Trace = append(gi.Trace, elem)
	return len(gi.Trace) - 1
}

func (gi *AdvocateRoutine) getElement(index int) string {
	return gi.Trace[index]
}

/*
 * Update an element in the trace of the current routine
 * Params:
 * 	index: the index of the element to update
 * 	elem: the new element
 */
func (gi *AdvocateRoutine) updateElement(index int, elem string) {
	if advocateTracingDisabled {
		return
	}

	if gi == nil {
		return
	}

	if gi.Trace == nil {
		panic("Tried to update element in nil trace")
	}

	if index >= len(gi.Trace) {
		panic("Tried to update element out of bounds")
	}

	gi.Trace[index] = elem
}

/*
 * Get the current routine
 * Return:
 * 	the current routine
 */
func currentGoRoutine() *AdvocateRoutine {
	return getg().goInfo
}

/*
 * GetRoutineID gets the id of the current routine
 * Return:
 * 	id of the current routine, 0 if current routine is nil
 */
func GetRoutineID() uint64 {
	if currentGoRoutine() == nil {
		return 0
	}
	return currentGoRoutine().id
}

// ADVOCATE-FILE-END
