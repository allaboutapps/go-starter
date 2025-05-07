package wellknown_test

import (
	"io"
	"net/http"
	"path/filepath"
	"testing"

	"allaboutapps.dev/aw/go-starter/internal/api"
	"allaboutapps.dev/aw/go-starter/internal/api/httperrors"
	"allaboutapps.dev/aw/go-starter/internal/config"
	"allaboutapps.dev/aw/go-starter/internal/test"
	"allaboutapps.dev/aw/go-starter/internal/util"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

func TestGetAndroidWellKnown(t *testing.T) {
	config := config.DefaultServiceConfigFromEnv()
	config.Paths.AndroidAssetlinksFile = filepath.Join(util.GetProjectRootDir(), "test", "testdata", "android-assetlinks.json")

	test.WithTestServerConfigurable(t, config, func(s *api.Server) {
		res := test.PerformRequest(t, s, "GET", "/.well-known/assetlinks.json", nil, nil)
		require.Equal(t, http.StatusOK, res.Result().StatusCode)

		result, err := io.ReadAll(res.Body)
		require.NoError(t, err)

		test.Snapshoter.SaveString(t, string(result))
	})
}

func TestGetAndroidWellKnownNotFound(t *testing.T) {
	config := config.DefaultServiceConfigFromEnv()
	config.Paths.AndroidAssetlinksFile = ""

	test.WithTestServerConfigurable(t, config, func(s *api.Server) {
		res := test.PerformRequest(t, s, "GET", "/.well-known/assetlinks.json", nil, nil)
		test.RequireHTTPError(t, res, httperrors.NewFromEcho(echo.ErrNotFound))
	})
}
