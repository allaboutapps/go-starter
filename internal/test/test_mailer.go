package test

import (
	"testing"

	"allaboutapps.dev/aw/go-starter/internal/api"
	"allaboutapps.dev/aw/go-starter/internal/config"
	"allaboutapps.dev/aw/go-starter/internal/mailer"
	"allaboutapps.dev/aw/go-starter/internal/mailer/transport"
	"github.com/jordan-wright/email"
)

const (
	TestMailerDefaultSender = "test@example.com"
)

func WithTestMailer(t *testing.T, closure func(m *mailer.Mailer)) {
	t.Helper()

	closure(NewTestMailer(t))
}

func WithMailer(t *testing.T, s *api.Server) {
	t.Helper()

	if err := s.InitMailer(); err != nil {
		t.Fatalf("Failed to init mailer: %v", err)
	}
}

func NewTestMailer(t *testing.T) *mailer.Mailer {
	t.Helper()

	config := config.DefaultServiceConfigFromEnv().Mailer
	config.DefaultSender = TestMailerDefaultSender

	m := mailer.New(config, transport.NewMock())

	if err := m.ParseTemplates(); err != nil {
		t.Fatal("Failed to parse mailer templates", err)
	}

	return m
}

func GetLastSentMail(t *testing.T, m *mailer.Mailer) *email.Email {
	t.Helper()
	mt, ok := m.Transport.(*transport.MockMailTransport)
	if !ok {
		t.Fatalf("invalid mailer transport type, got %T, want *transport.MockMailTransport", m.Transport)
	}

	return mt.GetLastSentMail()
}

func GetSentMails(t *testing.T, m *mailer.Mailer) []*email.Email {
	t.Helper()
	mt, ok := m.Transport.(*transport.MockMailTransport)
	if !ok {
		t.Fatalf("invalid mailer transport type, got %T, want *transport.MockMailTransport", m.Transport)
	}

	return mt.GetSentMails()
}
