package url

import (
	"fmt"
	"net/url"
	"path"

	"allaboutapps.dev/aw/go-starter/internal/config"
)

const (
	queryParamToken         = "token"
	accountConfirmationPath = "/api/v1/auth/register"
)

func PasswordResetDeeplinkURL(config config.Server, token string) (*url.URL, error) {
	u, err := url.Parse(config.Frontend.BaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse the base URL: %w", err)
	}

	u.Path = path.Join(u.Path, config.Frontend.PasswordResetEndpoint)

	q := u.Query()
	q.Set(queryParamToken, token)
	u.RawQuery = q.Encode()

	return u, nil
}

func ConfirmationDeeplinkURL(config config.Server, token string) (*url.URL, error) {
	u, err := url.Parse(config.Echo.BaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse the base URL: %w", err)
	}

	u.Path = path.Join(u.Path, accountConfirmationPath)

	q := u.Query()
	q.Set(queryParamToken, token)
	u.RawQuery = q.Encode()

	return u, nil
}

func ConfirmationRequestURL(config config.Server, token string) (*url.URL, error) {
	u, err := url.Parse(config.Echo.BaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse the base URL: %w", err)
	}

	u.Path = path.Join(u.Path, accountConfirmationPath, token)

	return u, nil
}
