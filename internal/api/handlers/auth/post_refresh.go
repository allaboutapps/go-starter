package auth

import (
	"database/sql"
	"net/http"
	"time"

	"allaboutapps.dev/aw/go-starter/internal/api"
	"allaboutapps.dev/aw/go-starter/internal/api/middleware"
	"allaboutapps.dev/aw/go-starter/internal/models"
	"allaboutapps.dev/aw/go-starter/internal/types"
	"allaboutapps.dev/aw/go-starter/internal/util"
	"allaboutapps.dev/aw/go-starter/internal/util/db"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/strfmt/conv"
	"github.com/go-openapi/swag"
	"github.com/labstack/echo/v4"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
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

		oldRefreshToken, err := models.RefreshTokens(
			models.RefreshTokenWhere.Token.EQ(body.RefreshToken.String()),
			qm.Load(models.RefreshTokenRels.User),
		).One(ctx, s.DB)
		if err != nil {
			if err == sql.ErrNoRows {
				log.Debug().Err(err).Msg("Refresh token not found")
				return echo.ErrUnauthorized
			}

			log.Debug().Err(err).Msg("Failed to load refresh token")
			return err
		}

		user := oldRefreshToken.R.User

		if !user.IsActive {
			log.Debug().Msg("User is deactivated, rejecting token refresh")
			return middleware.ErrForbiddenUserDeactivated
		}

		response := &types.PostLoginResponse{
			TokenType: swag.String(TokenTypeBearer),
			ExpiresIn: swag.Int64(int64(s.Config.Auth.AccessTokenValidity.Seconds())),
		}

		if err := db.WithTransaction(ctx, s.DB, func(tx boil.ContextExecutor) error {
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

			if _, err := oldRefreshToken.Delete(ctx, tx); err != nil {
				log.Debug().Err(err).Msg("Failed to delete old refresh token")
				return err
			}

			response.AccessToken = conv.UUID4(strfmt.UUID4(accessToken.Token))
			response.RefreshToken = conv.UUID4(strfmt.UUID4(refreshToken.Token))

			return nil
		}); err != nil {
			log.Debug().Err(err).Msg("Failed to refresh tokens")
			return err
		}

		log.Debug().Msg("Successfully refreshed tokens, returning new set of access and refresh tokens")

		return util.ValidateAndReturn(c, http.StatusOK, response)
	}
}
