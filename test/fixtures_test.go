package test

import (
	"context"
	"database/sql"
	"testing"

	. "allaboutapps.at/aw/go-mranftl-sample/models"
	_ "github.com/lib/pq"
	"github.com/volatiletech/null"
	"github.com/volatiletech/sqlboiler/boil"
)

func TestFixturesReload(t *testing.T) {

	t.Parallel()

	WithTestDatabase(t, func(db *sql.DB) {
		err := Fixtures().user1.Reload(context.Background(), db)

		if err != nil {
			t.Error("failed to reload")
		}

		// fmt.Println(user1)
	})

}

func TestInsert(t *testing.T) {

	t.Parallel()

	WithTestDatabase(t, func(db *sql.DB) {

		userNew := User{
			ID:       "6d00d09b-fab3-43d8-a163-279fe7ba533e",
			IsActive: true,
			Username: null.StringFrom("userNew@example.com"),
			Password: null.StringFrom("d79f8266fc2942c27af782e2daab181fb6948835fbef33ec4d6c1f3886b3d901381c048a4539bf4ae2ae01d31016e1ca7432917dabd11db27dbc28579dfbfee007f685746788e2bea36ff07fdc45806ff48e9551634b2675534d1169114d5564443c350476df985d6be0cec0fea0d7e3b248c7022a0d9cecbd10d481aecf7b0c97c5756a4f640138fa459729e0921e05c77e77c0a7b025007f2c6d780a1824fe404030f3b9117130818deb370465a9f5f9e32d5508ab06cb8867283b524c6e3a02e7002c95a69fca95f789f82c53d8ab3b2e435983ab25a8544ee9db3cd896c1fc2e9ba24c528754665a0785fe4f25d9d72e47807804d09360e38dc11d02febad845dd233c5227f8cd8c970ab3b117ce9d000b6af0aaf01a6ee6d82627e022ece143642b157f5f883e5fd930deaf4a253e5c5566fbd9ccbd2e9d8a3605dfe979856ca3e743b8d3027448718bba21d73baeb2d459f909b66f257c01bd76e23bf09da8ae1bbab5a6c6cef8b23ec260f25db4d9191a14142db762f4043896c49f6aca4b317f00394715612107b4b360d7120fda3022a7d67aa11380acaa676e43195f527824cf12346898040ed77bb6eb760c3b5861777809b6cfd0e57257584803e18ef580c07eeb0102f04603676702cc83c26aa30bc2bdca2cea51fe5937ff3616467ab9f40a8f1192669f06a387d2d9ed6831fd6371f4c6ad309e9730482a07"),
			Salt:     null.StringFrom("3f9dd9e6d89c7b545ad2f22cf4f68f25188fbf1b2c2fec30c2884c939eead8b5"),
		}

		err := userNew.Insert(context.Background(), db, boil.Infer())

		if err != nil {
			t.Error("failed to insert")
		}

		// fmt.Println(userNew)
	})

}

func TestUpdate(t *testing.T) {

	t.Parallel()

	WithTestDatabase(t, func(db *sql.DB) {

		originalUser1 := Fixtures().user1

		updatedUser1 := User(*originalUser1)

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
	WithTestDatabase(t, func(db *sql.DB) {

		originalUser1 := Fixtures().user1

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
