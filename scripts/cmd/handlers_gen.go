//go:build scripts

package cmd

import (
	"log"

	"allaboutapps.dev/aw/go-starter/scripts/internal/handlers"
	"github.com/spf13/cobra"
)

const (
	printOnlyFlag = "print-only"
)

var handlersGenCmd = &cobra.Command{
	Use:   "gen",
	Short: "Generate internal/api/handlers/handlers.go.",
	Long:  `Generates internal/api/handlers/handlers.go file based on the current implemented handlers.`,
	Run: func(cmd *cobra.Command, _ []string /* args */) {

		printOnly, err := cmd.Flags().GetBool(printOnlyFlag)
		if err != nil {
			log.Fatal(err)
		}

		err = handlers.GenHandlers(printOnly)
		if err != nil {
			log.Fatal(err)
		}
	},
}

//nolint:gochecknoinits
func init() {
	handlersCmd.AddCommand(handlersGenCmd)
	handlersGenCmd.Flags().Bool(printOnlyFlag, false, "Print all checked handlers regardless of errors.")
}
