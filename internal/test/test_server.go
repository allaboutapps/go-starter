package test

import (
	"context"
	"database/sql"
	"testing"

	"allaboutapps.dev/aw/go-starter/internal/api"
	"allaboutapps.dev/aw/go-starter/internal/api/router"
	"allaboutapps.dev/aw/go-starter/internal/config"
)

// Use this utility func to test with an full blown server (default server config)
func WithTestServer(t *testing.T, closure func(s *api.Server)) {
	t.Helper()
	defaultConfig := config.DefaultServiceConfigFromEnv()
	WithTestServerConfigurable(t, defaultConfig, closure)
}

// Use this utility func to test with an full blown server (default server config, database dump injectable)
func WithTestServerFromDump(t *testing.T, dumpFile string, closure func(s *api.Server)) {
	t.Helper()
	defaultConfig := config.DefaultServiceConfigFromEnv()
	WithTestServerConfigurableFromDump(t, defaultConfig, dumpFile, closure)
}

// Use this utility func to test with an full blown server (server env configurable)
func WithTestServerConfigurable(t *testing.T, config config.Server, closure func(s *api.Server)) {
	t.Helper()
	ctx := context.Background()
	WithTestServerConfigurableContext(ctx, t, config, closure)
}

// Use this utility func to test with an full blown server (server env configurable, context injectable).
func WithTestServerConfigurableContext(ctx context.Context, t *testing.T, config config.Server, closure func(s *api.Server)) {
	t.Helper()
	WithTestDatabaseContext(ctx, t, func(db *sql.DB) {
		t.Helper()
		execClosureNewTestServer(ctx, t, config, db, closure)
	})
}

// Use this utility func to test with an full blown server (server env configurable, database dump injectable).
func WithTestServerConfigurableFromDump(t *testing.T, config config.Server, dumpFile string, closure func(s *api.Server)) {
	t.Helper()
	ctx := context.Background()
	WithTestServerConfigurableFromDumpContext(ctx, t, config, dumpFile, closure)
}

// Use this utility func to test with an full blown server (server env configurable, database dump injectable, context injectable).
func WithTestServerConfigurableFromDumpContext(ctx context.Context, t *testing.T, config config.Server, dumpFile string, closure func(s *api.Server)) {
	t.Helper()
	WithTestDatabaseFromDump(t, dumpFile, func(db *sql.DB) {
		t.Helper()
		execClosureNewTestServer(ctx, t, config, db, closure)
	})
}

// Executes closure on a new test server with a pre-provided database
func execClosureNewTestServer(ctx context.Context, t *testing.T, config config.Server, db *sql.DB, closure func(s *api.Server)) {
	t.Helper()

	// https://stackoverflow.com/questions/43424787/how-to-use-next-available-port-in-http-listenandserve
	// You may use port 0 to indicate you're not specifying an exact port but you want a free, available port selected by the system
	config.Echo.ListenAddress = ":0"

	s := api.NewServer(config)

	// attach the already initalized db
	s.DB = db

	if err := s.InitMailer(); err != nil {
		t.Fatalf("Failed to init mailer: %v", err)
	}

	// attach any other mocks
	s.Push = NewTestPusher(t, db)

	router.Init(s)

	closure(s)

	// echo is managed and should close automatically after running the test
	if err := s.Echo.Shutdown(ctx); err != nil {
		t.Fatalf("failed to shutdown server: %v", err)
	}

	// disallow any further refs to managed object after running the test
	s = nil
}
