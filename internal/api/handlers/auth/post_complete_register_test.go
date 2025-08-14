package auth_test

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"allaboutapps.dev/aw/go-starter/internal/api"
	"allaboutapps.dev/aw/go-starter/internal/api/httperrors"
	"allaboutapps.dev/aw/go-starter/internal/models"
	"allaboutapps.dev/aw/go-starter/internal/test"
	"allaboutapps.dev/aw/go-starter/internal/test/fixtures"
	"allaboutapps.dev/aw/go-starter/internal/types"
	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPostCompleteRegister(t *testing.T) {
	test.WithTestServer(t, func(s *api.Server) {
		fix := fixtures.Fixtures()
		ctx := t.Context()

		res := test.PerformRequest(t, s, "POST", fmt.Sprintf("/api/v1/auth/register/%s", fix.UserRequiresConfirmationConfirmationToken.Token), nil, nil)
		require.Equal(t, http.StatusOK, res.Result().StatusCode)

		var response types.PostLoginResponse
		test.ParseResponseAndValidate(t, res, &response)

		assert.NotEmpty(t, response.AccessToken)
		assert.NotEmpty(t, response.RefreshToken)

		user, err := models.Users(
			models.UserWhere.ID.EQ(fix.UserRequiresConfirmationConfirmationToken.UserID),
		).One(ctx, s.DB)
		require.NoError(t, err)

		assert.True(t, user.IsActive)
		assert.False(t, user.RequiresConfirmation)

		// trying again with the same token should fail
		{
			res := test.PerformRequest(t, s, "POST", fmt.Sprintf("/api/v1/auth/register/%s", fix.UserRequiresConfirmationConfirmationToken.Token), nil, nil)
			test.RequireHTTPError(t, res, httperrors.NewFromEcho(echo.ErrUnauthorized))
		}
	})
}

func TestPostCompleteRegisterTokenNotFound(t *testing.T) {
	test.WithTestServer(t, func(s *api.Server) {
		res := test.PerformRequest(t, s, "POST", "/api/v1/auth/register/e45071b7-b9a0-4ed7-a5a0-16a16413d275", nil, nil)
		test.RequireHTTPError(t, res, httperrors.NewFromEcho(echo.ErrUnauthorized))
	})
}

func TestPostCompleteRegisterTokenExpired(t *testing.T) {
	test.WithTestServer(t, func(s *api.Server) {
		fix := fixtures.Fixtures()
		ctx := t.Context()

		fix.UserRequiresConfirmationConfirmationToken.ValidUntil = s.Clock.Now().Add(-1 * time.Second)
		_, err := fix.UserRequiresConfirmationConfirmationToken.Update(ctx, s.DB, boil.Whitelist(models.ConfirmationTokenColumns.ValidUntil))
		require.NoError(t, err)

		res := test.PerformRequest(t, s, "POST", fmt.Sprintf("/api/v1/auth/register/%s", fix.UserRequiresConfirmationConfirmationToken.Token), nil, nil)
		test.RequireHTTPError(t, res, httperrors.NewFromEcho(echo.ErrUnauthorized))
	})
}
