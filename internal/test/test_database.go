package test

import (
	"context"
	"database/sql"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
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
	client *integresql.Client

	// tracks IntegreSQL template(s) testDatabase initialization
	doOnceMap = &sync.Map{}
	hashMap   = make(map[string]string)

	// we will compute a db template hash over the following dirs/files
	migDir        = filepath.Join(pUtil.GetProjectRootDir(), "/migrations")
	fixFile       = filepath.Join(pUtil.GetProjectRootDir(), "/internal/test/fixtures.go")
	selfFile      = filepath.Join(pUtil.GetProjectRootDir(), "/internal/test/test_database.go")
	defaultPaths  = []string{migDir, fixFile, selfFile}
	defaultPoolID = strings.Join(defaultPaths[:], ",")
)

func init() {
	// initialize our default IntegreSQL template database hash (used by .WithTestDatabase and .WithTestServer)
	h, err := util.GetTemplateHash(defaultPaths...)

	if err != nil {
		panic(fmt.Sprintf("Failed to get default template hash: %#v", err))
	}

	hashMap[defaultPoolID] = h
	fmt.Printf("IntegreSQL default template hash: %v\n", h)
}

// Use this utility func to test with an isolated test database
func WithTestDatabase(t *testing.T, closure func(db *sql.DB)) {

	t.Helper()

	// new context derived from background
	ctx := context.Background()

	doOnce, _ := doOnceMap.LoadOrStore(hashMap[defaultPoolID], &sync.Once{})

	doOnce.(*sync.Once).Do(func() {
		t.Helper()
		initializeTestDatabaseTemplate(ctx, t)
	})

	testDatabase, err := client.GetTestDatabase(ctx, hashMap[defaultPoolID])

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

func WithTestDatabaseFromDump(t *testing.T, dumpFile string, closure func(db *sql.DB)) {

	t.Helper()

	// new context derived from background
	ctx := context.Background()

	doOnce, _ := doOnceMap.LoadOrStore(dumpFile, &sync.Once{})

	doOnce.(*sync.Once).Do(func() {
		t.Helper()
		initializeTestDatabaseTemplateFromDump(ctx, t, dumpFile)
	})

	testDatabase, err := client.GetTestDatabase(ctx, hashMap[dumpFile])

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

	// initTestDatabaseHash(t)

	initIntegresClient(t)

	if err := client.SetupTemplateWithDBClient(ctx, hashMap[defaultPoolID], func(db *sql.DB) error {

		t.Helper()

		err := applyDatabaseMigrations(t, db)

		if err != nil {
			return err
		}

		err = applyDatabaseFixtures(ctx, t, db)

		return err
	}); err != nil {

		// This error is exceptionally fatal as it hinders ANY future other
		// test execution with this hash as the template was *never* properly
		// setuped successfully. All GetTestDatabase will wait unti timeout
		// unless we interrupt them by discarding the base template...
		discardError := client.DiscardTemplate(ctx, hashMap[defaultPoolID])

		if discardError != nil {
			t.Fatalf("Failed to setup template database, also discarding failed for hash %q: %v, %v", hashMap[defaultPoolID], err, discardError)
		}

		t.Fatalf("Failed to setup template database (discarded) for hash %q: %v", hashMap[defaultPoolID], err)

	}
}

func initializeTestDatabaseTemplateFromDump(ctx context.Context, t *testing.T, dumpFile string) {

	t.Helper()

	initTestDatabaseHashFromDump(t, dumpFile)

	initIntegresClient(t)

	if err := client.SetupTemplateWithDBClient(ctx, hashMap[dumpFile], func(db *sql.DB) error {

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
		discardError := client.DiscardTemplate(ctx, hashMap[dumpFile])

		if discardError != nil {
			t.Fatalf("Failed to setup template database, also discarding failed for hash %q: %v, %v", hashMap[dumpFile], err, discardError)
		}

		t.Fatalf("Failed to setup template database (discarded) for hash %q: %v", hashMap[dumpFile], err)

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

func initTestDatabaseHashFromDump(t *testing.T, dumpFile string) {

	t.Helper()

	h, err := util.GetTemplateHash(dumpFile)
	if err != nil {
		t.Fatalf("Failed to get template hash: %#v", err)
	}

	hashMap[dumpFile] = h
}

// func initTestDatabaseHash(t *testing.T, identifiers ...string) {
// 	t.Helper()

// 	util.GetTemplateHash(identifiers...)

// 	h, err := util.GetTemplateHash(identifiers...)
// 	if err != nil {
// 		t.Fatalf("Failed to get template hash for %v: %#v", identifiers, err)
// 	}

// 	hashMap[getShortIdentifier(identifiers...)] = h
// }

// func getShortIdentifier(identifiers ...string) string {
// 	return strings.Join(identifiers[:], ",")
// }

func applyDatabaseMigrations(t *testing.T, db *sql.DB) error {

	t.Helper()

	migrations := &migrate.FileMigrationSource{Dir: migDir}
	n, err := migrate.Exec(db, "postgres", migrations, migrate.Up)
	if err != nil {
		return err
	}

	t.Logf("Applied %d migrations for hash %q", n, hashMap[defaultPoolID])

	return nil
}

func applyDatabaseFixtures(ctx context.Context, t *testing.T, db *sql.DB) error {

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

		t.Logf("Inserted %d fixtures for hash %q", len(inserts), hashMap[defaultPoolID])

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
