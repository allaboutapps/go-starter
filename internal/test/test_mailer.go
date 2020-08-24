package test

import (
	"testing"

	"allaboutapps.dev/aw/go-starter/internal/config"
	"allaboutapps.dev/aw/go-starter/internal/mailer"
	"allaboutapps.dev/aw/go-starter/internal/mailer/transport"
)

const (
	TestMailerDefaultSender = "test@example.com"
)

func NewTestMailer(t *testing.T) *mailer.Mailer {
	t.Helper()

	return newMailerWithTransporter(t, transport.NewMock())
}

func NewSMTPMailerFromDefaultEnv(t *testing.T) *mailer.Mailer {
	t.Helper()

	config := config.DefaultServiceConfigFromEnv().SMTP
	return newMailerWithTransporter(t, transport.NewSMTP(config))
}

func GetTestMailerMockTransport(t *testing.T, m *mailer.Mailer) *transport.MockMailTransport {
	t.Helper()
	mt, ok := m.Transport.(*transport.MockMailTransport)
	if !ok {
		t.Fatalf("invalid mailer transport type, got %T, want *transport.MockMailTransport", m.Transport)
	}

	return mt
}

func newMailerWithTransporter(t *testing.T, transporter transport.MailTransporter) *mailer.Mailer {
	t.Helper()

	config := config.DefaultServiceConfigFromEnv().Mailer
	config.DefaultSender = TestMailerDefaultSender

	m := mailer.New(config, transporter)

	if err := m.ParseTemplates(); err != nil {
		t.Fatal("Failed to parse mailer templates", err)
	}

	return m
}
