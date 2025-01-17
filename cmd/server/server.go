package server

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"allaboutapps.dev/aw/go-starter/cmd/db"
	"allaboutapps.dev/aw/go-starter/cmd/probe"
	"allaboutapps.dev/aw/go-starter/internal/api"
	"allaboutapps.dev/aw/go-starter/internal/api/router"
	"allaboutapps.dev/aw/go-starter/internal/config"
	"allaboutapps.dev/aw/go-starter/internal/util"
	"allaboutapps.dev/aw/go-starter/internal/util/command"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

type Flags struct {
	ProbeReadiness  bool
	ApplyMigrations bool
	SeedFixtures    bool
}

func New() *cobra.Command {
	var flags Flags

	cmd := &cobra.Command{
		Use:   "server",
		Short: "Starts the server",
		Long: `Starts the stateless RESTful JSON server
	
	Requires configuration through ENV and
	and a fully migrated PostgreSQL database.`,
		Run: func(_ *cobra.Command, _ []string) {
			runServer(flags)
		},
	}

	cmd.Flags().BoolVarP(&flags.ProbeReadiness, "probe", "p", false, "Probe readiness before startup.")
	cmd.Flags().BoolVarP(&flags.ApplyMigrations, "migrate", "m", false, "Apply migrations before startup.")
	cmd.Flags().BoolVarP(&flags.SeedFixtures, "seed", "s", false, "Seed fixtures into database before startup.")

	return cmd
}

func runServer(flags Flags) {
	err := command.WithServer(context.Background(), config.DefaultServiceConfigFromEnv(), func(ctx context.Context, s *api.Server) error {
		log := util.LogFromContext(ctx)

		if flags.ProbeReadiness {
			errs, err := probe.RunReadiness(ctx, s.Config, probe.ReadinessFlags{
				Verbose: true,
			})
			if err != nil {
				log.Fatal().Err(err).Msg("Failed to run readiness probes")
			}

			if len(errs) > 0 {
				log.Fatal().Errs("errs", errs).Msg("Unhealthy.")
			}

		}

		if flags.ApplyMigrations {
			_, err := db.ApplyMigrations(ctx, s.Config)
			if err != nil {
				log.Fatal().Err(err).Msg("Error while applying migrations")
			}
		}

		if flags.SeedFixtures {
			err := db.ApplySeedFixtures(ctx, s.Config)
			if err != nil {
				log.Fatal().Err(err).Msg("Error while applying seed fixtures")
			}
		}

		router.Init(s)

		go func() {
			if err := s.Start(); err != nil {
				if errors.Is(err, http.ErrServerClosed) {
					log.Info().Msg("Server closed")
				} else {
					log.Fatal().Err(err).Msg("Failed to start server")
				}
			}
		}()

		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
		<-quit

		return nil
	})
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to start server")
	}
}
