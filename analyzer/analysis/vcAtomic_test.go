package analysis

import (
	"analyzer/clock"
	"analyzer/trace"
	"testing"
)

func TestTraceElementAtomicUpdateVectorClock(t *testing.T) {
	nrRout := 4

	elemLoad := trace.TraceElementAtomic{
		routine: 1,
		tPost:   25,
		id:      2,
		opA:     LoadOp,
		vc:      clock.NewVectorClock(nrRout),
	}

	elemStore := TraceElementAtomic{
		routine: 1,
		tPost:   25,
		id:      2,
		opA:     StoreOp,
		vc:      clock.NewVectorClock(nrRout),
	}

	elemAdd := TraceElementAtomic{
		routine: 1,
		tPost:   25,
		id:      2,
		opA:     AddOp,
		vc:      clock.NewVectorClock(nrRout),
	}

	elemSwap := TraceElementAtomic{
		routine: 1,
		tPost:   25,
		id:      2,
		opA:     SwapOp,
		vc:      clock.NewVectorClock(nrRout),
	}

	elemComp := TraceElementAtomic{
		routine: 1,
		tPost:   25,
		id:      2,
		opA:     CompSwapOp,
		vc:      clock.NewVectorClock(nrRout),
	}

	currentVCHb[1] = clock.NewVectorClock(nrRout)
	analysis.lw[2] = clock.NewVectorClock(nrRout, map[int]int{1: 0, 2: 2, 3: 4, 4: 0})

	elemLoad.updateVectorClockAlt()
	elemStore.updateVectorClockAlt()
	elemAdd.updateVectorClockAlt()
	elemSwap.updateVectorClockAlt()
	elemComp.updateVectorClockAlt()

	expectVC("Load", elemLoad.GetVC(), clock.CreateVectorClock(4, map[int]int{1: 0, 2: 0, 3: 0, 4: 0}), t)
	expectVC("Store", elemStore.GetVC(), clock.CreateVectorClock(4, map[int]int{1: 1, 2: 0, 3: 0, 4: 0}), t)
	expectVC("Add", elemAdd.GetVC(), clock.CreateVectorClock(4, map[int]int{1: 2, 2: 0, 3: 0, 4: 0}), t)
	expectVC("Swap", elemSwap.GetVC(), clock.CreateVectorClock(4, map[int]int{1: 3, 2: 0, 3: 0, 4: 0}), t)
	expectVC("Comp", elemComp.GetVC(), clock.CreateVectorClock(4, map[int]int{1: 5, 2: 0, 3: 0, 4: 0}), t)
}

func TestTraceElementAtomicUpdateVectorClockAlt(t *testing.T) {
	nrRout := 4

	elemLoad := TraceElementAtomic{
		routine: 1,
		tPost:   25,
		id:      2,
		opA:     LoadOp,
		vc:      clock.NewVectorClock(nrRout),
	}

	elemStore := TraceElementAtomic{
		routine: 1,
		tPost:   25,
		id:      2,
		opA:     StoreOp,
		vc:      clock.NewVectorClock(nrRout),
	}

	elemAdd := TraceElementAtomic{
		routine: 1,
		tPost:   25,
		id:      2,
		opA:     AddOp,
		vc:      clock.NewVectorClock(nrRout),
	}

	elemSwap := TraceElementAtomic{
		routine: 1,
		tPost:   25,
		id:      2,
		opA:     SwapOp,
		vc:      clock.NewVectorClock(nrRout),
	}

	elemComp := TraceElementAtomic{
		routine: 1,
		tPost:   25,
		id:      2,
		opA:     CompSwapOp,
		vc:      clock.NewVectorClock(nrRout),
	}

	currentVCHb[1] = clock.NewVectorClock(nrRout)

	elemLoad.updateVectorClockAlt()
	elemStore.updateVectorClockAlt()
	elemAdd.updateVectorClockAlt()
	elemSwap.updateVectorClockAlt()
	elemComp.updateVectorClockAlt()

	expectVC("Load", elemLoad.GetVC(), clock.CreateVectorClock(4, map[int]int{1: 0, 2: 0, 3: 0, 4: 0}), t)
	expectVC("Store", elemStore.GetVC(), clock.CreateVectorClock(4, map[int]int{1: 1, 2: 0, 3: 0, 4: 0}), t)
	expectVC("Add", elemAdd.GetVC(), clock.CreateVectorClock(4, map[int]int{1: 2, 2: 0, 3: 0, 4: 0}), t)
	expectVC("Swap", elemSwap.GetVC(), clock.CreateVectorClock(4, map[int]int{1: 3, 2: 0, 3: 0, 4: 0}), t)
	expectVC("Comp", elemComp.GetVC(), clock.CreateVectorClock(4, map[int]int{1: 5, 2: 0, 3: 0, 4: 0}), t)
}

func expectVC(name string, got clock.VectorClock, expect clock.VectorClock, t *testing.T) {
	if got.ToString() != expect.ToString() {
		t.Errorf("%s: expected %s, got %s", name, expect.ToString(), got.ToString())
	}
}
