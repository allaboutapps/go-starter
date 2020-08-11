package auth

import (
	"database/sql"
	"errors"
	"net/http"

	"allaboutapps.dev/aw/go-starter/internal/api"
	"allaboutapps.dev/aw/go-starter/internal/api/auth"
	"allaboutapps.dev/aw/go-starter/internal/types"
	"allaboutapps.dev/aw/go-starter/internal/util"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/labstack/echo/v4"
)

func GetUserInfoRoute(s *api.Server) *echo.Route {
	return s.Router.APIV1Auth.GET("/userinfo", getUserInfoHandler(s))
}

func getUserInfoHandler(s *api.Server) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		user := auth.UserFromContext(ctx)
		log := util.LogFromContext(ctx)

		response := &types.GetUserInfoResponse{
			Sub:       swag.String(user.ID),
			UpdatedAt: swag.Int64(user.UpdatedAt.Unix()),
			Email:     strfmt.Email(user.Username.String),
			Scopes:    user.Scopes,
		}

		// if this user has an appUserProfile attached, add additional / modify props from there
		var err error
		appUserProfile, err := user.AppUserProfile().One(ctx, s.DB)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return util.ValidateAndReturn(c, http.StatusOK, response)
			}

			log.Debug().Err(err).Msg("Unknown error while getting appUserProfile information for user")
			return err
		}

		if appUserProfile.UpdatedAt.After(user.UpdatedAt) {
			response.UpdatedAt = swag.Int64(appUserProfile.UpdatedAt.Unix())
		}

		return util.ValidateAndReturn(c, http.StatusOK, response)
	}
}
