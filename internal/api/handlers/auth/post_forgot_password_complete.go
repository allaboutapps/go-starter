package auth

import (
	"net/http"

	"allaboutapps.dev/aw/go-starter/internal/api"
	"allaboutapps.dev/aw/go-starter/internal/data/dto"
	"allaboutapps.dev/aw/go-starter/internal/types"
	"allaboutapps.dev/aw/go-starter/internal/util"
	"github.com/go-openapi/swag"
	"github.com/labstack/echo/v4"
)

func PostForgotPasswordCompleteRoute(s *api.Server) *echo.Route {
	return s.Router.APIV1Auth.POST("/forgot-password/complete", postForgotPasswordCompleteHandler(s))
}

func postForgotPasswordCompleteHandler(s *api.Server) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		log := util.LogFromContext(ctx)

		var body types.PostForgotPasswordCompletePayload
		if err := util.BindAndValidateBody(c, &body); err != nil {
			return err
		}

		result, err := s.Auth.ResetPassword(ctx, dto.ResetPasswordRequest{
			ResetToken:  body.Token.String(),
			NewPassword: swag.StringValue(body.Password),
		})
		if err != nil {
			log.Debug().Err(err).Msg("Failed to reset password")
			return err
		}

		return util.ValidateAndReturn(c, http.StatusOK, result.ToTypes())
	}
}
