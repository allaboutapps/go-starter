// +build scripts

package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// handlersCmd represents the handlers command
// see handlers_*.go for sub_commands
var handlersCmd = &cobra.Command{
	Use:   "handlers <subcommand>",
	Short: "Handlers related subcommands",
	Run: func(cmd *cobra.Command, args []string) {
		if err := cmd.Help(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		os.Exit(0)
	},
}

func init() {
	rootCmd.AddCommand(handlersCmd)
}
