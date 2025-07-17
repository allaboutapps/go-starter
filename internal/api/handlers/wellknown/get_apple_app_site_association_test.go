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

func testGetWellKnown(t *testing.T, config config.Server, path string) {
	t.Helper()

	test.WithTestServerConfigurable(t, config, func(s *api.Server) {
		res := test.PerformRequest(t, s, "GET", path, nil, nil)
		require.Equal(t, http.StatusOK, res.Result().StatusCode)

		result, err := io.ReadAll(res.Body)
		require.NoError(t, err)

		test.Snapshoter.SaveString(t, string(result))
	})
}

func TestGetAppleWellKnown(t *testing.T) {
	config := config.DefaultServiceConfigFromEnv()
	config.Paths.AppleAppSiteAssociationFile = filepath.Join(util.GetProjectRootDir(), "test", "testdata", "apple-app-site-association.json")

	testGetWellKnown(t, config, "/.well-known/apple-app-site-association")
}

func TestGetAppleWellKnownNotFound(t *testing.T) {
	config := config.DefaultServiceConfigFromEnv()
	config.Paths.AppleAppSiteAssociationFile = ""

	test.WithTestServerConfigurable(t, config, func(s *api.Server) {
		res := test.PerformRequest(t, s, "GET", "/.well-known/apple-app-site-association", nil, nil)
		test.RequireHTTPError(t, res, httperrors.NewFromEcho(echo.ErrNotFound))
	})
}
