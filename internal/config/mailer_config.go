package config

type MailerTransporter string

var (
	MailerTransporterMock MailerTransporter = "mock"
	MailerTransporterSMTP MailerTransporter = "SMTP"
)

func (m MailerTransporter) String() string {
	return string(m)
}

type Mailer struct {
	DefaultSender               string
	Send                        bool
	WebTemplatesEmailBaseDirAbs string
	Transporter                 string
}
