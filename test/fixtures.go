package test

import (
	"context"

	. "allaboutapps.at/aw/go-mranftl-sample/models"
	"github.com/volatiletech/null"
	"github.com/volatiletech/sqlboiler/boil"
)

// A common interface for all model instances so they may be inserted via the Inserts() func
type Insertable interface {
	Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error
}

// The main definition which fixtures are available though Fixtures()
type FixtureMap struct {
	user1 *User
	user2 *User
}

// We return a function wrapping our fixtures, tests are allowed to manipulate those
// each test (which may run concurrently) can use a fresh copy
func Fixtures() FixtureMap {

	user1 := User{
		ID:       "f6ede5d8-e22a-4ca5-aa12-67821865a3e5",
		IsActive: true,
		Username: null.StringFrom("user1@example.com"),
		Password: null.StringFrom("$argon2id$v=19$m=65536,t=1,p=4$RFO8ulg2c2zloG0029pAUQ$2Po6NUIhVCMm9vivVDuzo7k5KVWfZzJJfeXzC+n+row"),
	}

	user2 := User{
		ID:       "76a79a2b-fbd8-45a0-b35b-671a28a87acf",
		IsActive: user1.IsActive,
		Username: null.StringFrom("user2@example.com"),
		Password: null.StringFrom("$argon2id$v=19$m=65536,t=1,p=4$RFO8ulg2c2zloG0029pAUQ$2Po6NUIhVCMm9vivVDuzo7k5KVWfZzJJfeXzC+n+row"),
	}

	return FixtureMap{
		&user1,
		&user2,
	}
}

// This function defines the order in which the fixtures will be inserted
// into the test database
func Inserts() []Insertable {
	fixtures := Fixtures()

	return []Insertable{
		fixtures.user1,
		fixtures.user2,
	}
}
