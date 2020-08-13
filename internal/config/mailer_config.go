package config

type Mailer struct {
	DefaultSender               string
	Send                        bool
	WebTemplatesEmailBaseDirAbs string
	UserMockTransporter         bool
}
