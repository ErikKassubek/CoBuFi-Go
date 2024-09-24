package main

import (
	"sync"
	"testing"
)

func Test57(t *testing.T) {
	n57()
}

// TN: No concurrent recv on same channel
func n57() {
	var m sync.Mutex

	m.Lock()
	m.Unlock()
}
