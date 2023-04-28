package router

import (
	"net/http"
	"runtime"
	"strings"

	"allaboutapps.dev/aw/go-starter/internal/api"
	"allaboutapps.dev/aw/go-starter/internal/api/handlers"
	"allaboutapps.dev/aw/go-starter/internal/api/middleware"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"

	// #nosec G108 - pprof handlers (conditionally made available via http.DefaultServeMux)
	"net/http/pprof"
)

func Init(s *api.Server) {
	s.Echo = echo.New()

	s.Echo.Debug = s.Config.Echo.Debug
	s.Echo.HideBanner = true
	s.Echo.Logger.SetOutput(&echoLogger{level: s.Config.Logger.RequestLevel, log: log.With().Str("component", "echo").Logger()})

	s.Echo.HTTPErrorHandler = HTTPErrorHandlerWithConfig(HTTPErrorHandlerConfig{
		HideInternalServerErrorDetails: s.Config.Echo.HideInternalServerErrorDetails,
	})

	// ---
	// General middleware
	if s.Config.Echo.EnableTrailingSlashMiddleware {
		s.Echo.Pre(echoMiddleware.RemoveTrailingSlash())
	} else {
		log.Warn().Msg("Disabling trailing slash middleware due to environment config")
	}

	if s.Config.Echo.EnableRecoverMiddleware {
		s.Echo.Use(echoMiddleware.Recover())
	} else {
		log.Warn().Msg("Disabling recover middleware due to environment config")
	}

	if s.Config.Echo.EnableSecureMiddleware {
		s.Echo.Use(echoMiddleware.SecureWithConfig(echoMiddleware.SecureConfig{
			Skipper:               echoMiddleware.DefaultSecureConfig.Skipper,
			XSSProtection:         s.Config.Echo.SecureMiddleware.XSSProtection,
			ContentTypeNosniff:    s.Config.Echo.SecureMiddleware.ContentTypeNosniff,
			XFrameOptions:         s.Config.Echo.SecureMiddleware.XFrameOptions,
			HSTSMaxAge:            s.Config.Echo.SecureMiddleware.HSTSMaxAge,
			HSTSExcludeSubdomains: s.Config.Echo.SecureMiddleware.HSTSExcludeSubdomains,
			ContentSecurityPolicy: s.Config.Echo.SecureMiddleware.ContentSecurityPolicy,
			CSPReportOnly:         s.Config.Echo.SecureMiddleware.CSPReportOnly,
			HSTSPreloadEnabled:    s.Config.Echo.SecureMiddleware.HSTSPreloadEnabled,
			ReferrerPolicy:        s.Config.Echo.SecureMiddleware.ReferrerPolicy,
		}))
	} else {
		log.Warn().Msg("Disabling secure middleware due to environment config")
	}

	if s.Config.Echo.EnableRequestIDMiddleware {
		s.Echo.Use(echoMiddleware.RequestID())
	} else {
		log.Warn().Msg("Disabling request ID middleware due to environment config")
	}

	if s.Config.Echo.EnableLoggerMiddleware {
		s.Echo.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
			Level:             s.Config.Logger.RequestLevel,
			LogRequestBody:    s.Config.Logger.LogRequestBody,
			LogRequestHeader:  s.Config.Logger.LogRequestHeader,
			LogRequestQuery:   s.Config.Logger.LogRequestQuery,
			LogResponseBody:   s.Config.Logger.LogResponseBody,
			LogResponseHeader: s.Config.Logger.LogResponseHeader,
			LogCaller:         s.Config.Logger.LogCaller,
			RequestBodyLogSkipper: func(req *http.Request) bool {
				// We skip all body logging for auth endpoints as these might contain credentials
				if strings.HasPrefix(req.URL.Path, "/api/v1/auth") {
					return true
				}

				return middleware.DefaultRequestBodyLogSkipper(req)
			},
			ResponseBodyLogSkipper: func(req *http.Request, res *echo.Response) bool {
				// We skip all body logging for auth endpoints as these might contain credentials
				if strings.HasPrefix(req.URL.Path, "/api/v1/auth") {
					return true
				}

				return middleware.DefaultResponseBodyLogSkipper(req, res)
			},
			Skipper: func(c echo.Context) bool {
				// We skip logging of readiness and liveness endpoints
				switch c.Path() {
				case "/-/ready", "/-/healthy":
					return true
				}
				return false
			},
		}))
	} else {
		log.Warn().Msg("Disabling logger middleware due to environment config")
	}

	if s.Config.Echo.EnableCORSMiddleware {
		s.Echo.Use(echoMiddleware.CORS())
	} else {
		log.Warn().Msg("Disabling CORS middleware due to environment config")
	}

	if s.Config.Echo.EnableCacheControlMiddleware {
		s.Echo.Use(middleware.CacheControl())
	} else {
		log.Warn().Msg("Disabling cache control middleware due to environment config")
	}

	if s.Config.Pprof.Enable {

		pprofAuthMiddleware := middleware.Noop()

		if s.Config.Pprof.EnableManagementKeyAuth {
			pprofAuthMiddleware = echoMiddleware.KeyAuthWithConfig(echoMiddleware.KeyAuthConfig{
				KeyLookup: "query:mgmt-secret",
				Validator: func(key string, c echo.Context) (bool, error) {
					return key == s.Config.Management.Secret, nil
				},
			})
		}

		s.Echo.GET("/debug/pprof", echo.WrapHandler(http.HandlerFunc(pprof.Index)), pprofAuthMiddleware)
		s.Echo.Any("/debug/pprof/*", echo.WrapHandler(http.DefaultServeMux), pprofAuthMiddleware)

		log.Warn().Bool("EnableManagementKeyAuth", s.Config.Pprof.EnableManagementKeyAuth).Msg("Pprof http handlers are available at /debug/pprof")

		if s.Config.Pprof.RuntimeBlockProfileRate != 0 {
			runtime.SetBlockProfileRate(s.Config.Pprof.RuntimeBlockProfileRate)
			log.Warn().Int("RuntimeBlockProfileRate", s.Config.Pprof.RuntimeBlockProfileRate).Msg("Pprof runtime.SetBlockProfileRate")
		}

		if s.Config.Pprof.RuntimeMutexProfileFraction != 0 {
			runtime.SetMutexProfileFraction(s.Config.Pprof.RuntimeMutexProfileFraction)
			log.Warn().Int("RuntimeMutexProfileFraction", s.Config.Pprof.RuntimeMutexProfileFraction).Msg("Pprof runtime.SetMutexProfileFraction")
		}
	}

	// Add your custom / additional middlewares here.
	// see https://echo.labstack.com/middleware

	// ---
	// Initialize our general groups and set middleware to use above them
	s.Router = &api.Router{
		Routes: nil, // will be populated by handlers.AttachAllRoutes(s)

		// Unsecured base group available at /**
		Root: s.Echo.Group(""),

		// Management endpoints, uncacheable, secured by key auth (query param), available at /-/**
		Management: s.Echo.Group("/-", echoMiddleware.KeyAuthWithConfig(echoMiddleware.KeyAuthConfig{
			KeyLookup: "query:mgmt-secret",
			Validator: func(key string, c echo.Context) (bool, error) {
				return key == s.Config.Management.Secret, nil
			},
			Skipper: func(c echo.Context) bool {
				switch c.Path() {
				case "/-/ready":
					return true
				}
				return false
			},
		}), middleware.NoCache()),

		// OAuth2, unsecured or secured by bearer auth, available at /api/v1/auth/**
		APIV1Auth: s.Echo.Group("/api/v1/auth", middleware.AuthWithConfig(middleware.AuthConfig{
			S:    s,
			Mode: middleware.AuthModeRequired,
			Skipper: func(c echo.Context) bool {
				switch c.Path() {
				case "/api/v1/auth/forgot-password",
					"/api/v1/auth/forgot-password/complete",
					"/api/v1/auth/login",
					"/api/v1/auth/refresh",
					"/api/v1/auth/register":
					return true
				}
				return false
			},
		})),

		// Your other endpoints, typically secured by bearer auth, available at /api/v1/**
		APIV1Push: s.Echo.Group("/api/v1/push", middleware.Auth(s)),
	}

	// ---
	// Finally attach our handlers
	handlers.AttachAllRoutes(s)
}
