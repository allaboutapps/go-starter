package auth

import (
	"database/sql"
	"net/http"
	"time"

	"allaboutapps.dev/aw/go-starter/internal/api"
	"allaboutapps.dev/aw/go-starter/internal/api/middleware"
	"allaboutapps.dev/aw/go-starter/internal/models"
	. "allaboutapps.dev/aw/go-starter/internal/types"
	"allaboutapps.dev/aw/go-starter/internal/util"
	"allaboutapps.dev/aw/go-starter/internal/util/db"
	"github.com/go-openapi/strfmt"
	"github.com/labstack/echo/v4"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

// swagger:route POST /api/v1/auth/refresh auth PostRefreshRoute
//
// Refresh tokens
//
// Return a fresh set of access and refresh tokens if a valid refresh token was provided.
// The old refresh token used to authenticate the request will be invalidated.
//
// Responses:
//   200: PostLoginResponse
//   400: body:HTTPValidationError
//   401: body:HTTPError
//   403: body:HTTPError HTTPError, type `USER_DEACTIVATED`
func PostRefreshRoute(s *api.Server) *echo.Route {
	return s.Router.APIV1Auth.POST("/refresh", postRefreshHandler(s))
}

func postRefreshHandler(s *api.Server) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		log := util.LogFromContext(ctx)

		var body PostRefreshPayload
		if err := util.BindAndValidate(c, &body); err != nil {
			return err
		}

		oldRefreshToken, err := models.RefreshTokens(
			models.RefreshTokenWhere.Token.EQ(body.RefreshToken.String()),
			qm.Load(models.RefreshTokenRels.User),
		).One(ctx, s.DB)
		if err != nil {
			if err == sql.ErrNoRows {
				log.Debug().Err(err).Msg("Refresh token not found")
			} else {
				log.Debug().Err(err).Msg("Failed to load refresh token")
			}

			return echo.ErrUnauthorized
		}

		user := oldRefreshToken.R.User

		if !user.IsActive {
			log.Debug().Msg("User is deactivated, rejecting token refresh")
			return middleware.ErrForbiddenUserDeactivated
		}

		response := &PostLoginResponse{
			TokenType: TokenTypeBearer,
			ExpiresIn: int(s.Config.Auth.AccessTokenValidity.Seconds()),
		}

		if err := db.WithTransaction(ctx, s.DB, func(tx boil.ContextExecutor) error {
			accessToken := models.AccessToken{
				ValidUntil: time.Now().Add(s.Config.Auth.AccessTokenValidity),
				UserID:     user.ID,
			}

			if err := accessToken.Insert(ctx, tx, boil.Infer()); err != nil {
				log.Debug().Err(err).Msg("Failed to insert access token")
				return echo.ErrUnauthorized
			}

			refreshToken := models.RefreshToken{
				UserID: user.ID,
			}

			if err := refreshToken.Insert(ctx, tx, boil.Infer()); err != nil {
				log.Debug().Err(err).Msg("Failed to insert refresh token")
				return echo.ErrUnauthorized
			}

			if _, err := oldRefreshToken.Delete(ctx, tx); err != nil {
				log.Debug().Err(err).Msg("Failed to delete old refresh token")
				return echo.ErrUnauthorized
			}

			response.AccessToken = strfmt.UUID4(accessToken.Token)
			response.RefreshToken = strfmt.UUID4(refreshToken.Token)

			return nil
		}); err != nil {
			log.Debug().Err(err).Msg("Failed to refresh tokens")
			return err
		}

		log.Debug().Msg("Successfully refreshed tokens, returning new set of access and refresh tokens")

		return util.ValidateAndReturn(c, http.StatusOK, response)
	}
}
