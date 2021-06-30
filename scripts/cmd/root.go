// +build scripts

package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"allaboutapps.dev/aw/go-starter/scripts/internal/util"
	"github.com/spf13/cobra"
)

var (
	projectRoot = util.GetProjectRootDir()
	modulePath  = filepath.Join(util.GetProjectRootDir(), "go.mod")
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gsdev",
	Short: "gsdev",
	Long: `go-starter development scripts
Utility commands while developing go-starter based projects.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
