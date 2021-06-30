package test

import (
	"fmt"
	"reflect"
	"strings"
)

// GetMapFromStructByTag returns a map of a given struct using a tag name as key and
// the string of the property value as value.
// inspired by: https://stackoverflow.com/questions/55879028/golang-get-structs-field-name-by-json-tag
func GetMapFromStructByTag(tag string, s interface{}) map[string]string {
	res := make(map[string]string)

	rt := reflect.TypeOf(s)
	if rt.Kind() != reflect.Struct {
		return res
	}

	rv := reflect.ValueOf(s)

	for i := 0; i < rt.NumField(); i++ {
		f := rt.Field(i)
		v := strings.Split(f.Tag.Get(tag), ",")[0] // use split to ignore tag "options" like omitempty, etc.
		if v == "" {
			continue
		}
		fs := rv.Field(i)
		if fs.Kind() == reflect.Ptr {
			if !fs.IsNil() {
				res[v] = fmt.Sprintf("%v", fs.Elem().Interface())
			}
		} else {
			res[v] = fmt.Sprintf("%v", fs.Interface())
		}
	}

	return res
}

func GetMapFromStruct(s interface{}) map[string]string {
	res := make(map[string]string)

	rt := reflect.TypeOf(s)
	if rt.Kind() != reflect.Struct {
		return res
	}

	rv := reflect.ValueOf(s)

	for i := 0; i < rt.NumField(); i++ {
		f := rt.Field(i)
		fs := rv.Field(i)
		if fs.Kind() == reflect.Ptr {
			if !fs.IsNil() {
				res[f.Name] = fmt.Sprintf("%v", fs.Elem().Interface())
			}
		} else {
			res[f.Name] = fmt.Sprintf("%v", fs.Interface())
		}
	}

	return res
}
