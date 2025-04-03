package auth

import (
	"time"

	"allaboutapps.dev/aw/go-starter/internal/data/dto"
)

type Result struct {
	Token      string
	User       *dto.User
	ValidUntil time.Time
	Scopes     []string
}
