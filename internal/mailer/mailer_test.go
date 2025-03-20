package mailer_test

import (
	"context"
	"testing"
	"time"

	"allaboutapps.dev/aw/go-starter/internal/api"
	"allaboutapps.dev/aw/go-starter/internal/config"
	"allaboutapps.dev/aw/go-starter/internal/test"
	"allaboutapps.dev/aw/go-starter/internal/test/fixtures"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMailerSendPasswordReset(t *testing.T) {
	ctx := context.Background()
	fix := fixtures.Fixtures()

	m := test.NewTestMailer(t)
	mt := test.GetTestMailerMockTransport(t, m)
	mt.Expect(1)

	//nolint:gosec
	passwordResetLink := "http://localhost/password/reset/12345"
	err := m.SendPasswordReset(ctx, fix.User1.Username.String, passwordResetLink)
	require.NoError(t, err)

	mt.WaitWithTimeout(time.Second)

	mail := mt.GetLastSentMail()
	mails := mt.GetSentMails()
	require.NotNil(t, mail)
	require.Len(t, mails, 1)
	assert.Equal(t, m.Config.DefaultSender, mail.From)
	assert.Len(t, mail.To, 1)
	assert.Equal(t, fix.User1.Username.String, mail.To[0])
	assert.Equal(t, test.TestMailerDefaultSender, mail.From)
	assert.Equal(t, "Password reset", mail.Subject)
	assert.Contains(t, string(mail.HTML), passwordResetLink)
}

func SkipTestMailerSendPasswordResetWithMailhog(t *testing.T) {
	t.Skip()
	ctx := context.Background()
	fix := fixtures.Fixtures()

	m := test.NewSMTPMailerFromDefaultEnv(t)

	//nolint:gosec
	passwordResetLink := "http://localhost/password/reset/12345"
	err := m.SendPasswordReset(ctx, fix.User1.Username.String, passwordResetLink)
	require.NoError(t, err)
}

func SkipTestMailerSendPasswordResetWithMailhogAndServer(t *testing.T) {
	t.Skip()
	ctx := context.Background()
	fix := fixtures.Fixtures()

	defaultConfig := config.DefaultServiceConfigFromEnv()
	defaultConfig.Mailer.Transporter = config.MailerTransporterSMTP.String()
	test.WithTestServerConfigurable(t, defaultConfig, func(s *api.Server) {
		//nolint:gosec
		passwordResetLink := "http://localhost/password/reset/12345"
		err := s.Mailer.SendPasswordReset(ctx, fix.User1.Username.String, passwordResetLink)
		require.NoError(t, err)
	})
}
