package provider

import (
	"errors"

	"allaboutapps.dev/aw/go-starter/internal/push"
	"github.com/rs/zerolog/log"
)

type Mock struct {
	Type push.ProviderType
}

func NewMock(providerType push.ProviderType) *Mock {
	return &Mock{
		Type: providerType,
	}
}

func (p *Mock) GetProviderType() push.ProviderType {
	return p.Type
}

func (p *Mock) Send(token string, title string, message string) push.ProviderSendResponse {
	valid := true
	var err error
	if len(token) < 40 {
		valid = false
		err = errors.New("invalid token")
	}

	if title == "other error" {
		err = errors.New("other error")
	}

	log.Info().Str("token", token).Str("title", title).Str("message", message).Msg("Mock Push Notification")

	return push.ProviderSendResponse{
		Token: token,
		Valid: valid,
		Err:   err,
	}
}

func (p *Mock) SendMulticast(tokens []string, title, message string) []push.ProviderSendResponse {
	return sendMulticastWithProvider(p, tokens, title, message)
}
