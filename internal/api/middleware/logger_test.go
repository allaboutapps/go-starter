package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"allaboutapps.dev/aw/go-starter/internal/api/middleware"
	"allaboutapps.dev/aw/go-starter/internal/util"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func logTestHandler(c echo.Context) error {
	log := util.LogFromEchoContext(c)
	log.Info().Msg("I'm here!")
	return nil
}

func TestLogWithCaller(t *testing.T) {
	cfg := middleware.DefaultLoggerConfig
	cfg.LogCaller = false

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	loggerMW := middleware.LoggerWithConfig(cfg, rec)

	e := echo.New()
	c := e.NewContext(req, rec)

	middlewareChain := loggerMW(logTestHandler)
	require.NoError(t, middlewareChain(c))

	bufSize := 2000
	loggedData := make([]byte, bufSize)
	n, err := rec.Body.Read(loggedData)
	require.NoError(t, err)
	assert.Less(t, n, bufSize)
	// LogCaller was set to false
	assert.NotContains(t, string(loggedData), "caller")

	// now log again with LogCaller set to true
	cfg.LogCaller = true
	loggerMW = middleware.LoggerWithConfig(cfg, rec)

	rec.Flush()

	middlewareChain = loggerMW(logTestHandler)
	require.NoError(t, middlewareChain(c))

	n, err = rec.Body.Read(loggedData)
	require.NoError(t, err)
	assert.Less(t, n, bufSize)
	// logger_test.go:X should match the line where log.Info() is placed in the logTestHandler
	assert.Contains(t, string(loggedData[:n]), `"caller":"/app/internal/api/middleware/logger_test.go:17"`)
}
