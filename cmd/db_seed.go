package cmd

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"allaboutapps.dev/aw/go-starter/internal/config"
	"allaboutapps.dev/aw/go-starter/internal/data"
	dbutil "allaboutapps.dev/aw/go-starter/internal/util/db"
	"github.com/spf13/cobra"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

// seedCmd represents the seed command
var seedCmd = &cobra.Command{
	Use:   "seed",
	Short: "Inserts or updates fixtures to the database.",
	Long:  `Uses upsert to add test data to the current environment.`,
	Run:   seedCmdFunc,
}

func init() {
	dbCmd.AddCommand(seedCmd)
}

func seedCmdFunc(cmd *cobra.Command, args []string) {
	if err := applyFixtures(); err != nil {
		fmt.Printf("Error while seeding fixtures: %v", err)
		os.Exit(1)
	}

	fmt.Println("Seeded all fixtures.")
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

	// insert fixtures in an auto-managed db transaction
	return dbutil.WithTransaction(ctx, db, func(tx boil.ContextExecutor) error {

		fixtures := data.Upserts()

		for _, fixture := range fixtures {
			if err := fixture.Upsert(ctx, tx, true, nil, boil.Infer(), boil.Infer()); err != nil {
				fmt.Printf("Failed to upsert fixture: %v\n", err)
				return err
			}
		}

		fmt.Printf("Upserted %d fixtures.\n", len(fixtures))
		return nil

	})
}
