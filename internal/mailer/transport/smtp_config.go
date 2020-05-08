package transport

import (
	"crypto/tls"
	"fmt"
	"strings"
)

type SMTPAuthType int

const (
	SMTPAuthTypeNone SMTPAuthType = iota
	SMTPAuthTypePlain
	SMTPAuthTypeCRAMMD5
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
	AuthType  SMTPAuthType
	Username  string
	Password  string
	UseTLS    bool
	TLSConfig *tls.Config
}
