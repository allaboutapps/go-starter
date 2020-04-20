package auth

import (
	"flag"
	"os"
	"testing"

	"allaboutapps.at/aw/go-mranftl-sample/test"
)

func TestMain(m *testing.M) {

	flag.Parse()

	// ensure the template database gets initialized before running our tests
	test.InitializeDatabaseTemplate()

	exit := m.Run()

	os.Exit(exit)
}
