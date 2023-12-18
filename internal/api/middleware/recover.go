package middleware

import (
	"allaboutapps.dev/aw/go-starter/internal/util"

	"github.com/labstack/echo/v4"
)

func LogErrorFuncWithRequestInfoFunc(c echo.Context, err error, stack []byte) error {
	log := util.LogFromContext(c.Request().Context())

	log.Error().Err(err).Bytes("stack", stack).Msg("PANIC RECOVER")

	return nil
}
