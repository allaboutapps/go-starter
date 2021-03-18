package test_test

import (
	"context"
	"database/sql"
	"testing"

	"allaboutapps.dev/aw/go-starter/internal/models"
	"allaboutapps.dev/aw/go-starter/internal/test"
	_ "github.com/lib/pq"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

func TestFixturesReload(t *testing.T) {
	test.WithTestDatabase(t, func(db *sql.DB) {
		err := test.Fixtures().User1.Reload(context.Background(), db)

		if err != nil {
			t.Error("failed to reload")
		}

		// fmt.Println(user1)
	})

}

func TestInsert(t *testing.T) {
	test.WithTestDatabase(t, func(db *sql.DB) {

		userNew := models.User{
			ID:       "6d00d09b-fab3-43d8-a163-279fe7ba533e",
			IsActive: true,
			Username: null.StringFrom("userNew@example.com"),
			Password: null.StringFrom("$argon2id$v=19$m=65536,t=1,p=4$RFO8ulg2c2zloG0029pAUQ$2Po6NUIhVCMm9vivVDuzo7k5KVWfZzJJfeXzC+n+row"),
			Scopes:   []string{"app"},
		}

		err := userNew.Insert(context.Background(), db, boil.Infer())

		if err != nil {
			t.Error("failed to insert")
		}
	})

}

func TestUpdate(t *testing.T) {
	test.WithTestDatabase(t, func(db *sql.DB) {

		originalUser1 := test.Fixtures().User1

		updatedUser1 := models.User(*originalUser1)

		updatedUser1.Username = null.StringFrom("user_updated@example.com")

		if updatedUser1.Username == originalUser1.Username {
			t.Fatalf("names match!")
		}

		_, err := updatedUser1.Update(context.Background(), db, boil.Infer())

		if err != nil {
			t.Error("failed to update")
		}

		// Attention, this actually mutates our user1 fixture!!!
		err = originalUser1.Reload(context.Background(), db)

		if err != nil {
			t.Error("failed to reload")
		}

		if updatedUser1.Username != originalUser1.Username {
			t.Fatalf("names don't match!")
		}

	})

	// with another testdatabase:
	test.WithTestDatabase(t, func(db *sql.DB) {

		originalUser1 := test.Fixtures().User1

		// ensure our fixture is the same again!
		if originalUser1.Username != null.StringFrom("user1@example.com") {

			err := originalUser1.Reload(context.Background(), db)

			if err != nil {
				t.Error("failed to reload")
			}

			if originalUser1.Username != null.StringFrom("user1@example.com") {
				t.Fatalf("fixture even not the same after reload!")
			}

			t.Fatalf("fixture was modified!")
		}

	})

}
