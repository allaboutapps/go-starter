package user

import (
	"net/http"

	"allaboutapps.at/aw/go-mranftl-sample/api"
	"allaboutapps.at/aw/go-mranftl-sample/models"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

func getUsersHandler(s *api.Server) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := log.Logger.WithContext(c.Request().Context())

		users, err := models.Users().All(ctx, s.DB)
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, users)
	}
}
