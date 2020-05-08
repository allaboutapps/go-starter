package api

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"allaboutapps.dev/aw/go-starter/internal/mailer"
	"allaboutapps.dev/aw/go-starter/internal/mailer/transport"
	"allaboutapps.dev/aw/go-starter/internal/util"
)

type DatabaseConfig struct {
	Host             string            `json:"host"`
	Port             int               `json:"port"`
	Username         string            `json:"username"`
	Password         string            `json:"password"`
	Database         string            `json:"database"`
	AdditionalParams map[string]string `json:"additionalParams,omitempty"` // Optional additional connection parameters mapped into the connection string
}

// Generates a connection string to be passed to sql.Open or equivalents, assuming Postgres syntax
func (c DatabaseConfig) ConnectionString() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s", c.Host, c.Port, c.Username, c.Password, c.Database))

	if _, ok := c.AdditionalParams["sslmode"]; !ok {
		b.WriteString(" sslmode=disable")
	}

	if len(c.AdditionalParams) > 0 {
		params := make([]string, 0, len(c.AdditionalParams))
		for param := range c.AdditionalParams {
			params = append(params, param)
		}

		sort.Strings(params)

		for _, param := range params {
			fmt.Fprintf(&b, " %s=%s", param, c.AdditionalParams[param])
		}
	}

	return b.String()
}

type EchoServerConfig struct {
	Debug         bool
	ListenAddress string
}

type AuthServerConfig struct {
	AccessTokenValidity          time.Duration
	PasswordResetTokenValidity   time.Duration
	DefaultUserScopes            []string
	LastAuthenticatedAtThreshold int
}

type FrontendServerConfig struct {	
	BaseURL               string
	PasswordResetEndpoint string
}

type ServerConfig struct {
	Database DatabaseConfig
	Echo     EchoServerConfig
	Auth     AuthServerConfig
	Mailer   mailer.MailerConfig
	SMTP     transport.SMTPMailTransportConfig
	Frontend FrontendServerConfig
}

func DefaultServiceConfigFromEnv() ServerConfig {
	return ServerConfig{
		Database: DatabaseConfig{
			Host:     util.GetEnv("PGHOST", "postgres"),
			Port:     util.GetEnvAsInt("PGPORT", 5432),
			Database: util.GetEnv("PGDATABASE", "development"),
			Username: util.GetEnv("PGUSER", "dbuser"),
			Password: util.GetEnv("PGPASSWORD", ""),
			AdditionalParams: map[string]string{
				"sslmode": util.GetEnv("PGSSLMODE", "disable"),
			},
		},
		Echo: EchoServerConfig{
			Debug:         util.GetEnvAsBool("SERVER_ECHO_DEBUG", false),
			ListenAddress: util.GetEnv("SERVER_ECHO_LISTEN_ADDRESS", ":8080"),
		},
		Auth: AuthServerConfig{
			AccessTokenValidity:          time.Second * time.Duration(util.GetEnvAsInt("SERVER_AUTH_ACCESS_TOKEN_VALIDITY", 86400)),
			PasswordResetTokenValidity:   time.Second * time.Duration(util.GetEnvAsInt("SERVER_AUTH_PASSWORD_RESET_TOKEN_VALIDITY", 900)),
			DefaultUserScopes:            util.GetEnvAsStringArr("SERVER_AUTH_DEFAULT_USER_SCOPES", []string{"app"}),
			LastAuthenticatedAtThreshold: util.GetEnvAsInt("SERVER_AUTH_LAST_AUTHENTICATED_AT_THRESHOLD", 900),
		},
		Mailer: mailer.MailerConfig{
			DefaultSender: util.GetEnv("SERVER_MAILER_DEFAULT_SENDER", "operations+go-starter-local@allaboutapps.at"),
			Send:          util.GetEnvAsBool("SERVER_MAILER_SEND", true),
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
		Frontend: FrontendServerConfig{
			BaseURL:               util.GetEnv("SERVER_FRONTEND_BASE_URL", "http://localhost:3000"),
			PasswordResetEndpoint: util.GetEnv("SERVER_FRONTEND_PASSWORD_RESET_ENDPOINT", "/set-new-password"),
		},
	}
}
