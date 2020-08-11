package common

import (
	"path/filepath"

	"allaboutapps.dev/aw/go-starter/internal/api"
	"github.com/labstack/echo/v4"
)

func GetSwaggerRoute(s *api.Server) *echo.Route {
	// ---
	// Serve generated swagger.yml file statically at /swagger.yml
	// hack: not attached to group - can go away after echo/group.go .File and .Static actually return the *echo.Route
	// see https://github.com/labstack/echo/issues/1595
	// return s.Router.Root.File("swagger.yml", filepath.Join(s.Config.Echo.APIBaseDirAbs, "swagger.yml"))

	return s.Echo.File("/swagger.yml", filepath.Join(s.Config.Paths.APIBaseDirAbs, "swagger.yml"))
}
