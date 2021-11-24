package test

import (
	"testing"

	"allaboutapps.dev/aw/go-starter/internal/config"
	"allaboutapps.dev/aw/go-starter/internal/i18n"
)

// TODO: add test

// NewTestMessages returns a i18n.Messages pointer created from the default service env
func NewTestMessages(t *testing.T) *i18n.Messages {
	t.Helper()
	messages, err := i18n.New(config.DefaultServiceConfigFromEnv().I18n)

	if err != nil {
		t.Fatal("Failed to setup i18n messages", err)
	}

	return messages
}
