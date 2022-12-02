package middleware

// Based on https://github.com/LYY/echo-middleware (MIT License, https://github.com/LYY/echo-middleware/blob/master/LICENSE)
// Ported from Goji's middleware, source:
// https://github.com/zenazn/goji/tree/master/web/middleware (MIT License, https://github.com/zenazn/goji/blob/master/LICENSE)

import (
	"time"

	"github.com/labstack/echo/v4"
	middleware "github.com/labstack/echo/v4/middleware"
)

type (
	// NoCacheConfig defines the config for nocache middleware.
	NoCacheConfig struct {
		// Skipper defines a function to skip middleware.
		Skipper middleware.Skipper
	}
)

var (
	// Unix epoch time
	epoch = time.Unix(0, 0).Format(time.RFC1123)

	// Taken from https://github.com/mytrile/nocache
	noCacheHeaders = map[string]string{
		"Expires":         epoch,
		"Cache-Control":   "no-cache, private, max-age=0",
		"Pragma":          "no-cache",
		"X-Accel-Expires": "0",
	}
	etagHeaders = []string{
		"ETag",
		"If-Modified-Since",
		"If-Match",
		"If-None-Match",
		"If-Range",
		"If-Unmodified-Since",
	}
	// DefaultNoCacheConfig is the default nocache middleware config.
	DefaultNoCacheConfig = NoCacheConfig{
		Skipper: middleware.DefaultSkipper,
	}
)

// NoCache is a simple piece of middleware that sets a number of HTTP headers to prevent
// a router (or subrouter) from being cached by an upstream proxy and/or client.
//
// As per http://wiki.nginx.org/HttpProxyModule - NoCache sets:
//
//	Expires: Thu, 01 Jan 1970 00:00:00 UTC
//	Cache-Control: no-cache, private, max-age=0
//	X-Accel-Expires: 0
//	Pragma: no-cache (for HTTP/1.0 proxies/clients)
func NoCache() echo.MiddlewareFunc {
	return NoCacheWithConfig(DefaultNoCacheConfig)
}

// NoCacheWithConfig returns a nocache middleware with config.
func NoCacheWithConfig(config NoCacheConfig) echo.MiddlewareFunc {
	// Defaults
	if config.Skipper == nil {
		config.Skipper = DefaultNoCacheConfig.Skipper
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			if config.Skipper(c) {
				return next(c)
			}

			req := c.Request()

			// Delete any ETag headers that may have been set
			for _, v := range etagHeaders {
				if req.Header.Get(v) != "" {
					req.Header.Del(v)
				}
			}

			// Set our NoCache headers
			res := c.Response()
			for k, v := range noCacheHeaders {
				res.Header().Set(k, v)
			}

			return next(c)
		}
	}
}
