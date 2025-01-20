package db

import (
	"context"
	"database/sql"

	"allaboutapps.dev/aw/go-starter/internal/api"
	"allaboutapps.dev/aw/go-starter/internal/config"
	"allaboutapps.dev/aw/go-starter/internal/data"
	"allaboutapps.dev/aw/go-starter/internal/util"
	"allaboutapps.dev/aw/go-starter/internal/util/command"
	dbutil "allaboutapps.dev/aw/go-starter/internal/util/db"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

func newSeed() *cobra.Command {
	return &cobra.Command{
		Use:   "seed",
		Short: "Inserts or updates fixtures to the database.",
		Long:  `Uses upsert to add test data to the current environment.`,
		Run: func(_ *cobra.Command, _ []string) {
			seedCmdFunc()
		},
	}
}

func seedCmdFunc() {
	err := command.WithServer(context.Background(), config.DefaultServiceConfigFromEnv(), func(ctx context.Context, s *api.Server) error {
		log := util.LogFromContext(ctx)

		err := ApplySeedFixtures(ctx, s.Config)
		if err != nil {
			log.Err(err).Msg("Error while applying seed fixtures")
			return err
		}

		log.Info().Msg("Successfully applied seed fixtures")

		return nil
	})
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to apply migrations")
	}
}

func ApplySeedFixtures(ctx context.Context, config config.Server) error {
	log := util.LogFromContext(ctx)

	db, err := sql.Open("postgres", config.Database.ConnectionString())
	if err != nil {
		return err
	}

	defer db.Close()

	if err := db.PingContext(ctx); err != nil {
		return err
	}

	// insert fixtures in an auto-managed db transaction
	return dbutil.WithTransaction(ctx, db, func(tx boil.ContextExecutor) error {

		fixtures := data.Upserts()

		for _, fixture := range fixtures {
			if err := fixture.Upsert(ctx, tx, true, nil, boil.Infer(), boil.Infer()); err != nil {
				log.Error().Err(err).Msg("Failed to upsert fixture")
				return err
			}
		}

		log.Info().Int("fixturesCount", len(fixtures)).Msg("Successfully upserted fixtures")
		return nil
	})
}
