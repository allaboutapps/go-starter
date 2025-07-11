package middleware

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"allaboutapps.dev/aw/go-starter/internal/util"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// RequestBodyLogSkipper defines a function to skip logging certain request bodies.
// Returning true skips logging the payload of the request.
type RequestBodyLogSkipper func(req *http.Request) bool

// DefaultRequestBodyLogSkipper returns true for all requests with Content-Type
// application/x-www-form-urlencoded or multipart/form-data as those might contain
// binary or URL-encoded file uploads unfit for logging purposes.
func DefaultRequestBodyLogSkipper(req *http.Request) bool {
	contentType := req.Header.Get(echo.HeaderContentType)
	switch {
	case strings.HasPrefix(contentType, echo.MIMEApplicationForm),
		strings.HasPrefix(contentType, echo.MIMEMultipartForm):
		return true
	default:
		return false
	}
}

// ResponseBodyLogSkipper defines a function to skip logging certain response bodies.
// Returning true skips logging the payload of the response.
type ResponseBodyLogSkipper func(req *http.Request, res *echo.Response) bool

// DefaultResponseBodyLogSkipper returns false for all responses with Content-Type
// application/json, preventing logging for all other types of payloads as those
// might contain binary or URL-encoded data unfit for logging purposes.
func DefaultResponseBodyLogSkipper(_ *http.Request, res *echo.Response) bool {
	contentType := res.Header().Get(echo.HeaderContentType)
	switch {
	case strings.HasPrefix(contentType, echo.MIMEApplicationJSON):
		return false
	default:
		return true
	}
}

// BodyLogReplacer defines a function to replace certain parts of a body before logging it,
// mainly used to strip sensitive information from a request or response payload.
// The []byte returned should contain a sanitized payload ready for logging.
type BodyLogReplacer func(body []byte) []byte

// DefaultBodyLogReplacer returns the body received without any modifications.
func DefaultBodyLogReplacer(body []byte) []byte {
	return body
}

// HeaderLogReplacer defines a function to replace certain parts of a header before logging it,
// mainly used to strip sensitive information from a request or response header.
// The http.Header returned should be a sanitized copy of the original header as not to modify
// the request or response while logging.
type HeaderLogReplacer func(header http.Header) http.Header

// DefaultHeaderLogReplacer replaces all Authorization, X-CSRF-Token and Proxy-Authorization
// header entries with a redacted string, indicating their presence without revealing actual,
// potentially sensitive values in the logs.
func DefaultHeaderLogReplacer(headers http.Header) http.Header {
	sanitizedHeader := http.Header{}

	for key, value := range headers {
		shouldRedact := strings.EqualFold(key, echo.HeaderAuthorization) ||
			strings.EqualFold(key, echo.HeaderXCSRFToken) ||
			strings.EqualFold(key, "Proxy-Authorization")

		for _, v := range value {
			if shouldRedact {
				sanitizedHeader.Add(key, "*****REDACTED*****")
			} else {
				sanitizedHeader.Add(key, v)
			}
		}
	}

	return sanitizedHeader
}

// QueryLogReplacer defines a function to replace certain parts of a URL query before logging it,
// mainly used to strip sensitive information from a request query.
// The url.Values returned should be a sanitized copy of the original query as not to modify the
// request while logging.
type QueryLogReplacer func(query url.Values) url.Values

// DefaultQueryLogReplacer returns the query received without any modifications.
func DefaultQueryLogReplacer(query url.Values) url.Values {
	return query
}

var (
	DefaultLoggerConfig = LoggerConfig{
		Skipper:                  middleware.DefaultSkipper,
		Level:                    zerolog.DebugLevel,
		LogRequestBody:           false,
		LogRequestHeader:         false,
		LogRequestQuery:          false,
		RequestBodyLogSkipper:    DefaultRequestBodyLogSkipper,
		RequestBodyLogReplacer:   DefaultBodyLogReplacer,
		RequestHeaderLogReplacer: DefaultHeaderLogReplacer,
		RequestQueryLogReplacer:  DefaultQueryLogReplacer,
		LogResponseBody:          false,
		LogResponseHeader:        false,
		ResponseBodyLogSkipper:   DefaultResponseBodyLogSkipper,
		ResponseBodyLogReplacer:  DefaultBodyLogReplacer,
	}
)

type LoggerConfig struct {
	Skipper                   middleware.Skipper
	Level                     zerolog.Level
	LogRequestBody            bool
	LogRequestHeader          bool
	LogRequestQuery           bool
	LogCaller                 bool
	RequestBodyLogSkipper     RequestBodyLogSkipper
	RequestBodyLogReplacer    BodyLogReplacer
	RequestHeaderLogReplacer  HeaderLogReplacer
	RequestQueryLogReplacer   QueryLogReplacer
	LogResponseBody           bool
	LogResponseHeader         bool
	ResponseBodyLogSkipper    ResponseBodyLogSkipper
	ResponseBodyLogReplacer   BodyLogReplacer
	ResponseHeaderLogReplacer HeaderLogReplacer
}

// Logger with default logger output and configuration
func Logger() echo.MiddlewareFunc {
	return LoggerWithConfig(DefaultLoggerConfig, nil)
}

// LoggerWithConfig returns a new MiddlewareFunc which creates a logger with the desired configuration.
// If output is set to nil, the default output is used. If more output params are provided, the first is being used.
func LoggerWithConfig(config LoggerConfig, output ...io.Writer) echo.MiddlewareFunc {
	if config.Skipper == nil {
		config.Skipper = DefaultLoggerConfig.Skipper
	}
	if config.RequestBodyLogSkipper == nil {
		config.RequestBodyLogSkipper = DefaultRequestBodyLogSkipper
	}
	if config.RequestBodyLogReplacer == nil {
		config.RequestBodyLogReplacer = DefaultBodyLogReplacer
	}
	if config.RequestHeaderLogReplacer == nil {
		config.RequestHeaderLogReplacer = DefaultHeaderLogReplacer
	}
	if config.RequestQueryLogReplacer == nil {
		config.RequestQueryLogReplacer = DefaultQueryLogReplacer
	}
	if config.ResponseBodyLogSkipper == nil {
		config.ResponseBodyLogSkipper = DefaultResponseBodyLogSkipper
	}
	if config.ResponseBodyLogReplacer == nil {
		config.ResponseBodyLogReplacer = DefaultBodyLogReplacer
	}
	if config.ResponseHeaderLogReplacer == nil {
		config.ResponseHeaderLogReplacer = DefaultHeaderLogReplacer
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config.Skipper(c) {
				return next(c)
			}

			req := c.Request()
			res := c.Response()

			requestID := req.Header.Get(echo.HeaderXRequestID)
			if len(requestID) == 0 {
				requestID = res.Header().Get(echo.HeaderXRequestID)
			}

			contentLength := req.Header.Get(echo.HeaderContentLength)
			if len(contentLength) == 0 {
				contentLength = "0"
			}

			logger := log.With().
				Dict("req", zerolog.Dict().
					Str("id", requestID).
					Str("host", req.Host).
					Str("method", req.Method).
					Str("url", req.URL.String()).
					Str("bytes_in", contentLength),
				).Logger()

			if len(output) > 0 {
				logger = logger.Output(output[0])
			}

			if config.LogCaller {
				// Caller uses https://pkg.go.dev/runtime#Caller underneath and might decrease the performance.
				logger = logger.With().Caller().Logger()
			}

			le := logger.WithLevel(config.Level)
			req = req.WithContext(logger.WithContext(context.WithValue(req.Context(), util.CTXKeyRequestID, requestID)))

			if config.LogRequestBody && !config.RequestBodyLogSkipper(req) {
				var reqBody []byte
				var err error
				if req.Body != nil {
					reqBody, err = io.ReadAll(req.Body)
					if err != nil {
						logger.Error().Err(err).Msg("Failed to read body while logging request")
						return fmt.Errorf("failed to read body while logging request: %w", err)
					}

					req.Body = io.NopCloser(bytes.NewBuffer(reqBody))
				}

				le = le.Bytes("req_body", config.RequestBodyLogReplacer(reqBody))
			}
			if config.LogRequestHeader {
				header := zerolog.Dict()
				for k, v := range config.RequestHeaderLogReplacer(req.Header) {
					header.Strs(k, v)
				}

				le = le.Dict("req_header", header)
			}
			if config.LogRequestQuery {
				query := zerolog.Dict()
				for k, v := range req.URL.Query() {
					query.Strs(k, v)
				}

				le = le.Dict("req_query", query)
			}

			le.Msg("Request received")

			c.SetRequest(req)

			var resBody bytes.Buffer
			if config.LogResponseBody {
				mw := io.MultiWriter(res.Writer, &resBody)
				writer := &bodyDumpResponseWriter{Writer: mw, ResponseWriter: res.Writer}
				res.Writer = writer
			}

			start := time.Now()
			err := next(c)
			if err != nil {
				c.Error(err)
			}
			stop := time.Now()

			// Retrieve logger from context again since other middlewares might have enhanced it
			ll := util.LogFromEchoContext(c)
			lle := ll.WithLevel(config.Level).
				Dict("res", zerolog.Dict().
					Int("status", res.Status).
					Int64("bytes_out", res.Size).
					TimeDiff("duration_ms", stop, start).
					Err(err),
				)

			if config.LogResponseBody && !config.ResponseBodyLogSkipper(req, res) {
				lle = lle.Bytes("res_body", config.ResponseBodyLogReplacer(resBody.Bytes()))
			}
			if config.LogResponseHeader {
				header := zerolog.Dict()
				for k, v := range config.ResponseHeaderLogReplacer(res.Header()) {
					header.Strs(k, v)
				}

				lle = lle.Dict("res_header", header)
			}

			lle.Msg("Response sent")

			return nil
		}
	}
}

type bodyDumpResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w *bodyDumpResponseWriter) WriteHeader(code int) {
	w.ResponseWriter.WriteHeader(code)
}

func (w *bodyDumpResponseWriter) Write(b []byte) (int, error) {
	n, err := w.Writer.Write(b)
	if err != nil {
		return 0, fmt.Errorf("failed to write response body: %w", err)
	}

	return n, nil
}

func (w *bodyDumpResponseWriter) Flush() {
	flusher, ok := w.ResponseWriter.(http.Flusher)
	if !ok {
		panic(fmt.Sprintf("failed to get flusher as http.Flusher, got %T", w.ResponseWriter))
	}

	flusher.Flush()
}

func (w *bodyDumpResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	hijacker, ok := w.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, fmt.Errorf("failed to get hijacker as http.Hijacker, got %T", w.ResponseWriter)
	}

	conn, rw, err := hijacker.Hijack()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to hijack connection: %w", err)
	}

	return conn, rw, nil
}
