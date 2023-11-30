package trace

import (
	"errors"
	"strconv"
)

// enum for operation
type opAtomic int

const (
	LoadOp opAtomic = iota
	StoreOp
	AddOp
	SwapOp
	CompSwapOp
)

/*
 * Struct to save an atomic event in the trace
 * Fields:
 *   routine (int): The routine id
 *   tpost (int): The timestamp of the event
 *   id (int): The id of the atomic variable
 *   operation (int, enum): The operation on the atomic variable
 */
type traceElementAtomic struct {
	routine int
	tpost   int
	id      int
	opA     opAtomic
}

/*
 * Create a new atomic trace element
 * Args:
 *   routine (int): The routine id
 *   tpost (string): The timestamp of the event
 *   id (string): The id of the atomic variable
 *   operation (string): The operation on the atomic variable
 */
func AddTraceElementAtomic(routine int, tpost string,
	id string, operation string) error {
	tPostInt, err := strconv.Atoi(tpost)
	if err != nil {
		return errors.New("tpost is not an integer")
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		return errors.New("tpost is not an integer")
	}

	var opAInt opAtomic
	switch operation {
	case "L":
		opAInt = LoadOp
	case "S":
		opAInt = StoreOp
	case "A":
		opAInt = AddOp
	case "W":
		opAInt = SwapOp
	case "C":
		opAInt = CompSwapOp
	default:
		return errors.New("operation is not a valid operation")
	}

	elem := traceElementAtomic{
		routine: routine,
		tpost:   tPostInt,
		id:      idInt,
		opA:     opAInt,
	}

	return addElementToTrace(&elem)
}

/*
 * Get the routine of the element
 * Returns:
 *   int: The routine of the element
 */
func (at *traceElementAtomic) getRoutine() int {
	return at.routine
}

/*
 * Get the tpost of the element. For atomic elements, tpre and tpost are the same
 * Returns:
 *   int: The tpost of the element
 */
func (at *traceElementAtomic) getTpre() int {
	return at.tpost
}

/*
 * Get the tpost of the element. For atomic elements, tpre and tpost are the same
 * Returns:
 *   int: The tpost of the element
 */
func (at *traceElementAtomic) getTpost() int {
	return at.tpost
}

/*
 * Get the timer, that is used for the sorting of the trace
 * Returns:
 *   int: The timer of the element
 */
func (at *traceElementAtomic) getTsort() int {
	return at.tpost
}

/*
 * Get the simple string representation of the element.
 * Returns:
 *   string: The simple string representation of the element
 */
func (at *traceElementAtomic) toString() string {
	return "A," + strconv.Itoa(at.tpost) + "," + strconv.Itoa(at.id) + "," +
		strconv.Itoa(int(at.opA))
}
