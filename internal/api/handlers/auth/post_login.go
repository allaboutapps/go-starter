package auth

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	"allaboutapps.dev/aw/go-starter/internal/api"
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

const (
	TokenTypeBearer = "bearer"
)

func PostLoginRoute(s *api.Server) *echo.Route {
	return s.Router.APIV1Auth.POST("/login", postLoginHandler(s))
}

func postLoginHandler(s *api.Server) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		log := util.LogFromContext(ctx)

		var body types.PostLoginPayload
		if err := util.BindAndValidateBody(c, &body); err != nil {
			return err
		}

		// enforce lowercase usernames, trim whitespaces
		username := util.ToUsernameFormat(body.Username.String())

		user, err := models.Users(models.UserWhere.Username.EQ(null.StringFrom(username))).One(ctx, s.DB)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				log.Debug().Err(err).Msg("User not found")
			} else {
				log.Debug().Err(err).Msg("Failed to load user")
			}

			return echo.ErrUnauthorized
		}

		if !user.IsActive {
			log.Debug().Msg("User is deactivated, rejecting authentication")
			return middleware.ErrForbiddenUserDeactivated
		}

		if !user.Password.Valid {
			log.Debug().Msg("User is missing password, forbidding authentication")
			return echo.ErrUnauthorized
		}

		match, err := hashing.ComparePasswordAndHash(*body.Password, user.Password.String)
		if err != nil {
			log.Debug().Err(err).Msg("Failed to compare password with stored hash")
			return echo.ErrUnauthorized
		}

		if !match {
			log.Debug().Msg("Provided password does not match stored hash")
			return echo.ErrUnauthorized
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

			user.LastAuthenticatedAt = null.TimeFrom(time.Now())
			if _, err := user.Update(ctx, tx, boil.Infer()); err != nil {
				log.Debug().Err(err).Msg("Failed to update user's last authenticated at timestamp")
				return err
			}

			response.AccessToken = conv.UUID4(strfmt.UUID4(accessToken.Token))
			response.RefreshToken = conv.UUID4(strfmt.UUID4(refreshToken.Token))

			return nil
		}); err != nil {
			log.Debug().Err(err).Msg("Failed to authenticate user")
			return err
		}

		log.Debug().Msg("Successfully authenticated user, returning new set of access and refresh tokens")

		return util.ValidateAndReturn(c, http.StatusOK, response)
	}
}
