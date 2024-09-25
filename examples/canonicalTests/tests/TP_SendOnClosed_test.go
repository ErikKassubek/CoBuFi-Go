package main

import (
	"advocate"
	"testing"
	"time"
)

func Test11(t *testing.T) {
	// ======= Preamble Start =======
  advocate.EnableReplay(2, true)
  defer advocate.WaitForReplayFinish()
  // ======= Preamble End =======
	n11()
}
// TP
func n11() {
	c := make(chan struct{}, 0)

	go func() {
		time.Sleep(200 * time.Millisecond) // prevent actual send on closed channel
		close(c)
	}()

	go func() {
		select {
		case c <- struct{}{}:
		default:
		}
	}()

	go func() {
		time.Sleep(100 * time.Millisecond)
		select {
		case <-c:
		default:
		}
	}()

	time.Sleep(300 * time.Millisecond) // make sure, that the default values are taken
}
