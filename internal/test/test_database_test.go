package test_test

import (
	"context"
	"database/sql"
	"path/filepath"
	"testing"

	"allaboutapps.dev/aw/go-starter/internal/test"
	pUtil "allaboutapps.dev/aw/go-starter/internal/util"
	"github.com/stretchr/testify/require"
)

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

	dumpFile := filepath.Join(pUtil.GetProjectRootDir(), "/test/testdata/minimal.sql")

	test.WithTestDatabaseFromDump(t, dumpFile, func(db1 *sql.DB) {
		test.WithTestDatabaseFromDump(t, dumpFile, func(db2 *sql.DB) {

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
			dumpFile := filepath.Join(pUtil.GetProjectRootDir(), "/test/testdata/minimal.sql")
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
