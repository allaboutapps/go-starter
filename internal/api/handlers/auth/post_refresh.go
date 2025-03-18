package auth

import (
	"net/http"

	"allaboutapps.dev/aw/go-starter/internal/api"
	"allaboutapps.dev/aw/go-starter/internal/data/dto"
	"allaboutapps.dev/aw/go-starter/internal/types"
	"allaboutapps.dev/aw/go-starter/internal/util"
	"github.com/labstack/echo/v4"
)

func PostRefreshRoute(s *api.Server) *echo.Route {
	return s.Router.APIV1Auth.POST("/refresh", postRefreshHandler(s))
}

func postRefreshHandler(s *api.Server) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		log := util.LogFromContext(ctx)

		var body types.PostRefreshPayload
		if err := util.BindAndValidateBody(c, &body); err != nil {
			return err
		}

		result, err := s.Auth.Refresh(ctx, dto.RefreshRequest{
			RefreshToken: body.RefreshToken.String(),
		})
		if err != nil {
			log.Debug().Err(err).Msg("Failed to refresh tokens")
			return err
		}

		return util.ValidateAndReturn(c, http.StatusOK, result.ToTypes())
	}
}
