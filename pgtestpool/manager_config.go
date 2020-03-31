package pgtestpool

import "os"

type ManagerConfig struct {
	DatabasePrefix            string
	TestDatabaseOwner         string
	TestDatabaseOwnerPassword string
	ManagerDatabaseConfig     DatabaseConfig
}

func DefaultManagerConfigFromEnv() ManagerConfig {
	return ManagerConfig{
		DatabasePrefix:            "test",
		TestDatabaseOwner:         os.Getenv("PSQL_USER"),
		TestDatabaseOwnerPassword: os.Getenv("PSQL_PASS"),
		ManagerDatabaseConfig: DatabaseConfig{
			Host:     os.Getenv("PSQL_HOST"),
			Port:     5432,
			Username: os.Getenv("PSQL_USER"),
			Password: os.Getenv("PSQL_PASS"),
			Database: "postgres",
		},
	}
}
