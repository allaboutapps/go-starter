package util

import "testing"

func TestParseCacheControlWhitespace(t *testing.T) {
	want := CacheControlDirectiveNoCache | CacheControlDirectiveNoStore
	for _, input := range []string{
		"no-cache,no-store",
		"no-cache, no-store",
		"no-store,\tno-cache",
		" no-cache , No-Store ",
		"max-age=0, no-cache, no-store",
	} {
		t.Run(input, func(t *testing.T) {
			if got := ParseCacheControlHeader(input); got != want {
				t.Fatalf("ParseCacheControlHeader(%q) = %d, want %d", input, got, want)
			}
		})
	}
}
