// +build scripts

package util

import "os"

// Note that VSCode/gopls currently spawns a "No packages found for open file: [...]" here.
// This is expected and will go away with gopls v1.0, see https://github.com/golang/go/issues/29202

// GetProjectRootDir returns the path as string to the project_root while **scripts generation**.
// Note: This function replaces the original util.GetProjectRootDir when go runs with the "script" build tag.
// https://stackoverflow.com/questions/43215655/building-multiple-binaries-using-different-packages-and-build-tags
// Should be in sync with "scripts/internal/util/get_project_root_dir.go"
func GetProjectRootDir() string {
	if val, ok := os.LookupEnv("PROJECT_ROOT_DIR"); ok {
		return val
	}

	return "/app"
}
