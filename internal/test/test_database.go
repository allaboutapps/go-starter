package test

import (
	"context"
	"crypto/md5" //nolint:gosec
	"database/sql"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	pUtil "allaboutapps.dev/aw/go-starter/internal/util"
	dbutil "allaboutapps.dev/aw/go-starter/internal/util/db"
	"github.com/allaboutapps/integresql-client-go"
	"github.com/allaboutapps/integresql-client-go/pkg/models"
	"github.com/allaboutapps/integresql-client-go/pkg/util"
	"github.com/pkg/errors"
	migrate "github.com/rubenv/sql-migrate"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

var (
	client *integresql.Client

	// tracks hash and IntegreSQL template(s) testDatabase initialization
	poolInitSyncMap = &sync.Map{}             // "poolID" -> *sync.Once
	poolHashMap     = make(map[string]string) // "poolID" -> "integresPoolHash"

	// we will compute a db template hash over the following dirs/files
	migDir           = filepath.Join(pUtil.GetProjectRootDir(), "/migrations")
	fixFile          = filepath.Join(pUtil.GetProjectRootDir(), "/internal/test/fixtures.go")
	selfFile         = filepath.Join(pUtil.GetProjectRootDir(), "/internal/test/test_database.go")
	defaultPoolPaths = []string{migDir, fixFile, selfFile}
)

func init() {
	// autoinitialize IntegreSQL client
	c, err := integresql.DefaultClientFromEnv()
	if err != nil {
		panic(errors.Wrap(err, "Failed to create new integresql-client"))
	}
	client = c
}

// Use this utility func to test with an isolated test database based on your current migrations and fixtures.
func WithTestDatabase(t *testing.T, closure func(db *sql.DB)) {
	t.Helper()
	ctx := context.Background()
	WithTestDatabaseContext(ctx, t, closure)
}

// Use this utility func to test with an isolated test database based on your current migrations and fixtures (context injectable).
func WithTestDatabaseContext(ctx context.Context, t *testing.T, closure func(db *sql.DB)) {
	t.Helper()

	poolID := strings.Join(defaultPoolPaths[:], ",")

	// Get a hold of the &sync.Once{} for this integresql pool in this pkg scope...
	doOnce, _ := poolInitSyncMap.LoadOrStore(poolID, &sync.Once{})
	doOnce.(*sync.Once).Do(func() {
		t.Helper()

		// enable poolHashMap[poolID] integresql hash lookup
		cachePoolIDToIntegresHash(t, poolID, defaultPoolPaths)

		// properly build up the template database once
		execClosureNewIntegresTemplate(ctx, t, poolHashMap[poolID], func(db *sql.DB) error {
			t.Helper()

			countMigrations, err := ApplyMigrations(t, db)
			if err != nil {
				t.Fatalf("Failed to apply migrations for %q: %v\n", poolHashMap[poolID], err)
				return err
			}
			t.Logf("Applied %d migrations for hash %q", countMigrations, poolHashMap[poolID])

			countFixtures, err := ApplyTestFixtures(ctx, t, db)
			if err != nil {
				t.Fatalf("Failed to apply test fixtures for %q: %v\n", poolHashMap[poolID], err)
				return err
			}
			t.Logf("Applied %d test fixtures for hash %q", countFixtures, poolHashMap[poolID])

			return nil
		})
	})

	// execute closure in a new IntegreSQL database build from above template
	execClosureNewIntegresDatabase(ctx, t, poolHashMap[poolID], func(testDatabase models.TestDatabase, db *sql.DB) {
		t.Helper()
		t.Logf("WithTestDatabase: %q", testDatabase.Config.Database)
		closure(db)
	})
}

type DatabaseDumpConfig struct {
	DumpFile          string // required, absolute path to dump file
	ApplyMigrations   bool   // optional, default false
	ApplyTestFixtures bool   // optional, default false
}

// Use this utility func to test with an isolated test database based on a dump file.
func WithTestDatabaseFromDump(t *testing.T, config DatabaseDumpConfig, closure func(db *sql.DB)) {
	t.Helper()
	ctx := context.Background()
	WithTestDatabaseFromDumpContext(ctx, t, config, closure)
}

// Use this utility func to test with an isolated test database based on a dump file (context injectable).
func WithTestDatabaseFromDumpContext(ctx context.Context, t *testing.T, config DatabaseDumpConfig, closure func(db *sql.DB)) {
	t.Helper()

	// DumpFile is mandadory.
	if config.DumpFile == "" {
		t.Fatal("DatabaseDumpConfig.DumpFile is mandadory and cannot be ''")
	}

	// poolID must incorporate additional config args in the final hash
	fragments := fmt.Sprintf("?migrations=%v&fixtures=%v", config.ApplyMigrations, config.ApplyTestFixtures)
	poolID := strings.Join([]string{config.DumpFile, selfFile}[:], ",") + fragments

	// Get a hold of the &sync.Once{} for this integresql pool in this pkg scope...
	doOnce, _ := poolInitSyncMap.LoadOrStore(poolID, &sync.Once{})
	doOnce.(*sync.Once).Do(func() {
		t.Helper()

		// enable poolHashMap[poolID] integresql hash lookup
		cachePoolIDToIntegresHash(t, poolID, []string{config.DumpFile, selfFile}, fragments)

		// properly build up the template database once
		execClosureNewIntegresTemplate(ctx, t, poolHashMap[poolID], func(db *sql.DB) error {
			t.Helper()

			if err := ApplyDump(ctx, t, db, config.DumpFile); err != nil {
				t.Fatalf("Failed to apply dumps for %q: %v\n", poolHashMap[poolID], err)
				return err
			}
			t.Logf("Applied dump for hash %q", poolHashMap[poolID])

			if config.ApplyMigrations {
				countMigrations, err := ApplyMigrations(t, db)
				if err != nil {
					t.Fatalf("Failed to apply migrations for %q: %v\n", poolHashMap[poolID], err)
					return err
				}
				t.Logf("Applied %d migrations for hash %q", countMigrations, poolHashMap[poolID])
			}

			if config.ApplyTestFixtures {
				countFixtures, err := ApplyTestFixtures(ctx, t, db)
				if err != nil {
					t.Fatalf("Failed to apply test fixtures for %q: %v\n", poolHashMap[poolID], err)
					return err
				}
				t.Logf("Applied %d test fixtures for hash %q", countFixtures, poolHashMap[poolID])
			}

			return nil
		})
	})

	execClosureNewIntegresDatabase(ctx, t, poolHashMap[poolID], func(testDatabase models.TestDatabase, db *sql.DB) {
		t.Helper()
		t.Logf("WithTestDatabaseFromDump: %q", testDatabase.Config.Database)
		closure(db)
	})
}

func WithTestDatabaseEmpty(t *testing.T, closure func(db *sql.DB)) {
	t.Helper()
	ctx := context.Background()
	WithTestDatabaseEmptyContext(ctx, t, closure)
}

func WithTestDatabaseEmptyContext(ctx context.Context, t *testing.T, closure func(db *sql.DB)) {
	t.Helper()

	poolID := selfFile

	// Get a hold of the &sync.Once{} for this integresql pool in this pkg scope...
	doOnce, _ := poolInitSyncMap.LoadOrStore(poolID, &sync.Once{})
	doOnce.(*sync.Once).Do(func() {
		t.Helper()

		// enable poolHashMap[poolID] integresql hash lookup
		cachePoolIDToIntegresHash(t, poolID, []string{selfFile})

		// properly build up the template database once (noop)
		execClosureNewIntegresTemplate(ctx, t, poolHashMap[poolID], func(db *sql.DB) error {
			t.Helper()
			return nil
		})
	})

	execClosureNewIntegresDatabase(ctx, t, poolHashMap[poolID], func(testDatabase models.TestDatabase, db *sql.DB) {
		t.Helper()
		t.Logf("WithTestDatabaseEmpty: %q", testDatabase.Config.Database)
		closure(db)
	})
}

// Adds poolID to poolHashMap pointing to the final integresql hash
// Expects hashPaths to be absolute paths to actual files or directories (its contents will be md5 hashed)
// Optional fragments can be used to further enhance the computed md5
func cachePoolIDToIntegresHash(t *testing.T, poolID string, hashPaths []string, fragments ...string) {
	t.Helper()

	// compute a new integreSQL pool hash
	integresPoolHash, err := util.GetTemplateHash(hashPaths...)
	if err != nil {
		t.Fatalf("Failed to create template hash for %v: %#v", poolID, err)
	}

	// update the hash with optional provided fragments
	if len(fragments) > 0 {
		integresPoolHash = fmt.Sprintf("%x", md5.Sum([]byte(integresPoolHash+strings.Join(fragments, ",")))) //nolint:gosec
	}

	// and point poolID to it (sideffect!)
	poolHashMap[poolID] = integresPoolHash // save it for all runners (not having this lock" yet)
}

// Executes closure on an integresql **template** database
func execClosureNewIntegresTemplate(ctx context.Context, t *testing.T, integresPoolHash string, closure func(db *sql.DB) error) {
	t.Helper()

	if err := client.SetupTemplateWithDBClient(ctx, integresPoolHash, closure); err != nil {

		// This error is exceptionally fatal as it hinders ANY future other
		// test execution with this hash as the template was *never* properly
		// setuped successfully. All GetTestDatabase will wait unti timeout
		// unless we interrupt them by discarding the base template...
		discardError := client.DiscardTemplate(ctx, integresPoolHash)

		if discardError != nil {
			t.Fatalf("Failed to setup template database, also discarding failed for integresPoolHash %q: %v, %v", integresPoolHash, err, discardError)
		}

		t.Fatalf("Failed to setup template database (discarded) for integresPoolHash %q: %v", integresPoolHash, err)

	}
}

// Executes closure on an integresql **test** database (scaffolded from a template)
func execClosureNewIntegresDatabase(ctx context.Context, t *testing.T, integresPoolHash string, closure func(testDatabase models.TestDatabase, db *sql.DB)) {
	t.Helper()

	testDatabase, err := client.GetTestDatabase(ctx, integresPoolHash)

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

	closure(testDatabase, db)

	// this database object is managed and should close automatically after running the test
	if err := db.Close(); err != nil {
		t.Fatalf("Failed to close db %q: %v", connectionString, err)
	}

	// disallow any further refs to managed object after running the test
	db = nil
}

// Applies all current database migrations to db
func ApplyMigrations(t *testing.T, db *sql.DB) (countMigrations int, err error) {
	t.Helper()

	migrations := &migrate.FileMigrationSource{Dir: migDir}
	countMigrations, err = migrate.Exec(db, "postgres", migrations, migrate.Up)
	if err != nil {
		return 0, err
	}

	return countMigrations, err
}

// Applies all current test fixtures (insert) to db
func ApplyTestFixtures(ctx context.Context, t *testing.T, db *sql.DB) (countFixtures int, err error) {
	t.Helper()

	inserts := Inserts()

	// insert test fixtures in an auto-managed db transaction
	err = dbutil.WithTransaction(ctx, db, func(tx boil.ContextExecutor) error {
		t.Helper()
		for _, fixture := range inserts {
			if err := fixture.Insert(ctx, tx, boil.Infer()); err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return 0, err
	}

	return len(inserts), nil
}

// Applies dumpFile (absolute path to .sql file) to db
func ApplyDump(ctx context.Context, t *testing.T, db *sql.DB, dumpFile string) error {
	t.Helper()

	// ensure file exists
	if _, err := os.Stat(dumpFile); err != nil {
		return err
	}

	// we need to get the db name before beeing able to do anything.
	var targetDB string
	if err := db.QueryRow("SELECT current_database();").Scan(&targetDB); err != nil {
		return err
	}

	cmd := exec.Command("bash", "-c", fmt.Sprintf("cat %q | psql %q", dumpFile, targetDB)) //nolint:gosec
	combinedOutput, err := cmd.CombinedOutput()

	if err != nil {
		return errors.Wrap(err, string(combinedOutput))
	}

	return err
}
