package auth_test

import (
	"context"
	"fmt"
	"net/http"
	"strings"
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

func TestPostLoginSuccess(t *testing.T) {
	test.WithTestServer(t, func(s *api.Server) {
		fixtures := test.Fixtures()
		payload := test.GenericPayload{
			"username": fixtures.User1.Username,
			"password": test.PlainTestUserPassword,
		}

		res := test.PerformRequest(t, s, "POST", "/api/v1/auth/login", payload, nil)
		assert.Equal(t, http.StatusOK, res.Result().StatusCode)

		var response types.PostLoginResponse
		test.ParseResponseAndValidate(t, res, &response)

		assert.NotEmpty(t, response.AccessToken)
		assert.NotEqual(t, fixtures.User1AccessToken1.Token, response.AccessToken)
		assert.NotEmpty(t, response.RefreshToken)
		assert.NotEqual(t, fixtures.User1RefreshToken1.Token, response.RefreshToken)
		assert.Equal(t, int64(s.Config.Auth.AccessTokenValidity.Seconds()), *response.ExpiresIn)
		assert.Equal(t, auth.TokenTypeBearer, *response.TokenType)
	})
}

func TestPostLoginSuccessLowercaseTrimWhitespaces(t *testing.T) {
	test.WithTestServer(t, func(s *api.Server) {
		fixtures := test.Fixtures()
		payload := test.GenericPayload{
			"username": fmt.Sprintf(" %s ", strings.ToUpper(fixtures.User1.Username.String)),
			"password": test.PlainTestUserPassword,
		}

		res := test.PerformRequest(t, s, "POST", "/api/v1/auth/login", payload, nil)
		assert.Equal(t, http.StatusOK, res.Result().StatusCode)

		var response types.PostLoginResponse
		test.ParseResponseAndValidate(t, res, &response)

		assert.NotEmpty(t, response.AccessToken)
		assert.NotEqual(t, fixtures.User1AccessToken1.Token, response.AccessToken)
		assert.NotEmpty(t, response.RefreshToken)
		assert.NotEqual(t, fixtures.User1RefreshToken1.Token, response.RefreshToken)
		assert.Equal(t, int64(s.Config.Auth.AccessTokenValidity.Seconds()), *response.ExpiresIn)
		assert.Equal(t, auth.TokenTypeBearer, *response.TokenType)
	})
}

func TestPostLoginInvalidCredentials(t *testing.T) {
	test.WithTestServer(t, func(s *api.Server) {
		fixtures := test.Fixtures()
		payload := test.GenericPayload{
			"username": fixtures.User1.Username,
			"password": "not my password",
		}

		res := test.PerformRequest(t, s, "POST", "/api/v1/auth/login", payload, nil)
		test.RequireHTTPError(t, res, httperrors.NewFromEcho(echo.ErrUnauthorized))
	})
}

func TestPostLoginUnknownUser(t *testing.T) {
	test.WithTestServer(t, func(s *api.Server) {
		payload := test.GenericPayload{
			"username": "definitelydoesnotexist@example.com",
			"password": test.PlainTestUserPassword,
		}

		res := test.PerformRequest(t, s, "POST", "/api/v1/auth/login", payload, nil)
		test.RequireHTTPError(t, res, httperrors.NewFromEcho(echo.ErrUnauthorized))
	})
}

func TestPostLoginDeactivatedUser(t *testing.T) {
	test.WithTestServer(t, func(s *api.Server) {
		fixtures := test.Fixtures()
		payload := test.GenericPayload{
			"username": fixtures.UserDeactivated.Username,
			"password": test.PlainTestUserPassword,
		}

		res := test.PerformRequest(t, s, "POST", "/api/v1/auth/login", payload, nil)
		test.RequireHTTPError(t, res, middleware.ErrForbiddenUserDeactivated)
	})
}

func TestPostLoginUserWithoutPassword(t *testing.T) {
	test.WithTestServer(t, func(s *api.Server) {
		fixtures := test.Fixtures()
		payload := test.GenericPayload{
			"username": fixtures.User2.Username,
			"password": test.PlainTestUserPassword,
		}

		fixtures.User2.Password = null.String{}
		rowsAff, err := fixtures.User2.Update(context.Background(), s.DB, boil.Infer())
		require.NoError(t, err)
		require.Equal(t, int64(1), rowsAff)

		res := test.PerformRequest(t, s, "POST", "/api/v1/auth/login", payload, nil)
		test.RequireHTTPError(t, res, httperrors.NewFromEcho(echo.ErrUnauthorized))
	})
}

func TestPostLoginBadRequest(t *testing.T) {
	test.WithTestServer(t, func(s *api.Server) {
		fixtures := test.Fixtures()

		tests := []struct {
			name    string
			payload test.GenericPayload
		}{
			{
				name: "InvalidUsername",
				payload: test.GenericPayload{
					"username": "definitely not an email",
					"password": test.PlainTestUserPassword,
				},
			},
			{
				name: "MissingUsername",
				payload: test.GenericPayload{
					"password": test.PlainTestUserPassword,
				},
			},
			{
				name: "MissingPassword",
				payload: test.GenericPayload{
					"username": fixtures.User1.Username,
				},
			},
			{
				name: "EmptyUsername",
				payload: test.GenericPayload{
					"username": "",
					"password": test.PlainTestUserPassword,
				},
			},
			{
				name: "EmptyPassword",
				payload: test.GenericPayload{
					"username": fixtures.User1.Username,
					"password": "",
				},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				res := test.PerformRequest(t, s, "POST", "/api/v1/auth/login", tt.payload, nil)
				assert.Equal(t, http.StatusBadRequest, res.Result().StatusCode)

				var response httperrors.HTTPValidationError
				test.ParseResponseAndValidate(t, res, &response)

				test.Snapshoter.Save(t, response)
			})
		}

	})
}
