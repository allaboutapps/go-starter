package common

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"allaboutapps.dev/aw/go-starter/internal/api"
	"github.com/labstack/echo/v4"
)

func GetHealthyRoute(s *api.Server) *echo.Route {
	return s.Router.Management.GET("/healthy", getHealthyHandler(s))
}

// Heathly check (= liveness)
// Returns an human readable string about the current service status.
// In addition to readiness probes, it performs actual write probes.
// Note that /-/healthy is private (shielded by the mgmt-secret) as it may expose sensitive information about your service.
// Structured upon https://prometheus.io/docs/prometheus/latest/management_api/
func getHealthyHandler(s *api.Server) echo.HandlerFunc {
	return func(c echo.Context) error {

		if !s.Ready() {
			// We use 521 to indicate an error state
			// same as Cloudflare: https://support.cloudflare.com/hc/en-us/articles/115003011431#521error
			return c.String(521, "Not ready.")
		}

		var str strings.Builder
		fmt.Fprintln(&str, "Ready.")

		// General Timeout and associated context.
		ctx, cancel := context.WithTimeout(c.Request().Context(), s.Config.Management.LivenessTimeout)
		defer cancel()

		healthyStr, errs := ProbeLiveness(ctx, s.DB, s.Config.Management.ProbeWriteablePathsAbs, s.Config.Management.ProbeWriteableTouchfile)
		str.WriteString(healthyStr)

		// Finally return the health status according to the seen states
		if ctx.Err() != nil || len(errs) != 0 {
			fmt.Fprintln(&str, "Probes failed.")
			// We use 521 to indicate this error state
			// same as Cloudflare: https://support.cloudflare.com/hc/en-us/articles/115003011431#521error
			return c.String(521, str.String())
		}

		fmt.Fprintln(&str, "Probes succeeded.")

		return c.String(http.StatusOK, str.String())
	}
}
