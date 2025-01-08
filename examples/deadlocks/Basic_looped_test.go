package main

import (
	"sync"
	"testing"
)

func TestBasicLooped(t *testing.T) {
	var x, y sync.Mutex
	repeat := 10

	go func() {
		x.Lock()
		y.Lock()
		y.Unlock()
		x.Unlock()
	}()

	go func() {
		for i := 0; i < repeat; i++ {
			y.Lock()
			x.Lock()
			x.Unlock()
			y.Unlock()
		}
	}()
}
