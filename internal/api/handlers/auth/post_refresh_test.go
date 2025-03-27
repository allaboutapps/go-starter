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
	"allaboutapps.dev/aw/go-starter/internal/test/fixtures"
	"allaboutapps.dev/aw/go-starter/internal/types"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestPostRefreshSuccess(t *testing.T) {
	test.WithTestServer(t, func(s *api.Server) {
		ctx := context.Background()
		fix := fixtures.Fixtures()
		payload := test.GenericPayload{
			"refresh_token": fix.User1RefreshToken1.Token,
		}

		res := test.PerformRequest(t, s, "POST", "/api/v1/auth/refresh", payload, nil)
		assert.Equal(t, http.StatusOK, res.Result().StatusCode)

		var response types.PostLoginResponse
		test.ParseResponseAndValidate(t, res, &response)

		assert.NotEmpty(t, response.AccessToken)
		assert.NotEqual(t, fix.User1AccessToken1.Token, response.AccessToken)
		assert.NotEmpty(t, response.RefreshToken)
		assert.NotEqual(t, fix.User1RefreshToken1.Token, response.RefreshToken)
		assert.Equal(t, int64(s.Config.Auth.AccessTokenValidity.Seconds()), *response.ExpiresIn)
		assert.Equal(t, auth.TokenTypeBearer, *response.TokenType)

		err := fix.User1RefreshToken1.Reload(ctx, s.DB)
		assert.ErrorIs(t, err, sql.ErrNoRows)
	})
}

func TestPostRefreshUnknownToken(t *testing.T) {
	test.WithTestServer(t, func(s *api.Server) {
		payload := test.GenericPayload{
			"refresh_token": "c094e933-e5f0-4ece-9c10-914f3122cdb6",
		}

		res := test.PerformRequest(t, s, "POST", "/api/v1/auth/refresh", payload, nil)
		test.RequireHTTPError(t, res, httperrors.NewFromEcho(echo.ErrUnauthorized))
	})
}

func TestPostRefreshDeactivatedUser(t *testing.T) {
	test.WithTestServer(t, func(s *api.Server) {
		ctx := context.Background()
		fix := fixtures.Fixtures()
		payload := test.GenericPayload{
			"refresh_token": fix.UserDeactivatedRefreshToken1.Token,
		}

		res := test.PerformRequest(t, s, "POST", "/api/v1/auth/refresh", payload, nil)
		test.RequireHTTPError(t, res, middleware.ErrForbiddenUserDeactivated)

		err := fix.UserDeactivatedRefreshToken1.Reload(ctx, s.DB)
		assert.NoError(t, err)
	})
}

func TestPostRefreshBadRequest(t *testing.T) {
	test.WithTestServer(t, func(s *api.Server) {
		tests := []struct {
			name    string
			payload test.GenericPayload
		}{
			{
				name:    "MissingRefreshToken",
				payload: test.GenericPayload{},
			},
			{
				name: "EmptyRefreshToken",
				payload: test.GenericPayload{
					"refresh_token": "",
				},
			},
			{
				name: "InvalidToken",
				payload: test.GenericPayload{
					"refresh_token": "not a valid token",
				},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				res := test.PerformRequest(t, s, "POST", "/api/v1/auth/refresh", tt.payload, nil)
				assert.Equal(t, http.StatusBadRequest, res.Result().StatusCode)

				var response httperrors.HTTPValidationError
				test.ParseResponseAndValidate(t, res, &response)

				test.Snapshoter.Save(t, response)
			})
		}
	})
}
