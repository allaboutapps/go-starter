package auth_test

import (
	"context"
	"database/sql"
	"net/http"
	"testing"
	"time"

	"allaboutapps.dev/aw/go-starter/internal/api"
	"allaboutapps.dev/aw/go-starter/internal/api/httperrors"
	"allaboutapps.dev/aw/go-starter/internal/api/middleware"
	"allaboutapps.dev/aw/go-starter/internal/models"
	"allaboutapps.dev/aw/go-starter/internal/test"
	"allaboutapps.dev/aw/go-starter/internal/test/fixtures"
	"allaboutapps.dev/aw/go-starter/internal/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

func TestPostForgotPasswordCompleteSuccess(t *testing.T) {
	test.WithTestServer(t, func(s *api.Server) {
		ctx := context.Background()
		fix := fixtures.Fixtures()

		passwordResetToken := models.PasswordResetToken{
			UserID:     fix.User1.ID,
			ValidUntil: s.Clock.Now().Add(s.Config.Auth.PasswordResetTokenValidity),
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
		assert.NotEqual(t, fix.User1AccessToken1.Token, response.AccessToken)
		assert.NotEmpty(t, response.RefreshToken)
		assert.NotEqual(t, fix.User1RefreshToken1.Token, *response.RefreshToken)
		test.Snapshoter.Skip([]string{"AccessToken", "RefreshToken"}).Save(t, response)

		err = fix.User1AccessToken1.Reload(ctx, s.DB)
		assert.ErrorIs(t, err, sql.ErrNoRows)
		err = fix.User1RefreshToken1.Reload(ctx, s.DB)
		assert.ErrorIs(t, err, sql.ErrNoRows)

		cnt, err := fix.User1.AccessTokens().Count(ctx, s.DB)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), cnt)

		cnt, err = fix.User1.RefreshTokens().Count(ctx, s.DB)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), cnt)

		err = fix.User1.Reload(ctx, s.DB)
		assert.NoError(t, err)
		assert.NotEqual(t, fixtures.HashedTestUserPassword, fix.User1.Password.String)
	})
}

func TestPostForgotPasswordCompleteUnknownToken(t *testing.T) {
	test.WithTestServer(t, func(s *api.Server) {
		ctx := context.Background()
		fix := fixtures.Fixtures()

		newPassword := "correct horse battery staple"
		payload := test.GenericPayload{
			"token":    "fd5c04ea-f39c-49e9-bb40-7f570ed1f66f",
			"password": newPassword,
		}

		res := test.PerformRequest(t, s, "POST", "/api/v1/auth/forgot-password/complete", payload, nil)
		test.RequireHTTPError(t, res, httperrors.ErrNotFoundTokenNotFound)

		err := fix.User1AccessToken1.Reload(ctx, s.DB)
		assert.NoError(t, err)
		err = fix.User1RefreshToken1.Reload(ctx, s.DB)
		assert.NoError(t, err)

		cnt, err := fix.User1.AccessTokens().Count(ctx, s.DB)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), cnt)

		cnt, err = fix.User1.RefreshTokens().Count(ctx, s.DB)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), cnt)

		err = fix.User1.Reload(ctx, s.DB)
		assert.NoError(t, err)
		assert.Equal(t, fixtures.HashedTestUserPassword, fix.User1.Password.String)
	})
}

func TestPostForgotPasswordCompleteExpiredToken(t *testing.T) {
	test.WithTestServer(t, func(s *api.Server) {
		ctx := context.Background()
		fix := fixtures.Fixtures()

		passwordResetToken := models.PasswordResetToken{
			UserID:     fix.User1.ID,
			ValidUntil: s.Clock.Now().Add(time.Minute * -10),
		}

		err := passwordResetToken.Insert(ctx, s.DB, boil.Infer())
		require.NoError(t, err)

		newPassword := "correct horse battery staple"
		payload := test.GenericPayload{
			"token":    passwordResetToken.Token,
			"password": newPassword,
		}

		res := test.PerformRequest(t, s, "POST", "/api/v1/auth/forgot-password/complete", payload, nil)
		test.RequireHTTPError(t, res, httperrors.ErrConflictTokenExpired)

		err = fix.User1AccessToken1.Reload(ctx, s.DB)
		assert.NoError(t, err)
		err = fix.User1RefreshToken1.Reload(ctx, s.DB)
		assert.NoError(t, err)

		cnt, err := fix.User1.AccessTokens().Count(ctx, s.DB)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), cnt)

		cnt, err = fix.User1.RefreshTokens().Count(ctx, s.DB)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), cnt)

		err = fix.User1.Reload(ctx, s.DB)
		assert.NoError(t, err)
		assert.Equal(t, fixtures.HashedTestUserPassword, fix.User1.Password.String)
	})
}

func TestPostForgotPasswordCompleteDeactivatedUser(t *testing.T) {
	test.WithTestServer(t, func(s *api.Server) {
		ctx := context.Background()
		fix := fixtures.Fixtures()

		passwordResetToken := models.PasswordResetToken{
			UserID:     fix.UserDeactivated.ID,
			ValidUntil: s.Clock.Now().Add(s.Config.Auth.PasswordResetTokenValidity),
		}

		err := passwordResetToken.Insert(ctx, s.DB, boil.Infer())
		require.NoError(t, err)

		newPassword := "correct horse battery staple"
		payload := test.GenericPayload{
			"token":    passwordResetToken.Token,
			"password": newPassword,
		}

		res := test.PerformRequest(t, s, "POST", "/api/v1/auth/forgot-password/complete", payload, nil)
		test.RequireHTTPError(t, res, middleware.ErrForbiddenUserDeactivated)

		err = fix.UserDeactivatedAccessToken1.Reload(ctx, s.DB)
		assert.NoError(t, err)
		err = fix.UserDeactivatedRefreshToken1.Reload(ctx, s.DB)
		assert.NoError(t, err)

		cnt, err := fix.UserDeactivated.AccessTokens().Count(ctx, s.DB)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), cnt)

		cnt, err = fix.UserDeactivated.RefreshTokens().Count(ctx, s.DB)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), cnt)

		err = fix.UserDeactivated.Reload(ctx, s.DB)
		assert.NoError(t, err)
		assert.Equal(t, fixtures.HashedTestUserPassword, fix.UserDeactivated.Password.String)
	})
}

func TestPostForgotPasswordCompleteUserWithoutPassword(t *testing.T) {
	test.WithTestServer(t, func(s *api.Server) {
		ctx := context.Background()
		fix := fixtures.Fixtures()

		passwordResetToken := models.PasswordResetToken{
			UserID:     fix.User2.ID,
			ValidUntil: s.Clock.Now().Add(s.Config.Auth.PasswordResetTokenValidity),
		}

		err := passwordResetToken.Insert(ctx, s.DB, boil.Infer())
		require.NoError(t, err)

		newPassword := "correct horse battery staple"
		payload := test.GenericPayload{
			"token":    passwordResetToken.Token,
			"password": newPassword,
		}

		fix.User2.Password = null.String{}
		rowsAff, err := fix.User2.Update(context.Background(), s.DB, boil.Infer())
		require.NoError(t, err)
		require.Equal(t, int64(1), rowsAff)

		res := test.PerformRequest(t, s, "POST", "/api/v1/auth/forgot-password/complete", payload, nil)
		test.RequireHTTPError(t, res, httperrors.ErrForbiddenNotLocalUser)

		err = fix.User2AccessToken1.Reload(ctx, s.DB)
		assert.NoError(t, err)
		err = fix.User2RefreshToken1.Reload(ctx, s.DB)
		assert.NoError(t, err)

		cnt, err := fix.User2.AccessTokens().Count(ctx, s.DB)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), cnt)

		cnt, err = fix.User2.RefreshTokens().Count(ctx, s.DB)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), cnt)

		err = fix.User2.Reload(ctx, s.DB)
		assert.NoError(t, err)
		assert.False(t, fix.User2.Password.Valid)
	})
}

func TestPostForgotPasswordCompleteBadRequest(t *testing.T) {
	tests := []struct {
		name    string
		payload test.GenericPayload
	}{
		{
			name: "MissingToken",
			payload: test.GenericPayload{
				"password": "correct horse battery stable",
			},
		},
		{
			name: "MissingPassword",
			payload: test.GenericPayload{
				"token": "7b6e2366-7806-421f-bd56-ffcb39d7b1ee",
			},
		},
		{
			name: "InvalidToken",
			payload: test.GenericPayload{
				"token":    "definitelydoesnotexist",
				"password": "correct horse battery stable",
			},
		},
		{
			name: "EmptyToken",
			payload: test.GenericPayload{
				"password": "correct horse battery stable",
				"token":    "",
			},
		},
		{
			name: "EmptyPassword",
			payload: test.GenericPayload{
				"token":    "42deb737-fa9c-4e9e-bdce-e33b829c72f7",
				"password": "",
			},
		},
	}

	test.WithTestServer(t, func(s *api.Server) {
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				res := test.PerformRequest(t, s, "POST", "/api/v1/auth/forgot-password/complete", tt.payload, nil)
				assert.Equal(t, http.StatusBadRequest, res.Result().StatusCode)

				var response httperrors.HTTPValidationError
				test.ParseResponseAndValidate(t, res, &response)

				test.Snapshoter.Save(t, response)
			})
		}
	})
}
