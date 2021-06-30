package middleware

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"time"

	"allaboutapps.dev/aw/go-starter/internal/api"
	"allaboutapps.dev/aw/go-starter/internal/api/auth"
	"allaboutapps.dev/aw/go-starter/internal/api/httperrors"
	"allaboutapps.dev/aw/go-starter/internal/models"
	"allaboutapps.dev/aw/go-starter/internal/util"
	"github.com/go-openapi/strfmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

var (
	ErrBadRequestMalformedToken                = httperrors.NewHTTPError(http.StatusBadRequest, "MALFORMED_TOKEN", "Auth token is malformed")
	ErrUnauthorizedLastAuthenticatedAtExceeded = httperrors.NewHTTPError(http.StatusUnauthorized, "LAST_AUTHENTICATED_AT_EXCEEDED", "LastAuthenticatedAt timestamp exceeds threshold, re-authentication required")
	ErrForbiddenUserDeactivated                = httperrors.NewHTTPError(http.StatusForbidden, "USER_DEACTIVATED", "User account is deactivated")
	ErrForbiddenMissingScopes                  = httperrors.NewHTTPError(http.StatusForbidden, "MISSING_SCOPES", "User is missing required scopes")
	ErrAuthTokenValidationFailed               = errors.New("auth token validation failed")
)

// AuthMode controls the type of authentication check performed for a specific route or group
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
	// AuthTokenSourceForm retrieves the auth token from a form parameter, specified by TokenSourceKey
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

type AuthTokenFormatValidator func(string) bool

func DefaultAuthTokenFormatValidator(token string) bool {
	return strfmt.IsUUID4(token)
}

type AuthTokenValidator func(c echo.Context, config AuthConfig, token string) (auth.AuthenticationResult, error)

func DefaultAuthTokenValidator(c echo.Context, config AuthConfig, token string) (auth.AuthenticationResult, error) {
	accessToken, err := models.AccessTokens(
		models.AccessTokenWhere.Token.EQ(token),
		qm.Load(models.AccessTokenRels.User),
	).One(c.Request().Context(), config.S.DB)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Trace().Err(err).Msg("Access token not found in database")
			return auth.AuthenticationResult{}, ErrAuthTokenValidationFailed
		}

		log.Error().Err(err).Msg("Failed to query for access token in database, aborting request")
		return auth.AuthenticationResult{}, echo.ErrInternalServerError
	}

	return auth.AuthenticationResult{
		Token:      accessToken.Token,
		User:       accessToken.R.User,
		ValidUntil: accessToken.ValidUntil,
	}, nil
}

var (
	DefaultAuthConfig = AuthConfig{
		Mode:            AuthModeRequired,
		FailureMode:     AuthFailureModeUnauthorized,
		TokenSource:     AuthTokenSourceHeader,
		TokenSourceKey:  echo.HeaderAuthorization,
		Scheme:          "Bearer",
		Skipper:         middleware.DefaultSkipper,
		FormatValidator: DefaultAuthTokenFormatValidator,
		TokenValidator:  DefaultAuthTokenValidator,
		Scopes:          []string{auth.AuthScopeApp.String()},
	}
)

type AuthConfig struct {
	S               *api.Server              // API server used for database and service access
	Mode            AuthMode                 // Controls type of authentication required (default: AuthModeRequired)
	FailureMode     AuthFailureMode          // Controls response on auth failure (default: AuthFailureModeUnauthorized)
	TokenSource     AuthTokenSource          // Sets source of auth token (default: AuthTokenSourceHeader)
	TokenSourceKey  string                   // Sets key for auth token source lookup (default: "Authorization")
	Scheme          string                   // Sets required token scheme (default: "Bearer")
	Skipper         middleware.Skipper       // Controls skipping of certain routes (default: no skipped routes)
	FormatValidator AuthTokenFormatValidator // Validates the format of the token retrieved
	TokenValidator  AuthTokenValidator       // Validates token retrieved and returns associated user (default: performs lookup in access_tokens table)
	Scopes          []string                 // List of scopes required to access endpoint (default: none required)
}

func (c AuthConfig) CheckLastAuthenticatedAt(user *models.User) bool {
	if c.Mode != AuthModeSecure {
		return true
	}

	if !user.LastAuthenticatedAt.Valid {
		return false
	}

	return time.Since(user.LastAuthenticatedAt.Time).Seconds() <= c.S.Config.Auth.LastAuthenticatedAtThreshold.Seconds()
}

func (c AuthConfig) CheckUserScopes(user *models.User) bool {
	if len(c.Scopes) == 0 {
		return true
	}

	if len(user.Scopes) == 0 {
		return false
	}

	for _, scope := range c.Scopes {
		for _, userScope := range user.Scopes {
			if scope == userScope {
				return true
			}
		}
	}

	return false
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

	if config.FormatValidator == nil {
		config.FormatValidator = DefaultAuthConfig.FormatValidator
	}

	if config.TokenValidator == nil {
		config.TokenValidator = DefaultAuthConfig.TokenValidator
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			log := util.LogFromEchoContext(c).With().Str("middleware", "auth").Str("auth_mode", config.Mode.String()).Logger()

			if config.Mode == AuthModeNone {
				log.Trace().Msg("No authentication required, allowing request")
				return next(c)
			}

			if config.Skipper(c) {
				log.Trace().Msg("Skipping auth middleware, allowing request")
				return next(c)
			}

			user := auth.UserFromEchoContext(c)
			if user != nil {
				if !config.CheckLastAuthenticatedAt(user) {
					log.Trace().
						Time("last_authenticated_at", user.LastAuthenticatedAt.Time).
						Dur("last_authenticated_at_threshold", config.S.Config.Auth.LastAuthenticatedAtThreshold).
						Msg("Authentication already performed, but last authenticated at time exceeds threshold, rejecting request")
					return ErrUnauthorizedLastAuthenticatedAtExceeded
				}

				if !config.CheckUserScopes(user) {
					log.Trace().
						Strs("scopes", config.Scopes).
						Strs("user_scopes", user.Scopes).
						Msg("Authentication already performed, but user does not have required scopes, rejecting request")
					return ErrForbiddenMissingScopes
				}

				log.Trace().Msg("Authentication already performed, allowing request")
				return next(c)
			}

			token, exists := config.TokenSource.Extract(c, config.TokenSourceKey, config.Scheme)
			if len(token) == 0 {
				if config.Mode == AuthModeRequired || config.Mode == AuthModeSecure || (exists && config.Mode == AuthModeOptional) {
					log.Trace().Bool("token_exists", exists).Msg("Request has missing or malformed token, rejecting")
					return config.FailureMode.Error()
				}

				log.Trace().Bool("token_exists", exists).Msg("Request does not have valid token, but auth mode permits access, allowing request")
				return next(c)
			}

			if !config.FormatValidator(token) {
				if config.Mode == AuthModeRequired || config.Mode == AuthModeSecure || config.Mode == AuthModeOptional {
					log.Trace().Msg("Request has malformed token, rejecting")
					return ErrBadRequestMalformedToken
				}

				log.Trace().Msg("Request have malformed token, but auth mode permits access, allowing request")
				return next(c)
			}

			res, err := config.TokenValidator(c, config, token)
			if err != nil {
				if errors.Is(err, ErrAuthTokenValidationFailed) {
					if config.Mode == AuthModeTry {
						log.Trace().Msg("Auth token validation failed, but auth mode permits access, allowing request")
						return next(c)
					}

					log.Trace().Msg("Auth token validation failed, rejecting request")
					return config.FailureMode.Error()
				}

				log.Trace().Err(err).Msg("Failed to validate auth token, aborting request")
				return echo.ErrInternalServerError
			}

			user = res.User

			if res.ValidUntil.IsZero() {
				log.Trace().Str("user_id", user.ID).Msg("Auth token has no expiry, allowing request")
			} else {
				if time.Now().After(res.ValidUntil) {
					if config.Mode == AuthModeTry {
						log.Trace().Time("valid_until", res.ValidUntil).Str("user_id", user.ID).Msg("Auth token is expired, but auth mode permits access, allowing request")
						return next(c)
					}

					log.Trace().Time("valid_until", res.ValidUntil).Str("user_id", user.ID).Msg("Auth token is expired, rejecting request")
					return config.FailureMode.Error()
				}
			}

			// ! User has been explicitly deactivated - we do not allow access here, even with AuthModeTry
			if !user.IsActive {
				log.Trace().Str("user_id", user.ID).Msg("User is deactivated, rejecting request")
				return ErrForbiddenUserDeactivated
			}

			if !config.CheckLastAuthenticatedAt(user) {
				log.Trace().
					Time("last_authenticated_at", user.LastAuthenticatedAt.Time).
					Dur("last_authenticated_at_threshold", config.S.Config.Auth.LastAuthenticatedAtThreshold).
					Msg("Authentication already performed, but last authenticated at time exceeds threshold, rejecting request")
				return ErrUnauthorizedLastAuthenticatedAtExceeded
			}

			if !config.CheckUserScopes(user) {
				log.Trace().
					Strs("scopes", config.Scopes).
					Strs("user_scopes", user.Scopes).
					Msg("Authentication already performed, but user does not have required scopes, rejecting request")
				return ErrForbiddenMissingScopes
			}

			auth.EnrichEchoContextWithCredentials(c, res)

			log.Trace().Str("user_id", user.ID).Msg("Auth token is valid, allowing request")

			return next(c)
		}
	}
}
