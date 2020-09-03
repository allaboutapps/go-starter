package config_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"allaboutapps.dev/aw/go-starter/internal/config"
)

func TestPrintServiceEnv(t *testing.T) {
	t.Parallel()

	config := config.DefaultServiceConfigFromEnv()

	c, err := json.MarshalIndent(config, "", "  ")

	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(string(c))
}
