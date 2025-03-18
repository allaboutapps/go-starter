package auth

import (
	"net/http"

	"allaboutapps.dev/aw/go-starter/internal/api"
	"allaboutapps.dev/aw/go-starter/internal/auth"
	"allaboutapps.dev/aw/go-starter/internal/util"
	"github.com/labstack/echo/v4"
)

func GetUserInfoRoute(s *api.Server) *echo.Route {
	return s.Router.APIV1Auth.GET("/userinfo", getUserInfoHandler(s))
}

func getUserInfoHandler(s *api.Server) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		user := auth.UserFromContext(ctx)
		log := util.LogFromContext(ctx)

		var err error
		user.Profile, err = s.Auth.GetAppUserProfileIfExists(ctx, user.ID)
		if err != nil {
			log.Debug().Err(err).Msg("Failed to get user profile")
			return err
		}

		return util.ValidateAndReturn(c, http.StatusOK, user.ToTypes())
	}
}
