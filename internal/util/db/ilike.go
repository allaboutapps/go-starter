package db

import (
	"fmt"
	"strings"

	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

// ILike returns a query mod containing a pre-formatted ILIKE clause.
// The value provided is applied directly - to perform a wildcard search,
// enclose the desired search value in `%` as desired before passing it
// to ILike.
// The path provided will be joined to construct the full SQL path used,
// allowing for filtering of values nested across multiple joins if needed.
func ILike(val string, path ...string) qm.QueryMod {
	// ! Attention: we **must** use ? instead of $1 or similar to bind query parameters here since
	// ! other parts of the query might have already defined $1, leading to incorrect parameters
	// ! being inserted. On the contrary to other parts using PG queries, ? actually works with qm.Where.
	return qm.Where(fmt.Sprintf("%s ILIKE ?", strings.Join(path, ".")), val)
}
