// +build tools

package tools

// Tooling dependencies
// https://github.com/golang/go/wiki/Modules#how-can-i-track-tool-dependencies-for-a-module
// https://github.com/go-modules-by-example/index/blob/master/010_tools/README.md

import (
	_ "github.com/go-swagger/go-swagger/cmd/swagger"
	_ "github.com/kyoh86/richgo"
	_ "github.com/rubenv/sql-migrate/sql-migrate"
	_ "github.com/volatiletech/sqlboiler"
	_ "github.com/volatiletech/sqlboiler/drivers/sqlboiler-psql"
)
