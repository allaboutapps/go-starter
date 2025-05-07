package wellknown

import (
	"allaboutapps.dev/aw/go-starter/internal/api"
	"github.com/labstack/echo/v4"
)

func GetAndroidDigitalAssetLinksRoute(s *api.Server) *echo.Route {
	return s.Router.WellKnown.GET("/assetlinks.json", getAndroidDigitalAssetLinksHandler(s))
}

func getAndroidDigitalAssetLinksHandler(s *api.Server) echo.HandlerFunc {
	return func(c echo.Context) error {
		if s.Config.Paths.AndroidAssetlinksFile == "" {
			return echo.ErrNotFound
		}

		c.Response().Header().Set("Cache-Control", "public, max-age=0, must-revalidate")
		c.Response().Header().Set("Content-Type", "application/json")

		return c.File(s.Config.Paths.AndroidAssetlinksFile)
	}
}
