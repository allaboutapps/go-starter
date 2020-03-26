package test

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"

	"allaboutapps.at/aw/go-mranftl-sample/models"
	_ "github.com/lib/pq"
	"github.com/volatiletech/sqlboiler/boil"
)

var (
	host            = os.Getenv("PSQL_HOST")
	port     int64  = 5432
	user            = os.Getenv("PSQL_USER")
	password string = os.Getenv("PSQL_PASS")
	dbname          = os.Getenv("PSQL_DBNAME")
)

func TestFixturesThroughSQLBoiler(t *testing.T) {

	fmt.Println("Connecting...")

	boil.DebugMode = true

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		panic(err)
	}

	defer db.Close()

	err = db.Ping()

	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")

	// trunc (only for now, will be useless when integrated with pgpool)
	models.Jets().DeleteAllP(context.Background(), db)
	models.PilotLanguages().DeleteAllP(context.Background(), db)
	models.Languages().DeleteAllP(context.Background(), db)
	models.Pilots().DeleteAllP(context.Background(), db)

	tx, err := db.BeginTx(context.TODO(), nil)

	if err != nil {
		t.Error("transaction fail")
	}

	for _, item := range fixtures {
		item.InsertP(context.Background(), db, boil.Infer())
	}

	// Rollback or commit
	err = tx.Commit()

	if err != nil {
		t.Error("transaction commit failed")
	}

	pilot1.ReloadP(context.TODO(), db)

	fmt.Println(pilot1)

}
