package auth

import (
	"database/sql"
	"net/http"
	"time"

	"allaboutapps.at/aw/go-mranftl-sample/api"
	"allaboutapps.at/aw/go-mranftl-sample/models"
	"allaboutapps.at/aw/go-mranftl-sample/pkg/auth/hashing"
	"allaboutapps.at/aw/go-mranftl-sample/pkg/util"
	. "allaboutapps.at/aw/go-mranftl-sample/types"
	"github.com/go-openapi/strfmt"
	"github.com/labstack/echo/v4"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries/qm"
)

var (
	SAMPLE_EXPORTED_PGK_CONST = "test"
)

// swagger:route POST /api/v1/auth/login PostLoginHandler
//
// Login with local user
//
// Returns AccessToken for the user if correct credentials are provided
//
// ---
// produces:
// - application/json
// parameters: PostLoginPayload
// responses:
//   200: PostLoginResponse
func PostLoginHandler(s *api.Server) echo.HandlerFunc {

	return func(c echo.Context) error {
		ctx := c.Request().Context()
		log := util.LogFromContext(ctx)

		var body PostLoginPayload
		if err := c.Bind(&body); err != nil {
			log.Debug().Err(err).Msg("Failed to bind payload")
			return err
		}

		if err := body.Validate(strfmt.Default); err != nil {
			log.Debug().Err(err).Msg("Failed to validate payload")
			return err
		}

		if len(body.Username) == 0 || len(body.Password) == 0 {
			log.Debug().Str("username", body.Username).Str("password", body.Password).Msg("Missing username or password")
			return echo.ErrBadRequest
		}

		user, err := models.Users(qm.Where("username = ?", body.Username)).One(ctx, s.DB)
		if err != nil {
			if err == sql.ErrNoRows {
				log.Debug().Err(err).Msg("User not found")
			} else {
				log.Debug().Err(err).Msg("Failed to load user")
			}

			return echo.ErrUnauthorized
		}

		if !user.Password.Valid {
			log.Debug().Msg("User is missing password")
			return echo.ErrForbidden
		}

		log.Debug().Str("userID", user.ID).Msg("Found user")

		match, err := hashing.ComparePasswordAndHash(body.Password, user.Password.String)
		if err != nil {
			log.Debug().Err(err).Msg("Failed to compare password with stored hash")
			return echo.ErrUnauthorized
		}

		if !match {
			log.Debug().Msg("Provided password does not match stored hash")
			return echo.ErrUnauthorized
		}

		tx, err := s.DB.BeginTx(ctx, nil)
		if err != nil {
			log.Debug().Err(err).Msg("Failed to start transaction")
			return echo.ErrUnauthorized
		}

		accessToken := models.AccessToken{
			ValidUntil: time.Now().Add(24 * time.Hour),
			UserID:     user.ID,
		}

		if err := accessToken.Insert(ctx, tx, boil.Infer()); err != nil {
			log.Debug().Err(err).Msg("Failed to insert access token")
			return echo.ErrUnauthorized
		}

		log.Debug().Str("token", accessToken.Token).Msg("Inserted access token")

		refreshToken := models.RefreshToken{
			UserID: user.ID,
		}

		if err := refreshToken.Insert(ctx, tx, boil.Infer()); err != nil {
			log.Debug().Err(err).Msg("Failed to insert refresh token")
			return echo.ErrUnauthorized
		}

		log.Debug().Str("token", refreshToken.Token).Msg("Inserted refresh token")

		if err := tx.Commit(); err != nil {
			log.Debug().Err(err).Msg("Failed to commit transaction")
			return echo.ErrUnauthorized
		}

		response := PostLoginResponse{
			AccessToken:  accessToken.Token,
			TokenType:    "bearer",
			ExpiresIn:    int((time.Hour * 24).Seconds()),
			RefreshToken: refreshToken.Token,
		}

		if err := response.Validate(strfmt.Default); err != nil {
			log.Debug().Err(err).Msg("Failed to validate response")
			return err
		}

		return c.JSON(http.StatusOK, &response)
	}
}
