package auth_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"allaboutapps.dev/aw/go-starter/internal/api"
	"allaboutapps.dev/aw/go-starter/internal/api/handlers/auth"
	"allaboutapps.dev/aw/go-starter/internal/api/httperrors"
	"allaboutapps.dev/aw/go-starter/internal/models"
	"allaboutapps.dev/aw/go-starter/internal/test"
	"allaboutapps.dev/aw/go-starter/internal/types"
	"github.com/go-openapi/strfmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func TestPostRegisterSuccess(t *testing.T) {
	test.WithTestServer(t, func(s *api.Server) {
		ctx := context.Background()

		username := "usernew@example.com"
		payload := test.GenericPayload{
			"username": username,
			"password": test.PlainTestUserPassword,
		}

		res := test.PerformRequest(t, s, "POST", "/api/v1/auth/register", payload, nil)

		require.Equal(t, http.StatusOK, res.Result().StatusCode)

		var response types.PostLoginResponse
		test.ParseResponseAndValidate(t, res, &response)

		assert.NotEmpty(t, response.AccessToken)
		assert.NotEmpty(t, response.RefreshToken)
		assert.Equal(t, int64(s.Config.Auth.AccessTokenValidity.Seconds()), *response.ExpiresIn)
		assert.Equal(t, auth.TokenTypeBearer, *response.TokenType)

		user, err := models.Users(
			models.UserWhere.Username.EQ(null.StringFrom(username)),
			qm.Load(models.UserRels.AppUserProfile),
			qm.Load(models.UserRels.AccessTokens),
			qm.Load(models.UserRels.RefreshTokens),
		).One(ctx, s.DB)
		assert.NoError(t, err)
		assert.Equal(t, null.StringFrom(username), user.Username)
		assert.Equal(t, true, user.LastAuthenticatedAt.Valid)
		assert.WithinDuration(t, time.Now(), user.LastAuthenticatedAt.Time, time.Second*10)
		assert.EqualValues(t, s.Config.Auth.DefaultUserScopes, user.Scopes)

		assert.NotNil(t, user.R.AppUserProfile)
		assert.Equal(t, false, user.R.AppUserProfile.LegalAcceptedAt.Valid)

		assert.Len(t, user.R.AccessTokens, 1)
		assert.Equal(t, strfmt.UUID4(user.R.AccessTokens[0].Token), *response.AccessToken)
		assert.Len(t, user.R.RefreshTokens, 1)
		assert.Equal(t, strfmt.UUID4(user.R.RefreshTokens[0].Token), *response.RefreshToken)

		res2 := test.PerformRequest(t, s, "POST", "/api/v1/auth/login", payload, nil)

		assert.Equal(t, http.StatusOK, res2.Result().StatusCode)

		var response2 types.PostLoginResponse
		test.ParseResponseAndValidate(t, res2, &response2)

		assert.NotEmpty(t, response2.AccessToken)
		assert.NotEqual(t, response.AccessToken, *response2.AccessToken)
		assert.NotEmpty(t, response2.RefreshToken)
		assert.NotEqual(t, response.RefreshToken, *response2.RefreshToken)
		assert.Equal(t, int64(s.Config.Auth.AccessTokenValidity.Seconds()), *response2.ExpiresIn)
		assert.Equal(t, auth.TokenTypeBearer, *response2.TokenType)
	})
}

func TestPostRegisterAlreadyExists(t *testing.T) {
	test.WithTestServer(t, func(s *api.Server) {
		ctx := context.Background()

		fixtures := test.Fixtures()
		payload := test.GenericPayload{
			"username": fixtures.User1.Username,
			"password": test.PlainTestUserPassword,
		}

		res := test.PerformRequest(t, s, "POST", "/api/v1/auth/register", payload, nil)

		assert.Equal(t, http.StatusConflict, res.Result().StatusCode)

		var response httperrors.HTTPError
		test.ParseResponseAndValidate(t, res, &response)

		assert.Equal(t, *httperrors.ErrConflictUserAlreadyExists.Code, *response.Code)
		assert.Equal(t, *httperrors.ErrConflictUserAlreadyExists.Type, *response.Type)
		assert.Equal(t, *httperrors.ErrConflictUserAlreadyExists.Title, *response.Title)
		assert.Empty(t, response.Detail)
		assert.Nil(t, response.Internal)
		assert.Nil(t, response.AdditionalData)

		user, err := models.Users(
			models.UserWhere.Username.EQ(fixtures.User1.Username),
			qm.Load(models.UserRels.AppUserProfile),
			qm.Load(models.UserRels.AccessTokens),
			qm.Load(models.UserRels.RefreshTokens),
		).One(ctx, s.DB)
		assert.NoError(t, err)
		assert.Equal(t, user.ID, fixtures.User1.ID)

		assert.NotNil(t, user.R.AppUserProfile)
		assert.Len(t, user.R.AccessTokens, 1)
		assert.Len(t, user.R.RefreshTokens, 1)
	})
}

func TestPostRegisterMissingUsername(t *testing.T) {
	test.WithTestServer(t, func(s *api.Server) {
		payload := test.GenericPayload{
			"password": test.PlainTestUserPassword,
		}

		res := test.PerformRequest(t, s, "POST", "/api/v1/auth/register", payload, nil)

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

func TestPostRegisterMissingPassword(t *testing.T) {
	test.WithTestServer(t, func(s *api.Server) {
		fixtures := test.Fixtures()
		payload := test.GenericPayload{
			"username": fixtures.User1.Username,
		}

		res := test.PerformRequest(t, s, "POST", "/api/v1/auth/register", payload, nil)

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

func TestPostRegisterInvalidUsername(t *testing.T) {
	test.WithTestServer(t, func(s *api.Server) {
		payload := test.GenericPayload{
			"username": "definitely not an email",
			"password": test.PlainTestUserPassword,
		}

		res := test.PerformRequest(t, s, "POST", "/api/v1/auth/register", payload, nil)

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

func TestPostRegisterEmptyUsername(t *testing.T) {
	test.WithTestServer(t, func(s *api.Server) {
		payload := test.GenericPayload{
			"username": "",
			"password": test.PlainTestUserPassword,
		}

		res := test.PerformRequest(t, s, "POST", "/api/v1/auth/register", payload, nil)

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

func TestPostRegisterEmptyPassword(t *testing.T) {
	test.WithTestServer(t, func(s *api.Server) {
		payload := test.GenericPayload{
			"username": "usernew@example.com",
			"password": "",
		}

		res := test.PerformRequest(t, s, "POST", "/api/v1/auth/register", payload, nil)

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

func TestPostRegisterSuccessLowercaseTrimWhitespaces(t *testing.T) {
	test.WithTestServer(t, func(s *api.Server) {
		ctx := context.Background()

		username := " USERNEW@example.com "
		usernameLowerTrimmed := "usernew@example.com"
		payload := test.GenericPayload{
			"username": username,
			"password": test.PlainTestUserPassword,
			"name":     "Trim Whitespaces",
		}

		res := test.PerformRequest(t, s, "POST", "/api/v1/auth/register", payload, nil)

		require.Equal(t, http.StatusOK, res.Result().StatusCode)

		var response types.PostLoginResponse
		test.ParseResponseAndValidate(t, res, &response)

		assert.NotEmpty(t, response.AccessToken)
		assert.NotEmpty(t, response.RefreshToken)
		assert.Equal(t, int64(s.Config.Auth.AccessTokenValidity.Seconds()), *response.ExpiresIn)
		assert.Equal(t, auth.TokenTypeBearer, *response.TokenType)

		user, err := models.Users(
			models.UserWhere.Username.EQ(null.StringFrom(usernameLowerTrimmed)),
			qm.Load(models.UserRels.AppUserProfile),
			qm.Load(models.UserRels.AccessTokens),
			qm.Load(models.UserRels.RefreshTokens),
		).One(ctx, s.DB)
		assert.NoError(t, err)
		assert.Equal(t, null.StringFrom(usernameLowerTrimmed), user.Username)
		assert.Equal(t, true, user.LastAuthenticatedAt.Valid)
		assert.WithinDuration(t, time.Now(), user.LastAuthenticatedAt.Time, time.Second*10)
		assert.EqualValues(t, s.Config.Auth.DefaultUserScopes, user.Scopes)

		assert.NotNil(t, user.R.AppUserProfile)
		assert.Equal(t, false, user.R.AppUserProfile.LegalAcceptedAt.Valid)

		assert.Len(t, user.R.AccessTokens, 1)
		assert.Equal(t, strfmt.UUID4(user.R.AccessTokens[0].Token), *response.AccessToken)
		assert.Len(t, user.R.RefreshTokens, 1)
		assert.Equal(t, strfmt.UUID4(user.R.RefreshTokens[0].Token), *response.RefreshToken)

		res2 := test.PerformRequest(t, s, "POST", "/api/v1/auth/login", payload, nil)

		assert.Equal(t, http.StatusOK, res2.Result().StatusCode)

		var response2 types.PostLoginResponse
		test.ParseResponseAndValidate(t, res2, &response2)

		assert.NotEmpty(t, response2.AccessToken)
		assert.NotEqual(t, response.AccessToken, *response2.AccessToken)
		assert.NotEmpty(t, response2.RefreshToken)
		assert.NotEqual(t, response.RefreshToken, *response2.RefreshToken)
		assert.Equal(t, int64(s.Config.Auth.AccessTokenValidity.Seconds()), *response2.ExpiresIn)
		assert.Equal(t, auth.TokenTypeBearer, *response2.TokenType)
	})
}
