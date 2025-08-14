package auth_test

import (
	"context"
	"net/http"
	"testing"

	"allaboutapps.dev/aw/go-starter/internal/api"
	"allaboutapps.dev/aw/go-starter/internal/api/httperrors"
	"allaboutapps.dev/aw/go-starter/internal/models"
	"allaboutapps.dev/aw/go-starter/internal/test"
	"allaboutapps.dev/aw/go-starter/internal/test/fixtures"
	"github.com/aarondl/null/v8"
	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

func assertUserAndRelatedData(ctx context.Context, t *testing.T, s *api.Server, userID string, expectExists bool) {
	t.Helper()

	userExists, err := models.Users(
		models.UserWhere.ID.EQ(userID),
	).Exists(ctx, s.DB)
	require.NoError(t, err)
	require.Equal(t, expectExists, userExists)

	appUserProfileExists, err := models.AppUserProfiles(
		models.AppUserProfileWhere.UserID.EQ(userID),
	).Exists(ctx, s.DB)
	require.NoError(t, err)
	require.Equal(t, expectExists, appUserProfileExists)

	accessTokenExists, err := models.AccessTokens(
		models.AccessTokenWhere.UserID.EQ(userID),
	).Exists(ctx, s.DB)
	require.NoError(t, err)
	require.Equal(t, expectExists, accessTokenExists)

	refreshTokenExists, err := models.RefreshTokens(
		models.RefreshTokenWhere.UserID.EQ(userID),
	).Exists(ctx, s.DB)
	require.NoError(t, err)
	require.Equal(t, expectExists, refreshTokenExists)

	pushTokenExists, err := models.PushTokens(
		models.PushTokenWhere.UserID.EQ(userID),
	).Exists(ctx, s.DB)
	require.NoError(t, err)
	require.Equal(t, expectExists, pushTokenExists)
}

func TestDeleteUserAccount(t *testing.T) {
	test.WithTestServer(t, func(s *api.Server) {
		ctx := t.Context()
		fix := fixtures.Fixtures()

		// expect the user to have a app user profile and different kinds of tokens (access, refresh, push, password reset)
		assertUserAndRelatedData(ctx, t, s, fix.User1.ID, true)

		payload := test.GenericPayload{
			"currentPassword": fixtures.PlainTestUserPassword,
		}

		res := test.PerformRequest(t, s, "DELETE", "/api/v1/auth/account", payload, test.HeadersWithAuth(t, fix.User1AccessToken1.Token))
		require.Equal(t, http.StatusNoContent, res.Result().StatusCode)

		// expect the user and all related data to be deleted
		assertUserAndRelatedData(ctx, t, s, fix.User1.ID, false)
	})
}

func TestDeleteUserAccountCurrentPasswordWrong(t *testing.T) {
	test.WithTestServer(t, func(s *api.Server) {
		fix := fixtures.Fixtures()

		payload := test.GenericPayload{
			"currentPassword": "wrongpassword",
		}

		res := test.PerformRequest(t, s, "DELETE", "/api/v1/auth/account", payload, test.HeadersWithAuth(t, fix.User1AccessToken1.Token))
		test.RequireHTTPError(t, res, httperrors.NewFromEcho(echo.ErrUnauthorized))
	})
}

func TestDeleteUserAccountMissingCurrentPassword(t *testing.T) {
	test.WithTestServer(t, func(s *api.Server) {
		fix := fixtures.Fixtures()

		res := test.PerformRequest(t, s, "DELETE", "/api/v1/auth/account", nil, test.HeadersWithAuth(t, fix.User1AccessToken1.Token))
		require.Equal(t, http.StatusBadRequest, res.Result().StatusCode)
	})
}

func TestDeleteUserAccountNoAuth(t *testing.T) {
	test.WithTestServer(t, func(s *api.Server) {
		res := test.PerformRequest(t, s, "DELETE", "/api/v1/auth/account", nil, nil)
		test.RequireHTTPError(t, res, httperrors.NewFromEcho(echo.ErrUnauthorized))
	})
}

func TestDeleteUserAccountUserNotActive(t *testing.T) {
	test.WithTestServer(t, func(s *api.Server) {
		fix := fixtures.Fixtures()
		ctx := t.Context()

		fix.User1.IsActive = false
		_, err := fix.User1.Update(ctx, s.DB, boil.Whitelist(models.UserColumns.IsActive))
		require.NoError(t, err)

		payload := test.GenericPayload{
			"currentPassword": "somepassword",
		}

		res := test.PerformRequest(t, s, "DELETE", "/api/v1/auth/account", payload, test.HeadersWithAuth(t, fix.User1AccessToken1.Token))
		test.RequireHTTPError(t, res, httperrors.ErrForbiddenUserDeactivated)
	})
}

func TestDeleteUserAccountUserNotLocal(t *testing.T) {
	test.WithTestServer(t, func(s *api.Server) {
		fix := fixtures.Fixtures()
		ctx := t.Context()

		fix.User1.Password = null.String{}
		_, err := fix.User1.Update(ctx, s.DB, boil.Whitelist(models.UserColumns.Password))
		require.NoError(t, err)

		payload := test.GenericPayload{
			"currentPassword": "wrongpassword",
		}

		res := test.PerformRequest(t, s, "DELETE", "/api/v1/auth/account", payload, test.HeadersWithAuth(t, fix.User1AccessToken1.Token))
		test.RequireHTTPError(t, res, httperrors.ErrForbiddenNotLocalUser)
	})
}
