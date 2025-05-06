package auth

import (
	"fmt"
	"net/http"

	"allaboutapps.dev/aw/go-starter/internal/api"
	"allaboutapps.dev/aw/go-starter/internal/api/handlers/constants"
	"allaboutapps.dev/aw/go-starter/internal/data/dto"
	"allaboutapps.dev/aw/go-starter/internal/types/auth"

	"allaboutapps.dev/aw/go-starter/internal/util"
	"github.com/labstack/echo/v4"
)

func PostCompleteRegisterRoute(s *api.Server) *echo.Route {
	return s.Router.APIV1Auth.POST(fmt.Sprintf("/register/:%s", constants.RegistrationTokenParam), postCompleteRegisterHandler(s))
}

func postCompleteRegisterHandler(s *api.Server) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		log := util.LogFromContext(ctx)

		params := auth.NewPostCompleteRegisterRouteParams()
		if err := util.BindAndValidatePathAndQueryParams(c, &params); err != nil {
			return err
		}

		result, err := s.Auth.CompleteRegister(ctx, dto.CompleteRegisterRequest{
			ConfirmationToken: params.RegistrationToken.String(),
		})
		if err != nil {
			log.Debug().Err(err).Msg("Failed to complete registration")
			return echo.ErrUnauthorized
		}

		return util.ValidateAndReturn(c, http.StatusOK, result.ToTypes())
	}
}
