package auth

import (
	"net/http"
	"time"

	"allaboutapps.dev/aw/go-starter/internal/api"
	"allaboutapps.dev/aw/go-starter/internal/api/auth"
	"allaboutapps.dev/aw/go-starter/internal/api/httperrors"
	"allaboutapps.dev/aw/go-starter/internal/api/middleware"
	"allaboutapps.dev/aw/go-starter/internal/models"
	"allaboutapps.dev/aw/go-starter/internal/types"
	"allaboutapps.dev/aw/go-starter/internal/util"
	"allaboutapps.dev/aw/go-starter/internal/util/db"
	"allaboutapps.dev/aw/go-starter/internal/util/hashing"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/strfmt/conv"
	"github.com/go-openapi/swag"
	"github.com/labstack/echo/v4"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

func PostChangePasswordRoute(s *api.Server) *echo.Route {
	return s.Router.APIV1Auth.POST("/change-password", postChangePasswordHandler(s))
}

func postChangePasswordHandler(s *api.Server) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		log := util.LogFromContext(ctx)

		var body types.PostChangePasswordPayload
		if err := util.BindAndValidateBody(c, &body); err != nil {
			return err
		}

		user := auth.UserFromEchoContext(c)
		if !user.IsActive {
			log.Debug().Msg("User is deactivated, rejecting password change")
			return middleware.ErrForbiddenUserDeactivated
		}

		if !user.Password.Valid {
			log.Debug().Msg("User is missing password, forbidding password change")
			return httperrors.ErrForbiddenNotLocalUser
		}

		match, err := hashing.ComparePasswordAndHash(*body.CurrentPassword, user.Password.String)
		if err != nil {
			log.Debug().Err(err).Msg("Failed to compare password with stored hash")
			return err
		}

		if !match {
			log.Debug().Msg("Provided password does not match stored hash")
			return echo.ErrUnauthorized
		}

		hash, err := hashing.HashPassword(*body.NewPassword, hashing.DefaultArgon2Params)
		if err != nil {
			log.Debug().Err(err).Msg("Failed to hash new password")
			return httperrors.ErrBadRequestInvalidPassword
		}

		response := &types.PostLoginResponse{
			TokenType: swag.String(TokenTypeBearer),
			ExpiresIn: swag.Int64(int64(s.Config.Auth.AccessTokenValidity.Seconds())),
		}

		if err := db.WithTransaction(ctx, s.DB, func(tx boil.ContextExecutor) error {
			user.Password = null.StringFrom(hash)

			if _, err := user.Update(ctx, tx, boil.Whitelist(models.UserColumns.Password)); err != nil {
				log.Debug().Err(err).Msg("Failed to update user")
				return err
			}

			if _, err := user.AccessTokens().DeleteAll(ctx, tx); err != nil {
				log.Debug().Err(err).Msg("Failed to delete existing access tokens")
				return err
			}

			if _, err := user.RefreshTokens().DeleteAll(ctx, tx); err != nil {
				log.Debug().Err(err).Msg("Failed to delete existing refresh tokens")
				return err
			}

			accessToken := models.AccessToken{
				ValidUntil: time.Now().Add(s.Config.Auth.AccessTokenValidity),
				UserID:     user.ID,
			}

			if err := accessToken.Insert(ctx, tx, boil.Infer()); err != nil {
				log.Debug().Err(err).Msg("Failed to insert access token")
				return err
			}

			refreshToken := models.RefreshToken{
				UserID: user.ID,
			}

			if err := refreshToken.Insert(ctx, tx, boil.Infer()); err != nil {
				log.Debug().Err(err).Msg("Failed to insert refresh token")
				return err
			}

			response.AccessToken = conv.UUID4(strfmt.UUID4(accessToken.Token))
			response.RefreshToken = conv.UUID4(strfmt.UUID4(refreshToken.Token))

			return nil
		}); err != nil {
			log.Debug().Err(err).Msg("Failed to change password")
			return err
		}

		log.Debug().Msg("Successfully changed password, returning new set of access and refresh tokens")

		return util.ValidateAndReturn(c, http.StatusOK, response)
	}
}
