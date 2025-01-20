package db

import (
	"context"
	"database/sql"
	"fmt"

	"allaboutapps.dev/aw/go-starter/internal/api"
	"allaboutapps.dev/aw/go-starter/internal/config"
	"allaboutapps.dev/aw/go-starter/internal/util"
	"allaboutapps.dev/aw/go-starter/internal/util/command"
	"github.com/rs/zerolog/log"
	migrate "github.com/rubenv/sql-migrate"
	"github.com/spf13/cobra"
)

func newMigrate() *cobra.Command {
	return &cobra.Command{
		Use:   "migrate",
		Short: "Executes all migrations which are not yet applied.",
		Run: func(_ *cobra.Command, _ []string) {
			migrateCmdFunc()
		},
	}
}

func init() {
	// pin migrate to use the globally defined `migrations` table identifier
	migrate.SetTable(config.DatabaseMigrationTable)
}

func migrateCmdFunc() {
	err := command.WithServer(context.Background(), config.DefaultServiceConfigFromEnv(), func(ctx context.Context, s *api.Server) error {
		log := util.LogFromContext(ctx)

		n, err := ApplyMigrations(ctx, s.Config)
		if err != nil {
			log.Err(err).Msg("Error while applying migrations")
			return err
		}

		log.Info().Int("appliedMigrationsCount", n).Msg("Successfully applied migrations")

		return nil
	})
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to apply migrations")
	}
}

func ApplyMigrations(ctx context.Context, serviceConfig config.Server) (int, error) {
	log := util.LogFromContext(ctx)

	db, err := sql.Open("postgres", serviceConfig.Database.ConnectionString())
	if err != nil {
		return 0, err
	}
	defer db.Close()

	if err := db.PingContext(ctx); err != nil {
		return 0, err
	}

	// In case an old default sql-migrate migration table (named "gorp_migrations") still exists we rename it to the new name equivalent
	// in sync with the settings in dbconfig.yml and config.DatabaseMigrationTable.
	if _, err := db.Exec(fmt.Sprintf("ALTER TABLE IF EXISTS gorp_migrations RENAME TO %s;", config.DatabaseMigrationTable)); err != nil {
		return 0, err
	}

	migrations := &migrate.FileMigrationSource{
		Dir: config.DatabaseMigrationFolder,
	}

	missingMigrations, _, err := migrate.PlanMigration(db, "postgres", migrations, migrate.Up, 0)
	if err != nil {
		log.Err(err).Msg("Error while planning migrations")
		return 0, err
	}

	var appliedMigrationsCount int
	for i := 0; i < len(missingMigrations); i++ {
		log.Info().Str("migrationId", missingMigrations[i].Id).Msg("Applying migration")

		n, err := migrate.ExecMax(db, "postgres", migrations, migrate.Up, 1)
		if err != nil {
			log.Err(err).Msg("Error while applying migration")
			return 0, err
		}

		log.Info().Int("appliedMigrationsCount", n).Msg("Applied migration")

		appliedMigrationsCount += n
		if n == 0 {
			break
		}
	}

	return appliedMigrationsCount, nil
}
