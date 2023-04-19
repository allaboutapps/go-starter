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

	// pin migrate to use the globally defined `migrations` table identifier
	migrate.SetTable(config.DatabaseMigrationTable)
}

func migrateCmdFunc(_ *cobra.Command, _ []string) {
	n, err := applyMigrations()
	if err != nil {
		fmt.Printf("Error while applying migrations: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Applied %d migrations.\n", n)
}

func applyMigrations() (int, error) {
	ctx := context.Background()
	serviceConfig := config.DefaultServiceConfigFromEnv()
	db, err := sql.Open("postgres", serviceConfig.Database.ConnectionString())
	if err != nil {
		return 0, err
	}
	defer db.Close()

	if err := db.PingContext(ctx); err != nil {
		return 0, err
	}

	// In case an old default sql-migrate migration table (named "gorp_migrations") still exists we rename it to the new name equivalent
	// in sync with the settings in dbconfig.yml and config.DatabaseMigrationTable.
	if _, err := db.Exec(fmt.Sprintf("ALTER TABLE IF EXISTS gorp_migrations RENAME TO %s;", config.DatabaseMigrationTable)); err != nil {
		return 0, err
	}

	migrations := &migrate.FileMigrationSource{
		Dir: config.DatabaseMigrationFolder,
	}
	n, err := migrate.Exec(db, "postgres", migrations, migrate.Up)
	if err != nil {
		return 0, err
	}

	return n, nil
}
