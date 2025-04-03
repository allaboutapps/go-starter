package url

import (
	"net/url"
	"path"

	"allaboutapps.dev/aw/go-starter/internal/config"
)

const (
	queryParamToken = "token"
)

func PasswordResetDeeplinkURL(config config.Server, token string) (*url.URL, error) {
	u, err := url.Parse(config.Frontend.BaseURL)
	if err != nil {
		return nil, err
	}

	u.Path = path.Join(u.Path, config.Frontend.PasswordResetEndpoint)

	q := u.Query()
	q.Set(queryParamToken, token)
	u.RawQuery = q.Encode()

	return u, nil
}
