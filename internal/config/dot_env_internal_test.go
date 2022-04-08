package config

import (
	"os"
	"path/filepath"
	"testing"

	"allaboutapps.dev/aw/go-starter/internal/util"
	"github.com/stretchr/testify/assert"
)

func TestDotEnvOverride(t *testing.T) {
	assert.Equal(t, "", os.Getenv("IS_THIS_A_TEST_ENV"))

	orgPsqlUser := os.Getenv("PSQL_USER")

	overrideEnv(filepath.Join(util.GetProjectRootDir(), "/internal/config/testdata/.env1.local"))
	assert.Equal(t, "yes", os.Getenv("IS_THIS_A_TEST_ENV"))
	assert.Equal(t, "dotenv_override_psql_user", os.Getenv("PSQL_USER"))
	assert.Equal(t, orgPsqlUser, os.Getenv("ORIGINAL_PSQL_USER"))

	// override works as expected?
	overrideEnv(filepath.Join(util.GetProjectRootDir(), "/internal/config/testdata/.env2.local"))
	assert.Equal(t, "yes still", os.Getenv("IS_THIS_A_TEST_ENV"))
	assert.NotEqual(t, "dotenv_override_psql_user", os.Getenv("PSQL_USER"))
	assert.Equal(t, orgPsqlUser, os.Getenv("PSQL_USER"), "Reset to original does not work!")
}

func TestNoopEnvNotFound(t *testing.T) {
	assert.NotPanics(t, assert.PanicTestFunc(func() {
		overrideEnv(filepath.Join(util.GetProjectRootDir(), "/internal/config/testdata/.env.does.not.exist"))
	}), "does not panic on file inexistance")
}

func TestEmptyEnv(t *testing.T) {
	assert.NotPanics(t, assert.PanicTestFunc(func() {
		overrideEnv(filepath.Join(util.GetProjectRootDir(), "/internal/config/testdata/.env.local.sample"))
	}), "does not panic on file inexistance")

	assert.Empty(t, os.Getenv("EMPTY_VARIABLE_INIT"), "should be empty")
}

func TestPanicsOnEnvMalform(t *testing.T) {
	assert.Panics(t, assert.PanicTestFunc(func() {
		overrideEnv(filepath.Join(util.GetProjectRootDir(), "/internal/config/testdata/.env.local.malformed"))
	}), "does panic on file malform")

	SetEnvFromFile(filepath.Join(util.GetProjectRootDir(), "/internal/config/testdata/.env.does.not.exist"), t.Setenv)
}
