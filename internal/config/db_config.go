package config

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"allaboutapps.dev/aw/go-starter/internal/util"
)

// The DatabaseMigrationTable name is baked into the binary
// This setting should always be in sync with dbconfig.yml, sqlboiler.toml and the live database (e.g. to be able to test producation dumps locally)
const DatabaseMigrationTable = "migrations"

// The DatabaseMigrationFolder (folder with all *.sql migrations).
// This settings should always be in sync with dbconfig.yaml and Dockerfile (the final app stage).
// It's expected that the migrations folder lives at the root of this project or right next to the app binary.
var DatabaseMigrationFolder = filepath.Join(util.GetProjectRootDir(), "/migrations")

type Database struct {
	Host             string
	Port             int
	Username         string
	Password         string `json:"-"` // sensitive
	Database         string
	AdditionalParams map[string]string `json:",omitempty"` // Optional additional connection parameters mapped into the connection string
	MaxOpenConns     int
	MaxIdleConns     int
	ConnMaxLifetime  time.Duration
}

// ConnectionString generates a connection string to be passed to sql.Open or equivalents, assuming Postgres syntax
func (c Database) ConnectionString() string {
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
