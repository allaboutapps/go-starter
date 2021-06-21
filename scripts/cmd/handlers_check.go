// +build scripts

package cmd

import (
	"log"

	"allaboutapps.dev/aw/go-starter/scripts/internal/handlers"
	"github.com/spf13/cobra"
)

var (
	printOnlyFlag string = "print-only"
)

var handlersCheckCmd = &cobra.Command{
	Use:   "check",
	Short: "Checks handlers vs. swagger spec.",
	Long:  `Checks currently implemented handlers against swagger spec.`,
	Run: func(cmd *cobra.Command, args []string) {

		printOnly, err := cmd.Flags().GetBool(printOnlyFlag)
		if err != nil {
			log.Fatal(err)
		}
		handlers.GenHandlers(printOnly)
	},
}

func init() {
	handlersCmd.AddCommand(handlersCheckCmd)
	handlersCheckCmd.Flags().Bool(printOnlyFlag, false, "Print only print the current implemented handlers, do not generate the file.")
}
