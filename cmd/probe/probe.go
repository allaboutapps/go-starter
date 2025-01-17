package probe

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const (
	verboseFlag string = "verbose"
)

func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "probe <subcommand>",
		Short: "Probe related subcommands",
		Run: func(cmd *cobra.Command, _ []string /* args */) {
			if err := cmd.Help(); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			os.Exit(0)
		},
	}

	cmd.AddCommand(newLiveness(), newReadiness())

	return cmd
}
