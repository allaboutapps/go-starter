package middleware

import (
	"context"

	"allaboutapps.dev/aw/go-starter/internal/util"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var (
	DefaultCacheControlConfig = CacheControlConfig{
		Skipper: middleware.DefaultSkipper,
	}
)

type CacheControlConfig struct {
	Skipper middleware.Skipper
}

func CacheControl() echo.MiddlewareFunc {
	return CacheControlWithConfig(DefaultCacheControlConfig)
}

func CacheControlWithConfig(config CacheControlConfig) echo.MiddlewareFunc {
	if config.Skipper == nil {
		config.Skipper = DefaultCacheControlConfig.Skipper
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config.Skipper(c) {
				return next(c)
			}

			cacheControl := c.Request().Header.Get(util.HTTPHeaderCacheControl)
			if len(cacheControl) > 0 {
				directive := util.ParseCacheControlHeader(cacheControl)

				ctx := c.Request().Context()

				l := util.LogFromContext(ctx).With().Str("cacheControl", directive.String()).Logger()
				ctx = l.WithContext(ctx)

				ctx = context.WithValue(ctx, util.CTXKeyCacheControl, directive)

				l.Trace().Msg("Setting cache control directive for request")

				c.SetRequest(c.Request().WithContext(ctx))
			}

			return next(c)
		}
	}
}
