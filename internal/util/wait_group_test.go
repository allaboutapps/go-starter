package util_test

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"allaboutapps.dev/aw/go-starter/internal/util"
	"github.com/stretchr/testify/assert"
)

func TestWaitTimeoutErr(t *testing.T) {
	var wg sync.WaitGroup
	var n int32

	wg.Add(1)
	go func() {
		atomic.AddInt32(&n, 1)
		time.Sleep(10 * time.Millisecond)
		atomic.AddInt32(&n, 1)
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		atomic.AddInt32(&n, 1)
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		time.Sleep(400 * time.Millisecond)
		atomic.AddInt32(&n, 1) // available after wg.Wait()
		wg.Done()
	}()

	// timeout reached.
	err := util.WaitTimeout(&wg, 200*time.Millisecond)
	assert.Equal(t, util.ErrWaitTimeout, err)
	assert.Equal(t, int32(3), atomic.LoadInt32(&n))

	// ok (after timeout).
	err = util.WaitTimeout(&wg, 800*time.Second)
	assert.Equal(t, nil, err)
	assert.Equal(t, int32(4), atomic.LoadInt32(&n))
}
