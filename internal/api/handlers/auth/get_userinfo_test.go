package auth_test

import (
	"net/http"
	"testing"

	"allaboutapps.dev/aw/go-starter/internal/api"
	"allaboutapps.dev/aw/go-starter/internal/models"
	"allaboutapps.dev/aw/go-starter/internal/test"
	"allaboutapps.dev/aw/go-starter/internal/test/fixtures"
	"allaboutapps.dev/aw/go-starter/internal/types"
	"github.com/go-openapi/strfmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetUserInfo(t *testing.T) {
	test.WithTestServer(t, func(s *api.Server) {
		ctx := t.Context()
		fix := fixtures.Fixtures()

		res := test.PerformRequest(t, s, "GET", "/api/v1/auth/userinfo", nil, test.HeadersWithAuth(t, fix.User1AccessToken1.Token))
		require.Equal(t, http.StatusOK, res.Result().StatusCode)

		var response types.GetUserInfoResponse
		test.ParseResponseAndValidate(t, res, &response)

		assert.Equal(t, fix.User1.ID, *response.Sub)
		assert.Equal(t, strfmt.Email(fix.User1.Username.String), response.Email)
		test.Snapshoter.Skip([]string{"UpdatedAt"}).Save(t, response)

		for _, scope := range fix.User1.Scopes {
			assert.Contains(t, response.Scopes, scope)
		}

		appUserProfile, err := models.FindAppUserProfile(ctx, s.DB, fix.User1.ID, models.AccessTokenColumns.UpdatedAt)
		require.NoError(t, err)

		assert.Equal(t, appUserProfile.UpdatedAt.Unix(), *response.UpdatedAt)
	})
}

func TestGetUserInfoMinimal(t *testing.T) {
	test.WithTestServer(t, func(s *api.Server) {
		ctx := t.Context()
		fix := fixtures.Fixtures()

		_, err := models.AppUserProfiles(models.AppUserProfileWhere.UserID.EQ(fix.User1.ID)).DeleteAll(ctx, s.DB)
		require.NoError(t, err)

		res := test.PerformRequest(t, s, "GET", "/api/v1/auth/userinfo", nil, test.HeadersWithAuth(t, fix.User1AccessToken1.Token))
		require.Equal(t, http.StatusOK, res.Result().StatusCode)

		var response types.GetUserInfoResponse
		test.ParseResponseAndValidate(t, res, &response)

		assert.Equal(t, fix.User1.ID, *response.Sub)
		assert.Equal(t, strfmt.Email(fix.User1.Username.String), response.Email)

		for _, scope := range fix.User1.Scopes {
			assert.Contains(t, response.Scopes, scope)
		}

		user, err := models.FindUser(ctx, s.DB, fix.User1.ID, models.UserColumns.UpdatedAt)
		require.NoError(t, err)

		assert.Equal(t, user.UpdatedAt.Unix(), *response.UpdatedAt)
	})
}
