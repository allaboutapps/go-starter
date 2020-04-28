package middleware

import (
	"database/sql"
	"fmt"
	"time"

	"allaboutapps.at/aw/go-mranftl-sample/api"
	"allaboutapps.at/aw/go-mranftl-sample/models"
	"allaboutapps.at/aw/go-mranftl-sample/pkg/auth"
	"allaboutapps.at/aw/go-mranftl-sample/pkg/util"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

// Controls the type of authentication check performed for a specific route or group
type AuthMode int

const (
	// AuthModeRequired requires an auth token to be present and valid in order to access the route or group
	AuthModeRequired AuthMode = iota
	// AuthModeSecure requires an auth token to be present and for the user to have recently re-confirmed their authentication in order to access the route or group
	AuthModeSecure
	// AuthModeOptional does not require an auth token to be present, however if it is, it must be valid in order to access the route or group
	AuthModeOptional
	// AuthModeTry does not require an auth token to be present in order to access the route or group and will process the request even if an invalid one has been provided
	AuthModeTry
	// AuthModeNone does not require an auth token to be present in order to access the route or group and will not attempt to parse any authentication provided
	AuthModeNone
)

func (m AuthMode) String() string {
	switch m {
	case AuthModeRequired:
		return "required"
	case AuthModeSecure:
		return "secure"
	case AuthModeOptional:
		return "optional"
	case AuthModeTry:
		return "try"
	case AuthModeNone:
		return "none"
	default:
		return fmt.Sprintf("unknown (%d)", m)
	}
}

type AuthFailureMode int

const (
	// AuthFailureModeUnauthorized returns a 401 Unauthorized response on missing or invalid authentication
	AuthFailureModeUnauthorized AuthFailureMode = iota
	// AuthFailureModeNotFound returns a 404 Not Found response on missing or invalid authentication
	AuthFailureModeNotFound
)

func (m AuthFailureMode) String() string {
	switch m {
	case AuthFailureModeUnauthorized:
		return "unauthorized"
	case AuthFailureModeNotFound:
		return "not_found"
	default:
		return fmt.Sprintf("unknown (%d)", m)
	}
}

func (m AuthFailureMode) Error() error {
	switch m {
	case AuthFailureModeUnauthorized:
		return echo.ErrUnauthorized
	case AuthFailureModeNotFound:
		return echo.ErrNotFound
	default:
		return echo.ErrInternalServerError
	}
}

type AuthTokenSource int

const (
	// AuthTokenSourceHeader retrieves the auth token from a header, specified by TokenSourceKey
	AuthTokenSourceHeader AuthTokenSource = iota
	// AuthTokenSourceQuery retrieves the auth token from a query parameter, specified by TokenSourceKey
	AuthTokenSourceQuery
	// AuthTOkenSourceForm retrieves the auth token from a form parameter, specified by TokenSourceKey
	AuthTokenSourceForm
)

func (s AuthTokenSource) String() string {
	switch s {
	case AuthTokenSourceHeader:
		return "header"
	case AuthTokenSourceQuery:
		return "query"
	case AuthTokenSourceForm:
		return "form"
	default:
		return fmt.Sprintf("unknown (%d)", s)
	}
}

func (s AuthTokenSource) Extract(c echo.Context, key string, scheme string) (token string, exists bool) {
	var t string

	switch s {
	case AuthTokenSourceHeader:
		t = c.Request().Header.Get(key)
	case AuthTokenSourceForm:
		t = c.FormValue(key)
	case AuthTokenSourceQuery:
		t = c.QueryParam(key)
	default:
		return "", false
	}

	if len(t) == 0 {
		return "", false
	}

	lenScheme := len(scheme)
	if lenScheme == 0 {
		return t, true
	}

	if len(t) < lenScheme+1 {
		return "", true
	}

	if t[:lenScheme] != scheme {
		return "", true
	}

	return t[lenScheme+1:], true
}

var (
	DefaultAuthConfig = AuthConfig{
		Mode:           AuthModeRequired,
		FailureMode:    AuthFailureModeUnauthorized,
		TokenSource:    AuthTokenSourceHeader,
		TokenSourceKey: echo.HeaderAuthorization,
		Scheme:         "Bearer",
		Skipper:        middleware.DefaultSkipper,
	}
)

type AuthConfig struct {
	S              *api.Server        // API server used for database and service access
	Mode           AuthMode           // Controls type of authentication required (default: AuthModeRequired)
	FailureMode    AuthFailureMode    // Controls response on auth failure (default: AuthFailureModeUnauthorized)
	TokenSource    AuthTokenSource    // Sets source of auth token (default: AuthTokenSourceHeader)
	TokenSourceKey string             // Sets key for auth token source lookup (default: "Authorization")
	Scheme         string             // Sets required token scheme (default: "Bearer")
	Skipper        middleware.Skipper // Controls skipping of certain routes (default: no skipped routes)
}

func Auth(s *api.Server) echo.MiddlewareFunc {
	c := DefaultAuthConfig
	c.S = s
	return AuthWithConfig(c)
}

func AuthWithConfig(config AuthConfig) echo.MiddlewareFunc {
	if config.S == nil {
		panic("auth middleware: server is required")
	}

	if len(config.TokenSourceKey) == 0 {
		config.TokenSourceKey = DefaultAuthConfig.TokenSourceKey
	}

	if len(config.Scheme) == 0 {
		config.Scheme = DefaultAuthConfig.Scheme
	}

	if config.Skipper == nil {
		config.Skipper = DefaultAuthConfig.Skipper
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			log := util.LogFromEchoContext(c)

			if config.Mode == AuthModeNone {
				log.Trace().Str("middleware", "auth").Str("auth_mode", config.Mode.String()).Msg("No authentication required, allowing request")
				return next(c)
			}

			if config.Skipper(c) {
				log.Trace().Str("middleware", "auth").Str("auth_mode", config.Mode.String()).Msg("Skipping auth middleware, allowing request")
				return next(c)
			}

			user := auth.UserFromEchoContext(c)
			if user != nil {
				// TODO perform additional validation for AuthModeSecure
				log.Trace().Str("middleware", "auth").Str("auth_mode", config.Mode.String()).Msg("Authentication already performed, allowing request")
				return next(c)
			}

			token, exists := config.TokenSource.Extract(c, config.TokenSourceKey, config.Scheme)
			if len(token) == 0 {
				if config.Mode == AuthModeRequired || config.Mode == AuthModeSecure || (exists && config.Mode == AuthModeOptional) {
					log.Trace().Str("middleware", "auth").Str("auth_mode", config.Mode.String()).Bool("token_exists", exists).Msg("Request has missing or malformed token, rejecting")
					return config.FailureMode.Error()
				}

				log.Trace().Str("middleware", "auth").Str("auth_mode", config.Mode.String()).Bool("token_exists", exists).Msg("Request does not have valid token, but auth mode permits access, allowing request")
				return next(c)
			}

			accessToken, err := models.AccessTokens(
				qm.Load(models.AccessTokenRels.User),
				qm.Where("token = ?", token),
			).One(c.Request().Context(), config.S.DB)
			if err != nil {
				if err == sql.ErrNoRows {
					if config.Mode == AuthModeTry {
						log.Trace().Str("middleware", "auth").Str("auth_mode", config.Mode.String()).Msg("Access token not found in database, but auth mode permits access, allowing request")
						return next(c)
					}

					log.Trace().Str("middleware", "auth").Str("auth_mode", config.Mode.String()).Msg("Access token not found in database, rejecting request")
					return config.FailureMode.Error()
				}

				log.Trace().Str("middleware", "auth").Str("auth_mode", config.Mode.String()).Err(err).Msg("Failed to query for access token in database, aborting request")
				return echo.ErrInternalServerError
			}

			// TODO scopes enum array on user model

			if time.Now().After(accessToken.ValidUntil) {
				if config.Mode == AuthModeTry {
					log.Trace().
						Str("middleware", "auth").
						Str("auth_mode", config.Mode.String()).
						Time("valid_until", accessToken.ValidUntil).
						Str("user_id", accessToken.R.User.ID).
						Msg("Access token is expired, but auth mode permits access, allowing request")
					return next(c)
				}

				log.Trace().
					Str("middleware", "auth").
					Str("auth_mode", config.Mode.String()).
					Time("valid_until", accessToken.ValidUntil).
					Str("user_id", accessToken.R.User.ID).
					Msg("Access token is expired, rejecting request")
				return echo.ErrUnauthorized
			}

			// ! User has been explicitly deactivated - we do not allow access here, even with AuthModeTry
			if !accessToken.R.User.IsActive {
				log.Trace().
					Str("middleware", "auth").
					Str("auth_mode", config.Mode.String()).
					Time("valid_until", accessToken.ValidUntil).
					Str("user_id", accessToken.R.User.ID).
					Msg("User is deactivated, rejecting request")
				return echo.ErrUnauthorized
			}

			// TODO perform additional validation for AuthModeSecure
			// TODO ACL check

			auth.EnrichEchoContextWithCredentials(c, accessToken.R.User, accessToken)

			log.Trace().
				Str("middleware", "auth").
				Str("auth_mode", config.Mode.String()).
				Time("valid_until", accessToken.ValidUntil).
				Str("user_id", accessToken.R.User.ID).
				Msg("Access token is valid, allowing request")

			return next(c)
		}
	}
}
