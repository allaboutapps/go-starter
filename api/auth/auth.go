package auth

import (
	"net/http"

	"allaboutapps.at/aw/go-mranftl-sample/api"
	"github.com/labstack/echo/v4"
)

func postLoginHandler(s *api.Server) echo.HandlerFunc {
	return func(c echo.Context) error {
		return echo.NewHTTPError(http.StatusNotImplemented, http.StatusText(http.StatusNotImplemented))
	}
}
