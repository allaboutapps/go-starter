package router

import (
	"allaboutapps.at/aw/go-mranftl-sample/api"
	"allaboutapps.at/aw/go-mranftl-sample/api/auth"
	"allaboutapps.at/aw/go-mranftl-sample/api/middleware"
	"allaboutapps.at/aw/go-mranftl-sample/api/user"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"
)

func Init(s *api.Server) {
	s.Echo = echo.New()

	s.Echo.Debug = s.Config.Echo.Debug
	s.Echo.HideBanner = true
	// TODO write proper wrapper for echo logger and zerolog/use library for proper logging support including levels and paylods
	s.Echo.Logger.SetOutput(log.With().Str("component", "echo").Str("level", "debug").Logger())

	s.Echo.Pre(echoMiddleware.RemoveTrailingSlash())

	s.Echo.Use(echoMiddleware.Recover())
	s.Echo.Use(echoMiddleware.RequestID())
	s.Echo.Use(middleware.Logger())
	s.Echo.Use(echoMiddleware.Gzip())

	auth.InitRoutes(s)
	user.InitRoutes(s)
}
