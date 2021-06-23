// +build scripts

package cmd

import (
	"fmt"
	"log"

	"allaboutapps.dev/aw/go-starter/scripts/internal/util"
	"github.com/spf13/cobra"
)

// moduleCmd represents the server command
var moduleCmd = &cobra.Command{
	Use:   "modulename",
	Short: "Prints the modulename",
	Long:  `Prints the currently applied go modulename of this project.`,
	Run: func(cmd *cobra.Command, args []string) {
		runModulename()
	},
}

var (
	PROJECT_ROOT  = util.GetProjectRootDir()
	PATH_MOD_FILE = PROJECT_ROOT + "/go.mod"
)

func init() {
	rootCmd.AddCommand(moduleCmd)
}

func runModulename() {
	baseModuleName, err := util.GetModuleName(PATH_MOD_FILE)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(baseModuleName)
}
