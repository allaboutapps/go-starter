package config

import (
	"os"
	"path/filepath"
	"runtime"
	"time"

	"allaboutapps.dev/aw/go-starter/internal/mailer/transport"
	"allaboutapps.dev/aw/go-starter/internal/push/provider"
	"allaboutapps.dev/aw/go-starter/internal/util"
	"github.com/rs/zerolog"
	"golang.org/x/text/language"
)

type EchoServer struct {
	Debug                          bool
	ListenAddress                  string
	HideInternalServerErrorDetails bool
	BaseURL                        string
	EnableCORSMiddleware           bool
	EnableLoggerMiddleware         bool
	EnableRecoverMiddleware        bool
	EnableRequestIDMiddleware      bool
	EnableTrailingSlashMiddleware  bool
	EnableSecureMiddleware         bool
	EnableCacheControlMiddleware   bool
	SecureMiddleware               EchoServerSecureMiddleware
}

type PprofServer struct {
	Enable                      bool
	EnableManagementKeyAuth     bool
	RuntimeBlockProfileRate     int
	RuntimeMutexProfileFraction int
}

// EchoServerSecureMiddleware represents a subset of echo's secure middleware config relevant to the app server.
// https://github.com/labstack/echo/blob/master/middleware/secure.go
type EchoServerSecureMiddleware struct {
	XSSProtection         string
	ContentTypeNosniff    string
	XFrameOptions         string
	HSTSMaxAge            int
	HSTSExcludeSubdomains bool
	ContentSecurityPolicy string
	CSPReportOnly         bool
	HSTSPreloadEnabled    bool
	ReferrerPolicy        string
}

type AuthServer struct {
	AccessTokenValidity          time.Duration
	PasswordResetTokenValidity   time.Duration
	DefaultUserScopes            []string
	LastAuthenticatedAtThreshold time.Duration
}

type PathsServer struct {
	APIBaseDirAbs string
	MntBaseDirAbs string
}

type ManagementServer struct {
	Secret                  string `json:"-"` // sensitive
	ReadinessTimeout        time.Duration
	LivenessTimeout         time.Duration
	ProbeWriteablePathsAbs  []string
	ProbeWriteableTouchfile string
}

type FrontendServer struct {
	BaseURL               string
	PasswordResetEndpoint string
}

type LoggerServer struct {
	Level              zerolog.Level
	RequestLevel       zerolog.Level
	LogRequestBody     bool
	LogRequestHeader   bool
	LogRequestQuery    bool
	LogResponseBody    bool
	LogResponseHeader  bool
	LogCaller          bool
	PrettyPrintConsole bool
}

type I18n struct {
	DefaultLanguage language.Tag
	BundleDirAbs    string
}

type Server struct {
	Database   Database
	Echo       EchoServer
	Pprof      PprofServer
	Paths      PathsServer
	Auth       AuthServer
	Management ManagementServer
	Mailer     Mailer
	SMTP       transport.SMTPMailTransportConfig
	Frontend   FrontendServer
	Logger     LoggerServer
	Push       PushService
	FCMConfig  provider.FCMConfig
	I18n       I18n
}

// DefaultServiceConfigFromEnv returns the server config as parsed from environment variables
// and their respective defaults defined below.
// We don't expect that ENV_VARs change while we are running our application or our tests
// (and it would be a bad thing to do anyways with parallel testing).
// Do NOT use os.Setenv / os.Unsetenv in tests utilizing DefaultServiceConfigFromEnv()!
func DefaultServiceConfigFromEnv() Server {

	// An `.env.local` file in your project root can override the currently set ENV variables.
	//
	// We never automatically apply `.env.local` when running "go test" as these ENV variables
	// may be sensitive (e.g. secrets to external APIs) and applying them modifies the process
	// global "os.Env" state (it should be applied via t.SetEnv instead).
	//
	// If you need dotenv ENV variables available in a test, do that explicitly within that
	// test before executing DefaultServiceConfigFromEnv (or test.WithTestServer).
	// See /internal/test/helper_dot_env.go: test.DotEnvLoadLocalOrSkipTest(t)
	if !util.RunningInTest() {
		DotEnvTryLoad(filepath.Join(util.GetProjectRootDir(), ".env.local"), os.Setenv)
	}

	return Server{
		Database: Database{
			Host:     util.GetEnv("PGHOST", "postgres"),
			Port:     util.GetEnvAsInt("PGPORT", 5432),
			Database: util.GetEnv("PGDATABASE", "development"),
			Username: util.GetEnv("PGUSER", "dbuser"),
			Password: util.GetEnv("PGPASSWORD", ""),
			AdditionalParams: map[string]string{
				"sslmode": util.GetEnv("PGSSLMODE", "disable"),
			},
			MaxOpenConns:    util.GetEnvAsInt("DB_MAX_OPEN_CONNS", runtime.NumCPU()*2),
			MaxIdleConns:    util.GetEnvAsInt("DB_MAX_IDLE_CONNS", 1),
			ConnMaxLifetime: time.Second * time.Duration(util.GetEnvAsInt("DB_CONN_MAX_LIFETIME_SEC", 60)),
		},
		Echo: EchoServer{
			Debug:                          util.GetEnvAsBool("SERVER_ECHO_DEBUG", false),
			ListenAddress:                  util.GetEnv("SERVER_ECHO_LISTEN_ADDRESS", ":8080"),
			HideInternalServerErrorDetails: util.GetEnvAsBool("SERVER_ECHO_HIDE_INTERNAL_SERVER_ERROR_DETAILS", true),
			BaseURL:                        util.GetEnv("SERVER_ECHO_BASE_URL", "http://localhost:8080"),
			EnableCORSMiddleware:           util.GetEnvAsBool("SERVER_ECHO_ENABLE_CORS_MIDDLEWARE", true),
			EnableLoggerMiddleware:         util.GetEnvAsBool("SERVER_ECHO_ENABLE_LOGGER_MIDDLEWARE", true),
			EnableRecoverMiddleware:        util.GetEnvAsBool("SERVER_ECHO_ENABLE_RECOVER_MIDDLEWARE", true),
			EnableRequestIDMiddleware:      util.GetEnvAsBool("SERVER_ECHO_ENABLE_REQUEST_ID_MIDDLEWARE", true),
			EnableTrailingSlashMiddleware:  util.GetEnvAsBool("SERVER_ECHO_ENABLE_TRAILING_SLASH_MIDDLEWARE", true),
			EnableSecureMiddleware:         util.GetEnvAsBool("SERVER_ECHO_ENABLE_SECURE_MIDDLEWARE", true),
			EnableCacheControlMiddleware:   util.GetEnvAsBool("SERVER_ECHO_ENABLE_CACHE_CONTROL_MIDDLEWARE", true),
			// see https://echo.labstack.com/middleware/secure
			// see https://github.com/labstack/echo/blob/master/middleware/secure.go
			SecureMiddleware: EchoServerSecureMiddleware{
				XSSProtection:         util.GetEnv("SERVER_ECHO_SECURE_MIDDLEWARE_XSS_PROTECTION", "1; mode=block"),
				ContentTypeNosniff:    util.GetEnv("SERVER_ECHO_SECURE_MIDDLEWARE_CONTENT_TYPE_NOSNIFF", "nosniff"),
				XFrameOptions:         util.GetEnv("SERVER_ECHO_SECURE_MIDDLEWARE_X_FRAME_OPTIONS", "SAMEORIGIN"),
				HSTSMaxAge:            util.GetEnvAsInt("SERVER_ECHO_SECURE_MIDDLEWARE_HSTS_MAX_AGE", 0),
				HSTSExcludeSubdomains: util.GetEnvAsBool("SERVER_ECHO_SECURE_MIDDLEWARE_HSTS_EXCLUDE_SUBDOMAINS", false),
				ContentSecurityPolicy: util.GetEnv("SERVER_ECHO_SECURE_MIDDLEWARE_CONTENT_SECURITY_POLICY", ""),
				CSPReportOnly:         util.GetEnvAsBool("SERVER_ECHO_SECURE_MIDDLEWARE_CSP_REPORT_ONLY", false),
				HSTSPreloadEnabled:    util.GetEnvAsBool("SERVER_ECHO_SECURE_MIDDLEWARE_HSTS_PRELOAD_ENABLED", false),
				ReferrerPolicy:        util.GetEnv("SERVER_ECHO_SECURE_MIDDLEWARE_REFERRER_POLICY", ""),
			},
		},
		Pprof: PprofServer{
			// https://golang.org/pkg/net/http/pprof/
			Enable:                      util.GetEnvAsBool("SERVER_PPROF_ENABLE", false),
			EnableManagementKeyAuth:     util.GetEnvAsBool("SERVER_PPROF_ENABLE_MANAGEMENT_KEY_AUTH", true),
			RuntimeBlockProfileRate:     util.GetEnvAsInt("SERVER_PPROF_RUNTIME_BLOCK_PROFILE_RATE", 0),
			RuntimeMutexProfileFraction: util.GetEnvAsInt("SERVER_PPROF_RUNTIME_MUTEX_PROFILE_FRACTION", 0),
		},
		Paths: PathsServer{
			// Please ALWAYS work with ABSOLUTE (ABS) paths from ENV_VARS (however you may resolve a project-relative to absolute for the default value)
			APIBaseDirAbs: util.GetEnv("SERVER_PATHS_API_BASE_DIR_ABS", filepath.Join(util.GetProjectRootDir(), "/api")),        // /app/api (swagger.yml)
			MntBaseDirAbs: util.GetEnv("SERVER_PATHS_MNT_BASE_DIR_ABS", filepath.Join(util.GetProjectRootDir(), "/assets/mnt")), // /app/assets/mnt (user-generated content)
		},
		Auth: AuthServer{
			AccessTokenValidity:          time.Second * time.Duration(util.GetEnvAsInt("SERVER_AUTH_ACCESS_TOKEN_VALIDITY", 86400)),
			PasswordResetTokenValidity:   time.Second * time.Duration(util.GetEnvAsInt("SERVER_AUTH_PASSWORD_RESET_TOKEN_VALIDITY", 900)),
			DefaultUserScopes:            util.GetEnvAsStringArr("SERVER_AUTH_DEFAULT_USER_SCOPES", []string{"app"}),
			LastAuthenticatedAtThreshold: time.Second * time.Duration(util.GetEnvAsInt("SERVER_AUTH_LAST_AUTHENTICATED_AT_THRESHOLD", 900)),
		},
		Management: ManagementServer{
			Secret:           util.GetMgmtSecret("SERVER_MANAGEMENT_SECRET"),
			ReadinessTimeout: time.Second * time.Duration(util.GetEnvAsInt("SERVER_MANAGEMENT_READINESS_TIMEOUT_SEC", 4)),
			LivenessTimeout:  time.Second * time.Duration(util.GetEnvAsInt("SERVER_MANAGEMENT_LIVENESS_TIMEOUT_SEC", 9)),
			ProbeWriteablePathsAbs: util.GetEnvAsStringArr("SERVER_MANAGEMENT_PROBE_WRITEABLE_PATHS_ABS", []string{
				filepath.Join(util.GetProjectRootDir(), "/assets/mnt")}, ","),
			ProbeWriteableTouchfile: util.GetEnv("SERVER_MANAGEMENT_PROBE_WRITEABLE_TOUCHFILE", ".healthy"),
		},
		Mailer: Mailer{
			DefaultSender:               util.GetEnv("SERVER_MAILER_DEFAULT_SENDER", "go-starter@example.com"),
			Send:                        util.GetEnvAsBool("SERVER_MAILER_SEND", true),
			WebTemplatesEmailBaseDirAbs: util.GetEnv("SERVER_MAILER_WEB_TEMPLATES_EMAIL_BASE_DIR_ABS", filepath.Join(util.GetProjectRootDir(), "/web/templates/email")), // /app/web/templates/email
			Transporter:                 util.GetEnvEnum("SERVER_MAILER_TRANSPORTER", MailerTransporterMock.String(), []string{MailerTransporterSMTP.String(), MailerTransporterMock.String()}),
		},
		SMTP: transport.SMTPMailTransportConfig{
			Host:      util.GetEnv("SERVER_SMTP_HOST", "mailhog"),
			Port:      util.GetEnvAsInt("SERVER_SMTP_PORT", 1025),
			Username:  util.GetEnv("SERVER_SMTP_USERNAME", ""),
			Password:  util.GetEnv("SERVER_SMTP_PASSWORD", ""),
			AuthType:  transport.SMTPAuthTypeFromString(util.GetEnv("SERVER_SMTP_AUTH_TYPE", transport.SMTPAuthTypeNone.String())),
			UseTLS:    util.GetEnvAsBool("SERVER_SMTP_USE_TLS", false),
			TLSConfig: nil,
		},
		Frontend: FrontendServer{
			BaseURL:               util.GetEnv("SERVER_FRONTEND_BASE_URL", "http://localhost:3000"),
			PasswordResetEndpoint: util.GetEnv("SERVER_FRONTEND_PASSWORD_RESET_ENDPOINT", "/set-new-password"),
		},
		Logger: LoggerServer{
			Level:              util.LogLevelFromString(util.GetEnv("SERVER_LOGGER_LEVEL", zerolog.DebugLevel.String())),
			RequestLevel:       util.LogLevelFromString(util.GetEnv("SERVER_LOGGER_REQUEST_LEVEL", zerolog.DebugLevel.String())),
			LogRequestBody:     util.GetEnvAsBool("SERVER_LOGGER_LOG_REQUEST_BODY", false),
			LogRequestHeader:   util.GetEnvAsBool("SERVER_LOGGER_LOG_REQUEST_HEADER", false),
			LogRequestQuery:    util.GetEnvAsBool("SERVER_LOGGER_LOG_REQUEST_QUERY", false),
			LogResponseBody:    util.GetEnvAsBool("SERVER_LOGGER_LOG_RESPONSE_BODY", false),
			LogResponseHeader:  util.GetEnvAsBool("SERVER_LOGGER_LOG_RESPONSE_HEADER", false),
			LogCaller:          util.GetEnvAsBool("SERVER_LOGGER_LOG_CALLER", false),
			PrettyPrintConsole: util.GetEnvAsBool("SERVER_LOGGER_PRETTY_PRINT_CONSOLE", false),
		},
		Push: PushService{
			UseFCMProvider:  util.GetEnvAsBool("SERVER_PUSH_USE_FCM", false),
			UseMockProvider: util.GetEnvAsBool("SERVER_PUSH_USE_MOCK", true),
		},
		FCMConfig: provider.FCMConfig{
			GoogleApplicationCredentials: util.GetEnv("GOOGLE_APPLICATION_CREDENTIALS", ""),
			ProjectID:                    util.GetEnv("SERVER_FCM_PROJECT_ID", "no-fcm-project-id-set"),
			ValidateOnly:                 util.GetEnvAsBool("SERVER_FCM_VALIDATE_ONLY", true),
		},
		I18n: I18n{
			DefaultLanguage: util.GetEnvAsLanguageTag("SERVER_I18N_DEFAULT_LANGUAGE", language.English),
			BundleDirAbs:    util.GetEnv("SERVER_I18N_BUNDLE_DIR_ABS", filepath.Join(util.GetProjectRootDir(), "/web/i18n")), // /app/web/i18n
		},
	}

}
