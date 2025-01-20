package command

import (
	"context"
	"time"

	"allaboutapps.dev/aw/go-starter/internal/api"
	"allaboutapps.dev/aw/go-starter/internal/config"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	LogKeyCmdExecutionID = "cmdExecutionId"
)

func WithServer(ctx context.Context, config config.Server, f func(ctx context.Context, s *api.Server) error) error {
	ctx = log.With().Str(LogKeyCmdExecutionID, uuid.New().String()).Logger().WithContext(ctx)

	zerolog.TimeFieldFormat = time.RFC3339Nano
	zerolog.SetGlobalLevel(config.Logger.Level)
	if config.Logger.PrettyPrintConsole {
		log.Logger = log.Output(zerolog.NewConsoleWriter(func(w *zerolog.ConsoleWriter) {
			w.TimeFormat = "15:04:05"
		}))
	}

	s := api.NewServer(config).InitCmd()

	start := time.Now()

	err := f(ctx, s)

	elapsed := time.Since(start)
	log.Info().Dur("duration", elapsed).Msg("Command execution finished")

	if err != nil {
		log.Error().Err(err).Msg("Command failed")
		return err
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if errs := s.Shutdown(shutdownCtx); len(errs) > 0 {
		log.Fatal().Errs("shutdownErrors", errs).Msg("Failed to gracefully shut down server")
	}

	return nil
}
