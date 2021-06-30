package transport

import (
	"crypto/tls"
	"fmt"
	"strings"
)

type SMTPAuthType int

const (
	// SMTPAuthTypeNone indicates no SMTP authentication should be performed.
	SMTPAuthTypeNone SMTPAuthType = iota
	// SMTPAuthTypePlain indicates SMTP authentication should be performed using the "AUTH PLAIN" protocol.
	SMTPAuthTypePlain
	// SMTPAuthTypeCRAMMD5 indicates SMTP authentication should be performed using the "CRAM-MD5" protocol.
	SMTPAuthTypeCRAMMD5
	// SMTPAuthTypeLogin indicates SMTP authentication should be performed using the "LOGIN" protocol.
	SMTPAuthTypeLogin
)

func (t SMTPAuthType) String() string {
	switch t {
	case SMTPAuthTypeNone:
		return "none"
	case SMTPAuthTypePlain:
		return "plain"
	case SMTPAuthTypeCRAMMD5:
		return "cram-md5"
	case SMTPAuthTypeLogin:
		return "login"
	default:
		return fmt.Sprintf("unknown (%d)", t)
	}
}

func SMTPAuthTypeFromString(s string) SMTPAuthType {
	switch strings.ToLower(s) {
	case "plain":
		return SMTPAuthTypePlain
	case "cram-md5":
		return SMTPAuthTypeCRAMMD5
	case "login":
		return SMTPAuthTypeLogin
	default:
		return SMTPAuthTypeNone
	}
}

type SMTPMailTransportConfig struct {
	Host      string
	Port      int
	AuthType  SMTPAuthType `json:"-"` // iota
	Username  string
	Password  string `json:"-"` // sensitive
	UseTLS    bool
	TLSConfig *tls.Config `json:"-"` // pointer
}
