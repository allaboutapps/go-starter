package api

import (
	"fmt"
	"sort"
	"strings"

	"allaboutapps.at/aw/go-mranftl-sample/internal/util"
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
	LastAuthenticatedAtThreshold int
}

type ServerConfig struct {
	Database DatabaseConfig
	Echo     EchoServerConfig
	Auth     AuthServerConfig
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
			LastAuthenticatedAtThreshold: util.GetEnvAsInt("SERVER_AUTH_LAST_AUTHENTICATED_AT_THRESHOLD", 900),
		},
	}
}
