package auth_test

import (
	"context"
	"database/sql"
	"net/http"
	"testing"
	"time"

	"allaboutapps.dev/aw/go-starter/internal/api"
	"allaboutapps.dev/aw/go-starter/internal/api/handlers/auth"
	"allaboutapps.dev/aw/go-starter/internal/api/httperrors"
	"allaboutapps.dev/aw/go-starter/internal/api/middleware"
	"allaboutapps.dev/aw/go-starter/internal/models"
	"allaboutapps.dev/aw/go-starter/internal/test"
	"allaboutapps.dev/aw/go-starter/internal/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

func TestPostForgotPasswordCompleteSuccess(t *testing.T) {
	test.WithTestServer(t, func(s *api.Server) {
		ctx := context.Background()
		fixtures := test.Fixtures()

		passwordResetToken := models.PasswordResetToken{
			UserID:     fixtures.User1.ID,
			ValidUntil: time.Now().Add(s.Config.Auth.PasswordResetTokenValidity),
		}

		err := passwordResetToken.Insert(ctx, s.DB, boil.Infer())
		require.NoError(t, err)

		newPassword := "correct horse battery staple"
		payload := test.GenericPayload{
			"token":    passwordResetToken.Token,
			"password": newPassword,
		}

		res := test.PerformRequest(t, s, "POST", "/api/v1/auth/forgot-password/complete", payload, nil)

		assert.Equal(t, http.StatusOK, res.Result().StatusCode)

		var response types.PostLoginResponse
		test.ParseResponseAndValidate(t, res, &response)

		assert.NotEmpty(t, response.AccessToken)
		assert.NotEqual(t, fixtures.User1AccessToken1.Token, response.AccessToken)
		assert.NotEmpty(t, response.RefreshToken)
		assert.NotEqual(t, fixtures.User1RefreshToken1.Token, *response.RefreshToken)
		assert.Equal(t, int64(s.Config.Auth.AccessTokenValidity.Seconds()), *response.ExpiresIn)
		assert.Equal(t, auth.TokenTypeBearer, *response.TokenType)
		test.Snapshoter.Skip([]string{"AccessToken", "RefreshToken"}).Save(t, response)

		err = fixtures.User1AccessToken1.Reload(ctx, s.DB)
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

func TestPostForgotPasswordCompleteUnknownToken(t *testing.T) {
	test.WithTestServer(t, func(s *api.Server) {
		ctx := context.Background()
		fixtures := test.Fixtures()

		newPassword := "correct horse battery staple"
		payload := test.GenericPayload{
			"token":    "fd5c04ea-f39c-49e9-bb40-7f570ed1f66f",
			"password": newPassword,
		}

		res := test.PerformRequest(t, s, "POST", "/api/v1/auth/forgot-password/complete", payload, nil)

		assert.Equal(t, http.StatusNotFound, res.Result().StatusCode)

		var response httperrors.HTTPError
		test.ParseResponseAndValidate(t, res, &response)

		assert.Equal(t, *httperrors.ErrNotFoundTokenNotFound.Code, *response.Code)
		assert.Equal(t, *httperrors.ErrNotFoundTokenNotFound.Type, *response.Type)
		assert.Equal(t, *httperrors.ErrNotFoundTokenNotFound.Title, *response.Title)
		assert.Empty(t, response.Detail)
		assert.Nil(t, response.Internal)
		assert.Nil(t, response.AdditionalData)

		err := fixtures.User1AccessToken1.Reload(ctx, s.DB)
		assert.NoError(t, err)
		err = fixtures.User1RefreshToken1.Reload(ctx, s.DB)
		assert.NoError(t, err)

		cnt, err := fixtures.User1.AccessTokens().Count(ctx, s.DB)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), cnt)

		cnt, err = fixtures.User1.RefreshTokens().Count(ctx, s.DB)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), cnt)

		err = fixtures.User1.Reload(ctx, s.DB)
		assert.NoError(t, err)
		assert.Equal(t, test.HashedTestUserPassword, fixtures.User1.Password.String)
	})
}

func TestPostForgotPasswordCompleteExpiredToken(t *testing.T) {
	test.WithTestServer(t, func(s *api.Server) {
		ctx := context.Background()
		fixtures := test.Fixtures()

		passwordResetToken := models.PasswordResetToken{
			UserID:     fixtures.User1.ID,
			ValidUntil: time.Now().Add(time.Minute * -10),
		}

		err := passwordResetToken.Insert(ctx, s.DB, boil.Infer())
		require.NoError(t, err)

		newPassword := "correct horse battery staple"
		payload := test.GenericPayload{
			"token":    passwordResetToken.Token,
			"password": newPassword,
		}

		res := test.PerformRequest(t, s, "POST", "/api/v1/auth/forgot-password/complete", payload, nil)

		assert.Equal(t, http.StatusConflict, res.Result().StatusCode)

		var response httperrors.HTTPError
		test.ParseResponseAndValidate(t, res, &response)

		assert.Equal(t, *httperrors.ErrConflictTokenExpired.Code, *response.Code)
		assert.Equal(t, *httperrors.ErrConflictTokenExpired.Type, *response.Type)
		assert.Equal(t, *httperrors.ErrConflictTokenExpired.Title, *response.Title)
		assert.Empty(t, response.Detail)
		assert.Nil(t, response.Internal)
		assert.Nil(t, response.AdditionalData)

		err = fixtures.User1AccessToken1.Reload(ctx, s.DB)
		assert.NoError(t, err)
		err = fixtures.User1RefreshToken1.Reload(ctx, s.DB)
		assert.NoError(t, err)

		cnt, err := fixtures.User1.AccessTokens().Count(ctx, s.DB)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), cnt)

		cnt, err = fixtures.User1.RefreshTokens().Count(ctx, s.DB)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), cnt)

		err = fixtures.User1.Reload(ctx, s.DB)
		assert.NoError(t, err)
		assert.Equal(t, test.HashedTestUserPassword, fixtures.User1.Password.String)
	})
}

func TestPostForgotPasswordCompleteDeactivatedUser(t *testing.T) {
	test.WithTestServer(t, func(s *api.Server) {
		ctx := context.Background()
		fixtures := test.Fixtures()

		passwordResetToken := models.PasswordResetToken{
			UserID:     fixtures.UserDeactivated.ID,
			ValidUntil: time.Now().Add(s.Config.Auth.PasswordResetTokenValidity),
		}

		err := passwordResetToken.Insert(ctx, s.DB, boil.Infer())
		require.NoError(t, err)

		newPassword := "correct horse battery staple"
		payload := test.GenericPayload{
			"token":    passwordResetToken.Token,
			"password": newPassword,
		}

		res := test.PerformRequest(t, s, "POST", "/api/v1/auth/forgot-password/complete", payload, nil)

		assert.Equal(t, http.StatusForbidden, res.Result().StatusCode)

		var response httperrors.HTTPError
		test.ParseResponseAndValidate(t, res, &response)

		assert.Equal(t, *middleware.ErrForbiddenUserDeactivated.Code, *response.Code)
		assert.Equal(t, *middleware.ErrForbiddenUserDeactivated.Type, *response.Type)
		assert.Equal(t, *middleware.ErrForbiddenUserDeactivated.Title, *response.Title)
		assert.Empty(t, response.Detail)
		assert.Nil(t, response.Internal)
		assert.Nil(t, response.AdditionalData)

		err = fixtures.UserDeactivatedAccessToken1.Reload(ctx, s.DB)
		assert.NoError(t, err)
		err = fixtures.UserDeactivatedRefreshToken1.Reload(ctx, s.DB)
		assert.NoError(t, err)

		cnt, err := fixtures.UserDeactivated.AccessTokens().Count(ctx, s.DB)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), cnt)

		cnt, err = fixtures.UserDeactivated.RefreshTokens().Count(ctx, s.DB)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), cnt)

		err = fixtures.UserDeactivated.Reload(ctx, s.DB)
		assert.NoError(t, err)
		assert.Equal(t, test.HashedTestUserPassword, fixtures.UserDeactivated.Password.String)
	})
}

func TestPostForgotPasswordCompleteUserWithoutPassword(t *testing.T) {
	test.WithTestServer(t, func(s *api.Server) {
		ctx := context.Background()
		fixtures := test.Fixtures()

		passwordResetToken := models.PasswordResetToken{
			UserID:     fixtures.User2.ID,
			ValidUntil: time.Now().Add(s.Config.Auth.PasswordResetTokenValidity),
		}

		err := passwordResetToken.Insert(ctx, s.DB, boil.Infer())
		require.NoError(t, err)

		newPassword := "correct horse battery staple"
		payload := test.GenericPayload{
			"token":    passwordResetToken.Token,
			"password": newPassword,
		}

		fixtures.User2.Password = null.NewString("", false)
		rowsAff, err := fixtures.User2.Update(context.Background(), s.DB, boil.Infer())
		require.NoError(t, err)
		require.Equal(t, int64(1), rowsAff)

		res := test.PerformRequest(t, s, "POST", "/api/v1/auth/forgot-password/complete", payload, nil)

		assert.Equal(t, http.StatusForbidden, res.Result().StatusCode)

		var response httperrors.HTTPError
		test.ParseResponseAndValidate(t, res, &response)

		assert.Equal(t, *httperrors.ErrForbiddenNotLocalUser.Code, *response.Code)
		assert.Equal(t, *httperrors.ErrForbiddenNotLocalUser.Type, *response.Type)
		assert.Equal(t, *httperrors.ErrForbiddenNotLocalUser.Title, *response.Title)
		assert.Empty(t, response.Detail)
		assert.Nil(t, response.Internal)
		assert.Nil(t, response.AdditionalData)

		err = fixtures.User2AccessToken1.Reload(ctx, s.DB)
		assert.NoError(t, err)
		err = fixtures.User2RefreshToken1.Reload(ctx, s.DB)
		assert.NoError(t, err)

		cnt, err := fixtures.User2.AccessTokens().Count(ctx, s.DB)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), cnt)

		cnt, err = fixtures.User2.RefreshTokens().Count(ctx, s.DB)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), cnt)

		err = fixtures.User2.Reload(ctx, s.DB)
		assert.NoError(t, err)
		assert.False(t, fixtures.User2.Password.Valid)
	})
}

func TestPostForgotPasswordCompleteValidation(t *testing.T) {
	tests := []struct {
		name        string
		payload     test.GenericPayload
		wantKey     string
		wantMessage string
	}{
		{
			name: "MissingToken",
			payload: test.GenericPayload{
				"password": "correct horse battery stable",
			},
			wantKey:     "token",
			wantMessage: "token in body is required",
		},
		{
			name: "MissingPassword",
			payload: test.GenericPayload{
				"token": "7b6e2366-7806-421f-bd56-ffcb39d7b1ee",
			},
			wantKey:     "password",
			wantMessage: "password in body is required",
		},
		{
			name: "InvalidToken",
			payload: test.GenericPayload{
				"token":    "definitelydoesnotexist",
				"password": "correct horse battery stable",
			},
			wantKey:     "token",
			wantMessage: "token in body must be of type uuid4: \"definitelydoesnotexist\"",
		},
		{
			name: "EmptyToken",
			payload: test.GenericPayload{
				"password": "correct horse battery stable",
			},
			wantKey:     "token",
			wantMessage: "token in body is required",
		},
		{
			name: "EmptyPassword",
			payload: test.GenericPayload{
				"token":    "42deb737-fa9c-4e9e-bdce-e33b829c72f7",
				"password": "",
			},
			wantKey:     "password",
			wantMessage: "password in body should be at least 1 chars long",
		},
	}

	test.WithTestServer(t, func(s *api.Server) {
		for _, tt := range tests {
			tt := tt
			t.Run(tt.name, func(t *testing.T) {
				res := test.PerformRequest(t, s, "POST", "/api/v1/auth/forgot-password/complete", tt.payload, nil)

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
				assert.Equal(t, tt.wantKey, *response.ValidationErrors[0].Key)
				assert.Equal(t, "body", *response.ValidationErrors[0].In)
				assert.Equal(t, tt.wantMessage, *response.ValidationErrors[0].Error)
			})
		}
	})
}
