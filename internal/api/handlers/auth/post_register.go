package auth

import (
	"net/http"
	"time"

	"allaboutapps.dev/aw/go-starter/internal/api"
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
	ErrBadRequestInvalidPassword = NewHTTPErrorWithDetail(http.StatusBadRequest, "INVALID_PASSWORD", "The password provided was invalid", "Password was either too weak or did not match other criteria")
	ErrConflictUserAlreadyExists = NewHTTPError(http.StatusConflict, "USER_ALREADY_EXISTS", "User with given username already exists")
)

// swagger:route POST /api/v1/auth/register auth PostRegisterRoute
//
// Registers a local user
//
// Returns an access and refresh token on successful registration
//
// Responses:
//   200: PostLoginResponse
//   400: body:HTTPValidationError HTTPValidationError, type `INVALID_PASSWORD`
//   409: body:HTTPError HTTPError, type `USER_ALREADY_EXISTS`
func PostRegisterRoute(s *api.Server) *echo.Route {
	return s.Router.APIV1Auth.POST("/register", postRegisterHandler(s))
}

func postRegisterHandler(s *api.Server) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		log := util.LogFromContext(ctx)

		var body PostRegisterPayload
		if err := util.BindAndValidate(c, &body); err != nil {
			return err
		}

		exists, err := models.Users(models.UserWhere.Username.EQ(null.StringFrom(body.Username.String()))).Exists(ctx, s.DB)
		if err != nil {
			log.Debug().Err(err).Str("username", body.Username.String()).Msg("Failed to check whether user exists")
			return err
		}

		if exists {
			log.Debug().Str("username", body.Username.String()).Msg("User with given username already exists")
			return ErrConflictUserAlreadyExists
		}

		hash, err := hashing.HashPassword(*body.Password, hashing.DefaultArgon2Params)
		if err != nil {
			log.Debug().Str("username", body.Username.String()).Err(err).Msg("Failed to hash user password")
			return ErrBadRequestInvalidPassword
		}

		response := &PostLoginResponse{
			TokenType: TokenTypeBearer,
			ExpiresIn: int(s.Config.Auth.AccessTokenValidity.Seconds()),
		}

		if err := db.WithTransaction(ctx, s.DB, func(tx boil.ContextExecutor) error {
			user := models.User{
				Username:            null.StringFrom(body.Username.String()),
				Password:            null.StringFrom(hash),
				LastAuthenticatedAt: null.TimeFrom(time.Now()),
				IsActive:            true,
				Scopes:              s.Config.Auth.DefaultUserScopes,
			}

			if err := user.Insert(ctx, tx, boil.Infer()); err != nil {
				log.Debug().Err(err).Msg("Failed to insert user")
				return echo.ErrInternalServerError
			}

			appUserProfile := models.AppUserProfile{
				UserID: user.ID,
			}

			if err := appUserProfile.Insert(ctx, tx, boil.Infer()); err != nil {
				log.Debug().Err(err).Msg("Failed to insert app user profile")
				return echo.ErrInternalServerError
			}

			accessToken := models.AccessToken{
				ValidUntil: time.Now().Add(s.Config.Auth.AccessTokenValidity),
				UserID:     user.ID,
			}

			if err := accessToken.Insert(ctx, tx, boil.Infer()); err != nil {
				log.Debug().Err(err).Msg("Failed to insert access token")
				return echo.ErrInternalServerError
			}

			refreshToken := models.RefreshToken{
				UserID: user.ID,
			}

			if err := refreshToken.Insert(ctx, tx, boil.Infer()); err != nil {
				log.Debug().Err(err).Msg("Failed to insert refresh token")
				return echo.ErrInternalServerError
			}

			response.AccessToken = strfmt.UUID4(accessToken.Token)
			response.RefreshToken = strfmt.UUID4(refreshToken.Token)

			return nil
		}); err != nil {
			log.Debug().Err(err).Msg("Failed to register user")
			return err
		}

		log.Debug().Msg("Successfully registered user, returning new set of access and refresh tokens")

		return util.ValidateAndReturn(c, http.StatusOK, response)
	}
}
