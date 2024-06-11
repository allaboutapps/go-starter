package db

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

var (
	likeQueryEscapeRegex     = regexp.MustCompile(`(%|_)`)
	likeQueryWhiteSpaceRegex = regexp.MustCompile(`\s+`)
)

// ILike returns a query mod containing a pre-formatted ILIKE clause.
// The value provided is applied directly - to perform a wildcard search,
// enclose the desired search value in `%` as desired before passing it
// to ILike.
// The path provided will be joined to construct the full SQL path used,
// allowing for filtering of values nested across multiple joins if needed.
func ILike(val string, path ...string) qm.QueryMod {
	// ! Attention: we **must** use ? instead of $1 or similar to bind query parameters here since
	// ! other parts of the query might have already defined $1, leading to incorrect parameters
	// ! being inserted. On the contrary to other parts using PG queries, ? actually works with qm.Where.
	return qm.Where(fmt.Sprintf("%s ILIKE ?", strings.Join(path, ".")), val)
}

// ILikeSearch returns a query mod with one or multiple ILIKE clauses in an
// AND expression.
// The query is split on whitespace characters and for each word an escaped
// ILIKE with prefix and suffix wildcard will be generated.
func ILikeSearch(query string, path ...string) qm.QueryMod {
	res := []qm.QueryMod{}

	terms := likeQueryWhiteSpaceRegex.Split(strings.TrimSpace(query), -1)
	for _, t := range terms {
		res = append(res, ILike("%"+EscapeLike(t)+"%", path...))
	}

	return qm.Expr(res...)
}

// EscapeLike escapes a string to be placed in an ILIKE query.
func EscapeLike(val string) string {
	return likeQueryEscapeRegex.ReplaceAllString(val, "\\$1")
}
