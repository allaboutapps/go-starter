// Scaffolding tool to auto-generate basic CRUD handlers for a given database model.

package main

import (
	"log"
	"os"
	"path/filepath"

	"allaboutapps.dev/aw/go-starter/internal/util"
	"allaboutapps.dev/aw/go-starter/scripts/scaffold"
)

var (
	projectRoot = util.GetProjectRootDir()
	modelPath   = filepath.Join(projectRoot, "internal/models")
	swaggerPath = filepath.Join(projectRoot, "api")
	handlerPath = filepath.Join(projectRoot, "internal/api/handlers")
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Please provide a valid resource name")
	}

	resourceName := os.Args[1]
	if resourceName == "" {
		log.Fatal("Please provide a valid resource name")
	}

	resourcePath := filepath.Join(modelPath, resourceName+".go")
	resource, err := scaffold.ParseModel(resourcePath)
	if err != nil {
		log.Fatalf("Failed to parse resource from model file: %v", err)
	}

	if err = scaffold.GenerateSwagger(resource, swaggerPath); err != nil {
		log.Fatalf("Failed generate Swagger spec: %v", err)
	}

	if err = scaffold.GenerateHandlers(resource, handlerPath); err != nil {
		log.Fatalf("Failed generate handlers: %v", err)
	}
}
