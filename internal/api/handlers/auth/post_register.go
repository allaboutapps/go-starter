package auth

import (
	"net/http"
	"time"

	"allaboutapps.dev/aw/go-starter/internal/api"
	"allaboutapps.dev/aw/go-starter/internal/api/httperrors"
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

		// enforce lowercase usernames, trim whitespaces
		username := util.ToUsernameFormat(body.Username.String())

		exists, err := models.Users(models.UserWhere.Username.EQ(null.StringFrom(username))).Exists(ctx, s.DB)
		if err != nil {
			log.Debug().Err(err).Str("username", username).Msg("Failed to check whether user exists")
			return err
		}

		if exists {
			log.Debug().Str("username", username).Msg("User with given username already exists")
			return httperrors.ErrConflictUserAlreadyExists
		}

		hash, err := hashing.HashPassword(*body.Password, hashing.DefaultArgon2Params)
		if err != nil {
			log.Debug().Str("username", username).Err(err).Msg("Failed to hash user password")
			return httperrors.ErrBadRequestInvalidPassword
		}

		response := &types.PostLoginResponse{
			TokenType: swag.String(TokenTypeBearer),
			ExpiresIn: swag.Int64(int64(s.Config.Auth.AccessTokenValidity.Seconds())),
		}

		if err := db.WithTransaction(ctx, s.DB, func(tx boil.ContextExecutor) error {
			user := &models.User{
				Username:            null.StringFrom(username),
				Password:            null.StringFrom(hash),
				LastAuthenticatedAt: null.TimeFrom(time.Now()),
				IsActive:            true,
				Scopes:              s.Config.Auth.DefaultUserScopes,
			}

			if err := user.Insert(ctx, tx, boil.Infer()); err != nil {
				log.Debug().Err(err).Msg("Failed to insert user")
				return err
			}

			appUserProfile := models.AppUserProfile{
				UserID: user.ID,
			}

			if err := appUserProfile.Insert(ctx, tx, boil.Infer()); err != nil {
				log.Debug().Err(err).Msg("Failed to insert app user profile")
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
			log.Debug().Err(err).Msg("Failed to register user")
			return err
		}

		log.Debug().Msg("Successfully registered user, returning new set of access and refresh tokens")

		return util.ValidateAndReturn(c, http.StatusOK, response)
	}
}
