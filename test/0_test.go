package test

import (
	"flag"
	"os"
	"testing"
)

func TestMain(m *testing.M) {

	flag.Parse()

	// ensure the template database gets initialized before running our tests
	InitializeDatabaseTemplate()

	exit := m.Run()

	os.Exit(exit)
}
