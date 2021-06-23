// +build scripts

package util

import "os"

func GetProjectRootDir() string {

	if val, ok := os.LookupEnv("PROJECT_ROOT_DIR"); ok {
		return val
	}

	return "/app"
}
