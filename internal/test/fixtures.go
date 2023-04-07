package test

import (
	"context"
	"fmt"
	"time"

	"allaboutapps.dev/aw/go-starter/internal/models"
	"allaboutapps.dev/aw/go-starter/internal/util"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

const (
	PlainTestUserPassword  = "password"
	HashedTestUserPassword = "$argon2id$v=19$m=65536,t=1,p=4$RFO8ulg2c2zloG0029pAUQ$2Po6NUIhVCMm9vivVDuzo7k5KVWfZzJJfeXzC+n+row" //nolint:gosec
)

// Insertable represents a common interface for all model instances so they may be inserted via the Inserts() func
type Insertable interface {
	Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error
}

// The main definition which fixtures are available through Fixtures().
// Mind the declaration order! The fields get inserted exactly in the order they are declared.
type FixtureMap struct {
	User1                         *models.User
	User1AppUserProfile           *models.AppUserProfile
	User1AccessToken1             *models.AccessToken
	User1RefreshToken1            *models.RefreshToken
	User2                         *models.User
	User2AppUserProfile           *models.AppUserProfile
	User2AccessToken1             *models.AccessToken
	User2RefreshToken1            *models.RefreshToken
	UserDeactivated               *models.User
	UserDeactivatedAppUserProfile *models.AppUserProfile
	UserDeactivatedAccessToken1   *models.AccessToken
	UserDeactivatedRefreshToken1  *models.RefreshToken
	User1PushToken                *models.PushToken
	User1PushTokenAPN             *models.PushToken
}

// Fixtures returns a function wrapping our fixtures, which tests are allowed to manipulate.
// Each test (which may run concurrently) receives a fresh copy, preventing side effects between test runs.
func Fixtures() FixtureMap {
	now := time.Now()
	f := FixtureMap{}

	f.User1 = &models.User{
		ID:       "f6ede5d8-e22a-4ca5-aa12-67821865a3e5",
		IsActive: true,
		Username: null.StringFrom("user1@example.com"),
		Password: null.StringFrom(HashedTestUserPassword),
		Scopes:   []string{"app"},
	}

	f.User1AppUserProfile = &models.AppUserProfile{
		UserID:          f.User1.ID,
		LegalAcceptedAt: null.TimeFrom(now.Add(time.Minute * -10)),
	}

	f.User1AccessToken1 = &models.AccessToken{
		Token:      "1cfc27d7-a178-4051-802b-f3ff3967c95c",
		ValidUntil: now.Add(10 * 365 * 24 * time.Hour),
		UserID:     f.User1.ID,
	}

	f.User1RefreshToken1 = &models.RefreshToken{
		Token:  "66412eaf-2b89-404d-bbb5-46c3b8bf1a53",
		UserID: f.User1.ID,
	}

	f.User2 = &models.User{
		ID:       "76a79a2b-fbd8-45a0-b35b-671a28a87acf",
		IsActive: true,
		Username: null.StringFrom("user2@example.com"),
		Password: null.StringFrom(HashedTestUserPassword),
		Scopes:   []string{"app"},
	}

	f.User2AppUserProfile = &models.AppUserProfile{
		UserID:          f.User2.ID,
		LegalAcceptedAt: null.TimeFrom(now.Add(time.Minute * -10)),
	}

	f.User2AccessToken1 = &models.AccessToken{
		Token:      "115d28c5-f585-4fb5-9656-fb321739fee5",
		ValidUntil: now.Add(10 * 365 * 24 * time.Hour),
		UserID:     f.User2.ID,
	}

	f.User2RefreshToken1 = &models.RefreshToken{
		Token:  "ea909c75-63d1-4348-a63c-4bcf8ab334a2",
		UserID: f.User2.ID,
	}

	f.UserDeactivated = &models.User{
		ID:       "d9c0dee9-239e-4323-979a-a5354d289627",
		IsActive: false,
		Username: null.StringFrom("userdeactivated@example.com"),
		Password: null.StringFrom(HashedTestUserPassword),
		Scopes:   []string{"app"},
	}

	f.UserDeactivatedAppUserProfile = &models.AppUserProfile{
		UserID:          f.UserDeactivated.ID,
		LegalAcceptedAt: null.Time{},
	}

	f.UserDeactivatedAccessToken1 = &models.AccessToken{
		Token:      "24d0b38d-387c-400c-80fc-a71d85031d4c",
		ValidUntil: now.Add(10 * 365 * 24 * time.Hour),
		UserID:     f.UserDeactivated.ID,
	}

	f.UserDeactivatedRefreshToken1 = &models.RefreshToken{
		Token:  "b6e13a88-7b18-4f17-b819-71b196be2444",
		UserID: f.UserDeactivated.ID,
	}

	f.User1PushToken = &models.PushToken{
		ID:       "98ad176b-af90-44b7-b991-d9ebfc5dd9a0",
		Token:    "cQ_Qk3ZCCZelUZ_K_Yn2BV:APA91bG4jst5srGYZqBAn_wRfiJUzAOQ4k8tV0sDcV4uas2ln5wNwkE_ebneR5Fqk7GvndZ-h3mWnjWaI8yZ4sVwo8qu_Aztotqup4mlEPNYgFGqTlJ5ltQrJG5oKp4RoYQ_0CeFaymn",
		UserID:   f.User1.ID,
		Provider: models.ProviderTypeFCM,
	}

	f.User1PushTokenAPN = &models.PushToken{
		ID:       "5909b472-86f8-4d15-bb63-d49f4fad41a3",
		Token:    "0a863a72-d391-4217-9f26-388801684744",
		UserID:   f.User1.ID,
		Provider: models.ProviderTypeApn,
	}

	return f
}

// Inserts defines the order in which the fixtures will be inserted
// into the test database
func Inserts() []Insertable {
	fix := Fixtures()
	insertableIfc := (*Insertable)(nil)
	inserts, err := util.GetFieldsImplementing(&fix, insertableIfc)
	if err != nil {
		panic(fmt.Errorf("failed to get insertable fixture fields: %w", err))
	}

	return inserts
}
