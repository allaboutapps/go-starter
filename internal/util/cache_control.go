package util

import (
	"context"
	"strings"
)

type CacheControlDirective uint8

const (
	CacheControlDirectiveNoCache CacheControlDirective = 1 << iota
	CacheControlDirectiveNoStore
)

func (d CacheControlDirective) HasDirective(dir CacheControlDirective) bool { return d&dir != 0 }
func (d *CacheControlDirective) AddDirective(dir CacheControlDirective)     { *d |= dir }
func (d *CacheControlDirective) ClearDirective(dir CacheControlDirective)   { *d &= ^dir }
func (d *CacheControlDirective) ToggleDirective(dir CacheControlDirective)  { *d ^= dir }

func (d CacheControlDirective) String() string {
	res := make([]string, 0)

	if d.HasDirective(CacheControlDirectiveNoCache) {
		res = append(res, "no-cache")
	}
	if d.HasDirective(CacheControlDirectiveNoStore) {
		res = append(res, "no-store")
	}

	return strings.Join(res, "|")
}

func ParseCacheControlDirective(d string) CacheControlDirective {
	parts := strings.Split(d, "=")
	switch strings.ToLower(parts[0]) {
	case "no-cache":
		return CacheControlDirectiveNoCache
	case "no-store":
		return CacheControlDirectiveNoStore
	default:
		return 0
	}
}

func ParseCacheControlHeader(val string) CacheControlDirective {
	res := CacheControlDirective(0)

	directives := strings.Split(val, ",")
	for _, dir := range directives {
		res = res | ParseCacheControlDirective(dir)
	}

	return CacheControlDirective(res)
}

func CacheControlDirectiveFromContext(ctx context.Context) CacheControlDirective {
	d := ctx.Value(CTXKeyCacheControl)
	if d == nil {
		return CacheControlDirective(0)
	}

	directive, ok := d.(CacheControlDirective)
	if !ok {
		return CacheControlDirective(0)
	}

	return directive
}
