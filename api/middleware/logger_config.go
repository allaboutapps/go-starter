package middleware

import (
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	DefaultLoggerConfig = LoggerConfig{
		Skipper: middleware.DefaultSkipper,
		Logger:  log.Logger,
	}
)

type LoggerConfig struct {
	Skipper middleware.Skipper
	Logger  zerolog.Logger
}
