package main

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"
)

// https://www.uber.com/en-DE/blog/leakprof-featherlight-in-production-goroutine-leak-detection/

func TestMain(t *testing.T) {
	unbalancedConditionalCommunication()
	timeOutBug()
	communicationContention()
	channelIterationMisuse()

	time.Sleep(0.5 * time.Second)
}

// Premature Function Return

/*

1. Goroutine leak due to premature return.
2. We can identify the leak but there is no potential partner visible from the trace.
3. This case requires a more detailed source code inspection.

*/

func unbalancedConditionalCommunication() {
	c := make(chan int)

	go func() {
		err := rand.Intn(2) == 1
		if err {
			c <- 1
		} else {
			c <- 0
		}

	}()

	return // early

	<-c

}

// The Timeout Leak
/*

1. Goroutine leak if timeout comes earlier.
2. We can identify the leak.
   PRE done?
3. There's a potential partner.
   PRE done! as part of the select case.
4. Via replay we can show that PRE done gets unstuck.
5. This is useful information to lead to the fix
   where we use buffered (done) channel.

*/

func timeOutBug() {

	done := make(chan int)

	go func() {
		time.Sleep(2 * time.Second)
		done <- 1
		fmt.Printf("done")
	}()

	select {
	case <-done:
	case <-time.After(1 * time.Second):

	}

}

// The NCast Leak

/*

1. Multiple concurrent sender but too few receivers.
2. We can identify the leak(s) (there are several blocked sender threads).
3. There's a single potential partner, the POST receive on data.
4. This should give a strong indication for the source of the problem
 (message contention)
   Fix is to have multiple receivers, increase buffer space of data channel, ...

*/

func communicationContention() {
	n := 5
	data := make(chan int)
	for i := 0; i < n; i++ {
		go func() {
			data <- i
		}()

	}

	<-data

}

// Channel Iteration Misuse

/*

1. Once no messages are sent the consumer will block.
2. There's a single blocked receiver and all potential senders are not stuck.
3. To resolve the leak, the solution is to close the channel.


*/

func channelIterationMisuse() {
	wg := sync.WaitGroup{}
	queueJobs := make(chan int, 1)

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			queueJobs <- 1
		}()
	}

	// Consumer
	go func() {
		for e := range queueJobs {
			fmt.Printf("%d", e)
			wg.Done()
		}

	}()

	wg.Wait()
	// close(queueJobs)
}
