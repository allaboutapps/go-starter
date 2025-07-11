package mailer_test

import (
	"testing"
	"time"

	"allaboutapps.dev/aw/go-starter/internal/test"
	"allaboutapps.dev/aw/go-starter/internal/test/fixtures"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMailerSendPasswordReset(t *testing.T) {
	ctx := t.Context()
	fix := fixtures.Fixtures()

	mailer := test.NewTestMailer(t)
	mailTransport := test.GetTestMailerMockTransport(t, mailer)
	mailTransport.Expect(1)

	//nolint:gosec
	passwordResetLink := "http://localhost/password/reset/12345"
	err := mailer.SendPasswordReset(ctx, fix.User1.Username.String, passwordResetLink)
	require.NoError(t, err)

	mailTransport.WaitWithTimeout(time.Second)

	mail := mailTransport.GetLastSentMail()
	mails := mailTransport.GetSentMails()
	require.NotNil(t, mail)
	require.Len(t, mails, 1)
	assert.Equal(t, mailer.Config.DefaultSender, mail.From)
	assert.Len(t, mail.To, 1)
	assert.Equal(t, fix.User1.Username.String, mail.To[0])
	assert.Equal(t, test.TestMailerDefaultSender, mail.From)
	assert.Equal(t, "Password reset", mail.Subject)
	assert.Contains(t, string(mail.HTML), passwordResetLink)
}
