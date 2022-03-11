package auth_test

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"allaboutapps.dev/aw/go-starter/internal/api"
	"allaboutapps.dev/aw/go-starter/internal/api/httperrors"
	"allaboutapps.dev/aw/go-starter/internal/mailer"
	"allaboutapps.dev/aw/go-starter/internal/mailer/transport"
	"allaboutapps.dev/aw/go-starter/internal/models"
	"allaboutapps.dev/aw/go-starter/internal/test"
	"github.com/jordan-wright/email"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

func getLastSentMail(t *testing.T, m *mailer.Mailer) *email.Email {
	t.Helper()

	mt, ok := m.Transport.(*transport.MockMailTransport)
	if !ok {
		t.Fatalf("invalid mailer transport type, got %T, want *transport.MockMailTransport", m.Transport)
	}

	return mt.GetLastSentMail()
}

func TestPostForgotPasswordSuccess(t *testing.T) {
	test.WithTestServer(t, func(s *api.Server) {
		ctx := context.Background()
		fixtures := test.Fixtures()
		payload := test.GenericPayload{
			"username": fixtures.User1.Username,
		}

		res := test.PerformRequest(t, s, "POST", "/api/v1/auth/forgot-password", payload, nil)

		assert.Equal(t, http.StatusNoContent, res.Result().StatusCode)

		passwordResetToken, err := fixtures.User1.PasswordResetTokens().One(ctx, s.DB)
		require.NoError(t, err)

		mail := getLastSentMail(t, s.Mailer)
		require.NotNil(t, mail)
		assert.Contains(t, string(mail.HTML), fmt.Sprintf("http://localhost:3000/set-new-password?token=%s", passwordResetToken.Token))
	})
}

func TestPostForgotPasswordUnknownUser(t *testing.T) {
	test.WithTestServer(t, func(s *api.Server) {
		ctx := context.Background()
		payload := test.GenericPayload{
			"username": "definitelydoesnotexist@example.com",
		}

		res := test.PerformRequest(t, s, "POST", "/api/v1/auth/forgot-password", payload, nil)

		assert.Equal(t, http.StatusNoContent, res.Result().StatusCode)

		cnt, err := models.PasswordResetTokens().Count(ctx, s.DB)
		assert.NoError(t, err)
		assert.Equal(t, int64(0), cnt)

		mail := getLastSentMail(t, s.Mailer)
		assert.Nil(t, mail)
	})
}

func TestPostForgotPasswordDeactivatedUser(t *testing.T) {
	test.WithTestServer(t, func(s *api.Server) {
		ctx := context.Background()
		fixtures := test.Fixtures()
		payload := test.GenericPayload{
			"username": fixtures.UserDeactivated.Username,
		}

		res := test.PerformRequest(t, s, "POST", "/api/v1/auth/forgot-password", payload, nil)

		assert.Equal(t, http.StatusNoContent, res.Result().StatusCode)

		cnt, err := models.PasswordResetTokens().Count(ctx, s.DB)
		assert.NoError(t, err)
		assert.Equal(t, int64(0), cnt)

		mail := getLastSentMail(t, s.Mailer)
		assert.Nil(t, mail)
	})
}

func TestPostForgotPasswordUserWithoutPassword(t *testing.T) {
	test.WithTestServer(t, func(s *api.Server) {
		ctx := context.Background()
		fixtures := test.Fixtures()
		payload := test.GenericPayload{
			"username": fixtures.User2.Username,
		}

		fixtures.User2.Password = null.NewString("", false)
		rowsAff, err := fixtures.User2.Update(context.Background(), s.DB, boil.Infer())
		require.NoError(t, err)
		require.Equal(t, int64(1), rowsAff)

		res := test.PerformRequest(t, s, "POST", "/api/v1/auth/forgot-password", payload, nil)

		assert.Equal(t, http.StatusNoContent, res.Result().StatusCode)

		cnt, err := models.PasswordResetTokens().Count(ctx, s.DB)
		assert.NoError(t, err)
		assert.Equal(t, int64(0), cnt)

		mail := getLastSentMail(t, s.Mailer)
		assert.Nil(t, mail)
	})
}

func TestPostForgotPasswordMissingUsername(t *testing.T) {
	test.WithTestServer(t, func(s *api.Server) {
		ctx := context.Background()
		payload := test.GenericPayload{}

		res := test.PerformRequest(t, s, "POST", "/api/v1/auth/forgot-password", payload, nil)

		assert.Equal(t, http.StatusBadRequest, res.Result().StatusCode)

		var response httperrors.HTTPValidationError
		test.ParseResponseAndValidate(t, res, &response)

		assert.Equal(t, int64(http.StatusBadRequest), *response.Code)
		assert.Equal(t, httperrors.HTTPErrorTypeGeneric, *response.Type)
		assert.Equal(t, http.StatusText(http.StatusBadRequest), *response.Title)
		assert.Empty(t, response.Detail)
		assert.Nil(t, response.Internal)
		assert.Nil(t, response.AdditionalData)
		assert.NotEmpty(t, response.ValidationErrors)
		assert.Equal(t, "username", *response.ValidationErrors[0].Key)
		assert.Equal(t, "body", *response.ValidationErrors[0].In)
		assert.Equal(t, "username in body is required", *response.ValidationErrors[0].Error)

		cnt, err := models.PasswordResetTokens().Count(ctx, s.DB)
		assert.NoError(t, err)
		assert.Equal(t, int64(0), cnt)

		mail := getLastSentMail(t, s.Mailer)
		assert.Nil(t, mail)
	})
}

func TestPostForgotPasswordEmptyUsername(t *testing.T) {
	test.WithTestServer(t, func(s *api.Server) {
		ctx := context.Background()
		payload := test.GenericPayload{
			"username": "",
		}

		res := test.PerformRequest(t, s, "POST", "/api/v1/auth/forgot-password", payload, nil)

		assert.Equal(t, http.StatusBadRequest, res.Result().StatusCode)

		var response httperrors.HTTPValidationError
		test.ParseResponseAndValidate(t, res, &response)

		assert.Equal(t, int64(http.StatusBadRequest), *response.Code)
		assert.Equal(t, httperrors.HTTPErrorTypeGeneric, *response.Type)
		assert.Equal(t, http.StatusText(http.StatusBadRequest), *response.Title)
		assert.Empty(t, response.Detail)
		assert.Nil(t, response.Internal)
		assert.Nil(t, response.AdditionalData)
		assert.NotEmpty(t, response.ValidationErrors)
		assert.Equal(t, "username", *response.ValidationErrors[0].Key)
		assert.Equal(t, "body", *response.ValidationErrors[0].In)
		assert.Equal(t, "username in body should be at least 1 chars long", *response.ValidationErrors[0].Error)

		cnt, err := models.PasswordResetTokens().Count(ctx, s.DB)
		assert.NoError(t, err)
		assert.Equal(t, int64(0), cnt)

		mail := getLastSentMail(t, s.Mailer)
		assert.Nil(t, mail)
	})
}

func TestPostForgotPasswordInvalidUsername(t *testing.T) {
	test.WithTestServer(t, func(s *api.Server) {
		ctx := context.Background()
		payload := test.GenericPayload{
			"username": "definitely not an email",
		}

		res := test.PerformRequest(t, s, "POST", "/api/v1/auth/forgot-password", payload, nil)

		assert.Equal(t, http.StatusBadRequest, res.Result().StatusCode)

		var response httperrors.HTTPValidationError
		test.ParseResponseAndValidate(t, res, &response)

		assert.Equal(t, int64(http.StatusBadRequest), *response.Code)
		assert.Equal(t, httperrors.HTTPErrorTypeGeneric, *response.Type)
		assert.Equal(t, http.StatusText(http.StatusBadRequest), *response.Title)
		assert.Empty(t, response.Detail)
		assert.Nil(t, response.Internal)
		assert.Nil(t, response.AdditionalData)
		assert.NotEmpty(t, response.ValidationErrors)
		assert.Equal(t, "username", *response.ValidationErrors[0].Key)
		assert.Equal(t, "body", *response.ValidationErrors[0].In)
		assert.Equal(t, "username in body must be of type email: \"definitely not an email\"", *response.ValidationErrors[0].Error)

		cnt, err := models.PasswordResetTokens().Count(ctx, s.DB)
		assert.NoError(t, err)
		assert.Equal(t, int64(0), cnt)

		mail := getLastSentMail(t, s.Mailer)
		assert.Nil(t, mail)
	})
}

func TestPostForgotPasswordSuccessLowercaseTrimWhitespaces(t *testing.T) {
	test.WithTestServer(t, func(s *api.Server) {
		ctx := context.Background()
		fixtures := test.Fixtures()
		payload := test.GenericPayload{
			"username": fmt.Sprintf(" %s ", strings.ToUpper(fixtures.User1.Username.String)),
		}

		res := test.PerformRequest(t, s, "POST", "/api/v1/auth/forgot-password", payload, nil)

		assert.Equal(t, http.StatusNoContent, res.Result().StatusCode)

		passwordResetToken, err := fixtures.User1.PasswordResetTokens().One(ctx, s.DB)
		require.NoError(t, err)

		mail := getLastSentMail(t, s.Mailer)
		require.NotNil(t, mail)
		assert.Contains(t, string(mail.HTML), fmt.Sprintf("http://localhost:3000/set-new-password?token=%s", passwordResetToken.Token))
	})
}
