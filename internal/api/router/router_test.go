package router_test

import (
	"fmt"
	"net/http"
	"testing"

	"allaboutapps.dev/aw/go-starter/internal/api"
	"allaboutapps.dev/aw/go-starter/internal/config"
	"allaboutapps.dev/aw/go-starter/internal/metrics/users"
	"allaboutapps.dev/aw/go-starter/internal/models"
	"allaboutapps.dev/aw/go-starter/internal/test"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPprofEnabled(t *testing.T) {
	config := config.DefaultServiceConfigFromEnv()

	// these are typically our default values, however we force set them here to ensure those are set while test execution.
	config.Pprof.Enable = true
	config.Pprof.EnableManagementKeyAuth = true

	test.WithTestServerConfigurable(t, config, func(s *api.Server) {
		// heap (test any)
		res := test.PerformRequest(t, s, "GET", "/debug/pprof/heap?mgmt-secret="+s.Config.Management.Secret, nil, nil)
		require.Equal(t, 200, res.Result().StatusCode)

		// index
		res = test.PerformRequest(t, s, "GET", "/debug/pprof?mgmt-secret="+s.Config.Management.Secret, nil, nil)
		require.Equal(t, 200, res.Result().StatusCode)

		res = test.PerformRequest(t, s, "GET", "/debug/pprof/heap?mgmt-secret=wrongsecret", nil, nil)
		require.Equal(t, 401, res.Result().StatusCode)
	})
}

func TestPprofEnabledNoAuth(t *testing.T) {
	config := config.DefaultServiceConfigFromEnv()

	// these are typically our default values, however we force set them here to ensure those are set while test execution.
	config.Pprof.Enable = true
	config.Pprof.EnableManagementKeyAuth = false

	test.WithTestServerConfigurable(t, config, func(s *api.Server) {
		res := test.PerformRequest(t, s, "GET", "/debug/pprof/heap?", nil, nil)
		require.Equal(t, 200, res.Result().StatusCode)
	})
}

func TestPprofDisabled(t *testing.T) {
	config := config.DefaultServiceConfigFromEnv()
	config.Pprof.Enable = false

	test.WithTestServerConfigurable(t, config, func(s *api.Server) {
		res := test.PerformRequest(t, s, "GET", "/debug/pprof/heap?mgmt-secret="+s.Config.Management.Secret, nil, nil)
		require.Equal(t, 404, res.Result().StatusCode)
	})
}

func TestMiddlewaresDisabled(t *testing.T) {
	// disable all
	config := config.DefaultServiceConfigFromEnv()
	config.Echo.EnableCORSMiddleware = false
	config.Echo.EnableLoggerMiddleware = false
	config.Echo.EnableRecoverMiddleware = false
	config.Echo.EnableRequestIDMiddleware = false
	config.Echo.EnableSecureMiddleware = false
	config.Echo.EnableTrailingSlashMiddleware = false

	test.WithTestServerConfigurable(t, config, func(s *api.Server) {
		res := test.PerformRequest(t, s, "GET", "/-/ready", nil, nil)
		require.Equal(t, 200, res.Result().StatusCode)
	})
}

func TestMetricsEnabled(t *testing.T) {
	config := config.DefaultServiceConfigFromEnv()
	config.Management.EnableMetrics = true

	test.WithTestServerConfigurable(t, config, func(s *api.Server) {
		res := test.PerformRequest(t, s, "GET", "/metrics", nil, nil)
		require.Equal(t, http.StatusOK, res.Result().StatusCode)

		result := res.Body.String()

		// expect custom metric for the total user count
		expectedTotalUserCount, err := models.Users().Count(t.Context(), s.DB)
		require.NoError(t, err)

		assert.Contains(t, result, fmt.Sprintf("%s %d", users.MetricNameTotalUserCount, expectedTotalUserCount))

		// expect sqlstats metrics
		assert.Contains(t, result, "go_sql_stats_connections")
	})
}

func TestMetricsDisabled(t *testing.T) {
	test.WithTestServer(t, func(s *api.Server) {
		res := test.PerformRequest(t, s, "GET", "/metrics", nil, nil)
		require.Equal(t, http.StatusNotFound, res.Result().StatusCode)
	})
}

func TestNotFound(t *testing.T) {
	test.WithTestServer(t, func(s *api.Server) {
		t.Run("AcceptApplicationJSON", func(t *testing.T) {
			headers := http.Header{}
			headers.Set(echo.HeaderAccept, echo.MIMEApplicationJSON)

			res := test.PerformRequest(t, s, "GET", "/api/v1/unknown-path", nil, headers)
			require.Equal(t, http.StatusNotFound, res.Result().StatusCode)

			test.Snapshoter.Save(t, res.Body.String())
		})

		t.Run("AcceptTextHTML", func(t *testing.T) {
			headers := http.Header{}
			headers.Set(echo.HeaderAccept, "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")

			res := test.PerformRequest(t, s, "GET", "/api/v1/unknown-path", nil, headers)
			require.Equal(t, http.StatusNotFound, res.Result().StatusCode)

			test.Snapshoter.Save(t, res.Body.String())
		})
	})
}
