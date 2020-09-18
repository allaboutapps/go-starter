package common

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"allaboutapps.dev/aw/go-starter/internal/api"
	"github.com/labstack/echo/v4"
	"golang.org/x/sys/unix"
)

func GetHealthyRoute(s *api.Server) *echo.Route {
	return s.Router.Management.GET("/healthy", getHealthyHandler(s))
}

// Health check
// Returns an human readable string about the current service status.
// Does additional checks apart from the general server ready state
// Structured upon https://prometheus.io/docs/prometheus/latest/management_api/
func getHealthyHandler(s *api.Server) echo.HandlerFunc {
	return func(c echo.Context) error {

		if !s.Ready() {
			return c.String(521, "Not ready.")
		}

		var str strings.Builder

		checksHaveErrored := false

		fmt.Fprintln(&str, "Ready.")

		ctx := c.Request().Context()

		// Check database is pingable...
		dbPingStart := time.Now()
		if err := s.DB.PingContext(ctx); err != nil {
			checksHaveErrored = true
			fmt.Fprintf(&str, "Database: Ping errored after %s, error=%v.\n", time.Since(dbPingStart), err.Error())
		} else {
			fmt.Fprintf(&str, "Database: Ping succeeded in %s.\n", time.Since(dbPingStart))
		}

		// Check database is writable...
		dbWriteStart := time.Now()
		var seqVal int
		if err := s.DB.QueryRowContext(ctx, "SELECT nextval('seq_health');").Scan(&seqVal); err != nil {
			checksHaveErrored = true
			fmt.Fprintf(&str, "Database: Next health sequence errored after %s, error=%v.\n", time.Since(dbWriteStart), err.Error())
		} else {
			fmt.Fprintf(&str, "Database: Next health sequence succeeded in %s, seq_health=%v.\n", time.Since(dbWriteStart), seqVal)
		}

		// Check mount is writeable...
		fsStart := time.Now()
		if err := unix.Access(s.Config.Paths.MntBaseDirAbs, unix.W_OK); err != nil {
			checksHaveErrored = true
			fmt.Fprintf(&str, "Mount '%s': Errored after %s, error=%v.\n", s.Config.Paths.MntBaseDirAbs, time.Since(fsStart), err.Error())
		} else {
			fmt.Fprintf(&str, "Mount '%s': Writeable check succeeded in %s.\n", s.Config.Paths.MntBaseDirAbs, time.Since(fsStart))
		}

		// Feel free to add additional checks here...

		if checksHaveErrored {
			fmt.Fprintln(&str, "Not healthy.")
			// We use 521 to indicate this error state
			// same as Cloudflare: https://support.cloudflare.com/hc/en-us/articles/115003011431#521error
			return c.String(521, str.String())
		}

		fmt.Fprintln(&str, "Healthy.")

		return c.String(http.StatusOK, str.String())
	}
}
