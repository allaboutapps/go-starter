package test

import (
	"fmt"
	"reflect"
	"strings"
)

// GetMapFromStructByTag returns a map of a given struct using a tag name as key and
// the string of the property value as value.
// inspired by: https://stackoverflow.com/questions/55879028/golang-get-structs-field-name-by-json-tag
func GetMapFromStructByTag(tag string, input any) map[string]string {
	res := make(map[string]string)

	inputType := reflect.TypeOf(input)
	if inputType.Kind() != reflect.Struct {
		return res
	}

	val := reflect.ValueOf(input)

	for i := 0; i < inputType.NumField(); i++ {
		f := inputType.Field(i)
		tagvalue := strings.Split(f.Tag.Get(tag), ",")[0] // use split to ignore tag "options" like omitempty, etc.
		if tagvalue == "" {
			continue
		}
		field := val.Field(i)
		if field.Kind() == reflect.Ptr {
			if !field.IsNil() {
				res[tagvalue] = fmt.Sprintf("%v", field.Elem().Interface())
			}
		} else {
			res[tagvalue] = fmt.Sprintf("%v", field.Interface())
		}
	}

	return res
}

func GetMapFromStruct(input any) map[string]string {
	res := make(map[string]string)

	inputType := reflect.TypeOf(input)
	if inputType.Kind() != reflect.Struct {
		return res
	}

	val := reflect.ValueOf(input)

	for i := 0; i < inputType.NumField(); i++ {
		field := inputType.Field(i)
		fieldValue := val.Field(i)
		if fieldValue.Kind() == reflect.Ptr {
			if !fieldValue.IsNil() {
				res[field.Name] = fmt.Sprintf("%v", fieldValue.Elem().Interface())
			}
		} else {
			res[field.Name] = fmt.Sprintf("%v", fieldValue.Interface())
		}
	}

	return res
}
