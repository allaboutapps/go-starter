package test_test

import (
	"database/sql"
	"testing"

	"allaboutapps.dev/aw/go-starter/internal/test"
	"github.com/stretchr/testify/require"
)

func TestWithTestDatabase(t *testing.T) {

	t.Parallel()

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
