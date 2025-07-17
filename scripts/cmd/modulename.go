//go:build scripts

package cmd

import (
	"allaboutapps.dev/aw/go-starter/scripts/internal/util"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// moduleCmd represents the server command
var moduleCmd = &cobra.Command{
	Use:   "modulename",
	Short: "Prints the modulename",
	Long:  `Prints the currently applied go modulename of this project.`,
	Run: func(_ *cobra.Command /* cmd */, _ []string /* args */) {
		runModulename()
	},
}

//nolint:gochecknoinits
func init() {
	rootCmd.AddCommand(moduleCmd)
}

func runModulename() {
	baseModuleName, err := util.GetModuleName(modulePath)

	if err != nil {
		log.Fatal().Err(err).Msg("Failed to get module name")
	}

	log.Info().Msg(baseModuleName)
}
