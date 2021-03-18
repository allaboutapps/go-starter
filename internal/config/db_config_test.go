package config_test

import (
	"testing"

	"allaboutapps.dev/aw/go-starter/internal/config"
)

// via https://github.com/allaboutapps/integresql/blob/master/pkg/manager/database_config_test.go
func TestDatabaseConnectionString(t *testing.T) {
	tests := []struct {
		name   string
		config config.Database
		want   string
	}{
		{
			name: "Simple",
			config: config.Database{
				Host:     "localhost",
				Port:     5432,
				Username: "simple",
				Password: "database_config",
				Database: "simple_database_config",
			},
			want: "host=localhost port=5432 user=simple password=database_config dbname=simple_database_config sslmode=disable",
		},
		{
			name: "SSLMode",
			config: config.Database{
				Host:     "localhost",
				Port:     5432,
				Username: "simple",
				Password: "database_config",
				Database: "simple_database_config",
				AdditionalParams: map[string]string{
					"sslmode": "prefer",
				},
			},
			want: "host=localhost port=5432 user=simple password=database_config dbname=simple_database_config sslmode=prefer",
		},
		{
			name: "Complex",
			config: config.Database{
				Host:     "localhost",
				Port:     5432,
				Username: "simple",
				Password: "database_config",
				Database: "simple_database_config",
				AdditionalParams: map[string]string{
					"connect_timeout": "10",
					"sslmode":         "verify-full",
					"sslcert":         "/app/certs/pg.pem",
					"sslkey":          "/app/certs/pg.key",
					"sslrootcert":     "/app/certs/pg_root.pem",
				},
			},
			want: "host=localhost port=5432 user=simple password=database_config dbname=simple_database_config connect_timeout=10 sslcert=/app/certs/pg.pem sslkey=/app/certs/pg.key sslmode=verify-full sslrootcert=/app/certs/pg_root.pem",
		},
	}

	for _, tt := range tests {
		tt := tt // NOTE: https://github.com/golang/go/wiki/CommonMistakes#using-goroutines-on-loop-iterator-variables
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.config.ConnectionString(); got != tt.want {
				t.Errorf("invalid connection string, got %q, want %q", got, tt.want)
			}
		})
	}
}
