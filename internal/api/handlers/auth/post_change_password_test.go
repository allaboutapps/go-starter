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
		assert.Equal(t, sql.ErrNoRows, err)
		err = fixtures.User1RefreshToken1.Reload(ctx, s.DB)
		assert.Equal(t, sql.ErrNoRows, err)

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

		assert.Equal(t, http.StatusUnauthorized, res.Result().StatusCode)

		var response httperrors.HTTPError
		test.ParseResponseAndValidate(t, res, &response)

		assert.Equal(t, int64(http.StatusUnauthorized), *response.Code)
		assert.Equal(t, httperrors.HTTPErrorTypeGeneric, *response.Type)
		assert.Equal(t, http.StatusText(http.StatusUnauthorized), *response.Title)
		assert.Empty(t, response.Detail)
		assert.Nil(t, response.Internal)
		assert.Nil(t, response.AdditionalData)

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

		assert.Equal(t, http.StatusForbidden, res.Result().StatusCode)

		var response httperrors.HTTPError
		test.ParseResponseAndValidate(t, res, &response)

		assert.Equal(t, *middleware.ErrForbiddenUserDeactivated.Code, *response.Code)
		assert.Equal(t, *middleware.ErrForbiddenUserDeactivated.Type, *response.Type)
		assert.Equal(t, *middleware.ErrForbiddenUserDeactivated.Title, *response.Title)
		assert.Empty(t, response.Detail)
		assert.Nil(t, response.Internal)
		assert.Nil(t, response.AdditionalData)

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

		fixtures.User2.Password = null.NewString("", false)
		rowsAff, err := fixtures.User2.Update(context.Background(), s.DB, boil.Infer())
		require.NoError(t, err)
		require.Equal(t, int64(1), rowsAff)

		res := test.PerformRequest(t, s, "POST", "/api/v1/auth/change-password", payload, test.HeadersWithAuth(t, fixtures.User2AccessToken1.Token))

		assert.Equal(t, http.StatusForbidden, res.Result().StatusCode)

		var response httperrors.HTTPError
		test.ParseResponseAndValidate(t, res, &response)

		assert.Equal(t, *httperrors.ErrForbiddenNotLocalUser.Code, *response.Code)
		assert.Equal(t, *httperrors.ErrForbiddenNotLocalUser.Type, *response.Type)
		assert.Equal(t, *httperrors.ErrForbiddenNotLocalUser.Title, *response.Title)
		assert.Empty(t, response.Detail)
		assert.Nil(t, response.Internal)
		assert.Nil(t, response.AdditionalData)

		err = fixtures.User1AccessToken1.Reload(ctx, s.DB)
		assert.NoError(t, err)
		err = fixtures.User1RefreshToken1.Reload(ctx, s.DB)
		assert.NoError(t, err)
	})
}

func TestPostChangePasswordMissingCurrentPassword(t *testing.T) {
	test.WithTestServer(t, func(s *api.Server) {
		ctx := context.Background()
		fixtures := test.Fixtures()

		newPassword := "correct horse battery staple"
		payload := test.GenericPayload{
			"newPassword": newPassword,
		}

		res := test.PerformRequest(t, s, "POST", "/api/v1/auth/change-password", payload, test.HeadersWithAuth(t, fixtures.User1AccessToken1.Token))

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
		assert.Equal(t, "currentPassword", *response.ValidationErrors[0].Key)
		assert.Equal(t, "body", *response.ValidationErrors[0].In)
		assert.Equal(t, "currentPassword in body is required", *response.ValidationErrors[0].Error)

		err := fixtures.User1AccessToken1.Reload(ctx, s.DB)
		assert.NoError(t, err)
		err = fixtures.User1RefreshToken1.Reload(ctx, s.DB)
		assert.NoError(t, err)
	})
}

func TestPostChangePasswordMissingNewPassword(t *testing.T) {
	test.WithTestServer(t, func(s *api.Server) {
		ctx := context.Background()
		fixtures := test.Fixtures()

		payload := test.GenericPayload{
			"currentPassword": test.PlainTestUserPassword,
		}

		res := test.PerformRequest(t, s, "POST", "/api/v1/auth/change-password", payload, test.HeadersWithAuth(t, fixtures.User1AccessToken1.Token))

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
		assert.Equal(t, "newPassword", *response.ValidationErrors[0].Key)
		assert.Equal(t, "body", *response.ValidationErrors[0].In)
		assert.Equal(t, "newPassword in body is required", *response.ValidationErrors[0].Error)

		err := fixtures.User1AccessToken1.Reload(ctx, s.DB)
		assert.NoError(t, err)
		err = fixtures.User1RefreshToken1.Reload(ctx, s.DB)
		assert.NoError(t, err)
	})
}

func TestPostChangePasswordEmptyCurrentPassword(t *testing.T) {
	test.WithTestServer(t, func(s *api.Server) {
		ctx := context.Background()
		fixtures := test.Fixtures()

		newPassword := "correct horse battery staple"
		payload := test.GenericPayload{
			"currentPassword": "",
			"newPassword":     newPassword,
		}

		res := test.PerformRequest(t, s, "POST", "/api/v1/auth/change-password", payload, test.HeadersWithAuth(t, fixtures.User1AccessToken1.Token))

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
		assert.Equal(t, "currentPassword", *response.ValidationErrors[0].Key)
		assert.Equal(t, "body", *response.ValidationErrors[0].In)
		assert.Equal(t, "currentPassword in body should be at least 1 chars long", *response.ValidationErrors[0].Error)

		err := fixtures.User1AccessToken1.Reload(ctx, s.DB)
		assert.NoError(t, err)
		err = fixtures.User1RefreshToken1.Reload(ctx, s.DB)
		assert.NoError(t, err)
	})
}

func TestPostChangePasswordEmptyNewPassword(t *testing.T) {
	test.WithTestServer(t, func(s *api.Server) {
		ctx := context.Background()
		fixtures := test.Fixtures()

		payload := test.GenericPayload{
			"currentPassword": test.PlainTestUserPassword,
			"newPassword":     "",
		}

		res := test.PerformRequest(t, s, "POST", "/api/v1/auth/change-password", payload, test.HeadersWithAuth(t, fixtures.User1AccessToken1.Token))

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
		assert.Equal(t, "newPassword", *response.ValidationErrors[0].Key)
		assert.Equal(t, "body", *response.ValidationErrors[0].In)
		assert.Equal(t, "newPassword in body should be at least 1 chars long", *response.ValidationErrors[0].Error)

		err := fixtures.User1AccessToken1.Reload(ctx, s.DB)
		assert.NoError(t, err)
		err = fixtures.User1RefreshToken1.Reload(ctx, s.DB)
		assert.NoError(t, err)
	})
}
