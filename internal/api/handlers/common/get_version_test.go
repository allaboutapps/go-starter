package common_test

import (
	"net/http"
	"testing"

	"allaboutapps.dev/aw/go-starter/internal/api"
	"allaboutapps.dev/aw/go-starter/internal/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetVersion(t *testing.T) {
	test.WithTestServer(t, func(s *api.Server) {
		res := test.PerformRequest(t, s, "GET", "/-/version?mgmt-secret="+s.Config.Management.Secret, nil, nil)
		require.Equal(t, http.StatusOK, res.Result().StatusCode)
		require.Equal(t, "build.local/misses/ldflags @ < 40 chars git commit hash via ldflags > (1970-01-01T00:00:00+00:00)", res.Body.String()) // build args are not injected during test time.

		// caching: ensure this call is always uncached for browsers
		assert.Equal(t, "no-cache, private, max-age=0", res.Header().Get("Cache-Control"))
		assert.Equal(t, "Thu, 01 Jan 1970 00:00:00 UTC", res.Header().Get("Expires"))
		assert.Equal(t, "0", res.Header().Get("X-Accel-Expires"))
		assert.Equal(t, "no-cache", res.Header().Get("Pragma"))

		// caching: unset
		assert.Empty(t, res.Header().Get("ETag"))
		assert.Empty(t, res.Header().Get("If-Modified-Since"))
		assert.Empty(t, res.Header().Get("If-Match"))
		assert.Empty(t, res.Header().Get("If-None-Match"))
		assert.Empty(t, res.Header().Get("If-Range"))
		assert.Empty(t, res.Header().Get("If-Unmodified-Since"))
	})
}
