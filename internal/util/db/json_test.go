package db_test

import (
	"testing"

	"allaboutapps.dev/aw/go-starter/internal/models"
	"allaboutapps.dev/aw/go-starter/internal/test"
	"allaboutapps.dev/aw/go-starter/internal/util/db"
	"github.com/stretchr/testify/require"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func TestWhereJSONStruct(t *testing.T) {
	age := 42
	filter := struct {
		First        string `json:"firstName"`
		MiddleName   string `json:"-"`
		Lastname     string `json:"lastName"`
		Country      string `json:"country"`
		City         string
		Scopes       []string   `json:"scopes"`
		Age          *int       `json:"age"`
		Height       *float32   `json:"height"`
		PhoneNumbers *[2]string `json:"phoneNumbers"`
		Addresses    []string   `json:"addresses"`
	}{
		First:      "Max",
		MiddleName: "Gustav",
		Lastname:   "Muster",
		Country:    "Austria",
		City:       "Vienna",
		Scopes:     []string{"app", "user_info"},
		Age:        &age,
		PhoneNumbers: &[2]string{
			"+1 206 555 0100",
			"+44 113 496 0000",
		},
	}

	sql, args := buildWhereJSONQuery(t, filter)

	test.Snapshoter.Label("SQL").Save(t, sql)
	test.Snapshoter.Label("Args").Save(t, args)
}

func TestWhereJSONStructComposition(t *testing.T) {
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

	sql, args := buildWhereJSONQuery(t, filter)

	test.Snapshoter.Label("SQL").Save(t, sql)
	test.Snapshoter.Label("Args").Save(t, args)
}

func TestWhereJSONString(t *testing.T) {
	sql, args := buildWhereJSONQuery(t, "https://example.org/users/123/profile")

	test.Snapshoter.Label("SQL").Save(t, sql)
	test.Snapshoter.Label("Args").Save(t, args)
}

func TestWhereJSONPanicEmptyResult(t *testing.T) {
	type privateName struct {
		First      string `json:"firstName"`
		MiddleName string `json:"-"`
		Lastname   string `json:"lastName"`
	}

	filter := struct {
		privateName
		City string
	}{
		privateName: privateName{
			First:      "Max",
			MiddleName: "Gustav",
			Lastname:   "Muster",
		},
		City: "Vienna",
	}

	panicFunc := func() {
		db.WhereJSON("users", "profile", filter)
	}

	require.PanicsWithError(t, "filter resulted in empty query", panicFunc)
}

func TestWhereJSONPanicInvalidFilterType(t *testing.T) {
	panicFunc := func() {
		db.WhereJSON("users", "profile", 1)
	}

	require.PanicsWithError(t, "invalid filter type int", panicFunc)
}

func TestWhereJSONPanicRecursion(t *testing.T) {
	type A struct {
		One string `json:"one"`
	}
	type B struct {
		A
		Two string `json:"two"`
	}
	type C struct {
		B
		Three string `json:"three"`
	}
	type D struct {
		C
		Four string `json:"four"`
	}
	type E struct {
		D
		Five string `json:"five"`
	}
	type F struct {
		E
		Six string `json:"six"`
	}
	type G struct {
		F
		Seven string `json:"seven"`
	}
	type H struct {
		G
		Eight string `json:"eight"`
	}
	type I struct {
		H
		Nine string `json:"nine"`
	}
	type J struct {
		I
		Ten string `json:"ten"`
	}

	filter := struct {
		J
		Country string `json:"country"`
	}{
		J: J{
			I: I{
				H: H{
					G: G{
						F: F{
							E: E{
								D: D{
									C: C{
										B: B{
											A: A{
												One: "1",
											},
											Two: "2",
										},
										Three: "3",
									},
									Four: "4",
								},
								Five: "5",
							},
							Six: "6",
						},
						Seven: "7",
					},
					Eight: "8",
				},
				Nine: "9",
			},
			Ten: "10",
		},
		Country: "Austria",
	}

	panicFunc := func() {
		db.WhereJSON("users", "profile", filter)
	}

	require.PanicsWithError(t, "whereJSON reached maximum recursion (10/10)", panicFunc)
}

func buildWhereJSONQuery(t *testing.T, filter interface{}) (string, []interface{}) {
	t.Helper()

	q := models.NewQuery(
		qm.Select("*"),
		qm.From("users"),
		db.WhereJSON("users", "profile", filter),
	)

	return queries.BuildQuery(q)
}
