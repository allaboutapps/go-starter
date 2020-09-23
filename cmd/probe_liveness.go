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

// livenessCmd represents the server command
var livenessCmd = &cobra.Command{
	Use:   "liveness",
	Short: "Runs liveness probes",
	Long: `Runs connection livenesss probes

This command triggers the same livenesss probes as in
/-/healthy (apart from the actual server.ready 
probe) and prints the results to stdout. Fails with
non zero exitcode on encountered errors.

A typical usecase of this command are liveness probes 
to take action if dependant services (e.g. DB, NFS 
mounts) become unstable. You may also use this to 
ensure all requirements are fulfilled before starting
the app server.`,
	Run: func(cmd *cobra.Command, args []string) {

		verbose, err := cmd.Flags().GetBool(verboseFlag)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to parse args")
		}
		runLiveness(verbose)
	},
}

func init() {
	probeCmd.AddCommand(livenessCmd)
	livenessCmd.Flags().BoolP(verboseFlag, "v", false, "Show verbose output.")
}

func runLiveness(verbose bool) {
	config := config.DefaultServiceConfigFromEnv()

	db, err := sql.Open("postgres", config.Database.ConnectionString())
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to the database")
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), config.Management.LivenessTimeout)
	defer cancel()

	str, errs := common.ProbeLiveness(ctx, db, config.Management.ProbeWriteablePathsAbs, config.Management.ProbeWriteableTouchfile)

	if verbose {
		fmt.Print(str)
	}

	if len(errs) > 0 {
		log.Fatal().Errs("errs", errs).Msg("Unhealthy.")
	}
}
