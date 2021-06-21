// +build scripts

package cmd

import (
	"log"

	"allaboutapps.dev/aw/go-starter/scripts/internal/handlers"
	"github.com/spf13/cobra"
)

var (
	printAllFlag string = "print-all"
)

var handlersGenCmd = &cobra.Command{
	Use:   "gen",
	Short: "Generate internal/api/handlers/handlers.go.",
	Long:  `Generates internal/api/handlers/handlers.go file based on the current implemented handlers.`,
	Run: func(cmd *cobra.Command, args []string) {

		printAll, err := cmd.Flags().GetBool(printAllFlag)
		if err != nil {
			log.Fatal(err)
		}
		handlers.CheckHandlers(printAll)
	},
}

func init() {
	handlersCmd.AddCommand(handlersGenCmd)
	handlersGenCmd.Flags().Bool(printAllFlag, false, "Print all checked handlers regardless of errors.")
}
