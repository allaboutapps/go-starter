package common

import (
	"fmt"
	"net/http"

	"allaboutapps.dev/aw/go-starter/internal/api"
	"github.com/labstack/echo/v4"
)

func GetReadyRoute(s *api.Server) *echo.Route {
	return s.Router.Management.GET("/ready", getReadyHandler(s))
}

// Readiness check
// This endpoint returns 200 when our Service is ready to serve traffic (i.e. respond to queries).
// Structured upon https://prometheus.io/docs/prometheus/latest/management_api/
func getReadyHandler(s *api.Server) echo.HandlerFunc {
	return func(c echo.Context) error {

		if !s.Ready() {
			// We use 521 to indicate an error state
			// same as Cloudflare: https://support.cloudflare.com/hc/en-us/articles/115003011431#521error
			return c.String(521, "Not ready.")
		}

		return c.String(http.StatusOK, "Ready.")
	}
}

func init() {
	fmt.Println("TEST")
}
