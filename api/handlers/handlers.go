package handlers

import (
	"allaboutapps.at/aw/go-mranftl-sample/api"
	"allaboutapps.at/aw/go-mranftl-sample/api/handlers/auth"
	"allaboutapps.at/aw/go-mranftl-sample/api/handlers/users"
)

func AttachHandlers(s *api.Server) {
	// attach our handlers
	// TODO: auto-generate those attaches via go generate
	auth.PostLoginHandler(s)
	auth.GetHashBenchmarkHandler(s)
	users.GetUsersHandler(s)
}
