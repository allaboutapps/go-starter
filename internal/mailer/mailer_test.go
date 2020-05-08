package mailer_test

import (
	"context"
	"testing"

	"allaboutapps.dev/aw/go-starter/internal/mailer"
	"allaboutapps.dev/aw/go-starter/internal/mailer/transport"
	"allaboutapps.dev/aw/go-starter/internal/test"
	"github.com/jordan-wright/email"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	mockMailerDefaultSender = "test@example.com"
)

func mockMailer(t *testing.T) *mailer.Mailer {
	t.Helper()

	m := mailer.New(mailer.MailerConfig{DefaultSender: mockMailerDefaultSender, Send: true}, transport.NewMock())

	if err := m.ParseTemplates(); err != nil {
		t.Fatalf("failed to parse email templates for mock mailer: %v", err)
	}

	return m
}

func getLastSentMail(t *testing.T, m *mailer.Mailer) *email.Email {
	t.Helper()

	mt, ok := m.Transport.(*transport.MockMailTransport)
	if !ok {
		t.Fatalf("invalid mailer transport type, got %T, want *transport.MockMailTransport", m.Transport)
	}

	return mt.GetLastSentMail()
}

func TestMailerSendPasswordReset(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	fixtures := test.Fixtures()

	m := mockMailer(t)

	passwordResetLink := "http://localhost/password/reset/12345"
	err := m.SendPasswordReset(ctx, fixtures.User1.Username.String, passwordResetLink)
	require.NoError(t, err)

	mail := getLastSentMail(t, m)
	require.NotNil(t, mail)
	assert.Equal(t, m.Config.DefaultSender, mail.From)
	assert.Len(t, mail.To, 1)
	assert.Equal(t, fixtures.User1.Username.String, mail.To[0])
	assert.Equal(t, mockMailerDefaultSender, mail.From)
	assert.Equal(t, "Password reset", mail.Subject)
	assert.Contains(t, string(mail.HTML), passwordResetLink)
}
