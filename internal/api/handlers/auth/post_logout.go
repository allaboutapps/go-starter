package auth

import (
	"database/sql"
	"errors"
	"net/http"

	"allaboutapps.dev/aw/go-starter/internal/api"
	"allaboutapps.dev/aw/go-starter/internal/api/auth"
	"allaboutapps.dev/aw/go-starter/internal/models"
	"allaboutapps.dev/aw/go-starter/internal/types"
	"allaboutapps.dev/aw/go-starter/internal/util"
	"allaboutapps.dev/aw/go-starter/internal/util/db"
	"github.com/labstack/echo/v4"
	"github.com/volatiletech/sqlboiler/v4/boil"
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

		token := auth.AccessTokenFromEchoContext(c)

		if err := db.WithTransaction(ctx, s.DB, func(tx boil.ContextExecutor) error {
			if _, err := models.AccessTokens(models.AccessTokenWhere.Token.EQ(*token)).DeleteAll(ctx, tx); err != nil {
				log.Debug().Err(err).Msg("Failed to delete access token")
				return err
			}

			if len(body.RefreshToken.String()) > 0 {
				refreshToken, err := models.FindRefreshToken(ctx, tx, body.RefreshToken.String())
				if err != nil {
					if errors.Is(err, sql.ErrNoRows) {
						log.Debug().Msg("Did not find provided refresh token, ignoring")
						return nil
					}

					log.Debug().Err(err).Msg("Failed to load refresh token")
					return err
				}

				if _, err := refreshToken.Delete(ctx, tx); err != nil {
					log.Debug().Err(err).Msg("Failed to delete refresh token")
					return err
				}
			}

			return nil
		}); err != nil {
			log.Debug().Err(err).Msg("Failed to process logout")
			return err
		}

		log.Debug().Msg("Successfully logged out user")

		return c.NoContent(http.StatusNoContent)
	}
}
