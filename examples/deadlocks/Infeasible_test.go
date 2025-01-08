package main

import (
	"sync"
	"testing"
)

func TestInfeasible(t *testing.T) {
	var x, y sync.Mutex

	x.Lock()
	y.Lock()
	y.Unlock()
	x.Unlock()

	go func() {
		y.Lock()
		x.Lock()
		x.Unlock()
		y.Unlock()
	}()
}
