package handlers

import (
	"allaboutapps.at/aw/go-mranftl-sample/api"
	"allaboutapps.at/aw/go-mranftl-sample/api/handlers/auth"
	"allaboutapps.at/aw/go-mranftl-sample/api/handlers/users"
)

func AttachAllRoutes(s *api.Server) {
	// attach our routes
	// TODO: auto-generate those attaches via go generate
	auth.PostLoginRoute(s)
	auth.GetHashBenchmarkRoute(s)
	users.GetUsersRoute(s)
}
