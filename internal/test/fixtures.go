package test

import (
	"context"
	"time"

	. "allaboutapps.dev/aw/go-starter/internal/models"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

const (
	PlainTestUserPassword  = "password"
	HashedTestUserPassword = "$argon2id$v=19$m=65536,t=1,p=4$RFO8ulg2c2zloG0029pAUQ$2Po6NUIhVCMm9vivVDuzo7k5KVWfZzJJfeXzC+n+row"
)

// A common interface for all model instances so they may be inserted via the Inserts() func
type Insertable interface {
	Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error
}

// The main definition which fixtures are available though Fixtures()
type FixtureMap struct {
	User1                         *User
	User1AppUserProfile           *AppUserProfile
	User1AccessToken1             *AccessToken
	User1RefreshToken1            *RefreshToken
	User2                         *User
	User2AppUserProfile           *AppUserProfile
	User2AccessToken1             *AccessToken
	User2RefreshToken1            *RefreshToken
	UserDeactivated               *User
	UserDeactivatedAppUserProfile *AppUserProfile
	UserDeactivatedAccessToken1   *AccessToken
	UserDeactivatedRefreshToken1  *RefreshToken
}

// We return a function wrapping our fixtures, tests are allowed to manipulate those
// each test (which may run concurrently) can use a fresh copy
func Fixtures() FixtureMap {

	user1 := User{
		ID:       "f6ede5d8-e22a-4ca5-aa12-67821865a3e5",
		IsActive: true,
		Username: null.StringFrom("user1@example.com"),
		Password: null.StringFrom(HashedTestUserPassword),
		Scopes:   []string{"app"},
	}

	user1AppUserProfile := AppUserProfile{
		UserID:          user1.ID,
		LegalAcceptedAt: null.TimeFrom(time.Now().Add(time.Minute * -10)),
		HasGDPROptOut:   false,
	}

	user1AccessToken1 := AccessToken{
		Token:      "1cfc27d7-a178-4051-802b-f3ff3967c95c",
		ValidUntil: time.Now().Add(10 * 365 * 24 * time.Hour),
		UserID:     user1.ID,
	}

	user1RefreshToken1 := RefreshToken{
		Token:  "66412eaf-2b89-404d-bbb5-46c3b8bf1a53",
		UserID: user1.ID,
	}

	user2 := User{
		ID:       "76a79a2b-fbd8-45a0-b35b-671a28a87acf",
		IsActive: true,
		Username: null.StringFrom("user2@example.com"),
		Password: null.StringFrom(HashedTestUserPassword),
		Scopes:   []string{"app"},
	}

	user2AppUserProfile := AppUserProfile{
		UserID:          user2.ID,
		LegalAcceptedAt: null.TimeFrom(time.Now().Add(time.Minute * -10)),
		HasGDPROptOut:   true,
	}

	user2AccessToken1 := AccessToken{
		Token:      "115d28c5-f585-4fb5-9656-fb321739fee5",
		ValidUntil: time.Now().Add(10 * 365 * 24 * time.Hour),
		UserID:     user2.ID,
	}

	user2RefreshToken1 := RefreshToken{
		Token:  "ea909c75-63d1-4348-a63c-4bcf8ab334a2",
		UserID: user2.ID,
	}

	userDeactivated := User{
		ID:       "d9c0dee9-239e-4323-979a-a5354d289627",
		IsActive: false,
		Username: null.StringFrom("userdeactivated@example.com"),
		Password: null.StringFrom(HashedTestUserPassword),
		Scopes:   []string{"app"},
	}

	userDeactivatedAppUserProfile := AppUserProfile{
		UserID:          userDeactivated.ID,
		LegalAcceptedAt: null.Time{},
		HasGDPROptOut:   false,
	}

	userDeactivatedAccessToken1 := AccessToken{
		Token:      "24d0b38d-387c-400c-80fc-a71d85031d4c",
		ValidUntil: time.Now().Add(10 * 365 * 24 * time.Hour),
		UserID:     userDeactivated.ID,
	}

	userDeactivatedRefreshToken1 := RefreshToken{
		Token:  "b6e13a88-7b18-4f17-b819-71b196be2444",
		UserID: userDeactivated.ID,
	}

	return FixtureMap{
		&user1,
		&user1AppUserProfile,
		&user1AccessToken1,
		&user1RefreshToken1,
		&user2,
		&user2AppUserProfile,
		&user2AccessToken1,
		&user2RefreshToken1,
		&userDeactivated,
		&userDeactivatedAppUserProfile,
		&userDeactivatedAccessToken1,
		&userDeactivatedRefreshToken1,
	}
}

// This function defines the order in which the fixtures will be inserted
// into the test database
func Inserts() []Insertable {
	fixtures := Fixtures()

	return []Insertable{
		fixtures.User1,
		fixtures.User1AppUserProfile,
		fixtures.User1AccessToken1,
		fixtures.User1RefreshToken1,
		fixtures.User2,
		fixtures.User2AppUserProfile,
		fixtures.User2AccessToken1,
		fixtures.User2RefreshToken1,
		fixtures.UserDeactivated,
		fixtures.UserDeactivatedAppUserProfile,
		fixtures.UserDeactivatedAccessToken1,
		fixtures.UserDeactivatedRefreshToken1,
	}
}
