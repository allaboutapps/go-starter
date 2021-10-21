package db

import (
	"fmt"

	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

// WhereIn is a copy from sqlboiler's WHERE IN query helpers since these don't get generated for nullable columns.
func WhereIn(tableName string, columnName string, slice []string) qm.QueryMod {
	values := make([]interface{}, 0, len(slice))
	for _, value := range slice {
		values = append(values, value)
	}
	return qm.WhereIn(fmt.Sprintf("%s.%s IN ?", tableName, columnName), values...)
}
