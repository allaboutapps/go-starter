package test_test

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"allaboutapps.dev/aw/go-starter/internal/test"
	"allaboutapps.dev/aw/go-starter/internal/util"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/strfmt/conv"
	"github.com/stretchr/testify/require"
)

func TestCompareFileHashes(t *testing.T) {
	tmpDir := t.TempDir()
	newFilePath := tmpDir + "example2.jpg"
	filePath := filepath.Join(util.GetProjectRootDir(), "test", "testdata", "example.jpg")

	file1, err := os.Open(filePath)
	require.NoError(t, err)
	defer file1.Close()

	file2, err := os.Create(newFilePath)
	require.NoError(t, err)
	defer file2.Close()

	_, err = io.Copy(file2, file1)
	require.NoError(t, err)
	require.FileExists(t, newFilePath)

	test.CompareFileHashes(t, filePath, newFilePath)
}

func TestCompareAllPayload(t *testing.T) {
	payload := test.GenericPayload{
		"A":   1,
		"B":   "b",
		"C":   2.3,
		"D":   true,
		"E":   "2020-02-01",
		"F":   conv.UUID4(strfmt.UUID4("0862573e-6ccb-4684-847d-276d3364e91e")),
		"X_Y": "skipped",
	}
	response := map[string]string{
		"A": "1",
		"B": "b",
		"C": "2.3",
		"D": "true",
		"E": util.Date(2020, 2, 1, time.UTC).String(),
		"F": "0862573e-6ccb-4684-847d-276d3364e91e",
	}

	toSkip := map[string]bool{
		"X_Y": true,
	}
	test.CompareAllPayload(t, payload, response, toSkip)

	payload = test.GenericPayload{
		"a":   1,
		"B":   "b",
		"C":   2.3,
		"d":   true,
		"e":   "2020-02-01",
		"F":   conv.UUID4(strfmt.UUID4("0862573e-6ccb-4684-847d-276d3364e91e")),
		"X_Y": "skipped",
	}
	test.CompareAllPayload(t, payload, response, toSkip, strings.ToUpper)
}

func TestCompareAll(t *testing.T) {
	payload := map[string]string{
		"A":   "1",
		"B":   "b",
		"C":   "2.3",
		"D":   "true",
		"E":   "2020-02-01",
		"F":   strfmt.UUID4("0862573e-6ccb-4684-847d-276d3364e91e").String(),
		"X_Y": "skipped",
	}
	response := map[string]string{
		"A": "1",
		"B": "b",
		"C": "2.3",
		"D": "true",
		"E": util.Date(2020, 2, 1, time.UTC).String(),
		"F": "0862573e-6ccb-4684-847d-276d3364e91e",
	}

	toSkip := map[string]bool{
		"X_Y": true,
	}
	test.CompareAll(t, payload, response, toSkip)
}
