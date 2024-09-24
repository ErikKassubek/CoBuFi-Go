package main

import (
	"sync"
	"testing"
	"time"
)

func Test58(t *testing.T) {
	n57()
}

// TN: No concurrent recv on same channel
func n58() {
	var m sync.Mutex

	go func() {
		m.Lock()
	}()

	time.Sleep(100 * time.Millisecond)
	m.Unlock()

	time.Sleep(200 * time.Millisecond)
}
