package test

import (
	"context"
	"time"

	. "allaboutapps.dev/aw/go-starter/internal/models"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

// A common interface for all model instances so they may be inserted via the Inserts() func
type Insertable interface {
	Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error
}

// The main definition which fixtures are available though Fixtures()
type FixtureMap struct {
	User1             *User
	User1AccessToken1 *AccessToken
	User2             *User
}

// We return a function wrapping our fixtures, tests are allowed to manipulate those
// each test (which may run concurrently) can use a fresh copy
func Fixtures() FixtureMap {

	user1 := User{
		ID:       "f6ede5d8-e22a-4ca5-aa12-67821865a3e5",
		IsActive: true,
		Username: null.StringFrom("user1@example.com"),
		Password: null.StringFrom("$argon2id$v=19$m=65536,t=1,p=4$RFO8ulg2c2zloG0029pAUQ$2Po6NUIhVCMm9vivVDuzo7k5KVWfZzJJfeXzC+n+row"),
		Scopes:   []string{"app"},
	}

	user1AccessToken1 := AccessToken{
		Token:      "1cfc27d7-a178-4051-802b-f3ff3967c95c",
		ValidUntil: time.Now().Add(10 * 365 * 24 * time.Hour),
		UserID:     user1.ID,
	}

	user2 := User{
		ID:       "76a79a2b-fbd8-45a0-b35b-671a28a87acf",
		IsActive: user1.IsActive,
		Username: null.StringFrom("user2@example.com"),
		Password: null.StringFrom("$argon2id$v=19$m=65536,t=1,p=4$RFO8ulg2c2zloG0029pAUQ$2Po6NUIhVCMm9vivVDuzo7k5KVWfZzJJfeXzC+n+row"),
		Scopes:   []string{"app"},
	}

	return FixtureMap{
		&user1,
		&user1AccessToken1,
		&user2,
	}
}

// This function defines the order in which the fixtures will be inserted
// into the test database
func Inserts() []Insertable {
	fixtures := Fixtures()

	return []Insertable{
		fixtures.User1,
		fixtures.User1AccessToken1,
		fixtures.User2,
	}
}
