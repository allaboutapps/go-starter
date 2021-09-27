//go:build scripts

package scaffold_test

import (
	"os"
	"path/filepath"
	"testing"

	"allaboutapps.dev/aw/go-starter/internal/test"
	"allaboutapps.dev/aw/go-starter/internal/util"
	"allaboutapps.dev/aw/go-starter/scripts/internal/scaffold"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	resourcePath    = "testdata/test_resource.txt"
	definitionsPath = "testdata/definitions"
	pathsPath       = "testdata/paths"
	handlerPath     = "testdata/testresource"
	modulePath      = filepath.Join(util.GetProjectRootDir(), "go.mod")

	defaultMethods = []string{"get-all", "get", "post", "put", "delete"}
)

func TestParseModel_Success(t *testing.T) {
	// Execute
	resource, err := scaffold.ParseModel(resourcePath)

	// Assert
	require.NoError(t, err)
	test.Snapshoter.Save(t, resource)
}

func TestGenerateSwagger_Success(t *testing.T) {
	// Setup
	resource, err := scaffold.ParseModel(resourcePath)
	require.NoError(t, err)

	// Execute
	err = scaffold.GenerateSwagger(resource, "testdata", true)

	// Assert
	require.NoError(t, err)

	assert.DirExists(t, definitionsPath, "Should create the definitions directory")
	assert.DirExists(t, pathsPath, "Should create the paths directory")
	assert.FileExists(t, filepath.Join(definitionsPath, "testresource.yml"), "Should generate the definition spec")
	assert.FileExists(t, filepath.Join(pathsPath, "testresource.yml"), "Should generate the paths spec")

	// Cleanup
	err = os.RemoveAll(definitionsPath)
	require.NoError(t, err)
	err = os.RemoveAll(pathsPath)
	require.NoError(t, err)
}

func TestGenerateHandlers_Success(t *testing.T) {
	// Setup
	resource, err := scaffold.ParseModel(resourcePath)
	require.NoError(t, err)
	err = scaffold.GenerateSwagger(resource, "testdata", true)
	require.NoError(t, err)

	// Execute
	err = scaffold.GenerateHandlers(resource, "testdata", modulePath, defaultMethods, true)

	// Assert
	require.NoError(t, err)

	assert.DirExists(t, handlerPath, "Should create the handler directory")
	assert.FileExists(t, filepath.Join(handlerPath, "get_testresource_list.go"), "Should create the GET list handler")
	assert.FileExists(t, filepath.Join(handlerPath, "get_testresource.go"), "Should create the GET handler")
	assert.FileExists(t, filepath.Join(handlerPath, "post_testresource.go"), "Should create the POST handler")
	assert.FileExists(t, filepath.Join(handlerPath, "put_testresource.go"), "Should create the PUT handler")
	assert.FileExists(t, filepath.Join(handlerPath, "delete_testresource.go"), "Should create the DELETE handler")

	// Cleanup
	err = os.RemoveAll(definitionsPath)
	require.NoError(t, err)
	err = os.RemoveAll(pathsPath)
	require.NoError(t, err)
	err = os.RemoveAll(handlerPath)
	require.NoError(t, err)
}
