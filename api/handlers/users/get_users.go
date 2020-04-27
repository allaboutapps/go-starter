package users

import (
	"net/http"

	"allaboutapps.at/aw/go-mranftl-sample/api"
	"allaboutapps.at/aw/go-mranftl-sample/models"
	"allaboutapps.at/aw/go-mranftl-sample/pkg/auth"
	"allaboutapps.at/aw/go-mranftl-sample/pkg/util"
	"github.com/labstack/echo/v4"
)

func GetUsersRoute(s *api.Server) *echo.Route {
	return s.Router.ApiV1Users.GET("", getUsersHandler(s))
}

func getUsersHandler(s *api.Server) echo.HandlerFunc {

	return func(c echo.Context) error {
		log := util.LogFromEchoContext(c)
		user := auth.UserFromEchoContext(c)
		if user != nil {
			log.Trace().Str("username", user.Username.String).Msg("Retrieved user from context")
		}

		users, err := models.Users().All(c.Request().Context(), s.DB)
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, users)
	}

}
