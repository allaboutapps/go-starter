package env

import (
	"encoding/json"
	"fmt"

	"allaboutapps.dev/aw/go-starter/internal/config"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func New() *cobra.Command {
	return &cobra.Command{
		Use:   "env",
		Short: "Prints the env",
		Long: `Prints the currently applied env
	
	You may use this cmd to get an overview about how 
	your ENV_VARS are bound by the server config.
	Please note that certain secrets are automatically
	removed from this output.`,
		Run: func(_ *cobra.Command /* cmd */, _ []string /* args */) {
			runEnv()
		},
	}
}

func runEnv() {
	config := config.DefaultServiceConfigFromEnv()

	result, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to marshal the env")
	}

	//nolint:forbidigo
	fmt.Println(string(result))
}
