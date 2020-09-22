package common

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"path"
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
		fmt.Fprintln(&str, "Ready.")

		ctx := c.Request().Context()

		// DB writeable?
		dbStr, dbErr := CheckHealthyWriteableDatabase(ctx, s.DB)
		str.WriteString(dbStr)

		// FS writeable?
		fsErrs := make([]error, 0, len(s.Config.Management.HealthyCheckWriteablePathsAbs))
		for _, writeablePath := range s.Config.Management.HealthyCheckWriteablePathsAbs {

			fsStr, fsErr := CheckHealthyWriteablePath(ctx, writeablePath, s.Config.Management.HealthyCheckWriteablePathsTouch)
			str.WriteString(fsStr)
			if fsErr != nil {
				fsErrs = append(fsErrs, fsErr)
			}
		}

		// Feel free to add additional checks here...

		// --
		// Finally return the health status according to the seen states
		if dbErr != nil || len(fsErrs) != 0 {
			fmt.Fprintln(&str, "Not healthy.")
			// We use 521 to indicate this error state
			// same as Cloudflare: https://support.cloudflare.com/hc/en-us/articles/115003011431#521error
			return c.String(521, str.String())
		}

		fmt.Fprintln(&str, "Healthy.")

		return c.String(http.StatusOK, str.String())
	}
}

func CheckHealthyWriteableDatabase(ctx context.Context, database *sql.DB) (string, error) {
	var str strings.Builder

	// Check database is pingable...
	dbPingStart := time.Now()
	if err := database.PingContext(ctx); err != nil {
		fmt.Fprintf(&str, "Database: Ping errored after %s, error=%v.\n", time.Since(dbPingStart), err.Error())
		return str.String(), err
	}

	fmt.Fprintf(&str, "Database: Ping succeeded in %s.\n", time.Since(dbPingStart))

	// Check database is writable...
	dbWriteStart := time.Now()
	var seqVal int
	if err := database.QueryRowContext(ctx, "SELECT nextval('seq_health');").Scan(&seqVal); err != nil {
		fmt.Fprintf(&str, "Database: Next health sequence errored after %s, error=%v.\n", time.Since(dbWriteStart), err.Error())
		return str.String(), err
	}

	fmt.Fprintf(&str, "Database: Next health sequence succeeded in %s, seq_health=%v.\n", time.Since(dbWriteStart), seqVal)
	return str.String(), nil
}

func CheckHealthyWriteablePath(ctx context.Context, writeablePath string, touch string) (string, error) {
	var str strings.Builder

	// Check Path is writeable...
	fsStart := time.Now()
	if err := unix.Access(writeablePath, unix.W_OK); err != nil {
		fmt.Fprintf(&str, "Path '%s': Writeable check errored after %s, error=%v.\n", writeablePath, time.Since(fsStart), err.Error())
		return str.String(), err
	}
	fmt.Fprintf(&str, "Path '%s': Writeable check succeeded in %s.\n", writeablePath, time.Since(fsStart))

	// Actually write a file...
	fsWriteStart := time.Now()
	fsNameAbs := path.Join(writeablePath, touch)
	modTime, err := touchFile(fsNameAbs)
	if err != nil {
		fmt.Fprintf(&str, "Touch '%s': Write errored after %s, error=%v.\n", fsNameAbs, time.Since(fsWriteStart), err.Error())
		return str.String(), err
	}
	fmt.Fprintf(&str, "Touch '%s': Write succeeded in %s, modTime=%v.\n", fsNameAbs, time.Since(fsWriteStart), modTime.Unix())

	return str.String(), nil
}

func touchFile(fileName string) (time.Time, error) {
	_, err := os.Stat(fileName)

	if os.IsNotExist(err) {
		file, err := os.Create(fileName)

		if err != nil {
			return time.Time{}, err
		}

		defer file.Close()

		stat, err := file.Stat()

		if err != nil {
			return time.Time{}, err
		}

		return stat.ModTime(), nil
	}

	currentTime := time.Now().Local()
	err = os.Chtimes(fileName, currentTime, currentTime)
	return currentTime, err
}
