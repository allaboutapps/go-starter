package handlers

import (
	"allaboutapps.at/aw/go-mranftl-sample/api"
	"allaboutapps.at/aw/go-mranftl-sample/api/handlers/auth"
	"allaboutapps.at/aw/go-mranftl-sample/api/handlers/users"
)

func InitRoutes(s *api.Server) {
	authGroup := s.Echo.Group("/api/v1/auth")
	authGroup.POST("/login", auth.PostLoginHandler(s))

	usersGroup := s.Echo.Group("/api/v1/users")
	usersGroup.GET("", users.GetUsersHandler(s))
}
