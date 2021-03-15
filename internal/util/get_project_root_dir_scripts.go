// +build scripts

package util

// Note that VSCode/gopls currently spawns a "No packages found for open file: [...]" here.
// This is expected and will go away with gopls v1.0, see https://github.com/golang/go/issues/29202

import (
	"allaboutapps.dev/aw/go-starter/scripts"
)

// GetProjectRootDir() Returns the path as string to the project_root while **scripts generation**.
// Note: This function replaces the original util.GetProjectRootDir when go runs with the "script" build tag.
// https://stackoverflow.com/questions/43215655/building-multiple-binaries-using-different-packages-and-build-tags
func GetProjectRootDir() string {
	// delegate to the scripts pkg
	return scripts.GetProjectRootDir()
}
