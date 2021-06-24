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
				if !fs.IsValid() {
					continue
				}
				if fs.Elem().Kind() == reflect.Array {
					isArray = true
				}
				v = fs.Elem().Interface()
			case reflect.Array,
				reflect.Slice:
				if !fs.IsValid() {
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
