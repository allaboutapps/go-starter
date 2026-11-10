package util_test

import (
	"testing"

	"allaboutapps.dev/aw/go-starter/internal/util"
	"github.com/stretchr/testify/assert"
)

func TestParseCacheControlWhitespace(t *testing.T) {
	want := util.CacheControlDirectiveNoCache | util.CacheControlDirectiveNoStore
	for _, input := range []string{
		"no-cache,no-store",
		"no-cache, no-store",
		"no-store,\tno-cache",
		" no-cache , No-Store ",
		"max-age=0, no-cache, no-store",
	} {
		t.Run(input, func(t *testing.T) {
			assert.Equal(t, want, util.ParseCacheControlHeader(input))
		})
	}
}
