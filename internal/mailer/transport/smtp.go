package transport

import (
	"fmt"
	"net"
	"net/smtp"
	"strconv"

	"github.com/jordan-wright/email"
)

type SMTPMailTransport struct {
	config SMTPMailTransportConfig
	addr   string
	auth   smtp.Auth
}

func NewSMTP(config SMTPMailTransportConfig) *SMTPMailTransport {
	mailTransport := &SMTPMailTransport{
		config: config,
		addr:   net.JoinHostPort(config.Host, strconv.Itoa(config.Port)),
		auth:   nil,
	}

	switch config.AuthType {
	case SMTPAuthTypePlain:
		mailTransport.auth = smtp.PlainAuth("", config.Username, config.Password, config.Host)
	case SMTPAuthTypeCRAMMD5:
		mailTransport.auth = smtp.CRAMMD5Auth(config.Username, config.Password)
	case SMTPAuthTypeLogin:
		mailTransport.auth = LoginAuth(config.Username, config.Password, config.Host)
	}

	return mailTransport
}

func (m *SMTPMailTransport) Send(mail *email.Email) error {
	var err error

	switch m.config.Encryption {
	case SMTPEncryptionNone:
		err = mail.Send(m.addr, m.auth)
	case SMTPEncryptionTLS:
		err = mail.SendWithTLS(m.addr, m.auth, m.config.TLSConfig)
	case SMTPEncryptionStartTLS:
		err = mail.SendWithStartTLS(m.addr, m.auth, m.config.TLSConfig)
	default:
		return fmt.Errorf("invalid SMTP encryption %q", m.config.Encryption)
	}

	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
