package test

// import (
// 	"context"
// 	"database/sql"
// 	"fmt"
// 	"os"
// 	"testing"

// 	"allaboutapps.at/aw/go-mranftl-sample/models"
// 	_ "github.com/lib/pq"
// 	"github.com/volatiletech/sqlboiler/boil"
// )

// var (
// 	host            = os.Getenv("PSQL_HOST")
// 	port     int64  = 5432
// 	user            = os.Getenv("PSQL_USER")
// 	password string = os.Getenv("PSQL_PASS")
// 	dbname          = os.Getenv("PSQL_DBNAME")
// )

// type Model interface {
// 	DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error)
// }

// func TestFixturesThroughSQLBoiler(t *testing.T) {

// 	fmt.Println("Connecting...")

// 	// boil.DebugMode = true

// 	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
// 		"password=%s dbname=%s sslmode=disable",
// 		host, port, user, password, dbname)

// 	db, err := sql.Open("postgres", psqlInfo)

// 	if err != nil {
// 		panic(err)
// 	}

// 	defer db.Close()

// 	err = db.Ping()

// 	if err != nil {
// 		panic(err)
// 	}

// 	fmt.Println("Successfully connected!")

// 	// trunc (only for now, will be useless when integrated with pgpool)
// 	for _, model := range []Model{models.Jets(), models.PilotLanguages(), models.Languages(), models.Pilots()} {
// 		_, err = model.DeleteAll(context.TODO(), db)
// 		if err != nil {
// 			t.Error("truncate fail", model)
// 		}
// 	}

// 	tx, err := db.BeginTx(context.TODO(), nil)

// 	if err != nil {
// 		t.Error("transaction fail")
// 	}

// 	for _, fixture := range fixtures {
// 		err = fixture.Insert(context.Background(), db, boil.Infer())

// 		if err != nil {
// 			t.Error("Failed to insert fixture", fixture)
// 		}
// 	}

// 	// Rollback or commit
// 	err = tx.Commit()

// 	if err != nil {
// 		t.Error("transaction commit failed")
// 	}

// 	err = pilot1.Reload(context.TODO(), db)

// 	if err != nil {
// 		t.Error("failed to reload")
// 	}

// 	fmt.Println(pilot1)

// }
