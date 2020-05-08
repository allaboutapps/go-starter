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
	"allaboutapps.dev/aw/go-starter/internal/util/hashing"
	"github.com/go-openapi/strfmt"
	"github.com/labstack/echo/v4"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

var (
	ErrForbiddenNotLocalUser = NewHTTPError(http.StatusForbidden, "NOT_LOCAL_USER", "User account is not valid for local authentication")
)

const (
	TokenTypeBearer = "bearer"
)

// swagger:route POST /api/v1/auth/login auth PostLoginRoute
//
// Login with local user
//
// Returns an access and refresh token on successful authentication
//
// Responses:
//   200: PostLoginResponse
//   400: body:HTTPValidationError
//   401: body:HTTPError
//   403: body:HTTPError HTTPError, type `USER_DEACTIVATED`/`NOT_LOCAL_USER`
func PostLoginRoute(s *api.Server) *echo.Route {
	return s.Router.APIV1Auth.POST("/login", postLoginHandler(s))
}

func postLoginHandler(s *api.Server) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		log := util.LogFromContext(ctx)

		var body PostLoginPayload
		if err := util.BindAndValidate(c, &body); err != nil {
			return err
		}

		user, err := models.Users(models.UserWhere.Username.EQ(null.StringFrom(body.Username.String()))).One(ctx, s.DB)
		if err != nil {
			if err == sql.ErrNoRows {
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
			return ErrForbiddenNotLocalUser
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

			user.LastAuthenticatedAt = null.TimeFrom(time.Now())
			if _, err := user.Update(ctx, tx, boil.Infer()); err != nil {
				log.Debug().Err(err).Msg("Failed to update user's last authenticated at timestamp")
				return echo.ErrUnauthorized
			}

			response.AccessToken = strfmt.UUID4(accessToken.Token)
			response.RefreshToken = strfmt.UUID4(refreshToken.Token)

			return nil
		}); err != nil {
			log.Debug().Err(err).Msg("Failed to authenticate user")
			return err
		}

		log.Debug().Msg("Successfully authenticated user, returning new set of access and refresh tokens")

		return util.ValidateAndReturn(c, http.StatusOK, response)
	}
}
