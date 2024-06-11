package util

import (
	"sort"

	"golang.org/x/text/collate"
	"golang.org/x/text/language"
)

// SortCollateStringSlice is used to sort a slice of strings if the language specific order of caracters is
// important for the order of the string.
// ! The slice passed will be changed.
func SortCollateStringSlice(slice []string, lang language.Tag, options ...collate.Option) {
	if len(options) == 0 {
		options = []collate.Option{collate.IgnoreCase, collate.IgnoreWidth}
	}
	coll := collate.New(lang, options...)

	sort.Slice(slice, func(i int, j int) bool {
		return coll.CompareString(slice[i], slice[j]) < 0
	})
}
