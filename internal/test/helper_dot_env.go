package test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"allaboutapps.dev/aw/go-starter/internal/config"
	"allaboutapps.dev/aw/go-starter/internal/util"
)

// DotEnvLoadLocalOrSkipTest tries to load the `.env.local` file in the projectroot
// overwriting ENV vars. If the file is not found, the test will be automatically skipped.
func DotEnvLoadLocalOrSkipTest(t *testing.T) {
	t.Helper()

	absolutePathToEnvFile := filepath.Join(util.GetProjectRootDir(), ".env.local")
	DotEnvLoadFileOrSkipTest(t, absolutePathToEnvFile)
}

// DotEnvLoadFileOrSkipTest tries to load the overgiven path to the dotenv file overwriting
// ENV vars. If the file is not found, the test will be automatically skipped.
func DotEnvLoadFileOrSkipTest(t *testing.T, absolutePathToEnvFile string) {
	t.Helper()

	// this test should be automatically skipped if no default `.env.local` file was found.
	err := config.DotEnvLoad(
		absolutePathToEnvFile,
		func(k string, v string) error { t.Setenv(k, v); return nil })

	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			t.Skip(absolutePathToEnvFile, "not found, skipping test.")
		} else {
			t.Fatal(err)
		}
	} else {
		t.Log(absolutePathToEnvFile, "override ENV variables!")
	}
}
