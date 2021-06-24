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

func PrepareTestFile(t *testing.T, fileName string) {
	t.Helper()

	src, err := os.Open(filepath.Join(util.GetProjectRootDir(), "test", "testdata", fileName))
	require.NoError(t, err)
	defer src.Close()

	path := filepath.Join(util.GetProjectRootDir(), "assets", "mnt", strings.ToLower(t.Name()), "documents", fileName)
	err = os.MkdirAll(filepath.Dir(path), 0755)
	require.NoError(t, err)

	dst, err := os.Create(path)
	require.NoError(t, err)
	defer dst.Close()

	_, err = io.Copy(dst, src)
	require.NoError(t, err)
}

func CleanupTestFiles(t *testing.T) {
	t.Helper()

	err := os.RemoveAll(filepath.Join(util.GetProjectRootDir(), "assets", "mnt", strings.ToLower(t.Name())))
	require.NoError(t, err)
}
