package auth_test

import (
	"context"
	"database/sql"
	"net/http"
	"testing"

	"allaboutapps.dev/aw/go-starter/internal/api"
	"allaboutapps.dev/aw/go-starter/internal/api/handlers/auth"
	"allaboutapps.dev/aw/go-starter/internal/api/httperrors"
	"allaboutapps.dev/aw/go-starter/internal/api/middleware"
	"allaboutapps.dev/aw/go-starter/internal/test"
	"allaboutapps.dev/aw/go-starter/internal/types"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

func TestPostChangePasswordSuccess(t *testing.T) {
	test.WithTestServer(t, func(s *api.Server) {
		ctx := context.Background()
		fixtures := test.Fixtures()

		newPassword := "correct horse battery staple"
		payload := test.GenericPayload{
			"currentPassword": test.PlainTestUserPassword,
			"newPassword":     newPassword,
		}

		res := test.PerformRequest(t, s, "POST", "/api/v1/auth/change-password", payload, test.HeadersWithAuth(t, fixtures.User1AccessToken1.Token))
		assert.Equal(t, http.StatusOK, res.Result().StatusCode)

		var response types.PostLoginResponse
		test.ParseResponseAndValidate(t, res, &response)

		assert.NotEmpty(t, response.AccessToken)
		assert.NotEqual(t, fixtures.User1AccessToken1.Token, *response.AccessToken)
		assert.NotEmpty(t, response.RefreshToken)
		assert.NotEqual(t, fixtures.User1RefreshToken1.Token, *response.RefreshToken)
		assert.Equal(t, int64(s.Config.Auth.AccessTokenValidity.Seconds()), *response.ExpiresIn)
		assert.Equal(t, auth.TokenTypeBearer, *response.TokenType)

		err := fixtures.User1AccessToken1.Reload(ctx, s.DB)
		assert.ErrorIs(t, err, sql.ErrNoRows)
		err = fixtures.User1RefreshToken1.Reload(ctx, s.DB)
		assert.ErrorIs(t, err, sql.ErrNoRows)

		cnt, err := fixtures.User1.AccessTokens().Count(ctx, s.DB)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), cnt)

		cnt, err = fixtures.User1.RefreshTokens().Count(ctx, s.DB)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), cnt)

		err = fixtures.User1.Reload(ctx, s.DB)
		assert.NoError(t, err)
		assert.NotEqual(t, test.HashedTestUserPassword, fixtures.User1.Password.String)
	})
}

func TestPostChangePasswordInvalidPassword(t *testing.T) {
	test.WithTestServer(t, func(s *api.Server) {
		ctx := context.Background()
		fixtures := test.Fixtures()

		newPassword := "correct horse battery staple"
		payload := test.GenericPayload{
			"currentPassword": "not my password",
			"newPassword":     newPassword,
		}

		res := test.PerformRequest(t, s, "POST", "/api/v1/auth/change-password", payload, test.HeadersWithAuth(t, fixtures.User1AccessToken1.Token))
		test.RequireHTTPError(t, res, httperrors.NewFromEcho(echo.ErrUnauthorized))

		err := fixtures.User1AccessToken1.Reload(ctx, s.DB)
		assert.NoError(t, err)
		err = fixtures.User1RefreshToken1.Reload(ctx, s.DB)
		assert.NoError(t, err)
	})
}

func TestPostChangePasswordDeactivatedUser(t *testing.T) {
	test.WithTestServer(t, func(s *api.Server) {
		ctx := context.Background()
		fixtures := test.Fixtures()

		newPassword := "correct horse battery staple"
		payload := test.GenericPayload{
			"currentPassword": test.PlainTestUserPassword,
			"newPassword":     newPassword,
		}

		res := test.PerformRequest(t, s, "POST", "/api/v1/auth/change-password", payload, test.HeadersWithAuth(t, fixtures.UserDeactivatedAccessToken1.Token))
		test.RequireHTTPError(t, res, middleware.ErrForbiddenUserDeactivated)

		err := fixtures.User1AccessToken1.Reload(ctx, s.DB)
		assert.NoError(t, err)
		err = fixtures.User1RefreshToken1.Reload(ctx, s.DB)
		assert.NoError(t, err)
	})
}

func TestPostChangePasswordUserWithoutPassword(t *testing.T) {
	test.WithTestServer(t, func(s *api.Server) {
		ctx := context.Background()
		fixtures := test.Fixtures()

		newPassword := "correct horse battery staple"
		payload := test.GenericPayload{
			"currentPassword": test.PlainTestUserPassword,
			"newPassword":     newPassword,
		}

		fixtures.User2.Password = null.String{}
		rowsAff, err := fixtures.User2.Update(context.Background(), s.DB, boil.Infer())
		require.NoError(t, err)
		require.Equal(t, int64(1), rowsAff)

		res := test.PerformRequest(t, s, "POST", "/api/v1/auth/change-password", payload, test.HeadersWithAuth(t, fixtures.User2AccessToken1.Token))
		test.RequireHTTPError(t, res, httperrors.ErrForbiddenNotLocalUser)

		err = fixtures.User1AccessToken1.Reload(ctx, s.DB)
		assert.NoError(t, err)
		err = fixtures.User1RefreshToken1.Reload(ctx, s.DB)
		assert.NoError(t, err)
	})
}

func TestPostChangePasswordBadRequest(t *testing.T) {
	test.WithTestServer(t, func(s *api.Server) {
		ctx := context.Background()
		fixtures := test.Fixtures()

		tests := []struct {
			name    string
			payload test.GenericPayload
		}{
			{
				name: "MissingCurrentPassword",
				payload: test.GenericPayload{
					"newPassword": "correct horse battery staple",
				},
			},
			{
				name: "MissingNewPassword",
				payload: test.GenericPayload{
					"currentPassword": test.PlainTestUserPassword,
				},
			},
			{
				name: "EmptyCurrentPassword",
				payload: test.GenericPayload{
					"currentPassword": "",
					"newPassword":     "correct horse battery staple",
				},
			},
			{
				name: "EmptyNewPassword",
				payload: test.GenericPayload{
					"currentPassword": test.PlainTestUserPassword,
					"newPassword":     "",
				},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				res := test.PerformRequest(t, s, "POST", "/api/v1/auth/change-password", tt.payload, test.HeadersWithAuth(t, fixtures.User1AccessToken1.Token))
				assert.Equal(t, http.StatusBadRequest, res.Result().StatusCode)

				var response httperrors.HTTPValidationError
				test.ParseResponseAndValidate(t, res, &response)

				test.Snapshoter.Save(t, response)

				err := fixtures.User1AccessToken1.Reload(ctx, s.DB)
				assert.NoError(t, err)
				err = fixtures.User1RefreshToken1.Reload(ctx, s.DB)
				assert.NoError(t, err)
			})
		}
	})
}
