package command

import (
	"context"
	"errors"
	"fmt"
	"time"

	"allaboutapps.dev/aw/go-starter/internal/api"
	"allaboutapps.dev/aw/go-starter/internal/config"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

const (
	LogKeyCmdExecutionID = "cmdExecutionId"
	shutdownTimeout      = 30 * time.Second
)

func WithServer(ctx context.Context, config config.Server, handler func(ctx context.Context, s *api.Server) error) error {
	ctx = log.With().Str(LogKeyCmdExecutionID, uuid.New().String()).Logger().WithContext(ctx)

	zerolog.TimeFieldFormat = time.RFC3339Nano
	zerolog.SetGlobalLevel(config.Logger.Level)
	if config.Logger.PrettyPrintConsole {
		log.Logger = log.Output(zerolog.NewConsoleWriter(func(w *zerolog.ConsoleWriter) {
			w.TimeFormat = "15:04:05"
		}))
	}

	s, err := api.InitNewServer(config)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize server")
	}

	start := s.Clock.Now()

	err = handler(ctx, s)

	elapsed := time.Since(start)
	log.Info().Dur("duration", elapsed).Msg("Command execution finished")

	if err != nil {
		log.Error().Err(err).Msg("Command failed")
		return err
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	if errs := s.Shutdown(shutdownCtx); len(errs) > 0 {
		log.Error().Errs("shutdownErrors", errs).Msg("Failed to gracefully shut down server")
		return errors.Join(errs...)
	}

	return nil
}

func NewSubcommandGroup(subcommand string, subcommands ...*cobra.Command) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("%s <subcommand>", subcommand),
		Short: fmt.Sprintf("%s related subcommands", subcommand),
		RunE: func(cmd *cobra.Command, _ []string) error {
			if err := cmd.Help(); err != nil {
				return fmt.Errorf("failed to print help: %w", err)
			}

			return nil
		},
	}

	cmd.AddCommand(subcommands...)

	return cmd
}
