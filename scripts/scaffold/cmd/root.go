package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"allaboutapps.dev/aw/go-starter/internal/util"
	"allaboutapps.dev/aw/go-starter/scripts/scaffold"
	"github.com/spf13/cobra"
)

const (
	methodsFlag = "methods"
	forceFlag   = "force"
)

var (
	projectRoot = util.GetProjectRootDir()
	modelPath   = filepath.Join(projectRoot, "internal/models")
	swaggerPath = filepath.Join(projectRoot, "api")
	handlerPath = filepath.Join(projectRoot, "internal/api/handlers")
	modulePath  = filepath.Join(util.GetProjectRootDir(), "go.mod")

	makeSwaggerCmd    = exec.Command("make", "swagger")
	makeGoGenerateCmd = exec.Command("make", "go-generate")

	defaultMethods = []string{"get-all", "get", "post", "put", "delete"}
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "scaffold [resource name]",
	Short: "Scaffolding tool for CRUD endpoints.",
	Long:  "Scaffolding tool to generate CRUD endpoint stubs from sqlboiler model definitions.",
	Run:   generate,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	rootCmd.Flags().StringSliceP(methodsFlag, "m", defaultMethods, "Specify HTTP methods to generate handlers for. Example: scaffold --methods get-all,get,delete")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func generate(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		log.Fatal("Please provide a valid resource name")
	}

	resourceName := args[0]
	if resourceName == "" {
		log.Fatal("Please provide a valid resource name")
	}

	methods, err := cmd.Flags().GetStringSlice(methodsFlag)
	if err != nil {
		log.Fatalf("Failed to read %s: %v", methodsFlag, err)
	}

	resourcePath := filepath.Join(modelPath, resourceName+".go")
	resource, err := scaffold.ParseModel(resourcePath)
	if err != nil {
		log.Fatalf("Failed to parse resource from model file: %v", err)
	}

	if err = scaffold.GenerateSwagger(resource, swaggerPath); err != nil {
		log.Fatalf("Failed to generate Swagger spec: %v", err)
	}

	if err = scaffold.GenerateHandlers(resource, handlerPath, modulePath, methods); err != nil {
		log.Fatalf("Failed to generate handlers: %v", err)
	}

	if err = makeSwaggerCmd.Run(); err != nil {
		log.Fatalf("Failed to run 'make swagger': %v", err)
	}

	if err = makeGoGenerateCmd.Run(); err != nil {
		log.Fatalf("Failed to run 'make go-generate': %v", err)
	}
}
