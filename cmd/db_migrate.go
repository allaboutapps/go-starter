package cmd

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"allaboutapps.dev/aw/go-starter/internal/config"
	migrate "github.com/rubenv/sql-migrate"
	"github.com/spf13/cobra"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Executes all migrations which are not yet applied.",
	Run:   migrateCmdFunc,
}

func init() {
	dbCmd.AddCommand(migrateCmd)
}

func migrateCmdFunc(cmd *cobra.Command, args []string) {
	n, err := applyMigrations()
	if err != nil {
		fmt.Printf("Error while applying migrations: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Applied %d migrations.\n", n)
}

func applyMigrations() (int, error) {
	ctx := context.Background()
	config := config.DefaultServiceConfigFromEnv()
	db, err := sql.Open("postgres", config.Database.ConnectionString())
	if err != nil {
		return 0, err
	}
	defer db.Close()

	if err := db.PingContext(ctx); err != nil {
		return 0, err
	}

	migrations := &migrate.FileMigrationSource{
		Dir: "migrations",
	}
	n, err := migrate.Exec(db, "postgres", migrations, migrate.Up)
	if err != nil {
		return 0, err
	}

	return n, nil
}
