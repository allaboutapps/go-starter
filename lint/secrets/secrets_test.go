package main_test

import (
	"path/filepath"
	"testing"

	lint "allaboutapps.dev/aw/secrests-linter"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestSecretLinter(t *testing.T) {
	testData := filepath.Join("app", "test", "testdata", "lint")
	t.Log(testData)
	analysistest.Run(t, testData, lint.SecretsAnalyzer, "lint")
}
