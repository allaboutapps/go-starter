package test

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"allaboutapps.dev/aw/go-starter/internal/util"
	"github.com/stretchr/testify/require"
)

// PrepareTestFile copies a test file by name from test/testdata/<filename> to a folder unique to the test.
// Used for tests with document to load reference files by fixtures.
func PrepareTestFile(t *testing.T, fileName string, destFileName ...string) {
	t.Helper()

	src, err := os.Open(filepath.Join(util.GetProjectRootDir(), "test", "testdata", fileName))
	require.NoError(t, err)
	defer src.Close()

	dest := fileName
	if len(destFileName) > 0 {
		dest = destFileName[0]
	}

	path := filepath.Join(util.GetProjectRootDir(), "assets", "mnt", strings.ToLower(t.Name()), "documents", dest)
	err = os.MkdirAll(filepath.Dir(path), 0755)
	require.NoError(t, err)

	dst, err := os.Create(path)
	require.NoError(t, err)
	defer dst.Close()

	_, err = io.Copy(dst, src)
	require.NoError(t, err)
}

// CleanupTestFiles removes folder unique to the test if exists
func CleanupTestFiles(t *testing.T) {
	t.Helper()

	err := os.RemoveAll(filepath.Join(util.GetProjectRootDir(), "assets", "mnt", strings.ToLower(t.Name())))
	require.NoError(t, err)
}

// WithTempDir creates a folder unique to the tests and ensures cleanup of the folder will be
// performed after the fn got called.
func WithTempDir(t *testing.T, fn func(localBasePath string, basePath string)) {
	t.Helper()

	localBasePath := filepath.Join(util.GetProjectRootDir(), "assets", "mnt", strings.ToLower(t.Name()))
	basePath := "/documents"
	path := filepath.Join(localBasePath, basePath)
	err := os.MkdirAll(filepath.Dir(path), 0755)
	require.NoError(t, err)

	defer CleanupTestFiles(t)

	fn(localBasePath, basePath)
}
