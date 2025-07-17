package test

import (
	"crypto/sha256"
	"fmt"
	"io"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"allaboutapps.dev/aw/go-starter/internal/api/httperrors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func CompareFileHashes(t *testing.T, expectedFilePath string, actualFilePath string) {
	t.Helper()

	expectedFile, err := os.Open(expectedFilePath)
	require.NoError(t, err)
	defer expectedFile.Close()

	expectedHash := sha256.New()
	_, err = io.Copy(expectedHash, expectedFile)
	require.NoError(t, err)

	actualFile, err := os.Open(actualFilePath)
	require.NoError(t, err)
	defer actualFile.Close()

	actualhash := sha256.New()
	_, err = io.Copy(actualhash, actualFile)
	require.NoError(t, err)

	assert.Equal(t, expectedHash.Sum(nil), actualhash.Sum(nil))
}

func CompareAllPayload(t *testing.T, base map[string]interface{}, toCheck map[string]string, skipKeys map[string]bool, keyConvertFunc ...func(string) string) {
	t.Helper()

	var keyFunc func(string) string
	if len(keyConvertFunc) > 0 {
		keyFunc = keyConvertFunc[0]
	} else {
		keyFunc = func(s string) string {
			return s
		}
	}
	for key, val := range base {
		if skipKeys[key] {
			continue
		}

		// checks with contains because of dateTime and null.String struct
		contains := strings.Contains(toCheck[keyFunc(key)], fmt.Sprintf("%v", val))
		assert.Truef(t, contains, "Expected for %s: '%s'. Got: '%s'", key, val, toCheck[keyFunc(key)])
	}
}

func CompareAll(t *testing.T, base map[string]string, toCheck map[string]string, skipKeys map[string]bool) {
	t.Helper()

	for key, val := range base {
		if skipKeys[key] {
			continue
		}

		// checks with contains because of dateTime and null.String struct
		contains := strings.Contains(toCheck[key], val)
		assert.Truef(t, contains, "Expected for %s: '%s'. Got: '%s'", key, val, toCheck[key])
	}
}

func RequireHTTPError(t *testing.T, res *httptest.ResponseRecorder, httpError *httperrors.HTTPError) httperrors.HTTPError {
	t.Helper()

	if httpError.Code != nil {
		require.Equal(t, int(*httpError.Code), res.Result().StatusCode)
	}

	var response httperrors.HTTPError
	ParseResponseAndValidate(t, res, &response)

	require.Equal(t, httpError, &response)

	return response
}
