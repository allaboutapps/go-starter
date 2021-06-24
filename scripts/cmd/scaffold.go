// +build scripts

package cmd

import (
	"log"
	"os/exec"
	"path/filepath"

	"allaboutapps.dev/aw/go-starter/scripts/internal/scaffold"
	"github.com/spf13/cobra"
)

const (
	methodsFlag = "methods"
	forceFlag   = "force"
)

var (
	modelPath   = filepath.Join(projectRoot, "internal/models")
	swaggerPath = filepath.Join(projectRoot, "api")
	handlerPath = filepath.Join(projectRoot, "internal/api/handlers")

	makeSwaggerCmd    = exec.Command("make", "swagger")
	makeGoGenerateCmd = exec.Command("make", "go-generate")

	defaultMethods = []string{"get-all", "get", "post", "put", "delete"}
)

// rootCmd represents the base command when called without any subcommands
var scaffoldCmd = &cobra.Command{
	Use:   "scaffold [resource name]",
	Short: "Scaffolding tool for CRUD endpoints.",
	Long:  "Scaffolding tool to generate CRUD endpoint stubs from sqlboiler model definitions.",
	Run:   generate,
}

func init() {
	rootCmd.AddCommand(scaffoldCmd)
	scaffoldCmd.Flags().StringSliceP(methodsFlag, "m", defaultMethods, "Specify HTTP methods to generate handlers for. Example: scaffold --methods get-all,get,delete")
	scaffoldCmd.Flags().BoolP(forceFlag, "f", false, "Forces the tool to overwrite existing files.")
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

	force, err := cmd.Flags().GetBool(forceFlag)
	if err != nil {
		log.Fatalf("Failed to read %s: %v", forceFlag, err)
	}

	resourcePath := filepath.Join(modelPath, resourceName+".go")
	resource, err := scaffold.ParseModel(resourcePath)
	if err != nil {
		log.Fatalf("Failed to parse resource from model file: %v", err)
	}

	if err = scaffold.GenerateSwagger(resource, swaggerPath, force); err != nil {
		log.Fatalf("Failed to generate Swagger spec: %v", err)
	}

	if err = scaffold.GenerateHandlers(resource, handlerPath, modulePath, methods, force); err != nil {
		log.Fatalf("Failed to generate handlers: %v", err)
	}

	if err = makeSwaggerCmd.Run(); err != nil {
		log.Fatalf("Failed to run 'make swagger': %v", err)
	}

	if err = makeGoGenerateCmd.Run(); err != nil {
		log.Fatalf("Failed to run 'make go-generate': %v", err)
	}
}
