package db

import "github.com/volatiletech/sqlboiler/v4/queries/qm"

// CombineWithOr receives a slice of query mods and returns a new slice with
// a single query mod, combining all other query mods into an OR expression.
func CombineWithOr(qms []qm.QueryMod) []qm.QueryMod {
	if len(qms) == 0 {
		return []qm.QueryMod{}
	}

	if len(qms) == 1 {
		return qms
	}

	q := []qm.QueryMod{qms[0]}
	for _, sq := range qms[1:] {
		q = append(q, qm.Or2(sq))
	}

	return []qm.QueryMod{qm.Expr(q...)}
}
