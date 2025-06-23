//go:build wireinject

package api

import (
	"database/sql"

	"allaboutapps.dev/aw/go-starter/internal/config"
	"allaboutapps.dev/aw/go-starter/internal/data/local"
	"allaboutapps.dev/aw/go-starter/internal/metrics"
	"allaboutapps.dev/aw/go-starter/internal/persistence"
	"github.com/google/wire"
)

// INJECTORS - https://github.com/google/wire/blob/main/docs/guide.md#injectors

// serviceSet groups the default set of providers that are required for initing a server
var serviceSet = wire.NewSet(
	newServerWithComponents,
	NewPush,
	NewMailer,
	NewI18N,
	NewAuthService,
	local.NewService,
	metrics.New,
	NewClock,
)

// InitNewServer returns a new Server instance.
func InitNewServer(
	_ config.Server,
) (*Server, error) {
	wire.Build(serviceSet, persistence.NewDB)
	return new(Server), nil
}

// InitNewServerWithDB returns a new Server instance with the given DB instance.
// All the other components are initialized via go wire according to the configuration.
func InitNewServerWithDB(
	_ config.Server,
	_ *sql.DB,
) (*Server, error) {
	wire.Build(serviceSet)
	return new(Server), nil
}
