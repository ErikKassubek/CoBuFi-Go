package main

import (
	"sync"
	"testing"
)

func TestReadWriteBase(t *testing.T) {
	var x, y sync.RWMutex

	go func() {
		x.Lock()
		y.Lock()
		y.Unlock()
		x.Unlock()
	}()

	go func() {
		y.Lock()
		x.Lock()
		x.Unlock()
		y.Unlock()
	}()
}
