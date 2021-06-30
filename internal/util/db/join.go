package db

import (
	"fmt"

	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

// InnerJoinWithFilter returns an InnerJoin QueryMod formatted using the provided join tables and columns including an
// additional filter condition. Omitting the optional filter table will use the provided join table as a base for the filter.
func InnerJoinWithFilter(baseTable string, baseColumn string, joinTable string, joinColumn string, filterColumn string, filterValue interface{}, optFilterTable ...string) qm.QueryMod {
	filterTable := joinTable
	if len(optFilterTable) > 0 {
		filterTable = optFilterTable[0]
	}

	return qm.InnerJoin(fmt.Sprintf("%s ON %s.%s=%s.%s AND %s.%s=$1",
		joinTable,
		joinTable,
		joinColumn,
		baseTable,
		baseColumn,
		filterTable,
		filterColumn), filterValue)
}

// InnerJoin returns an InnerJoin QueryMod formatted using the provided join tables and columns.
func InnerJoin(baseTable string, baseColumn string, joinTable string, joinColumn string) qm.QueryMod {
	return qm.InnerJoin(fmt.Sprintf("%s ON %s.%s=%s.%s",
		joinTable,
		joinTable,
		joinColumn,
		baseTable,
		baseColumn))
}

// LeftOuterJoin returns an LeftOuterJoin QueryMod formatted using the provided join tables and columns.
func LeftOuterJoin(baseTable string, baseColumn string, joinTable string, joinColumn string) qm.QueryMod {
	return qm.LeftOuterJoin(fmt.Sprintf("%s ON %s.%s=%s.%s",
		joinTable,
		joinTable,
		joinColumn,
		baseTable,
		baseColumn))
}

// LeftOuterJoinWithFilter returns an LeftOuterJoin QueryMod formatted using the provided join tables and columns including an
// additional filter condition. Omitting the optional filter table will use the provided join table as a base for the filter.
func LeftOuterJoinWithFilter(baseTable string, baseColumn string, joinTable string, joinColumn string, filterColumn string, filterValue interface{}, optFilterTable ...string) qm.QueryMod {
	filterTable := joinTable
	if len(optFilterTable) > 0 {
		filterTable = optFilterTable[0]
	}

	return qm.LeftOuterJoin(fmt.Sprintf("%s ON %s.%s=%s.%s AND %s.%s=$1",
		joinTable,
		joinTable,
		joinColumn,
		baseTable,
		baseColumn,
		filterTable,
		filterColumn), filterValue)
}
