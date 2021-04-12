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

func WithTestServerFromDump(t *testing.T, closure func(s *api.Server), dumpFile string) {

	t.Helper()

	defaultConfig := config.DefaultServiceConfigFromEnv()

	WithTestServerConfigurableFromDump(t, defaultConfig, dumpFile, closure)
}

// Use this utility func to test with an full blown server (server env configurable)
func WithTestServerConfigurable(t *testing.T, config config.Server, closure func(s *api.Server)) {
	t.Helper()

	WithTestDatabase(t, func(db *sql.DB) {

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
		if err := s.Echo.Shutdown(context.Background()); err != nil {
			t.Fatalf("failed to shutdown server: %v", err)
		}

		// disallow any further refs to managed object after running the test
		s = nil
	})
}

func WithTestServerConfigurableFromDump(t *testing.T, config config.Server, dumpFile string, closure func(s *api.Server)) {
	t.Helper()

	WithTestDatabaseFromDump(t, dumpFile, func(db *sql.DB) {

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
		if err := s.Echo.Shutdown(context.Background()); err != nil {
			t.Fatalf("failed to shutdown server: %v", err)
		}

		// disallow any further refs to managed object after running the test
		s = nil
	})
}
