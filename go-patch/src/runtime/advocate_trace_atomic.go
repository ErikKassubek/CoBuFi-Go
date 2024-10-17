package runtime

type AtomicOp string

const (
	LoadOp     AtomicOp = "L"
	StoreOp    AtomicOp = "S"
	AddOp      AtomicOp = "A"
	SwapOp     AtomicOp = "W"
	CompSwapOp AtomicOp = "C"
)

/*
 * Add an atomic operation to the trace
 * Args:
 * 	index: index of the atomic event in advocateAtomicMap
 */
func AdvocateAtomic[T any](addr *T, op AtomicOp) {
	timer := GetNextTimeStep()

	_, file, line, _ := Caller(2)

	index := pointerAddressAsString(addr, true)

	elem := "A," + uint64ToString(timer) + "," + index + "," + string(op) + "," + file + ":" + intToString(line)
	insertIntoTrace(elem)
}
