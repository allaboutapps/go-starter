package auth

import (
	"time"

	"allaboutapps.dev/aw/go-starter/internal/models"
)

type AuthenticationResult struct {
	Token      string
	User       *models.User
	ValidUntil time.Time
	Scopes     []string
}
