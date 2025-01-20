package db

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "db <subcommand>",
		Short: "Database related subcommands",
		Run: func(cmd *cobra.Command, _ []string /* args */) {
			if err := cmd.Help(); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			os.Exit(0)
		},
	}

	cmd.AddCommand(newMigrate(), newSeed())

	return cmd
}
