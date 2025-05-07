package wellknown

import (
	"allaboutapps.dev/aw/go-starter/internal/api"
	"github.com/labstack/echo/v4"
)

func GetAppleAppSiteAssociationRoute(s *api.Server) *echo.Route {
	return s.Router.WellKnown.GET("/apple-app-site-association", getAppleAppSiteAssociationHandler(s))
}

func getAppleAppSiteAssociationHandler(s *api.Server) echo.HandlerFunc {
	return func(c echo.Context) error {
		if s.Config.Paths.AppleAppSiteAssociationFile == "" {
			return echo.ErrNotFound
		}

		c.Response().Header().Set("Cache-Control", "public, max-age=0, must-revalidate")
		return c.File(s.Config.Paths.AppleAppSiteAssociationFile)
	}
}
