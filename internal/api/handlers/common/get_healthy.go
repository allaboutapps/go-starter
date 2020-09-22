package common

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"path"
	"strings"
	"sync"
	"time"

	"allaboutapps.dev/aw/go-starter/internal/api"
	"allaboutapps.dev/aw/go-starter/internal/util"
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

		// General Timeout and associated context.
		ctx, cancel := context.WithTimeout(c.Request().Context(), s.Config.Management.HealthyTimeout)
		defer cancel()

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
		if ctx.Err() != nil || dbErr != nil || len(fsErrs) != 0 {
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

	// FS calls may be blocking and thus need to run detached
	// However, we want them to timeout (e.g. useful for hard mounted NFS paths)
	// Typically a context will already have a deadline associated, if not we will define one here.
	ctxDeadline, hasDeadline := ctx.Deadline()
	if !hasDeadline {
		ctxDeadline = time.Now().Add(5 * time.Second)
	}

	{
		// ---
		// Check Path is writeable...
		fsWriteStart := time.Now()

		if ctx.Err() != nil {
			fmt.Fprintf(&str, "Path '%s': Writeable check cancelled after %s, error=%v.\n", writeablePath, time.Since(fsWriteStart), ctx.Err())
			return str.String(), ctx.Err()
		}

		var fsWriteWg sync.WaitGroup
		var fsWriteErr error
		fsWriteWg.Add(1)
		go func(wp string) {
			fsWriteErr = unix.Access(wp, unix.W_OK)
			fsWriteWg.Done()
		}(writeablePath)

		if err := util.WaitTimeout(&fsWriteWg, time.Until(ctxDeadline)); err != nil {
			fmt.Fprintf(&str, "Path '%s': Writeable check deadline after %s, error=%v.\n", writeablePath, time.Since(fsWriteStart), err)
			return str.String(), err
		}

		if fsWriteErr != nil {
			fmt.Fprintf(&str, "Path '%s': Writeable check errored after %s, error=%v.\n", writeablePath, time.Since(fsWriteStart), fsWriteErr.Error())
			return str.String(), fsWriteErr
		}

		fmt.Fprintf(&str, "Path '%s': Writeable check succeeded in %s.\n", writeablePath, time.Since(fsWriteStart))

	}

	{

		// ---
		// Actually write a file...
		fsTouchStart := time.Now()
		fsTouchNameAbs := path.Join(writeablePath, touch)

		if ctx.Err() != nil {
			fmt.Fprintf(&str, "Touch '%s': Write cancelled after %s, error=%v.\n", fsTouchNameAbs, time.Since(fsTouchStart), ctx.Err())
			return str.String(), ctx.Err()
		}

		var fsTouchWg sync.WaitGroup
		var fsTouchErr error
		var fsTouchModTime time.Time
		fsTouchWg.Add(1)
		go func(tn string) {
			fsTouchModTime, fsTouchErr = touchFile(tn)
			fsTouchWg.Done()
		}(fsTouchNameAbs)

		if err := util.WaitTimeout(&fsTouchWg, time.Until(ctxDeadline)); err != nil {
			fmt.Fprintf(&str, "Touch '%s': Write deadline after %s, error=%v.\n", fsTouchNameAbs, time.Since(fsTouchStart), err)
			return str.String(), err
		}

		if fsTouchErr != nil {
			fmt.Fprintf(&str, "Touch '%s': Write errored after %s, error=%v.\n", fsTouchNameAbs, time.Since(fsTouchStart), fsTouchErr.Error())
			return str.String(), fsTouchErr
		}

		fmt.Fprintf(&str, "Touch '%s': Write succeeded in %s, modTime=%v.\n", fsTouchNameAbs, time.Since(fsTouchStart), fsTouchModTime.Unix())

	}

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
