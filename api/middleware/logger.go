package middleware

import (
	"time"

	"github.com/labstack/echo/v4"
)

func Logger() echo.MiddlewareFunc {
	return LoggerWithConfig(DefaultLoggerConfig)
}

func LoggerWithConfig(config LoggerConfig) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config.Skipper(c) {
				return next(c)
			}

			req := c.Request()
			res := c.Response()

			start := time.Now()
			err := next(c)
			if err != nil {
				c.Error(err)
			}
			stop := time.Now()

			id := req.Header.Get(echo.HeaderXRequestID)
			if len(id) == 0 {
				id = res.Header().Get(echo.HeaderXRequestID)
			}

			in := req.Header.Get(echo.HeaderContentLength)
			if len(in) == 0 {
				in = "0"
			}

			config.Logger.
				Debug().
				Str("id", id).
				Str("host", req.Host).
				Str("method", req.Method).
				Str("uri", req.RequestURI).
				Int("status", res.Status).
				Str("bytes_in", in).
				Int64("bytes_out", res.Size).
				TimeDiff("duration", stop, start).
				Err(err).
				Send()

			return nil
		}
	}
}
