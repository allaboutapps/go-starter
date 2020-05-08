package transport

import "github.com/jordan-wright/email"

type MailTransporter interface {
	Send(mail *email.Email) error
}
