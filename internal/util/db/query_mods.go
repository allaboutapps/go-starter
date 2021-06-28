package db

import (
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

// QueryMods represents a slice of query mods, implementing the `queries.Applicator`
// interface to allow for usage with eager loading methods of models. Unfortunately,
// sqlboiler does not import the (identical) type used by the library, so we have to
// declare and "implemented" it ourselves...
type QueryMods []qm.QueryMod

// Apply applies the query mods to the query provided
func (m QueryMods) Apply(q *queries.Query) {
	qm.Apply(q, m...)
}
