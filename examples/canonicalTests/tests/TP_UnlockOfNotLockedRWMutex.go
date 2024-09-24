package main

import (
	"sync"
	"testing"
	"time"
)

func Test59(t *testing.T) {
	n57()
}

// TN: No concurrent recv on same channel
func n59() {
	var m sync.RWMutex

	go func() {
		m.RLock()
	}()

	go func() {
		m.RLock()
	}()

	go func() {
		time.Sleep(100 * time.Millisecond)
		m.RUnlock()
	}()

	time.Sleep(100 * time.Millisecond)
	m.RUnlock()

	time.Sleep(200 * time.Millisecond)
}
