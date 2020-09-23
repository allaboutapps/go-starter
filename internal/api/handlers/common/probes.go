package common

import (
	"context"
	"database/sql"
	"fmt"
	"path"
	"strings"
	"sync"
	"time"

	"allaboutapps.dev/aw/go-starter/internal/util"
	"golang.org/x/sys/unix"
)

func ProbeReadiness(ctx context.Context, database *sql.DB, writeablePaths []string) (string, []error) {
	var str strings.Builder

	// slice collects all errors from probes
	errs := make([]error, 0, 1+len(writeablePaths))

	// DB readable?
	dbPingStr, dbPingErr := probeDatabasePingable(ctx, database)
	str.WriteString(dbPingStr)

	if dbPingErr != nil {
		errs = append(errs, dbPingErr)
	}

	// FS (potentially) writeable?
	for _, writeablePath := range writeablePaths {

		fsPermStr, fsPermErr := probePathWriteablePermission(ctx, writeablePath)
		str.WriteString(fsPermStr)

		if fsPermErr != nil {
			errs = append(errs, fsPermErr)
		}
	}

	// Feel free to add additional probes here...

	return str.String(), errs
}

func ProbeLiveness(ctx context.Context, database *sql.DB, writeablePaths []string, touch string) (string, []error) {

	// fail immediately if any readiness probes above have already failed.
	readinessProbeStr, readinessProbeErrs := ProbeReadiness(ctx, database, writeablePaths)

	if len(readinessProbeErrs) != 0 {
		return readinessProbeStr, readinessProbeErrs
	}

	var str strings.Builder

	// include previous readiness probe results in final string
	str.WriteString(readinessProbeStr)

	// slice collects all errors from probes
	errs := make([]error, 0, 1+len(writeablePaths))

	// DB writeable?
	dbHealthStr, dbHealthErr := probeDatabaseNextHealthSequence(ctx, database)
	str.WriteString(dbHealthStr)

	if dbHealthErr != nil {
		errs = append(errs, dbHealthErr)
	}

	// FS writeable?
	for _, writeablePath := range writeablePaths {

		fsTouchStr, fsTouchErr := probePathWriteableTouch(ctx, writeablePath, touch)
		str.WriteString(fsTouchStr)
		if fsTouchErr != nil {
			errs = append(errs, fsTouchErr)
		}
	}

	// Feel free to add additional probes here...

	return str.String(), errs
}

// FS (especially hard mounted NFS paths) or PostgreSQL calls may be blocking or running for too long and thus need to run detached
// We additionally want them to timeout (e.g. useful for hard mounted NFS paths)
// Typically a any context used here will already have a deadline associated
// If not we will explicitly return a short one here.
func ensureProbeDeadlineFromContext(ctx context.Context) time.Time {
	ctxDeadline, hasDeadline := ctx.Deadline()
	if !hasDeadline {
		ctxDeadline = time.Now().Add(1 * time.Second)
	}

	return ctxDeadline
}

func probeDatabasePingable(ctx context.Context, database *sql.DB) (string, error) {
	var str strings.Builder
	ctxDeadline := ensureProbeDeadlineFromContext(ctx)

	dbPingStart := time.Now()

	var dbPingWg sync.WaitGroup
	var dbErr error

	dbPingWg.Add(1)
	go func() {
		dbErr = database.PingContext(ctx)
		dbPingWg.Done()
	}()

	if err := util.WaitTimeout(&dbPingWg, time.Until(ctxDeadline)); err != nil {
		fmt.Fprintf(&str, "Probe db: Ping deadline after %s, error=%v.\n", time.Since(dbPingStart), err.Error())
		return str.String(), err
	}

	if dbErr != nil {
		fmt.Fprintf(&str, "Probe db: Ping errored after %s, error=%v.\n", time.Since(dbPingStart), dbErr.Error())
		return str.String(), dbErr
	}

	fmt.Fprintf(&str, "Probe db: Ping succeeded in %s.\n", time.Since(dbPingStart))

	return str.String(), nil
}

func probeDatabaseNextHealthSequence(ctx context.Context, database *sql.DB) (string, error) {
	var str strings.Builder
	ctxDeadline := ensureProbeDeadlineFromContext(ctx)

	dbWriteStart := time.Now()

	var seqVal int
	var dbWriteWg sync.WaitGroup
	var dbErr error

	dbWriteWg.Add(1)
	go func() {
		dbErr = database.QueryRowContext(ctx, "SELECT nextval('seq_health');").Scan(&seqVal)
		dbWriteWg.Done()
	}()

	if err := util.WaitTimeout(&dbWriteWg, time.Until(ctxDeadline)); err != nil {
		fmt.Fprintf(&str, "Probe db: Next health sequence deadline after %s, error=%v.\n", time.Since(dbWriteStart), err.Error())
		return str.String(), err
	}

	if dbErr != nil {
		fmt.Fprintf(&str, "Probe db: Next health sequence errored after %s, error=%v.\n", time.Since(dbWriteStart), dbErr.Error())
		return str.String(), dbErr
	}

	fmt.Fprintf(&str, "Probe db: Next health sequence succeeded in %s, seq_health=%v.\n", time.Since(dbWriteStart), seqVal)

	return str.String(), nil
}

func probePathWriteablePermission(ctx context.Context, writeablePath string) (string, error) {
	var str strings.Builder
	ctxDeadline := ensureProbeDeadlineFromContext(ctx)

	fsWriteStart := time.Now()

	if ctx.Err() != nil {
		fmt.Fprintf(&str, "Probe path '%s': W_OK check cancelled after %s, error=%v.\n", writeablePath, time.Since(fsWriteStart), ctx.Err())
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
		fmt.Fprintf(&str, "Probe path '%s': W_OK check deadline after %s, error=%v.\n", writeablePath, time.Since(fsWriteStart), err)
		return str.String(), err
	}

	if fsWriteErr != nil {
		fmt.Fprintf(&str, "Probe path '%s': W_OK check errored after %s, error=%v.\n", writeablePath, time.Since(fsWriteStart), fsWriteErr.Error())
		return str.String(), fsWriteErr
	}

	fmt.Fprintf(&str, "Probe path '%s': W_OK check succeeded in %s.\n", writeablePath, time.Since(fsWriteStart))

	return str.String(), nil
}

func probePathWriteableTouch(ctx context.Context, writeablePath string, touch string) (string, error) {
	var str strings.Builder
	ctxDeadline := ensureProbeDeadlineFromContext(ctx)

	fsTouchStart := time.Now()
	fsTouchNameAbs := path.Join(writeablePath, touch)

	if ctx.Err() != nil {
		fmt.Fprintf(&str, "Probe path '%s': Touch cancelled after %s, error=%v.\n", fsTouchNameAbs, time.Since(fsTouchStart), ctx.Err())
		return str.String(), ctx.Err()
	}

	var fsTouchWg sync.WaitGroup
	var fsTouchErr error
	var fsTouchModTime time.Time
	fsTouchWg.Add(1)
	go func(tn string) {
		fsTouchModTime, fsTouchErr = util.TouchFile(tn)
		fsTouchWg.Done()
	}(fsTouchNameAbs)

	if err := util.WaitTimeout(&fsTouchWg, time.Until(ctxDeadline)); err != nil {
		fmt.Fprintf(&str, "Probe path '%s': Touch deadline after %s, error=%v.\n", fsTouchNameAbs, time.Since(fsTouchStart), err)
		return str.String(), err
	}

	if fsTouchErr != nil {
		fmt.Fprintf(&str, "Probe path '%s': Touch errored after %s, error=%v.\n", fsTouchNameAbs, time.Since(fsTouchStart), fsTouchErr.Error())
		return str.String(), fsTouchErr
	}

	fmt.Fprintf(&str, "Probe path '%s': Touch succeeded in %s, modTime=%v.\n", fsTouchNameAbs, time.Since(fsTouchStart), fsTouchModTime.Unix())

	return str.String(), nil
}
