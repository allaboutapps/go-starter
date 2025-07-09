package probe

import (
	"allaboutapps.dev/aw/go-starter/internal/util/command"
	"github.com/spf13/cobra"
)

const (
	verboseFlag string = "verbose"
)

func New() *cobra.Command {
	return command.NewSubcommandGroup("probe",
		newLiveness(),
		newReadiness(),
	)
}
