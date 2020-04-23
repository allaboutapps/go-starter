package handlers

import (
	"allaboutapps.at/aw/go-mranftl-sample/api"
	"allaboutapps.at/aw/go-mranftl-sample/api/handlers/auth"
	"allaboutapps.at/aw/go-mranftl-sample/api/handlers/users"
	"allaboutapps.at/aw/go-mranftl-sample/api/middleware"
)

func InitRoutes(s *api.Server) {
	authGroup := s.Echo.Group("/api/v1/auth")
	authGroup.GET("/hash/benchmark", auth.GetHashBenchmarkHandler(s))
	authGroup.POST("/login", auth.PostLoginHandler(s))

	usersGroup := s.Echo.Group("/api/v1/users")
	usersGroup.GET("", users.GetUsersHandler(s), middleware.AuthWithConfig(middleware.AuthConfig{S: s, Mode: middleware.AuthModeSecure}))
}
