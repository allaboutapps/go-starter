package test_test

import (
	"testing"
	"time"

	"allaboutapps.dev/aw/go-starter/internal/api"
	"allaboutapps.dev/aw/go-starter/internal/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTestClock(t *testing.T) {
	test.WithTestServer(t, func(s *api.Server) {
		now := time.Date(2025, 2, 5, 11, 42, 30, 0, time.UTC)

		test.SetMockClock(t, s, now)

		assert.Equal(t, now, s.Clock.Now())

		clock := test.GetMockClock(t, s.Clock)
		require.NotNil(t, clock)

		assert.Equal(t, now, clock.Now())
	})
}
