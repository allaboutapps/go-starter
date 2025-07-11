//go:build scripts

package cmd

import (
	"os"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// handlersCmd represents the handlers command
// see handlers_*.go for sub_commands
var handlersCmd = &cobra.Command{
	Use:   "handlers <subcommand>",
	Short: "Handlers related subcommands",
	Run: func(cmd *cobra.Command, _ []string /* args */) {
		if err := cmd.Help(); err != nil {
			log.Error().Err(err).Msg("Failed to print help")
			os.Exit(1)
		}
		os.Exit(0)
	},
}

//nolint:gochecknoinits
func init() {
	rootCmd.AddCommand(handlersCmd)
}
