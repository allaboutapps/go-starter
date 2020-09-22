package util_test

import (
	"os"
	"testing"

	"allaboutapps.dev/aw/go-starter/internal/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTouchfile(t *testing.T) {
	t.Parallel()

	err := os.Remove("/tmp/.touchfile-test")

	if err != nil {
		require.Equalf(t, true, os.IsNotExist(err), "Expected err to be of tpye os.IsNotExist(err), but is: %v", err)
	}

	ts1, err := util.TouchFile("/tmp/.touchfile-test")
	assert.NoError(t, err)

	ts2, err := util.TouchFile("/tmp/.touchfile-test")
	assert.NoError(t, err)
	require.NotEqual(t, ts1.UnixNano(), ts2.UnixNano())

	zeroTime, err := util.TouchFile("/this/path/does/not/exist/.touchfile-test")
	assert.Error(t, err)
	assert.True(t, zeroTime.IsZero(), "time.Time on error should be zero time")
}
