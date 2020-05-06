package middleware

import (
	"time"

	"allaboutapps.dev/aw/go-starter/internal/util"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	DefaultLoggerConfig = LoggerConfig{
		Skipper: middleware.DefaultSkipper,
		Level:   zerolog.DebugLevel,
	}
)

type LoggerConfig struct {
	Skipper middleware.Skipper
	Level   zerolog.Level
}

func Logger() echo.MiddlewareFunc {
	return LoggerWithConfig(DefaultLoggerConfig)
}

func LoggerWithConfig(config LoggerConfig) echo.MiddlewareFunc {
	if config.Skipper == nil {
		config.Skipper = DefaultLoggerConfig.Skipper
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config.Skipper(c) {
				return next(c)
			}

			req := c.Request()
			res := c.Response()

			id := req.Header.Get(echo.HeaderXRequestID)
			if len(id) == 0 {
				id = res.Header().Get(echo.HeaderXRequestID)
			}

			in := req.Header.Get(echo.HeaderContentLength)
			if len(in) == 0 {
				in = "0"
			}

			l := log.With().
				Dict("req", zerolog.Dict().
					Str("id", id).
					Str("host", req.Host).
					Str("method", req.Method).
					Str("url", req.URL.String()).
					Str("bytes_in", in),
				).Logger()
			req = req.WithContext(l.WithContext(req.Context()))

			c.SetRequest(req)

			start := time.Now()
			err := next(c)
			if err != nil {
				c.Error(err)
			}
			stop := time.Now()

			// Retrieve logger from context again since other middlewares might have enhanced it
			ll := util.LogFromEchoContext(c)
			ll.WithLevel(config.Level).
				Dict("res", zerolog.Dict().
					Int("status", res.Status).
					Int64("bytes_out", res.Size).
					TimeDiff("duration_ms", stop, start).
					Err(err),
				).Send()

			return nil
		}
	}
}
