//go:build scripts

package util

import (
	"fmt"
	"log"
	"os"

	"golang.org/x/mod/modfile"
)

// GetModuleName returns the current go module's name as defined in the go.mod file.
// https://stackoverflow.com/questions/53183356/api-to-get-the-module-name
// https://github.com/rogpeppe/go-internal
func GetModuleName(absolutePathToGoMod string) (string, error) {
	dat, err := os.ReadFile(absolutePathToGoMod)

	if err != nil {
		log.Fatal(err)
		return "", fmt.Errorf("failed to read go.mod file: %w", err)
	}

	mod := modfile.ModulePath(dat)

	return mod, nil
}
