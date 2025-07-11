//go:build scripts

package cmd

import (
	"log"

	"allaboutapps.dev/aw/go-starter/scripts/internal/handlers"
	"github.com/spf13/cobra"
)

const (
	printAllFlag = "print-all"
)

var handlersCheckCmd = &cobra.Command{
	Use:   "check",
	Short: "Checks handlers vs. swagger spec.",
	Long:  `Checks currently implemented handlers against swagger spec.`,
	Run: func(cmd *cobra.Command, _ []string /* args */) {

		printAll, err := cmd.Flags().GetBool(printAllFlag)
		if err != nil {
			log.Fatal(err)
		}
		err = handlers.CheckHandlers(printAll)
		if err != nil {
			log.Fatal(err)
		}
	},
}

//nolint:gochecknoinits
func init() {
	handlersCmd.AddCommand(handlersCheckCmd)
	handlersCheckCmd.Flags().Bool(printAllFlag, false, "Print only print the current implemented handlers, do not generate the file.")
}
