package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (s *Server) GetAdminTemplates() echo.HandlerFunc {
	return func(c echo.Context) error {
		return echo.NewHTTPError(http.StatusNotImplemented, http.StatusText(http.StatusNotImplemented))
	}
}
