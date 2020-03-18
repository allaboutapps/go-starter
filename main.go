//go:generate sqlboiler --wipe --no-hooks --add-panic-variants psql

package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"

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

func main() {

	fmt.Println("Connecting...")

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

	// jsut for the showcase, we will use the panic variants here currently...

	var p1 models.Pilot
	p1.Name = "Mario"

	p1.InsertP(context.Background(), db, boil.Infer())

	pilots := models.Pilots().AllP(context.Background(), db)

	fmt.Println(len(pilots), "pilots")

	for _, pilot := range pilots {
		fmt.Println(pilot)
	}

}
