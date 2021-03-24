package push_test

import (
	"context"
	"database/sql"
	"testing"

	"allaboutapps.dev/aw/go-starter/internal/models"
	"allaboutapps.dev/aw/go-starter/internal/push"
	"allaboutapps.dev/aw/go-starter/internal/push/provider"
	"allaboutapps.dev/aw/go-starter/internal/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

func TestSendMessageSuccess(t *testing.T) {
	test.WithTestPusher(t, func(p *push.Service, db *sql.DB) {
		ctx := context.Background()
		fixtures := test.Fixtures()

		user1 := fixtures.User1

		err := p.SendToUser(ctx, user1, "Hello", "World")
		assert.NoError(t, err)

		tokenCount, err2 := user1.PushTokens().Count(ctx, db)
		require.NoError(t, err2)
		assert.Equal(t, int64(2), tokenCount)
	})
}

func TestSendMessageSuccessWithGenericError(t *testing.T) {
	test.WithTestPusher(t, func(p *push.Service, db *sql.DB) {
		ctx := context.Background()
		fixtures := test.Fixtures()

		user1 := fixtures.User1

		// provoke error from mock provider
		err := p.SendToUser(ctx, user1, "other error", "World")
		assert.NoError(t, err)

		tokenCount, err2 := user1.PushTokens().Count(ctx, db)
		require.NoError(t, err2)
		assert.Equal(t, int64(2), tokenCount)
	})
}

func TestSendMessageWithInvalidToken(t *testing.T) {
	test.WithTestPusher(t, func(p *push.Service, db *sql.DB) {
		ctx := context.Background()
		fixtures := test.Fixtures()
		user1 := fixtures.User1

		user1InvalidPushToken := models.PushToken{
			ID:       "55c37bc8-f245-40b3-bdef-14dee35b10bd",
			Token:    "d5ded380-3285-4243-8a9c-72cc3f063fee",
			UserID:   user1.ID,
			Provider: models.ProviderTypeFCM,
		}
		err := user1InvalidPushToken.Insert(ctx, db, boil.Infer())
		require.NoError(t, err)

		tokenCount, err2 := user1.PushTokens().Count(ctx, db)
		require.NoError(t, err2)
		require.Equal(t, int64(3), tokenCount)

		err = p.SendToUser(ctx, user1, "Hello", "World")
		assert.NoError(t, err)

		tokenCount, err2 = user1.PushTokens().Count(ctx, db)
		require.NoError(t, err2)
		assert.Equal(t, int64(2), tokenCount)
	})
}

func TestSendMessageWithNoProvider(t *testing.T) {
	test.WithTestPusher(t, func(p *push.Service, db *sql.DB) {

		ctx := context.Background()
		fixtures := test.Fixtures()

		p.ResetProviders()
		require.Equal(t, 0, p.GetProviderCount())

		user1 := fixtures.User1

		err := p.SendToUser(ctx, user1, "Hello", "World")
		assert.Error(t, err)

		tokenCount, err2 := user1.PushTokens().Count(ctx, db)
		require.NoError(t, err2)
		assert.Equal(t, int64(2), tokenCount)
	})
}

func TestSendMessageWithMultipleProvider(t *testing.T) {
	test.WithTestPusher(t, func(p *push.Service, db *sql.DB) {
		ctx := context.Background()
		fixtures := test.Fixtures()

		p.ResetProviders()
		require.Equal(t, 0, p.GetProviderCount())

		mockProviderFCM := provider.NewMock(push.ProviderTypeFCM)
		mockProviderAPN := provider.NewMock(push.ProviderTypeAPN)
		p.RegisterProvider(mockProviderAPN)
		p.RegisterProvider(mockProviderFCM)
		user1 := fixtures.User1

		err := p.SendToUser(ctx, user1, "Hello", "World")
		assert.NoError(t, err)

		tokenCount, err2 := user1.PushTokens().Count(ctx, db)
		require.NoError(t, err2)
		assert.Equal(t, int64(1), tokenCount)
	})
}
