package data

import (
	"context"
	"fmt"

	"allaboutapps.dev/aw/go-starter/internal/util"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

// Live Service fixtures to be applied by manually running the CLI "app db seed"
// Note that these fixtures are not available while testing
// see the separate internal/test/fixtures.go file

type Upsertable interface {
	Upsert(ctx context.Context, exec boil.ContextExecutor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns) error
}

// Mind the declaration order! The fields get upserted exactly in the order they are declared.
type FixtureMap struct{}

func Fixtures() FixtureMap {
	return FixtureMap{}
}

func Upserts() []Upsertable {
	fix := Fixtures()
	upsertableIfc := (*Upsertable)(nil)
	upserts, err := util.GetFieldsImplementing(&fix, upsertableIfc)
	if err != nil {
		panic(fmt.Errorf("failed to get upsertable fixture fields: %w", err))
	}

	return upserts
}
