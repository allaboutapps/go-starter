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

type SMTPEncryption string

const (
	SMTPEncryptionNone     SMTPEncryption = "none"
	SMTPEncryptionTLS      SMTPEncryption = "tls"
	SMTPEncryptionStartTLS SMTPEncryption = "starttls"
)

func (e SMTPEncryption) String() string {
	return string(e)
}

func SMTPEncryptionFromString(s string) SMTPEncryption {
	switch strings.ToLower(s) {
	case "tls":
		return SMTPEncryptionTLS
	case "starttls":
		return SMTPEncryptionStartTLS
	default:
		return SMTPEncryptionNone
	}
}

type SMTPMailTransportConfig struct {
	Host       string
	Port       int
	AuthType   SMTPAuthType `json:"-"` // iota
	Username   string
	Password   string         `json:"-"` // sensitive
	Encryption SMTPEncryption `json:"-"` // iota
	TLSConfig  *tls.Config    `json:"-"` // pointer
	UseTLS     bool           // ! deprecated since 2021-10-25, use Encryption type 'SMTPEncryptionTLS' instead
}
