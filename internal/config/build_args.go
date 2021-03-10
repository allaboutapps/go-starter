package config

import "fmt"

// The following vars are auto-injected via -ldflags during make build (see Makefile)
var (
	Commit    string = "unknown" // e.g. "59cb7684dd0b0f38d68cd7db657cb614feba8f7e"
	BuildDate string = "unknown" // e.g. "1970-01-01T00:00:00+00:00"
)

func GetBuildArgVersion() string {
	return fmt.Sprintf("@ %v (%v)", Commit, BuildDate)
}
