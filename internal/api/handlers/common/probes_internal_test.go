package common

import (
	"context"
	"database/sql"
	"os"
	"testing"
	"time"

	"allaboutapps.dev/aw/go-starter/internal/util"
	"github.com/stretchr/testify/assert"
)

func TestEnsureDeadline(t *testing.T) {
	t.Parallel()

	deadline := time.Now().Add(1 * time.Second)

	ctx, cancel := context.WithDeadline(context.Background(), deadline)
	defer cancel()

	receivedDeadline := ensureProbeDeadlineFromContext(ctx)
	assert.Equal(t, deadline, receivedDeadline)
}

func TestDummyDeadlineWithinOneSec(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	receivedDeadline := ensureProbeDeadlineFromContext(ctx)
	assert.WithinDuration(t, time.Now().Add(1*time.Second), receivedDeadline, 100*time.Millisecond)

}

func TestProbeDatabasePingableDeadline(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithDeadline(context.Background(), time.Now())
	defer cancel()

	_, err := probeDatabasePingable(ctx, &sql.DB{})
	assert.Equal(t, util.ErrWaitTimeout, err)
}

func TestProbeDatabaseNextHealthSequenceDeadline(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithDeadline(context.Background(), time.Now())
	defer cancel()

	_, err := probeDatabaseNextHealthSequence(ctx, &sql.DB{})
	assert.Equal(t, util.ErrWaitTimeout, err)
}

func TestProbePathWriteablePermissionContextDeadline(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithDeadline(context.Background(), time.Now())
	defer cancel()

	_, err := probePathWriteablePermission(ctx, "/any/thing")
	assert.Equal(t, context.DeadlineExceeded, err)
}

func TestProbePathWriteableTouchContextDeadline(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithDeadline(context.Background(), time.Now())
	defer cancel()

	_, err := probePathWriteableTouch(ctx, "/any/thing", ".touch")
	assert.Equal(t, context.DeadlineExceeded, err)
}

func TestProbePathWriteableTouchInaccessable(t *testing.T) {
	t.Parallel()

	_, err := probePathWriteableTouch(context.Background(), "/this/path/does/not/exist", ".touch")
	assert.True(t, os.IsNotExist(err))
}
