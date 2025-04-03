package auth

import (
	"net/http"

	"allaboutapps.dev/aw/go-starter/internal/api"
	"allaboutapps.dev/aw/go-starter/internal/auth"
	"allaboutapps.dev/aw/go-starter/internal/data/dto"
	"allaboutapps.dev/aw/go-starter/internal/types"
	"allaboutapps.dev/aw/go-starter/internal/util"
	"github.com/labstack/echo/v4"
	"github.com/volatiletech/null/v8"
)

func PostLogoutRoute(s *api.Server) *echo.Route {
	return s.Router.APIV1Auth.POST("/logout", postLogoutHandler(s))
}

func postLogoutHandler(s *api.Server) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		log := util.LogFromContext(ctx)

		var body types.PostLogoutPayload
		if err := util.BindAndValidateBody(c, &body); err != nil {
			return err
		}

		request := dto.LogoutRequest{
			AccessToken: *auth.AccessTokenFromEchoContext(c),
		}

		if len(body.RefreshToken.String()) > 0 {
			request.RefreshToken = null.StringFrom(body.RefreshToken.String())
		}

		err := s.Auth.Logout(ctx, request)
		if err != nil {
			log.Debug().Err(err).Msg("Failed to logout user")
			return err
		}

		return c.NoContent(http.StatusNoContent)
	}
}
