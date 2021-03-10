package common_test

import (
	"net/http"
	"testing"

	"allaboutapps.dev/aw/go-starter/internal/api"
	"allaboutapps.dev/aw/go-starter/internal/test"
	"github.com/stretchr/testify/require"
)

func TestGetVersion(t *testing.T) {
	t.Parallel()

	test.WithTestServer(t, func(s *api.Server) {
		res := test.PerformRequest(t, s, "GET", "/-/version?mgmt-secret="+s.Config.Management.Secret, nil, nil)
		require.Equal(t, http.StatusOK, res.Result().StatusCode)
		require.Equal(t, res.Body.String(), "@ unknown (unknown)") // build args are not injected during test time.
	})
}
