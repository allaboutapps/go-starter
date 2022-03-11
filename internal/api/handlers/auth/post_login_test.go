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

func TestPostLoginInvalidCredentials(t *testing.T) {
	test.WithTestServer(t, func(s *api.Server) {
		fixtures := test.Fixtures()
		payload := test.GenericPayload{
			"username": fixtures.User1.Username,
			"password": "not my password",
		}

		res := test.PerformRequest(t, s, "POST", "/api/v1/auth/login", payload, nil)

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

func TestPostLoginUnknownUser(t *testing.T) {
	test.WithTestServer(t, func(s *api.Server) {
		payload := test.GenericPayload{
			"username": "definitelydoesnotexist@example.com",
			"password": test.PlainTestUserPassword,
		}

		res := test.PerformRequest(t, s, "POST", "/api/v1/auth/login", payload, nil)

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

func TestPostLoginDeactivatedUser(t *testing.T) {
	test.WithTestServer(t, func(s *api.Server) {
		fixtures := test.Fixtures()
		payload := test.GenericPayload{
			"username": fixtures.UserDeactivated.Username,
			"password": test.PlainTestUserPassword,
		}

		res := test.PerformRequest(t, s, "POST", "/api/v1/auth/login", payload, nil)

		assert.Equal(t, http.StatusForbidden, res.Result().StatusCode)

		var response httperrors.HTTPError
		test.ParseResponseAndValidate(t, res, &response)

		assert.Equal(t, *middleware.ErrForbiddenUserDeactivated.Code, *response.Code)
		assert.Equal(t, *middleware.ErrForbiddenUserDeactivated.Type, *response.Type)
		assert.Equal(t, *middleware.ErrForbiddenUserDeactivated.Title, *response.Title)
		assert.Empty(t, response.Detail)
		assert.Nil(t, response.Internal)
		assert.Nil(t, response.AdditionalData)
	})
}

func TestPostLoginUserWithoutPassword(t *testing.T) {
	test.WithTestServer(t, func(s *api.Server) {
		fixtures := test.Fixtures()
		payload := test.GenericPayload{
			"username": fixtures.User2.Username,
			"password": test.PlainTestUserPassword,
		}

		fixtures.User2.Password = null.NewString("", false)
		rowsAff, err := fixtures.User2.Update(context.Background(), s.DB, boil.Infer())
		require.NoError(t, err)
		require.Equal(t, int64(1), rowsAff)

		res := test.PerformRequest(t, s, "POST", "/api/v1/auth/login", payload, nil)

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

func TestPostLoginInvalidUsername(t *testing.T) {
	test.WithTestServer(t, func(s *api.Server) {
		payload := test.GenericPayload{
			"username": "definitely not an email",
			"password": test.PlainTestUserPassword,
		}

		res := test.PerformRequest(t, s, "POST", "/api/v1/auth/login", payload, nil)

		assert.Equal(t, http.StatusBadRequest, res.Result().StatusCode)

		var response httperrors.HTTPValidationError
		test.ParseResponseAndValidate(t, res, &response)

		assert.Equal(t, int64(http.StatusBadRequest), *response.Code)
		assert.Equal(t, httperrors.HTTPErrorTypeGeneric, *response.Type)
		assert.Equal(t, http.StatusText(http.StatusBadRequest), *response.Title)
		assert.Empty(t, response.Detail)
		assert.Nil(t, response.Internal)
		assert.Nil(t, response.AdditionalData)
		assert.NotEmpty(t, response.ValidationErrors)
		assert.Equal(t, "username", *response.ValidationErrors[0].Key)
		assert.Equal(t, "body", *response.ValidationErrors[0].In)
		assert.Equal(t, "username in body must be of type email: \"definitely not an email\"", *response.ValidationErrors[0].Error)
	})
}

func TestPostLoginMissingUsername(t *testing.T) {
	test.WithTestServer(t, func(s *api.Server) {
		payload := test.GenericPayload{
			"password": test.PlainTestUserPassword,
		}

		res := test.PerformRequest(t, s, "POST", "/api/v1/auth/login", payload, nil)

		assert.Equal(t, http.StatusBadRequest, res.Result().StatusCode)

		var response httperrors.HTTPValidationError
		test.ParseResponseAndValidate(t, res, &response)

		assert.Equal(t, int64(http.StatusBadRequest), *response.Code)
		assert.Equal(t, httperrors.HTTPErrorTypeGeneric, *response.Type)
		assert.Equal(t, http.StatusText(http.StatusBadRequest), *response.Title)
		assert.Empty(t, response.Detail)
		assert.Nil(t, response.Internal)
		assert.Nil(t, response.AdditionalData)
		assert.NotEmpty(t, response.ValidationErrors)
		assert.Equal(t, "username", *response.ValidationErrors[0].Key)
		assert.Equal(t, "body", *response.ValidationErrors[0].In)
		assert.Equal(t, "username in body is required", *response.ValidationErrors[0].Error)
	})
}

func TestPostLoginMissingPassword(t *testing.T) {
	test.WithTestServer(t, func(s *api.Server) {
		fixtures := test.Fixtures()
		payload := test.GenericPayload{
			"username": fixtures.User1.Username,
		}

		res := test.PerformRequest(t, s, "POST", "/api/v1/auth/login", payload, nil)

		assert.Equal(t, http.StatusBadRequest, res.Result().StatusCode)

		var response httperrors.HTTPValidationError
		test.ParseResponseAndValidate(t, res, &response)

		assert.Equal(t, int64(http.StatusBadRequest), *response.Code)
		assert.Equal(t, httperrors.HTTPErrorTypeGeneric, *response.Type)
		assert.Equal(t, http.StatusText(http.StatusBadRequest), *response.Title)
		assert.Empty(t, response.Detail)
		assert.Nil(t, response.Internal)
		assert.Nil(t, response.AdditionalData)
		assert.NotEmpty(t, response.ValidationErrors)
		assert.Equal(t, "password", *response.ValidationErrors[0].Key)
		assert.Equal(t, "body", *response.ValidationErrors[0].In)
		assert.Equal(t, "password in body is required", *response.ValidationErrors[0].Error)
	})
}

func TestPostLoginEmptyUsername(t *testing.T) {
	test.WithTestServer(t, func(s *api.Server) {
		payload := test.GenericPayload{
			"username": "",
			"password": test.PlainTestUserPassword,
		}

		res := test.PerformRequest(t, s, "POST", "/api/v1/auth/login", payload, nil)

		assert.Equal(t, http.StatusBadRequest, res.Result().StatusCode)

		var response httperrors.HTTPValidationError
		test.ParseResponseAndValidate(t, res, &response)

		assert.Equal(t, int64(http.StatusBadRequest), *response.Code)
		assert.Equal(t, httperrors.HTTPErrorTypeGeneric, *response.Type)
		assert.Equal(t, http.StatusText(http.StatusBadRequest), *response.Title)
		assert.Empty(t, response.Detail)
		assert.Nil(t, response.Internal)
		assert.Nil(t, response.AdditionalData)
		assert.NotEmpty(t, response.ValidationErrors)
		assert.Equal(t, "username", *response.ValidationErrors[0].Key)
		assert.Equal(t, "body", *response.ValidationErrors[0].In)
		assert.Equal(t, "username in body should be at least 1 chars long", *response.ValidationErrors[0].Error)
	})
}

func TestPostLoginEmptyPassword(t *testing.T) {
	test.WithTestServer(t, func(s *api.Server) {
		fixtures := test.Fixtures()
		payload := test.GenericPayload{
			"username": fixtures.User1.Username,
			"password": "",
		}

		res := test.PerformRequest(t, s, "POST", "/api/v1/auth/login", payload, nil)

		assert.Equal(t, http.StatusBadRequest, res.Result().StatusCode)

		var response httperrors.HTTPValidationError
		test.ParseResponseAndValidate(t, res, &response)

		assert.Equal(t, int64(http.StatusBadRequest), *response.Code)
		assert.Equal(t, httperrors.HTTPErrorTypeGeneric, *response.Type)
		assert.Equal(t, http.StatusText(http.StatusBadRequest), *response.Title)
		assert.Empty(t, response.Detail)
		assert.Nil(t, response.Internal)
		assert.Nil(t, response.AdditionalData)
		assert.NotEmpty(t, response.ValidationErrors)
		assert.Equal(t, "password", *response.ValidationErrors[0].Key)
		assert.Equal(t, "body", *response.ValidationErrors[0].In)
		assert.Equal(t, "password in body should be at least 1 chars long", *response.ValidationErrors[0].Error)
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
