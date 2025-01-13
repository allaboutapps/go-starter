package probe

import (
	"context"
	"database/sql"
	"fmt"

	"allaboutapps.dev/aw/go-starter/internal/api"
	"allaboutapps.dev/aw/go-starter/internal/api/handlers/common"
	"allaboutapps.dev/aw/go-starter/internal/config"
	"allaboutapps.dev/aw/go-starter/internal/util"
	"allaboutapps.dev/aw/go-starter/internal/util/command"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

type LivenessFlags struct {
	Verbose bool
}

func newLiveness() *cobra.Command {
	var flags LivenessFlags

	cmd := &cobra.Command{
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
		Run: func(_ *cobra.Command, _ []string) {
			livenessCmdFunc(flags)
		},
	}

	cmd.Flags().BoolVarP(&flags.Verbose, verboseFlag, "v", false, "Show verbose output.")

	return cmd
}

func livenessCmdFunc(flags LivenessFlags) {
	err := command.WithServer(context.Background(), config.DefaultServiceConfigFromEnv(), func(ctx context.Context, s *api.Server) error {
		log := util.LogFromContext(ctx)

		errs, err := runLiveness(ctx, s.Config, flags)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to run liveness probes")
		}

		if len(errs) > 0 {
			log.Fatal().Errs("errs", errs).Msg("Unhealthy.")
		}

		return nil
	})
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to run liveness probes")
	}
}

func runLiveness(ctx context.Context, config config.Server, flags LivenessFlags) ([]error, error) {
	log := util.LogFromContext(ctx)

	db, err := sql.Open("postgres", config.Database.ConnectionString())
	if err != nil {
		log.Error().Err(err).Msg("Failed to connect to the database")
		return nil, err
	}
	defer db.Close()

	livenessCtx, cancel := context.WithTimeout(context.Background(), config.Management.LivenessTimeout)
	defer cancel()

	str, errs := common.ProbeLiveness(livenessCtx, db, config.Management.ProbeWriteablePathsAbs, config.Management.ProbeWriteableTouchfile)

	if flags.Verbose {
		fmt.Print(str)
	}

	return errs, nil
}
