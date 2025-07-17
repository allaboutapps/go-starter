package common

import (
	"context"
	"net/http"

	"allaboutapps.dev/aw/go-starter/internal/api"
	"github.com/labstack/echo/v4"
)

func GetReadyRoute(s *api.Server) *echo.Route {
	return s.Router.Management.GET("/ready", getReadyHandler(s))
}

// Readiness check
// This endpoint returns 200 when our Service is ready to serve traffic (i.e. respond to queries).
// Does read-only probes apart from the general server ready state.
// Note that /-/ready is typically public (and not shielded by a mgmt-secret), we thus prevent information leakage here and only return `"Ready."`.
// Structured upon https://prometheus.io/docs/prometheus/latest/management_api/
func getReadyHandler(s *api.Server) echo.HandlerFunc {
	return func(c echo.Context) error {
		if !s.Ready() {
			return c.String(httpStatusDown, "Not ready.")
		}

		// General Timeout and associated context.
		ctx, cancel := context.WithTimeout(c.Request().Context(), s.Config.Management.ReadinessTimeout)
		defer cancel()

		_, errs := ProbeReadiness(ctx, s.DB, s.Config.Management.ProbeWriteablePathsAbs)

		// Finally return the health status according to the seen states
		if ctx.Err() != nil || len(errs) != 0 {
			return c.String(httpStatusDown, "Not ready.")
		}

		return c.String(http.StatusOK, "Ready.")
	}
}
