package config

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

type Database struct {
	Host             string            `json:"host"`
	Port             int               `json:"port"`
	Username         string            `json:"username"`
	Password         string            `json:"-"` // sensitive, do not emit if this struct gets marshalled.
	Database         string            `json:"database"`
	AdditionalParams map[string]string `json:"additionalParams,omitempty"` // Optional additional connection parameters mapped into the connection string
	MaxOpenConns     int               `json:"maxOpenConns"`
	MaxIdleConns     int               `json:"maxIdleConns"`
	ConnMaxLifetime  time.Duration     `json:"connMaxLifetime"`
}

// Generates a connection string to be passed to sql.Open or equivalents, assuming Postgres syntax
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
