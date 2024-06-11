package db

import (
	"fmt"

	"github.com/lib/pq"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

// WhereIn was a copy from sqlboiler's WHERE IN query helpers since these don't get generated for nullable columns.
// Since sqlboilers IN query helpers will set a param for earch element in the slice we reccomment using this packages IN.
func WhereIn(tableName string, columnName string, slice []string) qm.QueryMod {
	return IN(fmt.Sprintf("%s.%s", tableName, columnName), slice)
}

// IN is a replacement for sqlboilers IN query mod. sqlboilers IN will set a param for
// each element in the slice and we do not reccomend to use this, because it will run into driver and
// database limits. While the sqlboiler IN fails at about ~10000 params this was tested with over 1000000.
func IN(path string, slice []string) qm.QueryMod {
	return qm.Where(fmt.Sprintf("%s = any(?)", path), pq.StringArray(slice))
}

func NIN(path string, slice []string) qm.QueryMod {
	return qm.Where(fmt.Sprintf("%s <> all(?)", path), pq.StringArray(slice))
}
