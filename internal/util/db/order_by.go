package db

import (
	"fmt"
	"strings"

	"allaboutapps.dev/aw/go-starter/internal/types"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func OrderBy(orderDir types.OrderDir, path ...string) qm.QueryMod {
	return qm.OrderBy(fmt.Sprintf("%s %s", strings.Join(path, "."), strings.ToUpper(string(orderDir))))
}

func OrderByLower(orderDir types.OrderDir, path ...string) qm.QueryMod {
	return qm.OrderBy(fmt.Sprintf("LOWER(%s) %s", strings.Join(path, "."), strings.ToUpper(string(orderDir))))
}

type OrderByNulls string

const (
	OrderByNullsFirst OrderByNulls = "FIRST"
	OrderByNullsLast  OrderByNulls = "LAST"
)

func OrderByWithNulls(orderDir types.OrderDir, orderByNulls OrderByNulls, path ...string) qm.QueryMod {
	return qm.OrderBy(fmt.Sprintf("%s %s NULLS %s", strings.Join(path, "."), strings.ToUpper(string(orderDir)), orderByNulls))
}

func OrderByLowerWithNulls(orderDir types.OrderDir, orderByNulls OrderByNulls, path ...string) qm.QueryMod {
	return qm.OrderBy(fmt.Sprintf("LOWER(%s) %s NULLS %s", strings.Join(path, "."), strings.ToUpper(string(orderDir)), orderByNulls))
}
