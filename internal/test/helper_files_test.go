package test_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"allaboutapps.dev/aw/go-starter/internal/test"
	"github.com/stretchr/testify/assert"
)

func TestPrepareTestFile(t *testing.T) {
	var path string
	test.WithTempDir(t, func(localBasePath, basePath string) {
		assert.True(t, strings.HasSuffix(localBasePath, strings.ToLower(t.Name())))
		assert.NotEmpty(t, basePath)

		fileName := "example.jpg"
		test.PrepareTestFile(t, fileName)

		path = filepath.Join(localBasePath, basePath, fileName)
		_, err := os.Stat(path)
		assert.NoError(t, err)
	})

	_, err := os.Stat(path)
	assert.Error(t, err)
	assert.ErrorIs(t, err, os.ErrNotExist)
}
