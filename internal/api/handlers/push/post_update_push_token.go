package push

import (
	"database/sql"
	"net/http"

	"allaboutapps.dev/aw/go-starter/internal/api"
	"allaboutapps.dev/aw/go-starter/internal/api/auth"
	"allaboutapps.dev/aw/go-starter/internal/api/httperrors"
	"allaboutapps.dev/aw/go-starter/internal/models"
	"allaboutapps.dev/aw/go-starter/internal/types"
	"allaboutapps.dev/aw/go-starter/internal/util"
	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

func PostUpdatePushTokenRoute(s *api.Server) *echo.Route {
	return s.Router.APIV1Push.PUT("/token", postUpdatePushTokenHandler(s))
}

func postUpdatePushTokenHandler(s *api.Server) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		log := util.LogFromContext(ctx)

		var body types.PostUpdatePushTokenPayload
		if err := util.BindAndValidateBody(c, &body); err != nil {
			return err
		}

		user := auth.UserFromEchoContext(c)

		// insert new token
		newToken := models.PushToken{
			UserID:   user.ID,
			Token:    *body.NewToken,
			Provider: *body.Provider,
		}
		if err := newToken.Insert(ctx, s.DB, boil.Infer()); err != nil {
			log.Debug().Str("user_id", user.ID).Err(err).Msg("Failed to insert push token.")

			// check for unique_violation on token column, 23505 == unique_violation
			var pqErr *pq.Error
			if errors.As(err, &pqErr) {
				if pqErr.Code == "23505" && pqErr.Constraint == "push_tokens_token_key" {
					return httperrors.ErrConflictPushToken
				}
			}

			return err
		}

		// delete old token if present in request
		if body.OldToken != nil {
			oldToken, err := models.PushTokens(models.PushTokenWhere.Token.EQ(*body.OldToken), models.PushTokenWhere.UserID.EQ(user.ID)).One(ctx, s.DB)
			if err != nil {
				log.Debug().Str("user_id", user.ID).Err(err).Msg("Old token to delete not found or not assigned to user.")
				if errors.Is(err, sql.ErrNoRows) {
					return httperrors.ErrNotFoundOldPushToken
				}

				return err
			}

			if _, err := oldToken.Delete(ctx, s.DB); err != nil {
				log.Debug().Str("user_id", user.ID).Err(err).Msg("Failed to delete old push token.")
				return err
			}
		}

		log.Debug().Str("user_id", user.ID).Msg("Successfully updated push token.")

		return c.String(http.StatusOK, "Success")
	}
}
