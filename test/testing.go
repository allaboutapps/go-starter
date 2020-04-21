package test

import (
	"context"
	"database/sql"
	"sync"
	"testing"

	"allaboutapps.at/aw/go-mranftl-sample/api"
	"allaboutapps.at/aw/go-mranftl-sample/api/router"
	"github.com/allaboutapps/integresql-client-go"
	"github.com/allaboutapps/integresql-client-go/pkg/util"
	migrate "github.com/rubenv/sql-migrate"
	"github.com/volatiletech/sqlboiler/boil"
)

var (
	client *integresql.Client
	hash   string

	// tracks template testDatabase initialization
	doOnce sync.Once

	// ! TODO golang does not support relative paths in files properly
	// It's only possible to supply this by
	// Use ENV var to specify app-root
	migDir  = "/app/migrations"
	fixFile = "/app/test/fixtures.go"
)

// Use this utility func to test with an isolated test database
func WithTestDatabase(t *testing.T, closure func(db *sql.DB)) {

	// mark as helper
	t.Helper()

	// new context derived from background
	ctx := context.Background()

	doOnce.Do(func() {
		initializeTestDatabaseTemplate(ctx, t)
	})

	testDatabase, err := client.GetTestDatabase(ctx, hash)

	if err != nil {
		t.Fatalf("Failed to obtain TestDatabase: %v", err)
	}

	connectionString := testDatabase.Config.ConnectionString()

	db, err := sql.Open("postgres", connectionString)

	if err != nil {
		t.Fatalf("Failed to setup testdatabase for connectionString %q: %v", connectionString, err)
	}

	// this database object is managed and should close automatically after running the test
	defer db.Close()

	if err := db.PingContext(ctx); err != nil {
		t.Fatalf("Failed to ping testdatabase for connectionString %q: %v", connectionString, err)
	}

	closure(db)
}

// Use this utility func to test with an full blown server
func WithTestServer(t *testing.T, closure func(s *api.Server)) {

	// mark as helper
	t.Helper()

	WithTestDatabase(t, func(db *sql.DB) {

		defaultConfig := api.DefaultServiceConfigFromEnv()

		// https://stackoverflow.com/questions/43424787/how-to-use-next-available-port-in-http-listenandserve
		// You may use port 0 to indicate you're not specifying an exact port but you want a free, available port selected by the system
		defaultConfig.Echo.ListenAddress = ":0"

		s := api.NewServer(defaultConfig)

		// attach the already initalized db
		s.DB = db

		router.Init(s)

		// no need to actually start echo!
		// see https://github.com/labstack/echo/issues/659

		closure(s)
	})
}

// main private function to properly build up the template database
// ensure it is called once once per pkg scope.
func initializeTestDatabaseTemplate(ctx context.Context, t *testing.T) {

	// mark as helper
	t.Helper()

	initTestDatabaseHash(t)

	initIntegresClient(t)

	if err := client.SetupTemplateWithDBClient(ctx, hash, func(db *sql.DB) error {

		err := applyMigrations(t, db)

		if err != nil {
			return err
		}

		err = insertFixtures(ctx, t, db)

		return err
	}); err != nil {
		t.Fatalf("Failed to setup template database for hash %q: %v", hash, err)
	}
}

func initIntegresClient(t *testing.T) {

	// mark as helper
	t.Helper()

	c, err := integresql.DefaultClientFromEnv()
	if err != nil {
		t.Fatalf("Failed to create new integresql-client: %v", err)
	}

	client = c
}

func initTestDatabaseHash(t *testing.T) {

	// mark as helper
	t.Helper()

	h, err := util.GetTemplateHash(migDir, fixFile)
	if err != nil {
		t.Fatalf("Failed to get template hash: %#v", err)
	}

	hash = h
}

func applyMigrations(t *testing.T, db *sql.DB) error {

	// mark as helper
	t.Helper()

	migrations := &migrate.FileMigrationSource{Dir: migDir}
	n, err := migrate.Exec(db, "postgres", migrations, migrate.Up)
	if err != nil {
		return err
	}

	t.Logf("Applied %d migrations for hash %q", n, hash)

	return nil
}

func insertFixtures(ctx context.Context, t *testing.T, db *sql.DB) error {

	// mark as helper
	t.Helper()

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	for _, fixture := range fixtures {
		if err := fixture.Insert(ctx, db, boil.Infer()); err != nil {
			if err := tx.Rollback(); err != nil {
				return err
			}

			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	t.Logf("Inserted %d fixtures for hash %q", len(fixtures), hash)

	return nil
}
