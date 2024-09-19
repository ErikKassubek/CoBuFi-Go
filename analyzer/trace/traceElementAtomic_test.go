package trace

import (
	"testing"
)

func TestTraceElementAtomicStruct(t *testing.T) {
	AddTraceElementAtomic(1, "2", "3", "L")
	AddTraceElementAtomic(1, "2", "3", "S")
	AddTraceElementAtomic(1, "2", "3", "A")
	AddTraceElementAtomic(1, "2", "3", "W")
	AddTraceElementAtomic(1, "2", "3", "C")

	if len(traces[1]) != 5 {
		t.Errorf("Expected 5 trace elements for routine 1, got %d", len(traces[1]))
	}

	expected := []string{
		"A,2,3,L",
		"A,2,3,S",
		"A,2,3,A",
		"A,2,3,W",
		"A,2,3,C",
	}

	for i, elem := range traces[1] {
		if elem.ToString() != expected[i] {
			t.Errorf("Expected %s, got %s", expected[i], elem.ToString())
		}
	}
}
