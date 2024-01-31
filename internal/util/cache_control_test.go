package util_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"allaboutapps.dev/aw/go-starter/internal/api"
	"allaboutapps.dev/aw/go-starter/internal/api/middleware"
	"allaboutapps.dev/aw/go-starter/internal/test"
	"allaboutapps.dev/aw/go-starter/internal/util"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCacheControl(t *testing.T) {
	test.WithTestServer(t, func(s *api.Server) {
		path := "/testing-1c2cad5f-7545-4177-9bfe-8dc7ed368b33"

		s.Echo.GET(path, func(c echo.Context) error {
			cache := util.CacheControlDirectiveFromContext(c.Request().Context())
			if cache.HasDirective(util.CacheControlDirectiveNoCache) &&
				cache.HasDirective(util.CacheControlDirectiveNoStore) {
				return c.JSON(http.StatusOK, "no-cache,no-store")
			} else if cache.HasDirective(util.CacheControlDirectiveNoCache) {
				return c.JSON(http.StatusOK, "no-cache")
			} else if cache.HasDirective(util.CacheControlDirectiveNoStore) {
				return c.JSON(http.StatusOK, "no-store")
			}
			return c.NoContent(http.StatusNoContent)

		}, middleware.CacheControl())

		header := http.Header{}
		header.Set(util.HTTPHeaderCacheControl, fmt.Sprintf("%s,%s", util.CacheControlDirectiveNoStore, util.CacheControlDirectiveNoCache))

		res := test.PerformRequest(t, s, "GET", path, nil, header)
		require.Equal(t, http.StatusOK, res.Result().StatusCode)

		var resp string
		err := json.NewDecoder(res.Result().Body).Decode(&resp)
		require.NoError(t, err)
		assert.Equal(t, "no-cache,no-store", resp)

		header.Set(util.HTTPHeaderCacheControl, util.CacheControlDirectiveNoCache.String())

		res = test.PerformRequest(t, s, "GET", path, nil, header)
		require.Equal(t, http.StatusOK, res.Result().StatusCode)

		err = json.NewDecoder(res.Result().Body).Decode(&resp)
		require.NoError(t, err)
		assert.Equal(t, "no-cache", resp)

		header.Set(util.HTTPHeaderCacheControl, util.CacheControlDirectiveNoStore.String())

		res = test.PerformRequest(t, s, "GET", path, nil, header)
		require.Equal(t, http.StatusOK, res.Result().StatusCode)

		err = json.NewDecoder(res.Result().Body).Decode(&resp)
		require.NoError(t, err)
		assert.Equal(t, "no-store", resp)

		res = test.PerformRequest(t, s, "GET", path, nil, nil)
		assert.Equal(t, http.StatusNoContent, res.Result().StatusCode)

		header.Set(util.HTTPHeaderCacheControl, "gunther")

		res = test.PerformRequest(t, s, "GET", path, nil, header)
		assert.Equal(t, http.StatusNoContent, res.Result().StatusCode)
	})
}
