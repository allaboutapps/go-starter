package db_test

import (
	"fmt"

	"allaboutapps.dev/aw/go-starter/internal/models"
	"allaboutapps.dev/aw/go-starter/internal/util/db"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type PublicName struct {
	First string `json:"firstName"`
}

type Name struct {
	PublicName
	MiddleName string `json:"-"`
	Lastname   string `json:"lastName"`
}

type UserFilter struct {
	Name
	Country string `json:"country"`
	City    string
	Scopes  []string `json:"scopes"`
	Age     *int     `json:"age"`
	Height  *float32 `json:"height"`
}

func ExampleWhereJSON() {
	age := 42
	filter := UserFilter{
		Name: Name{
			PublicName: PublicName{
				First: "Max",
			},
			MiddleName: "Gustav",
			Lastname:   "Muster",
		},
		Country: "Austria",
		City:    "Vienna",
		Scopes:  []string{"app", "user_info"},
		Age:     &age,
	}

	q := models.NewQuery(
		qm.Select("*"),
		qm.From("users"),
		db.WhereJSON("users", "profile", filter),
	)

	sql, args := queries.BuildQuery(q)

	fmt.Println(sql)
	fmt.Print("[")
	for i := range args {
		if i < len(args)-1 {
			fmt.Printf("%v, ", args[i])
		} else {
			fmt.Printf("%v", args[i])
		}
	}
	fmt.Println("]")

	// Output:
	// SELECT * FROM "users" WHERE (users.profile->>'firstName' = $1 AND users.profile->>'lastName' = $2 AND users.profile->>'country' = $3 AND users.profile->'scopes' <@ to_jsonb($4::text[]) AND users.profile->>'age' = $5);
	// [Max, Muster, Austria, &[app user_info], 42]
}
