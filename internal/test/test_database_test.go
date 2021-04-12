package test_test

import (
	"context"
	"database/sql"
	"path/filepath"
	"testing"

	"allaboutapps.dev/aw/go-starter/internal/models"
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
	test.WithTestDatabaseFromDump(t, filepath.Join(pUtil.GetProjectRootDir(), "/test/testdata/test.db"), func(db1 *sql.DB) {
		test.WithTestDatabaseFromDump(t, filepath.Join(pUtil.GetProjectRootDir(), "/test/testdata/test.db"), func(db2 *sql.DB) {

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

			userCount1, err := models.Users().Count(context.Background(), db1)
			require.NoError(t, err)
			require.Equal(t, int64(1), userCount1)

			userCount2, err := models.Users().Count(context.Background(), db2)
			require.NoError(t, err)
			require.Equal(t, int64(1), userCount2)
		})
	})
}
