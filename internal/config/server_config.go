package config

import (
	"path/filepath"
	"runtime"
	"time"

	"allaboutapps.dev/aw/go-starter/internal/mailer/transport"
	"allaboutapps.dev/aw/go-starter/internal/push/provider"
	"allaboutapps.dev/aw/go-starter/internal/util"
	"github.com/rs/zerolog"
)

type EchoServer struct {
	Debug                          bool
	ListenAddress                  string
	HideInternalServerErrorDetails bool
	BaseURL                        string
}

type AuthServer struct {
	AccessTokenValidity          time.Duration
	PasswordResetTokenValidity   time.Duration
	DefaultUserScopes            []string
	LastAuthenticatedAtThreshold time.Duration
}

type PathsServer struct {
	APIBaseDirAbs        string
	MntBaseDirAbs        string
	TestAssetsBaseDirAbs string
}

type ManagementServer struct {
	Secret string
}

type FrontendServer struct {
	BaseURL               string
	PasswordResetEndpoint string
}

type LoggerServer struct {
	Level             zerolog.Level
	RequestLevel      zerolog.Level
	LogRequestBody    bool
	LogRequestHeader  bool
	LogRequestQuery   bool
	LogResponseBody   bool
	LogResponseHeader bool
}

type Server struct {
	Database   Database
	Echo       EchoServer
	Paths      PathsServer
	Auth       AuthServer
	Management ManagementServer
	Mailer     Mailer
	SMTP       transport.SMTPMailTransportConfig
	Frontend   FrontendServer
	Logger     LoggerServer
	Push       PushService
	FCMConfig  provider.FCMConfig
}

func DefaultServiceConfigFromEnv() Server {
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
		},
		Paths: PathsServer{
			// Please ALWAYS work with ABSOLUTE (ABS) paths from ENV_VARS (however you may resolve a project-relative to absolute for the default value)
			APIBaseDirAbs:        util.GetEnv("SERVER_PATHS_API_BASE_DIR_ABS", filepath.Join(util.GetProjectRootDir(), "/api")),                   // /app/api (swagger.yml)
			MntBaseDirAbs:        util.GetEnv("SERVER_PATHS_MNT_BASE_DIR_ABS", filepath.Join(util.GetProjectRootDir(), "/assets/mnt")),            // /app/assets/mnt (user-generated content)
			TestAssetsBaseDirAbs: util.GetEnv("SERVER_PATHS_TEST_ASSETS_BASE_DIR_ABS", filepath.Join(util.GetProjectRootDir(), "/test/testdata")), // /app/test/testdata (files used in tests and test snapshots)
		},
		Auth: AuthServer{
			AccessTokenValidity:          time.Second * time.Duration(util.GetEnvAsInt("SERVER_AUTH_ACCESS_TOKEN_VALIDITY", 86400)),
			PasswordResetTokenValidity:   time.Second * time.Duration(util.GetEnvAsInt("SERVER_AUTH_PASSWORD_RESET_TOKEN_VALIDITY", 900)),
			DefaultUserScopes:            util.GetEnvAsStringArr("SERVER_AUTH_DEFAULT_USER_SCOPES", []string{"app"}),
			LastAuthenticatedAtThreshold: time.Second * time.Duration(util.GetEnvAsInt("SERVER_AUTH_LAST_AUTHENTICATED_AT_THRESHOLD", 900)),
		},
		Management: ManagementServer{
			Secret: util.GetEnv("SERVER_MANAGEMENT_SECRET", "mgmt-pass"),
		},
		Mailer: Mailer{
			DefaultSender:               util.GetEnv("SERVER_MAILER_DEFAULT_SENDER", "go-starter@example.com"),
			Send:                        util.GetEnvAsBool("SERVER_MAILER_SEND", true),
			WebTemplatesEmailBaseDirAbs: util.GetEnv("SERVER_MAILER_WEB_TEMPLATES_EMAIL_BASE_DIR_ABS", filepath.Join(util.GetProjectRootDir(), "/web/templates/email")), // /app/web/templates/email
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
			Level:             util.LogLevelFromString(util.GetEnv("SERVER_LOGGER_LEVEL", zerolog.DebugLevel.String())),
			RequestLevel:      util.LogLevelFromString(util.GetEnv("SERVER_LOGGER_REQUEST_LEVEL", zerolog.DebugLevel.String())),
			LogRequestBody:    util.GetEnvAsBool("SERVER_LOGGER_LOG_REQUEST_BODY", false),
			LogRequestHeader:  util.GetEnvAsBool("SERVER_LOGGER_LOG_REQUEST_HEADER", false),
			LogRequestQuery:   util.GetEnvAsBool("SERVER_LOGGER_LOG_REQUEST_QUERY", false),
			LogResponseBody:   util.GetEnvAsBool("SERVER_LOGGER_LOG_RESPONSE_BODY", false),
			LogResponseHeader: util.GetEnvAsBool("SERVER_LOGGER_LOG_RESPONSE_HEADER", false),
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
	}
}
