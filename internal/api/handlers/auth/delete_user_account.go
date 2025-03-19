package auth

import (
	"net/http"

	"allaboutapps.dev/aw/go-starter/internal/api"
	"allaboutapps.dev/aw/go-starter/internal/auth"
	"allaboutapps.dev/aw/go-starter/internal/data/dto"
	"allaboutapps.dev/aw/go-starter/internal/types"
	"allaboutapps.dev/aw/go-starter/internal/util"
	"github.com/go-openapi/swag"
	"github.com/labstack/echo/v4"
)

func DeleteUserAccountRoute(s *api.Server) *echo.Route {
	return s.Router.APIV1Auth.DELETE("/account", deleteUserAccountHandler(s))
}

func deleteUserAccountHandler(s *api.Server) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		user := auth.UserFromContext(ctx)
		log := util.LogFromContext(ctx)

		var body types.DeleteUserAccountPayload
		if err := util.BindAndValidateBody(c, &body); err != nil {
			return err
		}

		err := s.Auth.DeleteUserAccount(ctx, dto.DeleteUserAccountRequest{
			User:            *user,
			CurrentPassword: swag.StringValue(body.CurrentPassword),
		})
		if err != nil {
			log.Debug().Err(err).Msg("Failed to delete user")
			return err
		}

		return c.NoContent(http.StatusNoContent)
	}
}
