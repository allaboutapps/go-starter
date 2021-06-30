package db

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/lib/pq"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

const (
	whereJSONMaxLevel = 10
)

// WhereJSON constructs a QueryMod for querying a JSONB column.
//
// The filter interface provided is inspected using reflection, all fields with
// a (non-empty) `json` tag will be added to the query and combined using `AND` -
// fields tagged with `json:"-"` will be ignored as well. Alternatively, a string
// can be provided, performing a string comparison with the database value (the
// stored JSON value does not necessarily have to be a string, but could be an
// integer or similar). The `json` tag's (first) value will be used as the "key"
// for the query, allowing for field renaming or different capitalizations.
//
// At the moment, the root level `filter` value must either be a struct or a string.
// WhereJSON will panic should it encounter a type it cannot process or the filter
// provided results in an empty QueryMod - this allows for easier call chaining
// at the expense of panics in case of incorrect filters being passed.
//
// WhereJSON should support all basic types as well as pointers and array/slices
// of those out of the box, given the Postgres driver can handle their serialization.
// nil pointers are skipped automatically.
// At the moment, struct fields are only supported for composition purposes: if a
// struct is encountered, WhereJSON recursively traverses it (up to 10 levels deep)
// and adds all eligible fields to the top level query.
// Should an array or slice be encountered, their values will be added using the
// `<@` JSONB operator, checking whether all entries existx at the top level within
// the JSON column.
// At the time of writing, no support for special database/HTTP types such as the
// `null` or `strfmt` packages exists - use their respective base types instead.
//
// Whilst WhereJSON was designed to be used with Postgres' JSONB column type, the
// current implementation also supports the JSON type as long as the filter struct
// does not contain any arrays or slices. Note that this compatibility might change
// at some point in the future, so it is advised to use the JSONB data type unless
// your requirements do not allow for it.
func WhereJSON(table string, column string, filter interface{}) qm.QueryMod {
	qms := whereJSON(table, column, filter, 0)
	if len(qms) == 0 {
		panic(errors.New("filter resulted in empty query"))
	}
	return qm.Expr(qms...)
}

func whereJSON(table string, column string, filter interface{}, level int) []qm.QueryMod {
	if level >= whereJSONMaxLevel {
		panic(fmt.Errorf("whereJSON reached maximum recursion (%d/%d)", level, whereJSONMaxLevel))
	}

	qms := make([]qm.QueryMod, 0)

	rt := reflect.TypeOf(filter)
	switch rt.Kind() {
	case reflect.Struct:
		rv := reflect.ValueOf(filter)
		for i := 0; i < rt.NumField(); i++ {
			f := rt.Field(i)

			// skip unexported fields as we cannot retrieve their values
			if len(f.PkgPath) != 0 {
				continue
			}

			k := strings.Split(f.Tag.Get("json"), ",")[0]
			if k == "-" {
				continue
			}

			fs := rv.Field(i)
			if fs.Kind() != reflect.Struct && k == "" {
				continue
			}

			isArray := false
			var v interface{}
			switch fs.Kind() {
			case reflect.Struct:
				qms = append(qms, whereJSON(table, column, fs.Interface(), level+1)...)
				continue
			case reflect.Ptr:
				if !fs.IsValid() || fs.IsNil() {
					continue
				}
				if fs.Elem().Kind() == reflect.Array ||
					fs.Elem().Kind() == reflect.Slice {
					isArray = true
				}
				v = fs.Elem().Interface()
			case reflect.Array,
				reflect.Slice:
				if !fs.IsValid() || fs.IsNil() {
					continue
				}
				isArray = true
				v = fs.Interface()
			default:
				v = fs.Interface()
			}

			if isArray {
				qms = append(qms, qm.Where(fmt.Sprintf("%s.%s->'%s' <@ to_jsonb(?::text[])", table, column, k), pq.Array(v)))
			} else {
				qms = append(qms, qm.Where(fmt.Sprintf("%s.%s->>'%s' = ?", table, column, k), v))
			}
		}
	case reflect.String:
		qms = append(qms, qm.Where(fmt.Sprintf("%s.%s::text = ?", table, column), filter))
	default:
		panic(fmt.Errorf("invalid filter type %v", rt.Kind()))
	}

	return qms
}
