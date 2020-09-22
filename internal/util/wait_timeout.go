package util

import (
	"errors"
	"sync"
	"time"
)

var (
	ErrWaitTimeout = errors.New("WaitGroup has timed out")
)

// WaitTimeout waits for the waitgroup for the specified max timeout.
// Returns nil on completion or ErrWaitTimeout if waiting timed out.
// See https://stackoverflow.com/questions/32840687/timeout-for-waitgroup-wait
// Note that the spawned goroutine to wg.Wait() gets leaked and will continue running detached
func WaitTimeout(wg *sync.WaitGroup, timeout time.Duration) error {
	c := make(chan struct{})
	go func() {
		defer close(c)
		wg.Wait()
	}()
	select {
	case <-c:
		return nil // completed normally
	case <-time.After(timeout):
		return ErrWaitTimeout // timed out
	}
}
