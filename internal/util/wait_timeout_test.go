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
	t.Parallel()

	var wg sync.WaitGroup
	var n int32

	wg.Add(1)
	go func() {
		atomic.AddInt32(&n, 1) // pass
		time.Sleep(25 * time.Millisecond)
		atomic.AddInt32(&n, 1) // pass
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		atomic.AddInt32(&n, 1) // pass
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		time.Sleep(350 * time.Millisecond)
		atomic.AddInt32(&n, 1) // available after wg.Wait()
		wg.Done()
	}()

	// timeout reached.
	err := util.WaitTimeout(&wg, 100*time.Millisecond)
	assert.Equal(t, util.ErrWaitTimeout, err)
	assert.Equal(t, int32(3), atomic.LoadInt32(&n))

	// ok (after timeout).
	err = util.WaitTimeout(&wg, 1*time.Second)
	assert.Equal(t, nil, err)
	assert.Equal(t, int32(4), atomic.LoadInt32(&n))
}
