package db

import (
	"fmt"
	"strings"

	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func ILike(val string, path ...string) qm.QueryMod {
	// ! Attention: we **must** use ? instead of $1 or similar to bind query parameters here since
	// ! other parts of the query might have already defined $1, leading to incorrect parameters
	// ! being inserted. On the contrary to other parts using PG queries, ? actually works with qm.Where.
	return qm.Where(fmt.Sprintf("%s ILIKE ?", strings.Join(path, ".")), val)
}
