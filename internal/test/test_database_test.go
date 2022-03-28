package test_test

import (
	"context"
	"database/sql"
	"fmt"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"allaboutapps.dev/aw/go-starter/internal/config"
	"allaboutapps.dev/aw/go-starter/internal/test"
	pUtil "allaboutapps.dev/aw/go-starter/internal/util"
	"github.com/stretchr/testify/require"
)

func TestWithTestDatabaseConcurrentUsage(t *testing.T) {
	wg := sync.WaitGroup{}
	wg.Add(4)

	go func() {
		test.WithTestDatabase(t, func(db1 *sql.DB) {
			wg.Done()
		})
	}()

	go func() {
		test.WithTestDatabaseEmpty(t, func(db2 *sql.DB) {
			wg.Done()
		})
	}()

	go func() {
		test.WithTestDatabaseFromDump(t, test.DatabaseDumpConfig{DumpFile: filepath.Join(pUtil.GetProjectRootDir(), "/test/testdata/plain.sql")}, func(db3 *sql.DB) {
			wg.Done()
		})
	}()

	go func() {
		test.WithTestDatabaseFromDump(t, test.DatabaseDumpConfig{DumpFile: filepath.Join(pUtil.GetProjectRootDir(), "/test/testdata/users.sql")}, func(db4 *sql.DB) {
			wg.Done()
		})
	}()

	// the above will concurrently write to the database pool maps,
	wg.Wait()
}

func TestWithTestDatabase(t *testing.T) {
	test.WithTestDatabase(t, func(db1 *sql.DB) {
		test.WithTestDatabase(t, func(db2 *sql.DB) {

			var db1Name string
			err := db1.QueryRow("SELECT current_database();").Scan(&db1Name)
			if err != nil {
				t.Fatal(err)
			}

			var db2Name string
			err = db2.QueryRow("SELECT current_database();").Scan(&db2Name)
			if err != nil {
				t.Fatal(err)
			}

			require.NotEqual(t, db1Name, db2Name)
		})
	})
}

func TestWithTestDatabaseFromDump(t *testing.T) {

	dumpFile := filepath.Join(pUtil.GetProjectRootDir(), "/test/testdata/users.sql")

	test.WithTestDatabaseFromDump(t, test.DatabaseDumpConfig{DumpFile: dumpFile}, func(db1 *sql.DB) {
		test.WithTestDatabaseFromDump(t, test.DatabaseDumpConfig{DumpFile: dumpFile}, func(db2 *sql.DB) {

			var db1Name string
			if err := db1.QueryRow("SELECT current_database();").Scan(&db1Name); err != nil {
				t.Fatal(err)
			}

			var db2Name string
			if err := db2.QueryRow("SELECT current_database();").Scan(&db2Name); err != nil {
				t.Fatal(err)
			}

			require.NotEqual(t, db1Name, db2Name)

			if _, err := db2.Exec("DELETE FROM users WHERE true;"); err != nil {
				t.Fatal(err)
			}

			var userCount1 int
			if err := db1.QueryRow("SELECT count(id) FROM users;").Scan(&userCount1); err != nil {
				t.Fatal(err)
			}
			require.Equal(t, 3, userCount1)

			var userCount2 int
			if err := db2.QueryRow("SELECT count(id) FROM users;").Scan(&userCount2); err != nil {
				t.Fatal(err)
			}
			require.Equal(t, 0, userCount2)
		})
	})
}

func TestWithTestDatabaseFromDumpAutoMigrateAndTestFixtures(t *testing.T) {
	dumpFile := filepath.Join(pUtil.GetProjectRootDir(), "/test/testdata/plain.sql")

	test.WithTestDatabaseFromDump(t, test.DatabaseDumpConfig{DumpFile: dumpFile}, func(db0 *sql.DB) {
		test.WithTestDatabaseFromDump(t, test.DatabaseDumpConfig{DumpFile: dumpFile, ApplyMigrations: true}, func(db1 *sql.DB) {
			test.WithTestDatabaseFromDump(t, test.DatabaseDumpConfig{DumpFile: dumpFile, ApplyMigrations: true, ApplyTestFixtures: true}, func(db2 *sql.DB) {

				// db0: has only a plain dump
				// db1: has migrations
				// db2: has migrations and testFixtures

				var db0Name string
				if err := db0.QueryRow("SELECT current_database();").Scan(&db0Name); err != nil {
					t.Fatal(err)
				}

				var db1Name string
				if err := db1.QueryRow("SELECT current_database();").Scan(&db1Name); err != nil {
					t.Fatal(err)
				}

				var db2Name string
				if err := db2.QueryRow("SELECT current_database();").Scan(&db2Name); err != nil {
					t.Fatal(err)
				}

				require.NotEqual(t, db0Name, db1Name)
				require.NotEqual(t, db1Name, db2Name)
				require.NotEqual(t, db2Name, db0Name)

				// expect hash to be different for all 3 databases!
				db0Hash := strings.Split(strings.Join(strings.Split(db0Name, "integresql_test_"), ""), "_")[0]
				db1Hash := strings.Split(strings.Join(strings.Split(db1Name, "integresql_test_"), ""), "_")[0]
				db2Hash := strings.Split(strings.Join(strings.Split(db2Name, "integresql_test_"), ""), "_")[0]

				require.NotEqual(t, db0Hash, db1Hash)
				require.NotEqual(t, db1Hash, db2Hash)
				require.NotEqual(t, db2Hash, db0Hash)
			})
		})
	})
}

func TestWithTestDatabaseFromDumpGorp(t *testing.T) {
	dumpFile := filepath.Join(pUtil.GetProjectRootDir(), "/test/testdata/uuid_extension_only.sql")

	// migrate transforms gorp_migratons to migrations
	test.WithTestDatabaseFromDump(t, test.DatabaseDumpConfig{DumpFile: dumpFile, ApplyMigrations: true}, func(db *sql.DB) {

		// check that we properly renamed the "gorp_migrations" migration tracking table to config.DatabaseMigrationTable
		var migrationsTableName string
		if err := db.QueryRow(fmt.Sprintf("SELECT table_name FROM information_schema.tables WHERE table_schema = 'public' AND table_name = '%s';", config.DatabaseMigrationTable)).Scan(&migrationsTableName); err != nil {
			t.Fatal(err)
		}

		require.Equal(t, config.DatabaseMigrationTable, migrationsTableName)

		// check that gorp_migrations does not exist!
		var gorp string
		err := db.QueryRow("SELECT table_name FROM information_schema.tables WHERE table_schema = 'public' AND table_name = 'gorp_migrations';").Scan(&gorp)

		require.Error(t, err)

		// we can expect that the '20200428064736-install-extension-uuid.sql' migration within the uuid_extension_only.sql is a very stable migrations.
		// let's check if other migrations were applied...
		var migrationsCount int
		err = db.QueryRow("SELECT COUNT(table_name) FROM information_schema.tables WHERE table_schema = 'public';").Scan(&migrationsCount)
		require.NoError(t, err)

		require.Greater(t, migrationsCount, 1)
	})

	// with no migrate, gorp_migrations still exists:
	test.WithTestDatabaseFromDump(t, test.DatabaseDumpConfig{DumpFile: dumpFile, ApplyMigrations: false}, func(db *sql.DB) {

		// check that "gorp_migrations" migration tracking table still exists (as we have not set ApplyMigrations to true!)
		var migrationsTableName string
		if err := db.QueryRow("SELECT table_name FROM information_schema.tables WHERE table_schema = 'public' AND table_name = 'gorp_migrations';").Scan(&migrationsTableName); err != nil {
			t.Fatal(err)
		}

		require.Equal(t, "gorp_migrations", migrationsTableName)
	})
}

func TestWithTestDatabaseEmpty(t *testing.T) {
	test.WithTestDatabaseEmpty(t, func(db1 *sql.DB) {
		test.WithTestDatabaseEmpty(t, func(db2 *sql.DB) {

			var db1Name string
			err := db1.QueryRow("SELECT current_database();").Scan(&db1Name)
			if err != nil {
				t.Fatal(err)
			}

			var db2Name string
			err = db2.QueryRow("SELECT current_database();").Scan(&db2Name)
			if err != nil {
				t.Fatal(err)
			}
			require.NotEqual(t, db1Name, db2Name)

			// test apply migrations + fixtures to a empty database 1
			_, err = test.ApplyMigrations(t, db1)
			require.NoError(t, err)

			_, err = test.ApplyTestFixtures(context.Background(), t, db1)
			require.NoError(t, err)

			// test apply dump to a empty database 2
			dumpFile := filepath.Join(pUtil.GetProjectRootDir(), "/test/testdata/users.sql")
			err = test.ApplyDump(context.Background(), t, db2, dumpFile)
			require.NoError(t, err)

			// check user count
			var usrCount int
			if err := db1.QueryRow("SELECT count(id) FROM users;").Scan(&usrCount); err != nil {
				t.Fatal(err)
			}

			require.Equal(t, 3, usrCount)
		})
	})
}
