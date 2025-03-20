package test_test

import (
	"context"
	"testing"

	"allaboutapps.dev/aw/go-starter/internal/mailer/transport"
	"allaboutapps.dev/aw/go-starter/internal/test"
	"allaboutapps.dev/aw/go-starter/internal/test/fixtures"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWithTestMailer(t *testing.T) {
	ctx := context.Background()
	fix := fixtures.Fixtures()
	//nolint:gosec
	passwordResetLink := "http://localhost/password/reset/12345"

	m1 := test.NewTestMailer(t)
	m2 := test.NewTestMailer(t)
	err := m1.SendPasswordReset(ctx, fix.User1.Username.String, passwordResetLink)
	require.NoError(t, err)

	sender2 := "test2@example.com"
	m2.Config.DefaultSender = sender2
	err = m2.SendPasswordReset(ctx, fix.User1.Username.String, passwordResetLink)
	require.NoError(t, err)

	mt1 := test.GetTestMailerMockTransport(t, m1)
	mail := mt1.GetLastSentMail()
	mails := mt1.GetSentMails()

	assert.Equal(t, mail, test.GetLastSentMail(t, m1))
	assert.Equal(t, mails, test.GetSentMails(t, m1))

	require.NotNil(t, mail)
	require.Len(t, mails, 1)
	assert.Equal(t, m1.Config.DefaultSender, mail.From)
	assert.Len(t, mail.To, 1)
	assert.Equal(t, fix.User1.Username.String, mail.To[0])
	assert.Equal(t, test.TestMailerDefaultSender, mail.From)
	assert.Equal(t, "Password reset", mail.Subject)
	assert.Contains(t, string(mail.HTML), passwordResetLink)

	mt2 := test.GetTestMailerMockTransport(t, m2)
	mail = mt2.GetLastSentMail()
	mails = mt2.GetSentMails()

	assert.Equal(t, mail, test.GetLastSentMail(t, m2))
	assert.Equal(t, mails, test.GetSentMails(t, m2))

	require.NotNil(t, mail)
	require.Len(t, mails, 1)
	assert.Equal(t, m2.Config.DefaultSender, mail.From)
	assert.Len(t, mail.To, 1)
	assert.Equal(t, fix.User1.Username.String, mail.To[0])
	assert.Equal(t, sender2, mail.From)
	assert.Equal(t, "Password reset", mail.Subject)
	assert.Contains(t, string(mail.HTML), passwordResetLink)
}

func TestWithSMTPMailerFromDefaultEnv(t *testing.T) {
	m := test.NewSMTPMailerFromDefaultEnv(t)
	require.NotNil(t, m)
	require.NotEmpty(t, m.Transport)
	assert.IsType(t, &transport.SMTPMailTransport{}, m.Transport)
}
