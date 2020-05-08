package transport

import (
	"errors"
	"fmt"
	"net/smtp"
)

type loginAuth struct {
	username []byte
	password []byte
	host     string
}

func LoginAuth(username string, password string, host ...string) smtp.Auth {
	a := &loginAuth{
		username: []byte(username),
		password: []byte(password),
		host:     "",
	}

	if len(host) > 0 {
		a.host = host[0]
	}

	return a
}

func (a *loginAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	if len(a.host) > 0 && server.Name != a.host {
		return "", nil, errors.New("wrong host name")
	}

	return "LOGIN", a.username, nil
}

func (a *loginAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	if more {
		switch string(fromServer) {
		case "Username:":
			return a.username, nil
		case "Password:":
			return a.password, nil
		default:
			return nil, fmt.Errorf("Unknown challenge received from server: %q", fromServer)
		}
	}

	return nil, nil
}
