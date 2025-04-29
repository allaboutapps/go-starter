package auth_test

import (
	"net/http"
	"testing"
	"time"

	"allaboutapps.dev/aw/go-starter/internal/api"
	"allaboutapps.dev/aw/go-starter/internal/api/httperrors"
	"allaboutapps.dev/aw/go-starter/internal/auth"
	"allaboutapps.dev/aw/go-starter/internal/models"
	"allaboutapps.dev/aw/go-starter/internal/test"
	"allaboutapps.dev/aw/go-starter/internal/test/fixtures"
	"allaboutapps.dev/aw/go-starter/internal/types"
	"github.com/go-openapi/strfmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func TestPostRegisterSuccess(t *testing.T) {
	test.WithTestServer(t, func(s *api.Server) {
		ctx := t.Context()

		now := time.Date(2025, 2, 5, 11, 42, 30, 0, time.UTC)
		test.SetMockClock(t, s, now)

		username := "usernew@example.com"
		payload := test.GenericPayload{
			"username": username,
			"password": fixtures.PlainTestUserPassword,
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
		assert.Equal(t, now, user.LastAuthenticatedAt.Time)
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

func TestPostRegisterSuccessLowercaseTrimWhitespaces(t *testing.T) {
	test.WithTestServer(t, func(s *api.Server) {
		ctx := t.Context()

		username := " USERNEW@example.com "
		usernameLowerTrimmed := "usernew@example.com"
		payload := test.GenericPayload{
			"username": username,
			"password": fixtures.PlainTestUserPassword,
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
		assert.WithinDuration(t, s.Clock.Now(), user.LastAuthenticatedAt.Time, time.Second*10)
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
		ctx := t.Context()

		fix := fixtures.Fixtures()
		payload := test.GenericPayload{
			"username": fix.User1.Username,
			"password": fixtures.PlainTestUserPassword,
		}

		res := test.PerformRequest(t, s, "POST", "/api/v1/auth/register", payload, nil)
		test.RequireHTTPError(t, res, httperrors.ErrConflictUserAlreadyExists)

		user, err := models.Users(
			models.UserWhere.Username.EQ(fix.User1.Username),
			qm.Load(models.UserRels.AppUserProfile),
			qm.Load(models.UserRels.AccessTokens),
			qm.Load(models.UserRels.RefreshTokens),
		).One(ctx, s.DB)
		assert.NoError(t, err)
		assert.Equal(t, user.ID, fix.User1.ID)

		assert.NotNil(t, user.R.AppUserProfile)
		assert.Len(t, user.R.AccessTokens, 1)
		assert.Len(t, user.R.RefreshTokens, 1)
	})
}

func TestPostRegisterBadRequest(t *testing.T) {
	test.WithTestServer(t, func(s *api.Server) {
		fix := fixtures.Fixtures()

		tests := []struct {
			name    string
			payload test.GenericPayload
		}{
			{
				name: "MissingUsername",
				payload: test.GenericPayload{
					"password": fixtures.PlainTestUserPassword,
				},
			},
			{
				name: "MissingPassword",
				payload: test.GenericPayload{
					"username": fix.User1.Username,
				},
			},
			{
				name: "InvalidUsername",
				payload: test.GenericPayload{
					"username": "definitely not an email",
					"password": fixtures.PlainTestUserPassword,
				},
			},
			{
				name: "EmptyUsername",
				payload: test.GenericPayload{
					"username": "",
					"password": fixtures.PlainTestUserPassword,
				},
			},
			{
				name: "EmptyPassword",
				payload: test.GenericPayload{
					"username": fix.User1.Username,
					"password": "",
				},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				res := test.PerformRequest(t, s, "POST", "/api/v1/auth/register", tt.payload, nil)
				assert.Equal(t, http.StatusBadRequest, res.Result().StatusCode)

				var response httperrors.HTTPValidationError
				test.ParseResponseAndValidate(t, res, &response)

				test.Snapshoter.Save(t, response)
			})
		}
	})
}
