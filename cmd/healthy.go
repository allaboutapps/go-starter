package cmd

import (
	"context"
	"database/sql"
	"fmt"

	"allaboutapps.dev/aw/go-starter/internal/api/handlers/common"
	"allaboutapps.dev/aw/go-starter/internal/config"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var (
	verboseFlag string = "verbose"
)

// healthyCmd represents the server command
var healthyCmd = &cobra.Command{
	Use:   "healthy",
	Short: "Runs conn healthy checks",
	Long: `Runs connection healthy checks

This command triggers the same healthy check as in
/-/healthy (apart from the actual server readiness 
probe) and prints the results to stdout. Fails with
non zero exitcode on encountered errors.

A typical usecase of this command are live-/readiness
probes to take action if dependant services (e.g. DB,
NFS mounts) become unstable. You may also use this to
ensure all requirements are fulfilled before starting
the app server.`,
	Run: func(cmd *cobra.Command, args []string) {

		verbose, err := cmd.Flags().GetBool(verboseFlag)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to parse args")
		}
		runHealthy(verbose)
	},
}

func init() {
	rootCmd.AddCommand(healthyCmd)
	healthyCmd.Flags().BoolP(verboseFlag, "v", false, "Show verbose output.")
}

func runHealthy(verbose bool) {
	config := config.DefaultServiceConfigFromEnv()

	db, err := sql.Open("postgres", config.Database.ConnectionString())
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to the database")
	}

	ctx, cancel := context.WithTimeout(context.Background(), config.Management.HealthyTimeout)
	defer cancel()

	str, errs := common.CheckHealthy(ctx, db, config.Management.HealthyCheckWriteablePathsAbs, config.Management.HealthyCheckWriteablePathsTouch)

	if verbose {
		fmt.Print(str)
	}

	if len(errs) > 0 {
		log.Fatal().Errs("errs", errs).Msg("Unhealthy.")
	}
}
