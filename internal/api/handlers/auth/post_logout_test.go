package auth_test

import (
	"context"
	"database/sql"
	"net/http"
	"testing"

	"allaboutapps.dev/aw/go-starter/internal/api"
	"allaboutapps.dev/aw/go-starter/internal/api/httperrors"
	"allaboutapps.dev/aw/go-starter/internal/api/middleware"
	"allaboutapps.dev/aw/go-starter/internal/test"
	"allaboutapps.dev/aw/go-starter/internal/test/fixtures"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestPostLogoutSuccess(t *testing.T) {
	test.WithTestServer(t, func(s *api.Server) {
		ctx := context.Background()
		fix := fixtures.Fixtures()

		res := test.PerformRequest(t, s, "POST", "/api/v1/auth/logout", nil, test.HeadersWithAuth(t, fix.User1AccessToken1.Token))
		assert.Equal(t, http.StatusNoContent, res.Result().StatusCode)

		err := fix.User1AccessToken1.Reload(ctx, s.DB)
		assert.ErrorIs(t, err, sql.ErrNoRows)

		err = fix.User1RefreshToken1.Reload(ctx, s.DB)
		assert.NoError(t, err)
	})
}

func TestPostLogoutSuccessWithRefreshToken(t *testing.T) {
	test.WithTestServer(t, func(s *api.Server) {
		ctx := context.Background()
		fix := fixtures.Fixtures()

		payload := test.GenericPayload{
			"refresh_token": fix.User1RefreshToken1.Token,
		}

		res := test.PerformRequest(t, s, "POST", "/api/v1/auth/logout", payload, test.HeadersWithAuth(t, fix.User1AccessToken1.Token))
		assert.Equal(t, http.StatusNoContent, res.Result().StatusCode)

		err := fix.User1AccessToken1.Reload(ctx, s.DB)
		assert.ErrorIs(t, err, sql.ErrNoRows)

		err = fix.User1RefreshToken1.Reload(ctx, s.DB)
		assert.ErrorIs(t, err, sql.ErrNoRows)
	})
}

func TestPostLogoutSuccessWithUnknownRefreshToken(t *testing.T) {
	test.WithTestServer(t, func(s *api.Server) {
		ctx := context.Background()
		fix := fixtures.Fixtures()

		payload := test.GenericPayload{
			"refresh_token": "93d8ccd0-be30-4661-a428-cbe74e1a3ffe",
		}

		res := test.PerformRequest(t, s, "POST", "/api/v1/auth/logout", payload, test.HeadersWithAuth(t, fix.User1AccessToken1.Token))
		assert.Equal(t, http.StatusNoContent, res.Result().StatusCode)

		err := fix.User1AccessToken1.Reload(ctx, s.DB)
		assert.ErrorIs(t, err, sql.ErrNoRows)

		err = fix.User1RefreshToken1.Reload(ctx, s.DB)
		assert.NoError(t, err)
	})
}

func TestPostLogoutInvalidRefreshToken(t *testing.T) {
	test.WithTestServer(t, func(s *api.Server) {
		ctx := context.Background()
		fix := fixtures.Fixtures()
		payload := test.GenericPayload{
			"refresh_token": "not my refresh token",
		}

		res := test.PerformRequest(t, s, "POST", "/api/v1/auth/logout", payload, test.HeadersWithAuth(t, fix.User1AccessToken1.Token))
		assert.Equal(t, http.StatusBadRequest, res.Result().StatusCode)

		var response httperrors.HTTPValidationError
		test.ParseResponseAndValidate(t, res, &response)

		test.Snapshoter.Save(t, response)

		err := fix.User1AccessToken1.Reload(ctx, s.DB)
		assert.NoError(t, err)

		err = fix.User1RefreshToken1.Reload(ctx, s.DB)
		assert.NoError(t, err)
	})
}

func TestPostLogoutError(t *testing.T) {
	test.WithTestServer(t, func(s *api.Server) {
		tests := []struct {
			name          string
			expectedError *httperrors.HTTPError
			headers       http.Header
		}{
			{
				name:          "InvalidAuthToken",
				expectedError: middleware.ErrBadRequestMalformedToken,
				headers:       test.HeadersWithAuth(t, "not my auth token"),
			},
			{
				name:          "UnknownAuthToken",
				expectedError: httperrors.NewFromEcho(echo.ErrUnauthorized),
				headers:       test.HeadersWithAuth(t, "25e8630e-9a41-4f38-8339-373f0c203cef"),
			},
			{
				name:          "MissingAuthToken",
				expectedError: httperrors.NewFromEcho(echo.ErrUnauthorized),
				headers:       nil,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				res := test.PerformRequest(t, s, "POST", "/api/v1/auth/logout", nil, tt.headers)
				test.RequireHTTPError(t, res, tt.expectedError)
			})
		}
	})
}
