package persistence_test

import (
	"database/sql"
	"testing"

	"allaboutapps.dev/aw/go-starter/internal/config"
	"allaboutapps.dev/aw/go-starter/internal/test"
	migrate "github.com/rubenv/sql-migrate"
	"github.com/stretchr/testify/require"
)

func TestMigrations(t *testing.T) {
	test.WithTestDatabase(t, func(db *sql.DB) {
		migrate.SetTable(config.DatabaseMigrationTable)

		// use the migrations from the migrations folder
		migrationSource := &migrate.FileMigrationSource{Dir: config.DatabaseMigrationFolder}

		// expect the migrations to already be applied to the database
		missingUpMigrations, _, err := migrate.PlanMigration(db, "postgres", migrationSource, migrate.Up, 0)
		require.NoError(t, err)
		require.Empty(t, missingUpMigrations)

		// run all down migrations
		downMigrations, _, err := migrate.PlanMigration(db, "postgres", migrationSource, migrate.Down, 0)
		require.NoError(t, err)
		require.NotEmpty(t, downMigrations)

		for _, migration := range downMigrations {
			n, err := migrate.ExecMax(db, "postgres", migrationSource, migrate.Down, 1)
			require.NoError(t, err, "failed to apply down migration %s", migration.Id)
			require.Equal(t, 1, n, "expected 1 migration to be applied for %s, got %d", migration.Id, n)
		}

		// run all up migrations again to test if the down migrations cleanup the database correctly
		upMigrations, _, err := migrate.PlanMigration(db, "postgres", migrationSource, migrate.Up, 0)
		require.NoError(t, err)
		require.NotEmpty(t, upMigrations)

		for _, migration := range upMigrations {
			n, err := migrate.ExecMax(db, "postgres", migrationSource, migrate.Up, 1)
			require.NoError(t, err, "failed to apply up migration %s after down migrations", migration.Id)
			require.Equal(t, 1, n, "expected 1 migration to be applied for %s after down migrations, got %d", migration.Id, n)
		}
	})
}
