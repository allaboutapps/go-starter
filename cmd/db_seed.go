package cmd

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"allaboutapps.dev/aw/go-starter/internal/config"
	"allaboutapps.dev/aw/go-starter/internal/data"
	"github.com/spf13/cobra"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

// seedCmd represents the seed command
var seedCmd = &cobra.Command{
	Use:   "seed",
	Short: "Inserts or updates fixtures to the database.",
	Long:  `Uses upsert to add test data to the current environment.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := applyFixtures(); err != nil {
			fmt.Printf("Error while applying fixtures: %v", err)
			os.Exit(1)
		}

		fmt.Println("Applied all fixtures")
	},
}

func init() {
	dbCmd.AddCommand(seedCmd)
}

func applyFixtures() error {
	ctx := context.Background()
	config := config.DefaultServiceConfigFromEnv()
	db, err := sql.Open("postgres", config.Database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	if err := db.PingContext(ctx); err != nil {
		return err
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	fixtures := data.Upserts()

	for _, fixture := range fixtures {
		if err := fixture.Upsert(ctx, db, true, nil, boil.Infer(), boil.Infer()); err != nil {
			if err := tx.Rollback(); err != nil {
				return err
			}

			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
