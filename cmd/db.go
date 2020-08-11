package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// dbCmd represents the db command
// see db_*.go for sub_commands
var dbCmd = &cobra.Command{
	Use:   "db <subcommand>",
	Short: "Database related subcommands",
	Run: func(cmd *cobra.Command, args []string) {
		if err := cmd.Help(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		os.Exit(0)
	},
}

func init() {
	rootCmd.AddCommand(dbCmd)
}
