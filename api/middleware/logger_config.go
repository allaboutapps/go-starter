package middleware

import (
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
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
