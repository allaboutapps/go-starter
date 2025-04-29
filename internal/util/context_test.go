package util_test

import (
	"context"
	"testing"
	"time"

	"allaboutapps.dev/aw/go-starter/internal/util"
	"github.com/stretchr/testify/assert"
)

type contextKey string

func TestDetachContextWithCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())

	var key contextKey = "test"
	val := 42
	ctx2 := context.WithValue(ctx, key, val)
	detachedContext := util.DetachContext(ctx2)

	cancel()

	select {
	case <-ctx.Done():
		t.Log("Context cancelled")
	default:
		t.Error("Context is not canceled")
	}

	select {
	case <-ctx2.Done():
		t.Log("Context with value cancelled")
	default:
		t.Error("Context with value is not canceled")
	}

	select {
	case <-detachedContext.Done():
		t.Error("Detached context is cancelled")
	default:
		t.Log("Detached context is not cancelled")
	}

	res := detachedContext.Value(key).(int)
	assert.Equal(t, val, res)
}

func TestDetachContextWithDeadline(t *testing.T) {
	ctx, cancel := context.WithTimeout(t.Context(), time.Second*1)
	defer cancel()

	var key contextKey = "test"
	val := 42
	ctx2 := context.WithValue(ctx, key, val)
	detachedContext := util.DetachContext(ctx2)

	time.Sleep(time.Second * 2)

	select {
	case <-ctx.Done():
		t.Log("Context cancelled")
	default:
		t.Error("Context is not canceled")
	}

	select {
	case <-ctx2.Done():
		t.Log("Context with value cancelled")
	default:
		t.Error("Context with value is not canceled")
	}

	select {
	case <-detachedContext.Done():
		t.Error("Detached context is cancelled")
	default:
		t.Log("Detached context is not cancelled")
	}

	res := detachedContext.Value(key).(int)
	assert.Equal(t, val, res)
}
