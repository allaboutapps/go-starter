package users

import (
	"net/http"

	"allaboutapps.at/aw/go-mranftl-sample/api"
	"allaboutapps.at/aw/go-mranftl-sample/models"
	"allaboutapps.at/aw/go-mranftl-sample/util"
	"github.com/labstack/echo/v4"
)

func GetUsersHandler(s *api.Server) echo.HandlerFunc {
	return func(c echo.Context) error {
		log := util.LogFromEchoContext(c)
		log.Trace().Msg("Loading all users")

		users, err := models.Users().All(c.Request().Context(), s.DB)
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, users)
	}
}
