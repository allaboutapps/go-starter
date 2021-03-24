package router_test

import (
	"testing"

	"allaboutapps.dev/aw/go-starter/internal/api"
	"allaboutapps.dev/aw/go-starter/internal/config"
	"allaboutapps.dev/aw/go-starter/internal/test"
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
