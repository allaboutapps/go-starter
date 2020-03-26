package pgtestpool

type ManagerConfig struct {
	TemplateDatabaseBaseName string           // Optional base name to be used for template databases, will have hash appended. Defaults to manager database name with "_template" appended if empty
	TestDatabaseBaseName     string           // Optional base name to be used for test databases, will have continuous numeric ID appended. Defaults to managater database name if empty
	DatabaseConfig           ConnectionConfig // Config for manager to connect to database, will require privileged user for creating test databases and managing templates
}
