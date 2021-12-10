package test

import (
	"context"
	"database/sql"
	"testing"

	"allaboutapps.dev/aw/go-starter/internal/api"
	"allaboutapps.dev/aw/go-starter/internal/api/router"
	"allaboutapps.dev/aw/go-starter/internal/config"
)

// WithTestServer returns a fully configured server (using the default server config).
func WithTestServer(t *testing.T, closure func(s *api.Server)) {
	t.Helper()
	defaultConfig := config.DefaultServiceConfigFromEnv()
	WithTestServerConfigurable(t, defaultConfig, closure)
}

// WithTestServerFromDump returns a fully configured server (using the default server config) and allows for a database dump to be injected.
func WithTestServerFromDump(t *testing.T, dumpConfig DatabaseDumpConfig, closure func(s *api.Server)) {
	t.Helper()
	defaultConfig := config.DefaultServiceConfigFromEnv()
	WithTestServerConfigurableFromDump(t, defaultConfig, dumpConfig, closure)
}

// WithTestServerConfigurable returns a fully configured server, allowing for configuration using the provided server config.
func WithTestServerConfigurable(t *testing.T, config config.Server, closure func(s *api.Server)) {
	t.Helper()
	ctx := context.Background()
	WithTestServerConfigurableContext(ctx, t, config, closure)
}

// WithTestServerConfigurableContext returns a fully configured server, allowing for configuration using the provided server config.
// The provided context will be used during setup (instead of the default background context).
func WithTestServerConfigurableContext(ctx context.Context, t *testing.T, config config.Server, closure func(s *api.Server)) {
	t.Helper()
	WithTestDatabaseContext(ctx, t, func(db *sql.DB) {
		t.Helper()
		execClosureNewTestServer(ctx, t, config, db, closure)
	})
}

// WithTestServerConfigurableFromDump returns a fully configured server, allowing for configuration using the provided server config and a database dump to be injected.
func WithTestServerConfigurableFromDump(t *testing.T, config config.Server, dumpConfig DatabaseDumpConfig, closure func(s *api.Server)) {
	t.Helper()
	ctx := context.Background()
	WithTestServerConfigurableFromDumpContext(ctx, t, config, dumpConfig, closure)
}

// WithTestServerConfigurableFromDumpContext returns a fully configured server, allowing for configuration using the provided server config and a database dump to be injected.
// The provided context will be used during setup (instead of the default background context).
func WithTestServerConfigurableFromDumpContext(ctx context.Context, t *testing.T, config config.Server, dumpConfig DatabaseDumpConfig, closure func(s *api.Server)) {
	t.Helper()
	WithTestDatabaseFromDump(t, dumpConfig, func(db *sql.DB) {
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

	// attach the already initialized db
	s.DB = db

	if err := s.InitMailer(); err != nil {
		t.Fatalf("Failed to init mailer: %v", err)
	}

	// attach any other mocks
	s.Push = NewTestPusher(t, db)

	if err := s.InitI18n(); err != nil {
		t.Fatalf("Failed to init i18n service: %v", err)
	}

	router.Init(s)

	closure(s)

	// echo is managed and should close automatically after running the test
	if err := s.Echo.Shutdown(ctx); err != nil {
		t.Fatalf("failed to shutdown server: %v", err)
	}

	// disallow any further refs to managed object after running the test
	s = nil
}
