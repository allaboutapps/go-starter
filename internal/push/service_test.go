package push_test

import (
	"database/sql"
	"testing"

	"allaboutapps.dev/aw/go-starter/internal/data/mapper"
	"allaboutapps.dev/aw/go-starter/internal/models"
	"allaboutapps.dev/aw/go-starter/internal/push"
	"allaboutapps.dev/aw/go-starter/internal/push/provider"
	"allaboutapps.dev/aw/go-starter/internal/test"
	"allaboutapps.dev/aw/go-starter/internal/test/fixtures"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

func TestSendMessageSuccess(t *testing.T) {
	test.WithTestPusher(t, func(service *push.Service, db *sql.DB) {
		ctx := t.Context()
		fix := fixtures.Fixtures()

		err := service.SendToUser(ctx, mapper.LocalUserToDTO(fix.User1).Ptr(), "Hello", "World")
		require.NoError(t, err)

		tokenCount, err2 := fix.User1.PushTokens().Count(ctx, db)
		require.NoError(t, err2)
		assert.Equal(t, int64(2), tokenCount)
	})
}

func TestSendMessageSuccessWithGenericError(t *testing.T) {
	test.WithTestPusher(t, func(service *push.Service, db *sql.DB) {
		ctx := t.Context()
		fix := fixtures.Fixtures()

		// provoke error from mock provider
		err := service.SendToUser(ctx, mapper.LocalUserToDTO(fix.User1).Ptr(), "other error", "World")
		require.NoError(t, err)

		tokenCount, err2 := fix.User1.PushTokens().Count(ctx, db)
		require.NoError(t, err2)
		assert.Equal(t, int64(2), tokenCount)
	})
}

func TestSendMessageWithInvalidToken(t *testing.T) {
	test.WithTestPusher(t, func(service *push.Service, db *sql.DB) {
		ctx := t.Context()
		fix := fixtures.Fixtures()

		user1InvalidPushToken := models.PushToken{
			ID:       "55c37bc8-f245-40b3-bdef-14dee35b10bd",
			Token:    "d5ded380-3285-4243-8a9c-72cc3f063fee",
			UserID:   fix.User1.ID,
			Provider: models.ProviderTypeFCM,
		}
		err := user1InvalidPushToken.Insert(ctx, db, boil.Infer())
		require.NoError(t, err)

		tokenCount, err2 := fix.User1.PushTokens().Count(ctx, db)
		require.NoError(t, err2)
		require.Equal(t, int64(3), tokenCount)

		err = service.SendToUser(ctx, mapper.LocalUserToDTO(fix.User1).Ptr(), "Hello", "World")
		require.NoError(t, err)

		tokenCount, err2 = fix.User1.PushTokens().Count(ctx, db)
		require.NoError(t, err2)
		assert.Equal(t, int64(2), tokenCount)
	})
}

func TestSendMessageWithNoProvider(t *testing.T) {
	test.WithTestPusher(t, func(service *push.Service, db *sql.DB) {
		ctx := t.Context()
		fix := fixtures.Fixtures()

		service.ResetProviders()
		require.Equal(t, 0, service.GetProviderCount())

		err := service.SendToUser(ctx, mapper.LocalUserToDTO(fix.User1).Ptr(), "Hello", "World")
		require.Error(t, err)

		tokenCount, err2 := fix.User1.PushTokens().Count(ctx, db)
		require.NoError(t, err2)
		assert.Equal(t, int64(2), tokenCount)
	})
}

func TestSendMessageWithMultipleProvider(t *testing.T) {
	test.WithTestPusher(t, func(service *push.Service, db *sql.DB) {
		ctx := t.Context()
		fix := fixtures.Fixtures()

		service.ResetProviders()
		require.Equal(t, 0, service.GetProviderCount())

		mockProviderFCM := provider.NewMock(push.ProviderTypeFCM)
		mockProviderAPN := provider.NewMock(push.ProviderTypeAPN)
		service.RegisterProvider(mockProviderAPN)
		service.RegisterProvider(mockProviderFCM)

		err := service.SendToUser(ctx, mapper.LocalUserToDTO(fix.User1).Ptr(), "Hello", "World")
		require.NoError(t, err)

		tokenCount, err2 := fix.User1.PushTokens().Count(ctx, db)
		require.NoError(t, err2)
		assert.Equal(t, int64(1), tokenCount)
	})
}
