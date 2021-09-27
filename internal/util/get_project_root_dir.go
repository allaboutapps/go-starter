//go:build !scripts

package util

import (
	"os"
	"path/filepath"
	"sync"

	"github.com/rs/zerolog/log"
)

var (
	projectRootDir string
	dirOnce        sync.Once
)

// GetProjectRootDir returns the path as string to the project_root for a **running application**.
// Note: This function should not be used for generation targets (go generate, make go-generate).
// Thus it's explicitly excluded from the build tag scripts, see instead:
// * /scripts/get_project_root_dir.go
// * ./get_project_root_dir_scripts.go (delegates to above)
// https://stackoverflow.com/questions/43215655/building-multiple-binaries-using-different-packages-and-build-tags
func GetProjectRootDir() string {
	dirOnce.Do(func() {
		ex, err := os.Executable()
		if err != nil {
			log.Panic().Err(err).Msg("Failed to get executable path while retrieving project root directory")
		}

		projectRootDir = GetEnv("PROJECT_ROOT_DIR", filepath.Dir(ex))
	})

	return projectRootDir
}
