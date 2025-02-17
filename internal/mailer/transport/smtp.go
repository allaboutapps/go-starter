package transport

import (
	"fmt"
	"net"
	"net/smtp"

	"github.com/jordan-wright/email"
	"github.com/rs/zerolog/log"
)

type SMTPMailTransport struct {
	config SMTPMailTransportConfig
	addr   string
	auth   smtp.Auth
}

func NewSMTP(config SMTPMailTransportConfig) *SMTPMailTransport {
	m := &SMTPMailTransport{
		config: config,
		addr:   net.JoinHostPort(config.Host, fmt.Sprintf("%d", config.Port)),
		auth:   nil,
	}

	switch config.AuthType {
	case SMTPAuthTypePlain:
		m.auth = smtp.PlainAuth("", config.Username, config.Password, config.Host)
	case SMTPAuthTypeCRAMMD5:
		m.auth = smtp.CRAMMD5Auth(config.Username, config.Password)
	case SMTPAuthTypeLogin:
		m.auth = LoginAuth(config.Username, config.Password, config.Host)
	}

	return m
}

func (m *SMTPMailTransport) Send(mail *email.Email) error {
	if m.config.UseTLS {
		log.Warn().Msg("Enabling TLS with the UseTLS flag is *DEPRECATED* and will be removed in a future release")
		return mail.SendWithTLS(m.addr, m.auth, m.config.TLSConfig)
	}

	switch m.config.Encryption {
	case SMTPEncryptionNone:
		return mail.Send(m.addr, m.auth)
	case SMTPEncryptionTLS:
		return mail.SendWithTLS(m.addr, m.auth, m.config.TLSConfig)
	case SMTPEncryptionStartTLS:
		return mail.SendWithStartTLS(m.addr, m.auth, m.config.TLSConfig)
	default:
		return fmt.Errorf("invalid SMTP encryption %q", m.config.Encryption)
	}
}
