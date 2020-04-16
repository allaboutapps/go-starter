package util

import (
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// LogFromContext returns a request-specific zerolog instance using the echo.Context of the request.
// The returned logger will have the request ID as well as some other value predefined.
func LogFromContext(c echo.Context) *zerolog.Logger {
	return log.Ctx(c.Request().Context())
}
