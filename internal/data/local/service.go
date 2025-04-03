package local

import (
	"database/sql"

	"allaboutapps.dev/aw/go-starter/internal/config"
	"github.com/dropbox/godropbox/time2"
)

type Service struct {
	config config.Server
	db     *sql.DB
	clock  time2.Clock
}

func NewService(config config.Server, db *sql.DB, clock time2.Clock) *Service {
	return &Service{
		config: config,
		db:     db,
		clock:  clock,
	}
}
