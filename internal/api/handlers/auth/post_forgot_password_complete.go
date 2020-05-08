package auth

import (
	"database/sql"
	"net/http"
	"time"

	"allaboutapps.dev/aw/go-starter/internal/api"
	"allaboutapps.dev/aw/go-starter/internal/api/middleware"
	"allaboutapps.dev/aw/go-starter/internal/models"
	"allaboutapps.dev/aw/go-starter/internal/types"
	. "allaboutapps.dev/aw/go-starter/internal/types"
	"allaboutapps.dev/aw/go-starter/internal/util"
	"allaboutapps.dev/aw/go-starter/internal/util/db"
	"allaboutapps.dev/aw/go-starter/internal/util/hashing"
	"github.com/go-openapi/strfmt"
	"github.com/labstack/echo/v4"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

var (
	ErrNotFoundTokenNotFound = types.NewHTTPError(http.StatusNotFound, "TOKEN_NOT_FOUND", "Provided token was not found")
	ErrConflictTokenExpired  = types.NewHTTPError(http.StatusConflict, "TOKEN_EXPIRED", "Provided token has expired and is no longer valid")
)

// swagger:route POST /api/v1/auth/forgot-password/complete auth PostForgotPasswordCompleteRoute
//
// Completes password reset for local user
//
// Completes a password reset for a local user, using the password reset token sent via email
// to confirm user access, setting the new password if successful. All current access and refresh
// tokens are invalidated and a new set of auth tokens is returned
//
// Responses:
//   200: PostLoginResponse
//   400: body:HTTPValidationError HTTPValidationError, type `INVALID_PASSWORD`
//   403: body:HTTPError HTTPError, type `USER_DEACTIVATED`/`NOT_LOCAL_USER`
//   404: body:HTTPError HTTPError, type `TOKEN_NOT_FOUND`
//   409: body:HTTPError HTTPError, type `TOKEN_EXPIRED`
func PostForgotPasswordCompleteRoute(s *api.Server) *echo.Route {
	return s.Router.APIV1Auth.POST("/forgot-password/complete", postForgotPasswordCompleteHandler(s))
}

func postForgotPasswordCompleteHandler(s *api.Server) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		log := util.LogFromContext(ctx)

		var body PostForgotPasswordCompletePayload
		if err := util.BindAndValidate(c, &body); err != nil {
			return err
		}

		passwordResetToken, err := models.PasswordResetTokens(
			models.PasswordResetTokenWhere.Token.EQ(body.Token.String()),
			qm.Load(models.PasswordResetTokenRels.User),
		).One(ctx, s.DB)
		if err != nil {
			if err == sql.ErrNoRows {
				log.Debug().Err(err).Msg("Password reset token not found")
				return ErrNotFoundTokenNotFound
			}

			log.Debug().Msg("Failed to load password reset token")
			return err
		}

		user := passwordResetToken.R.User

		if time.Now().After(passwordResetToken.ValidUntil) {
			log.Debug().
				Str("user_id", user.ID).
				Time("valid_until", passwordResetToken.ValidUntil).
				Msg("Password reset token is no longer valid, rejecting password reset")
			return ErrConflictTokenExpired
		}

		if !user.IsActive {
			log.Debug().Str("user_id", user.ID).Msg("User is deactivated, rejecting password reset")
			return middleware.ErrForbiddenUserDeactivated
		}

		if !user.Password.Valid {
			log.Debug().Str("user_id", user.ID).Msg("User is missing password, forbidding password reset")
			return ErrForbiddenNotLocalUser
		}

		hash, err := hashing.HashPassword(*body.Password, hashing.DefaultArgon2Params)
		if err != nil {
			log.Debug().Str("user_id", user.ID).Err(err).Msg("Failed to hash new password")
			return ErrBadRequestInvalidPassword
		}

		response := &PostLoginResponse{
			TokenType: TokenTypeBearer,
			ExpiresIn: int(s.Config.Auth.AccessTokenValidity.Seconds()),
		}

		if err := db.WithTransaction(ctx, s.DB, func(tx boil.ContextExecutor) error {
			user.Password = null.StringFrom(hash)

			if _, err := user.Update(ctx, tx, boil.Whitelist(models.UserColumns.Password)); err != nil {
				log.Debug().Str("user_id", user.ID).Err(err).Msg("Failed to update user")
				return err
			}

			if _, err := user.AccessTokens().DeleteAll(ctx, tx); err != nil {
				log.Debug().Str("user_id", user.ID).Err(err).Msg("Failed to delete existing access tokens")
				return err
			}

			if _, err := user.RefreshTokens().DeleteAll(ctx, tx); err != nil {
				log.Debug().Str("user_id", user.ID).Err(err).Msg("Failed to delete existing refresh tokens")
				return err
			}

			accessToken := models.AccessToken{
				ValidUntil: time.Now().Add(s.Config.Auth.AccessTokenValidity),
				UserID:     user.ID,
			}

			if err := accessToken.Insert(ctx, tx, boil.Infer()); err != nil {
				log.Debug().Str("user_id", user.ID).Err(err).Msg("Failed to insert access token")
				return err
			}

			refreshToken := models.RefreshToken{
				UserID: user.ID,
			}

			if err := refreshToken.Insert(ctx, tx, boil.Infer()); err != nil {
				log.Debug().Str("user_id", user.ID).Err(err).Msg("Failed to insert refresh token")
				return err
			}

			if _, err := passwordResetToken.Delete(ctx, tx); err != nil {
				log.Debug().Str("user_id", user.ID).Err(err).Msg("Failed to delete password reset token")
				return err
			}

			response.AccessToken = strfmt.UUID4(accessToken.Token)
			response.RefreshToken = strfmt.UUID4(refreshToken.Token)

			return nil
		}); err != nil {
			log.Debug().Str("user_id", user.ID).Err(err).Msg("Failed to complete password reset")
			return err
		}

		log.Debug().Str("user_id", user.ID).Msg("Successfully completed password reset, returning new set of access and refresh tokens")

		return util.ValidateAndReturn(c, http.StatusOK, response)
	}
}
