package auth

import (
	"net/http"

	"allaboutapps.dev/aw/go-starter/internal/api"
	"allaboutapps.dev/aw/go-starter/internal/data/dto"
	"allaboutapps.dev/aw/go-starter/internal/types"
	"allaboutapps.dev/aw/go-starter/internal/util"
	"allaboutapps.dev/aw/go-starter/internal/util/url"
	"github.com/go-openapi/swag/conv"
	"github.com/labstack/echo/v4"
)

func PostRegisterRoute(s *api.Server) *echo.Route {
	return s.Router.APIV1Auth.POST("/register", postRegisterHandler(s))
}

func postRegisterHandler(s *api.Server) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		log := util.LogFromContext(ctx)

		var body types.PostRegisterPayload
		if err := util.BindAndValidateBody(c, &body); err != nil {
			return err
		}

		username := dto.NewUsername(body.Username.String())

		result, err := s.Auth.Register(ctx, dto.RegisterRequest{
			Username: username,
			Password: conv.Value(body.Password),
		})
		if err != nil {
			log.Debug().Err(err).Msg("Failed to register user")
			return err
		}

		if !result.RequiresConfirmation {
			loginResult, err := s.Auth.Login(ctx, dto.LoginRequest{
				Username: username,
				Password: conv.Value(body.Password),
			})
			if err != nil {
				log.Debug().Err(err).Msg("Failed to authenticate user after registration")
				return err
			}

			return util.ValidateAndReturn(c, http.StatusOK, loginResult.ToTypes())
		}

		if result.ConfirmationToken.Valid {
			confirmationLink, err := url.ConfirmationDeeplinkURL(s.Config, result.ConfirmationToken.String)
			if err != nil {
				log.Debug().Err(err).Msg("Failed to generate confirmation link")
				return err
			}

			if err := s.Mailer.SendAccountConfirmation(ctx, username.String(), dto.ConfirmatioNotificationPayload{
				ConfirmationLink: confirmationLink.String(),
			}); err != nil {
				log.Debug().Err(err).Msg("Failed to send confirmation email")
				return err
			}
		}

		return util.ValidateAndReturn(c, http.StatusAccepted, &types.RegisterResponse{
			RequiresConfirmation: conv.Pointer(result.RequiresConfirmation),
		})
	}
}
