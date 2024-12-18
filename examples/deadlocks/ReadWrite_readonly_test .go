package main

import (
	"sync"
	"testing"
)

func TestReadWriteReadonly(t *testing.T) {
	var x, y sync.RWMutex

	go func() {
		x.RLock()
		y.RLock()
		y.RUnlock()
		x.RUnlock()
	}()

	go func() {
		y.RLock()
		x.RLock()
		x.RUnlock()
		y.RUnlock()
	}()
}
