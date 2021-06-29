package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const (
	verboseFlag string = "verbose"
)

// probeCmd represents the probe command
// see probe_*.go for sub_commands
var probeCmd = &cobra.Command{
	Use:   "probe <subcommand>",
	Short: "Probe related subcommands",
	Run: func(cmd *cobra.Command, args []string) {
		if err := cmd.Help(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		os.Exit(0)
	},
}

func init() {
	rootCmd.AddCommand(probeCmd)
}
