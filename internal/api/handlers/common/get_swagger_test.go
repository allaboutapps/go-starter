package common_test

import (
	"net/http"
	"testing"

	"allaboutapps.dev/aw/go-starter/internal/api"
	"allaboutapps.dev/aw/go-starter/internal/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSwaggerYAMLRetrieval(t *testing.T) {
	test.WithTestServer(t, func(s *api.Server) {
		res := test.PerformRequest(t, s, "GET", "/swagger.yml", nil, nil)
		require.Equal(t, http.StatusOK, res.Result().StatusCode)

		// caching: ensure this call is always uncached for browsers
		assert.Equal(t, "no-cache, private, max-age=0", res.Header().Get("Cache-Control"))
		assert.Equal(t, "Thu, 01 Jan 1970 00:00:00 UTC", res.Header().Get("Expires"))
		assert.Equal(t, "0", res.Header().Get("X-Accel-Expires"))
		assert.Equal(t, "no-cache", res.Header().Get("Pragma"))

		// caching: unset
		assert.Equal(t, "", res.Header().Get("ETag"))
		assert.Equal(t, "", res.Header().Get("If-Modified-Since"))
		assert.Equal(t, "", res.Header().Get("If-Match"))
		assert.Equal(t, "", res.Header().Get("If-None-Match"))
		assert.Equal(t, "", res.Header().Get("If-Range"))
		assert.Equal(t, "", res.Header().Get("If-Unmodified-Since"))
	})
}
