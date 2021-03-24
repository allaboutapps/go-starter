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
	"github.com/stretchr/testify/assert"
)

func TestPostLogoutSuccess(t *testing.T) {
	test.WithTestServer(t, func(s *api.Server) {
		ctx := context.Background()
		fixtures := test.Fixtures()

		res := test.PerformRequest(t, s, "POST", "/api/v1/auth/logout", nil, test.HeadersWithAuth(t, fixtures.User1AccessToken1.Token))

		assert.Equal(t, http.StatusNoContent, res.Result().StatusCode)

		err := fixtures.User1AccessToken1.Reload(ctx, s.DB)
		assert.Equal(t, sql.ErrNoRows, err)

		err = fixtures.User1RefreshToken1.Reload(ctx, s.DB)
		assert.NoError(t, err)
	})
}

func TestPostLogoutSuccessWithRefreshToken(t *testing.T) {
	test.WithTestServer(t, func(s *api.Server) {
		ctx := context.Background()
		fixtures := test.Fixtures()
		payload := test.GenericPayload{
			"refresh_token": fixtures.User1RefreshToken1.Token,
		}

		res := test.PerformRequest(t, s, "POST", "/api/v1/auth/logout", payload, test.HeadersWithAuth(t, fixtures.User1AccessToken1.Token))

		assert.Equal(t, http.StatusNoContent, res.Result().StatusCode)

		err := fixtures.User1AccessToken1.Reload(ctx, s.DB)
		assert.Equal(t, sql.ErrNoRows, err)

		err = fixtures.User1RefreshToken1.Reload(ctx, s.DB)
		assert.Equal(t, sql.ErrNoRows, err)
	})
}

func TestPostLogoutSuccessWithUnknownRefreshToken(t *testing.T) {
	test.WithTestServer(t, func(s *api.Server) {
		ctx := context.Background()
		fixtures := test.Fixtures()
		payload := test.GenericPayload{
			"refresh_token": "93d8ccd0-be30-4661-a428-cbe74e1a3ffe",
		}

		res := test.PerformRequest(t, s, "POST", "/api/v1/auth/logout", payload, test.HeadersWithAuth(t, fixtures.User1AccessToken1.Token))

		assert.Equal(t, http.StatusNoContent, res.Result().StatusCode)

		err := fixtures.User1AccessToken1.Reload(ctx, s.DB)
		assert.Equal(t, sql.ErrNoRows, err)

		err = fixtures.User1RefreshToken1.Reload(ctx, s.DB)
		assert.NoError(t, err)
	})
}

func TestPostLogoutInvalidRefreshToken(t *testing.T) {
	test.WithTestServer(t, func(s *api.Server) {
		ctx := context.Background()
		fixtures := test.Fixtures()
		payload := test.GenericPayload{
			"refresh_token": "not my refresh token",
		}

		res := test.PerformRequest(t, s, "POST", "/api/v1/auth/logout", payload, test.HeadersWithAuth(t, fixtures.User1AccessToken1.Token))

		assert.Equal(t, http.StatusBadRequest, res.Result().StatusCode)

		var response httperrors.HTTPError
		test.ParseResponseAndValidate(t, res, &response)

		assert.Equal(t, int64(http.StatusBadRequest), *response.Code)
		assert.Equal(t, httperrors.HTTPErrorTypeGeneric, *response.Type)
		assert.Equal(t, http.StatusText(http.StatusBadRequest), *response.Title)
		assert.Empty(t, response.Detail)
		assert.Nil(t, response.Internal)
		assert.Nil(t, response.AdditionalData)

		err := fixtures.User1AccessToken1.Reload(ctx, s.DB)
		assert.NoError(t, err)

		err = fixtures.User1RefreshToken1.Reload(ctx, s.DB)
		assert.NoError(t, err)
	})
}

func TestPostLogoutInvalidAuthToken(t *testing.T) {
	test.WithTestServer(t, func(s *api.Server) {
		res := test.PerformRequest(t, s, "POST", "/api/v1/auth/logout", nil, test.HeadersWithAuth(t, "not my auth token"))

		assert.Equal(t, http.StatusBadRequest, res.Result().StatusCode)

		var response httperrors.HTTPError
		test.ParseResponseAndValidate(t, res, &response)

		assert.Equal(t, *middleware.ErrBadRequestMalformedToken.Code, *response.Code)
		assert.Equal(t, *middleware.ErrBadRequestMalformedToken.Type, *response.Type)
		assert.Equal(t, *middleware.ErrBadRequestMalformedToken.Title, *response.Title)
		assert.Empty(t, response.Detail)
		assert.Nil(t, response.Internal)
		assert.Nil(t, response.AdditionalData)
	})
}

func TestPostLogoutUnknownAuthToken(t *testing.T) {
	test.WithTestServer(t, func(s *api.Server) {
		res := test.PerformRequest(t, s, "POST", "/api/v1/auth/logout", nil, test.HeadersWithAuth(t, "25e8630e-9a41-4f38-8339-373f0c203cef"))

		assert.Equal(t, http.StatusUnauthorized, res.Result().StatusCode)

		var response httperrors.HTTPError
		test.ParseResponseAndValidate(t, res, &response)

		assert.Equal(t, int64(http.StatusUnauthorized), *response.Code)
		assert.Equal(t, httperrors.HTTPErrorTypeGeneric, *response.Type)
		assert.Equal(t, http.StatusText(http.StatusUnauthorized), *response.Title)
		assert.Empty(t, response.Detail)
		assert.Nil(t, response.Internal)
		assert.Nil(t, response.AdditionalData)
	})
}

func TestPostLogoutMissingAuthToken(t *testing.T) {
	test.WithTestServer(t, func(s *api.Server) {
		res := test.PerformRequest(t, s, "POST", "/api/v1/auth/logout", nil, nil)

		assert.Equal(t, http.StatusUnauthorized, res.Result().StatusCode)

		var response httperrors.HTTPError
		test.ParseResponseAndValidate(t, res, &response)

		assert.Equal(t, int64(http.StatusUnauthorized), *response.Code)
		assert.Equal(t, httperrors.HTTPErrorTypeGeneric, *response.Type)
		assert.Equal(t, http.StatusText(http.StatusUnauthorized), *response.Title)
		assert.Empty(t, response.Detail)
		assert.Nil(t, response.Internal)
		assert.Nil(t, response.AdditionalData)
	})
}
