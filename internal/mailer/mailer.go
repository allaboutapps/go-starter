package mailer

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"html/template"
	"os"
	"path/filepath"

	"allaboutapps.dev/aw/go-starter/internal/config"
	"allaboutapps.dev/aw/go-starter/internal/data/dto"
	"allaboutapps.dev/aw/go-starter/internal/mailer/transport"
	"allaboutapps.dev/aw/go-starter/internal/util"
	"github.com/jordan-wright/email"
	"github.com/rs/zerolog/log"
)

var (
	ErrEmailTemplateNotFound         = errors.New("email template not found")
	emailTemplatePasswordReset       = "password_reset"       // /app/templates/email/password_reset/**.
	emailTemplateAccountConfirmation = "account_confirmation" // /app/templates/email/account_confirmation/**
)

type Mailer struct {
	Config    config.Mailer
	Transport transport.MailTransporter
	Templates map[string]*template.Template
}

func New(config config.Mailer, transport transport.MailTransporter) *Mailer {
	return &Mailer{
		Config:    config,
		Transport: transport,
		Templates: map[string]*template.Template{},
	}
}

func NewWithConfig(cfg config.Mailer, smtpConfig transport.SMTPMailTransportConfig) (*Mailer, error) {
	var mailer *Mailer

	switch config.MailerTransporter(cfg.Transporter) {
	case config.MailerTransporterMock:
		log.Warn().Msg("Initializing mock mailer")
		mailer = New(cfg, transport.NewMock())
	case config.MailerTransporterSMTP:
		mailer = New(cfg, transport.NewSMTP(smtpConfig))
	default:
		return nil, fmt.Errorf("unsupported mail transporter: %s", cfg.Transporter)
	}

	if err := mailer.ParseTemplates(); err != nil {
		return nil, fmt.Errorf("failed to parse mailer templates: %w", err)
	}

	return mailer, nil
}

func (m *Mailer) ParseTemplates() error {
	files, err := os.ReadDir(m.Config.WebTemplatesEmailBaseDirAbs)
	if err != nil {
		log.Error().Str("dir", m.Config.WebTemplatesEmailBaseDirAbs).Err(err).Msg("Failed to read email templates directory while parsing templates")
		return fmt.Errorf("failed to read email templates directory while parsing templates: %w", err)
	}

	for _, file := range files {
		if !file.IsDir() {
			continue
		}

		tmpl, err := template.ParseGlob(filepath.Join(m.Config.WebTemplatesEmailBaseDirAbs, file.Name(), "**"))
		if err != nil {
			log.Error().Str("template", file.Name()).Err(err).Msg("Failed to parse email template files as glob")
			return fmt.Errorf("failed to parse email template files as glob: %w", err)
		}

		m.Templates[file.Name()] = tmpl
	}

	return nil
}

func (m *Mailer) SendPasswordReset(ctx context.Context, to string, passwordResetLink string) error {
	log := util.LogFromContext(ctx).With().Str("component", "mailer").Str("email_template", emailTemplatePasswordReset).Logger()

	tmpl, ok := m.Templates[emailTemplatePasswordReset]
	if !ok {
		log.Error().Msg("Password reset email template not found")
		return ErrEmailTemplateNotFound
	}

	data := map[string]interface{}{
		"passwordResetLink": passwordResetLink,
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		log.Error().Err(err).Msg("Failed to execute password reset email template")
		return fmt.Errorf("failed to execute password reset email template: %w", err)
	}

	mail := email.NewEmail()

	mail.From = m.Config.DefaultSender
	mail.To = []string{to}
	mail.Subject = "Password reset"
	mail.HTML = buf.Bytes()

	if !m.Config.Send {
		log.Warn().Str("to", to).Str("passwordResetLink", passwordResetLink).Msg("Sending has been disabled in mailer config, skipping password reset email")
		return nil
	}

	if err := m.Transport.Send(mail); err != nil {
		log.Debug().Err(err).Msg("Failed to send password reset email")
		return fmt.Errorf("failed to send password reset email: %w", err)
	}

	log.Debug().Msg("Successfully sent password reset email")

	return nil
}

func (m *Mailer) SendAccountConfirmation(ctx context.Context, to string, payload dto.ConfirmatioNotificationPayload) error {
	log := util.LogFromContext(ctx)

	tmpl, ok := m.Templates[emailTemplateAccountConfirmation]
	if !ok {
		log.Error().Msg("Account confirmation email template not found")
		return ErrEmailTemplateNotFound
	}

	data := map[string]interface{}{
		"confirmationLink": payload.ConfirmationLink,
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		log.Error().Err(err).Msg("Failed to execute account confirmation email template")
		return fmt.Errorf("failed to execute account confirmation email template: %w", err)
	}

	mail := email.NewEmail()

	mail.From = m.Config.DefaultSender
	mail.To = []string{to}
	mail.Subject = "Account confirmation"
	mail.HTML = buf.Bytes()

	if !m.Config.Send {
		log.Warn().Str("to", to).Msg("Sending has been disabled in mailer config, skipping account confirmation email")
		return nil
	}

	if err := m.Transport.Send(mail); err != nil {
		log.Debug().Err(err).Msg("Failed to send account confirmation email")
		return fmt.Errorf("failed to send account confirmation email: %w", err)
	}

	return nil
}
