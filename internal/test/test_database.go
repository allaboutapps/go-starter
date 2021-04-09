package test

import (
	"context"
	"database/sql"
	"io/ioutil"
	"path/filepath"
	"sync"
	"testing"

	pUtil "allaboutapps.dev/aw/go-starter/internal/util"
	dbutil "allaboutapps.dev/aw/go-starter/internal/util/db"
	"github.com/allaboutapps/integresql-client-go"
	"github.com/allaboutapps/integresql-client-go/pkg/util"
	migrate "github.com/rubenv/sql-migrate"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

var (
	client   *integresql.Client
	hash     string
	hashDump string

	// tracks template testDatabase initialization
	doOnce     sync.Once
	doOnceDump sync.Once

	// we will compute a db template hash over the following dirs/files
	migDir   = filepath.Join(pUtil.GetProjectRootDir(), "/migrations")
	fixFile  = filepath.Join(pUtil.GetProjectRootDir(), "/internal/test/fixtures.go")
	selfFile = filepath.Join(pUtil.GetProjectRootDir(), "/internal/test/test_database.go")
)

// Use this utility func to test with an isolated test database
func WithTestDatabase(t *testing.T, closure func(db *sql.DB)) {

	t.Helper()

	// new context derived from background
	ctx := context.Background()

	doOnce.Do(func() {

		t.Helper()
		initializeTestDatabaseTemplate(ctx, t)
	})

	testDatabase, err := client.GetTestDatabase(ctx, hash)

	if err != nil {
		t.Fatalf("Failed to obtain test database: %v", err)
	}

	connectionString := testDatabase.Config.ConnectionString()

	db, err := sql.Open("postgres", connectionString)

	if err != nil {
		t.Fatalf("Failed to setup test database for connectionString %q: %v", connectionString, err)
	}

	if err := db.PingContext(ctx); err != nil {
		t.Fatalf("Failed to ping test database for connectionString %q: %v", connectionString, err)
	}

	t.Logf("WithTestDatabase: %q", testDatabase.Config.Database)

	closure(db)

	// this database object is managed and should close automatically after running the test
	if err := db.Close(); err != nil {
		t.Fatalf("Failed to close db %q: %v", connectionString, err)
	}

	// disallow any further refs to managed object after running the test
	db = nil
}

func WithTestDatabaseFromDump(t *testing.T, closure func(db *sql.DB), dumpFile string) {

	t.Helper()

	// new context derived from background
	ctx := context.Background()

	doOnceDump.Do(func() {

		t.Helper()
		initializeTestDatabaseTemplateFromDump(ctx, t, dumpFile)
	})

	testDatabase, err := client.GetTestDatabase(ctx, hashDump)

	if err != nil {
		t.Fatalf("Failed to obtain test database: %v", err)
	}

	connectionString := testDatabase.Config.ConnectionString()

	db, err := sql.Open("postgres", connectionString)

	if err != nil {
		t.Fatalf("Failed to setup test database for connectionString %q: %v", connectionString, err)
	}

	if err := db.PingContext(ctx); err != nil {
		t.Fatalf("Failed to ping test database for connectionString %q: %v", connectionString, err)
	}

	t.Logf("WithTestDatabase: %q", testDatabase.Config.Database)

	closure(db)

	// this database object is managed and should close automatically after running the test
	if err := db.Close(); err != nil {
		t.Fatalf("Failed to close db %q: %v", connectionString, err)
	}

	// disallow any further refs to managed object after running the test
	db = nil
}

// main private function to properly build up the template database
// ensure it is called once once per pkg scope.
func initializeTestDatabaseTemplate(ctx context.Context, t *testing.T) {

	t.Helper()

	initTestDatabaseHash(t)

	initIntegresClient(t)

	if err := client.SetupTemplateWithDBClient(ctx, hash, func(db *sql.DB) error {

		t.Helper()

		err := applyMigrations(t, db)

		if err != nil {
			return err
		}

		err = insertFixtures(ctx, t, db)

		return err
	}); err != nil {

		// This error is exceptionally fatal as it hinders ANY future other
		// test execution with this hash as the template was *never* properly
		// setuped successfully. All GetTestDatabase will wait unti timeout
		// unless we interrupt them by discarding the base template...
		discardError := client.DiscardTemplate(ctx, hash)

		if discardError != nil {
			t.Fatalf("Failed to setup template database, also discarding failed for hash %q: %v, %v", hash, err, discardError)
		}

		t.Fatalf("Failed to setup template database (discarded) for hash %q: %v", hash, err)

	}
}

func initializeTestDatabaseTemplateFromDump(ctx context.Context, t *testing.T, dumpFile string) {

	t.Helper()

	initTestDatabaseHashFromDump(t, dumpFile)

	initIntegresClient(t)

	if err := client.SetupTemplateWithDBClient(ctx, hashDump, func(db *sql.DB) error {

		t.Helper()

		err := applyDumpFile(t, db, dumpFile)

		if err != nil {
			return err
		}

		return err
	}); err != nil {

		// This error is exceptionally fatal as it hinders ANY future other
		// test execution with this hash as the template was *never* properly
		// setuped successfully. All GetTestDatabase will wait unti timeout
		// unless we interrupt them by discarding the base template...
		discardError := client.DiscardTemplate(ctx, hashDump)

		if discardError != nil {
			t.Fatalf("Failed to setup template database, also discarding failed for hash %q: %v, %v", hash, err, discardError)
		}

		t.Fatalf("Failed to setup template database (discarded) for hash %q: %v", hash, err)

	}
}

func initIntegresClient(t *testing.T) {

	t.Helper()

	c, err := integresql.DefaultClientFromEnv()
	if err != nil {
		t.Fatalf("Failed to create new integresql-client: %v", err)
	}

	client = c
}

func initTestDatabaseHash(t *testing.T) {

	t.Helper()

	h, err := util.GetTemplateHash(migDir, fixFile, selfFile)
	if err != nil {
		t.Fatalf("Failed to get template hash: %#v", err)
	}

	hash = h
}

func initTestDatabaseHashFromDump(t *testing.T, dumpFile string) {

	t.Helper()

	h, err := util.GetTemplateHash(selfFile, dumpFile)
	if err != nil {
		t.Fatalf("Failed to get template hash: %#v", err)
	}

	hashDump = h
}

func applyMigrations(t *testing.T, db *sql.DB) error {

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

	t.Helper()

	// insert test fixtures in an auto-managed db transaction
	return dbutil.WithTransaction(ctx, db, func(tx boil.ContextExecutor) error {
		inserts := Inserts()

		for _, fixture := range inserts {
			if err := fixture.Insert(ctx, tx, boil.Infer()); err != nil {
				t.Errorf("Failed to upsert test fixture: %v\n", err)
				return err
			}
		}

		t.Logf("Inserted %d fixtures for hash %q", len(inserts), hash)

		return nil
	})
}

func applyDumpFile(t *testing.T, db *sql.DB, dumpFile string) error {

	t.Helper()

	c, err := ioutil.ReadFile(dumpFile)
	if err != nil {
		t.Errorf("Failed to read dumpfile: %v\n", err)
		return err
	}
	sql := string(c)
	_, err = db.Exec(sql)
	if err != nil {
		t.Errorf("Failed to execute dumpfile: %v\n", err)
		return err
	}

	return nil
}
