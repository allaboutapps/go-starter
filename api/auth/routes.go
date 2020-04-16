package auth

import (
	"allaboutapps.at/aw/go-mranftl-sample/api"
)

func InitRoutes(s *api.Server) {
	g := s.Echo.Group("/api/v1/auth")

	g.POST("/login", postLoginHandler(s))
}
