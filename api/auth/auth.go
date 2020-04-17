package auth

import (
	"crypto/sha512"
	"crypto/subtle"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"allaboutapps.at/aw/go-mranftl-sample/api"
	"allaboutapps.at/aw/go-mranftl-sample/models"
	"allaboutapps.at/aw/go-mranftl-sample/util"
	"github.com/labstack/echo/v4"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries/qm"
	"golang.org/x/crypto/pbkdf2"
)

func postLoginHandler(s *api.Server) echo.HandlerFunc {
	type payload struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	type response struct {
		AccessToken  string `json:"access_token"`
		TokenType    string `json:"token_type"`
		ExpiresIn    int    `json:"expires_in"`
		RefreshToken string `json:"refresh_token"`
	}

	return func(c echo.Context) error {
		ctx := c.Request().Context()
		log := util.LogFromContext(ctx)

		var body payload
		if err := c.Bind(&body); err != nil {
			log.Debug().Err(err).Msg("Failed to bind payload")
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

		if !user.Password.Valid || !user.Salt.Valid {
			log.Debug().Bool("passwordValid", user.Password.Valid).Bool("saltValid", user.Salt.Valid).Msg("User is missing password or salt")
			return echo.ErrForbidden
		}

		log.Debug().Str("userID", user.ID).Msg("Found user")

		hash := pbkdf2.Key([]byte(body.Password), []byte(user.Salt.String), 12000, 512, sha512.New)

		if subtle.ConstantTimeCompare([]byte(fmt.Sprintf("%x", hash)), []byte(user.Password.String)) == 0 {
			log.Debug().Msg("Invalid password")
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

		return c.JSON(http.StatusOK, response{
			AccessToken:  accessToken.Token,
			TokenType:    "bearer",
			ExpiresIn:    int((time.Hour * 24).Seconds()),
			RefreshToken: refreshToken.Token,
		})
	}
}
