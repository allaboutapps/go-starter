package data

import (
	"context"

	"github.com/volatiletech/sqlboiler/v4/boil"
)

// Live Service fixtures to be applied by manually running the CLI "app db seed"
// Note that these fixtures are not available while testing
// see the separate internal/test/fixtures.go file

type Upsertable interface {
	Upsert(ctx context.Context, exec boil.ContextExecutor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns) error
}

type FixtureMap struct{}

func Fixtures() FixtureMap {
	return FixtureMap{}
}

func Upserts() []Upsertable {
	return []Upsertable{}
}
