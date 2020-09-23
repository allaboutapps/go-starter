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

// readinessCmd represents the server command
var readinessCmd = &cobra.Command{
	Use:   "readiness",
	Short: "Runs readiness probes",
	Long: `Runs connection readinesss probes

This command triggers the same readinesss probes as in
/-/ready (apart from the actual server.ready 
probe) and prints the results to stdout. Fails with
non zero exitcode on encountered errors.

A typical usecase of this command are readiness probes 
to take action if dependant services (e.g. DB, NFS 
mounts) become unstable. You may also use this to 
ensure all requirements are fulfilled before starting
the app server.`,
	Run: func(cmd *cobra.Command, args []string) {

		verbose, err := cmd.Flags().GetBool(verboseFlag)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to parse args")
		}
		runReadiness(verbose)
	},
}

func init() {
	probeCmd.AddCommand(readinessCmd)
	readinessCmd.Flags().BoolP(verboseFlag, "v", false, "Show verbose output.")
}

func runReadiness(verbose bool) {
	config := config.DefaultServiceConfigFromEnv()

	db, err := sql.Open("postgres", config.Database.ConnectionString())
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to the database")
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), config.Management.ReadinessTimeout)
	defer cancel()

	str, errs := common.ProbeReadiness(ctx, db, config.Management.ProbeWriteablePathsAbs)

	if verbose {
		fmt.Print(str)
	}

	if len(errs) > 0 {
		log.Fatal().Errs("errs", errs).Msg("Unhealthy.")
	}
}
