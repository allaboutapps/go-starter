package push_test

import (
	"context"
	"database/sql"
	"net/http"
	"testing"

	"allaboutapps.dev/aw/go-starter/internal/api"
	"allaboutapps.dev/aw/go-starter/internal/api/httperrors"
	"allaboutapps.dev/aw/go-starter/internal/models"
	"allaboutapps.dev/aw/go-starter/internal/test"
	"allaboutapps.dev/aw/go-starter/internal/test/fixtures"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

func TestPutUpdatePushTokenSuccess(t *testing.T) {
	test.WithTestServer(t, func(s *api.Server) {
		ctx := context.Background()
		fix := fixtures.Fixtures()

		//nolint:gosec
		testToken := "869f6deb-73e6-4691-9d40-2a2a794006cf"
		testProvider := "fcm"

		payload := test.GenericPayload{
			"newToken": testToken,
			"provider": testProvider,
		}

		res := test.PerformRequest(t, s, "PUT", "/api/v1/push/token", payload, test.HeadersWithAuth(t, fix.User1AccessToken1.Token))
		assert.Equal(t, http.StatusOK, res.Result().StatusCode)

		newToken, err := models.PushTokens(models.PushTokenWhere.Token.EQ(testToken)).One(ctx, s.DB)
		require.NoError(t, err)
		assert.NotEmpty(t, newToken.ID)
		assert.Equal(t, testToken, newToken.Token)
		assert.Equal(t, testProvider, newToken.Provider)
		assert.Equal(t, fix.User1.ID, newToken.UserID)
	})
}

func TestPutUpdatePushTokenSuccessWithOldToken(t *testing.T) {
	test.WithTestServer(t, func(s *api.Server) {
		ctx := context.Background()
		fix := fixtures.Fixtures()

		//nolint:gosec
		oldToken := "6803ccb4-c91d-47b2-960e-291afa5e29cd"

		oldPushToken := models.PushToken{
			Token:    oldToken,
			Provider: "fcm",
			UserID:   fix.User1.ID,
		}
		err := oldPushToken.Insert(ctx, s.DB, boil.Infer())
		require.NoError(t, err)

		//nolint:gosec
		testToken := "af55b6cf-1fb0-4bb7-960c-25268a5ce7c3"
		testProvider := "fcm"

		payload := test.GenericPayload{
			"newToken": testToken,
			"provider": testProvider,
			"oldToken": oldToken,
		}

		res := test.PerformRequest(t, s, "PUT", "/api/v1/push/token", payload, test.HeadersWithAuth(t, fix.User1AccessToken1.Token))
		assert.Equal(t, http.StatusOK, res.Result().StatusCode)

		newToken, err := models.PushTokens(models.PushTokenWhere.Token.EQ(testToken)).One(ctx, s.DB)
		require.NoError(t, err)
		assert.NotEmpty(t, newToken.ID)
		assert.Equal(t, testToken, newToken.Token)
		assert.Equal(t, testProvider, newToken.Provider)
		assert.Equal(t, fix.User1.ID, newToken.UserID)

		err = oldPushToken.Reload(ctx, s.DB)
		assert.ErrorIs(t, err, sql.ErrNoRows)
	})
}

func TestPutUpdatePushTokenWithDuplicateToken(t *testing.T) {
	test.WithTestServer(t, func(s *api.Server) {
		ctx := context.Background()
		fix := fixtures.Fixtures()

		//nolint:gosec
		oldToken := "6803ccb4-c91d-47b2-960e-291afa5e29cd"

		oldPushToken := models.PushToken{
			Token:    oldToken,
			Provider: "fcm",
			UserID:   fix.User1.ID,
		}
		err := oldPushToken.Insert(ctx, s.DB, boil.Infer())
		require.NoError(t, err)

		testProvider := "fcm"
		payload := test.GenericPayload{
			"newToken": oldToken,
			"provider": testProvider,
			"oldToken": oldToken,
		}

		oldCnt, err := fix.User1.PushTokens().Count(ctx, s.DB)
		assert.NoError(t, err)

		res := test.PerformRequest(t, s, "PUT", "/api/v1/push/token", payload, test.HeadersWithAuth(t, fix.User1AccessToken1.Token))
		test.RequireHTTPError(t, res, httperrors.ErrConflictPushToken)

		err = oldPushToken.Reload(ctx, s.DB)
		assert.NoError(t, err)

		cnt, err := fix.User1.PushTokens().Count(ctx, s.DB)
		assert.NoError(t, err)
		assert.Equal(t, oldCnt, cnt)
	})
}

func TestPutUpdatePushTokenWithOldTokenNotfound(t *testing.T) {
	test.WithTestServer(t, func(s *api.Server) {
		ctx := context.Background()
		fix := fixtures.Fixtures()

		//nolint:gosec
		oldToken := "cc08624a-b40d-4b8e-bbfe-f62aabb47592"

		oldPushToken := models.PushToken{
			Token:    oldToken,
			Provider: "fcm",
			UserID:   fix.User1.ID,
		}
		err := oldPushToken.Insert(ctx, s.DB, boil.Infer())
		require.NoError(t, err)

		oldCnt, err := fix.User1.PushTokens().Count(ctx, s.DB)
		assert.NoError(t, err)

		//nolint:gosec
		testToken := "8e4ad85f-cbb6-4ef3-a455-d9d8bd8917b3"
		testProvider := "fcm"

		payload := test.GenericPayload{
			"newToken": testToken,
			"provider": testProvider,
			"oldToken": "3199aa21-eb41-47dd-9287-338e9e88a5ae",
		}

		res := test.PerformRequest(t, s, "PUT", "/api/v1/push/token", payload, test.HeadersWithAuth(t, fix.User1AccessToken1.Token))
		test.RequireHTTPError(t, res, httperrors.ErrNotFoundOldPushToken)

		newToken, err := models.PushTokens(models.PushTokenWhere.Token.EQ(testToken)).One(ctx, s.DB)
		require.NoError(t, err)
		assert.NotEmpty(t, newToken.ID)
		assert.Equal(t, testToken, newToken.Token)
		assert.Equal(t, testProvider, newToken.Provider)
		assert.Equal(t, fix.User1.ID, newToken.UserID)

		err = oldPushToken.Reload(ctx, s.DB)
		assert.NoError(t, err)

		cnt, err := fix.User1.PushTokens().Count(ctx, s.DB)
		assert.NoError(t, err)
		assert.Equal(t, oldCnt+1, cnt)
	})
}
