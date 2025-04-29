package common

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"testing"
	"time"

	"allaboutapps.dev/aw/go-starter/internal/util"
	"github.com/stretchr/testify/assert"
)

func TestEnsureDeadline(t *testing.T) {
	deadline := time.Now().Add(1 * time.Second)

	ctx, cancel := context.WithDeadline(t.Context(), deadline)
	defer cancel()

	receivedDeadline := ensureProbeDeadlineFromContext(ctx)
	assert.Equal(t, deadline, receivedDeadline)
}

func TestDummyDeadlineWithinOneSec(t *testing.T) {
	ctx := t.Context()

	receivedDeadline := ensureProbeDeadlineFromContext(ctx)
	assert.WithinDuration(t, time.Now().Add(1*time.Second), receivedDeadline, 100*time.Millisecond)

}

func TestProbeDatabasePingableDeadline(t *testing.T) {
	ctx, cancel := context.WithDeadline(t.Context(), time.Now())
	defer cancel()

	_, err := probeDatabasePingable(ctx, &sql.DB{})
	assert.Truef(t, errors.Is(err, util.ErrWaitTimeout) || errors.Is(err, context.DeadlineExceeded), "err must be util.ErrWaitTimeout or context.DeadlineExceeded but is %v", err)
}

func TestProbeDatabaseNextHealthSequenceDeadline(t *testing.T) {
	ctx, cancel := context.WithDeadline(t.Context(), time.Now())
	defer cancel()

	_, err := probeDatabaseNextHealthSequence(ctx, &sql.DB{})
	assert.Truef(t, errors.Is(err, util.ErrWaitTimeout) || errors.Is(err, context.DeadlineExceeded), "err must be util.ErrWaitTimeout or context.DeadlineExceeded but is %v", err)
}

func TestProbePathWriteablePermissionContextDeadline(t *testing.T) {
	ctx, cancel := context.WithDeadline(t.Context(), time.Now())
	defer cancel()

	_, err := probePathWriteablePermission(ctx, "/any/thing")
	assert.Truef(t, errors.Is(err, util.ErrWaitTimeout) || errors.Is(err, context.DeadlineExceeded), "err must be util.ErrWaitTimeout or context.DeadlineExceeded but is %v", err)
}

func TestProbePathWriteableTouchContextDeadline(t *testing.T) {
	ctx, cancel := context.WithDeadline(t.Context(), time.Now())
	defer cancel()

	_, err := probePathWriteableTouch(ctx, "/any/thing", ".touch")
	assert.Truef(t, errors.Is(err, util.ErrWaitTimeout) || errors.Is(err, context.DeadlineExceeded), "err must be util.ErrWaitTimeout or context.DeadlineExceeded but is %v", err)
}

func TestProbePathWriteableTouchInaccessable(t *testing.T) {
	_, err := probePathWriteableTouch(t.Context(), "/this/path/does/not/exist", ".touch")
	assert.True(t, os.IsNotExist(err))
}
