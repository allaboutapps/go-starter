package cmd

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"allaboutapps.dev/aw/go-starter/internal/api"
	"allaboutapps.dev/aw/go-starter/internal/api/router"
	"allaboutapps.dev/aw/go-starter/internal/config"
	"allaboutapps.dev/aw/go-starter/internal/util"
	"allaboutapps.dev/aw/go-starter/internal/util/command"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

type ServerFlags struct {
	ProbeReadiness  bool
	ApplyMigrations bool
	SeedFixtures    bool
}

var serverFlags ServerFlags

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Starts the server",
	Long: `Starts the stateless RESTful JSON server

Requires configuration through ENV and
and a fully migrated PostgreSQL database.`,
	Run: func(_ *cobra.Command, _ []string) {
		runServer(serverFlags)
	},
}

func init() {
	serverCmd.Flags().BoolVarP(&serverFlags.ProbeReadiness, "probe", "p", false, "Probe readiness before startup.")
	serverCmd.Flags().BoolVarP(&serverFlags.ApplyMigrations, "migrate", "m", false, "Apply migrations before startup.")
	serverCmd.Flags().BoolVarP(&serverFlags.SeedFixtures, "seed", "s", false, "Seed fixtures into database before startup.")
	rootCmd.AddCommand(serverCmd)
}

func runServer(flags ServerFlags) {
	err := command.WithServer(context.Background(), config.DefaultServiceConfigFromEnv(), func(ctx context.Context, s *api.Server) error {
		log := util.LogFromContext(ctx)

		if flags.ProbeReadiness {
			errs, err := runReadiness(ctx, s.Config, ReadinessFlags{
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
			_, err := applyMigrations(ctx, s.Config)
			if err != nil {
				log.Fatal().Err(err).Msg("Error while applying migrations")
			}
		}

		if flags.SeedFixtures {
			err := applySeedFixtures(ctx, s.Config)
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

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := s.Shutdown(ctx); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal().Err(err).Msg("Failed to gracefully shut down server")
		}

		return nil
	})
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to start server")
	}
}
