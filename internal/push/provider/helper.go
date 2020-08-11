package provider

import "allaboutapps.dev/aw/go-starter/internal/push"

func sendMulticastWithProvider(p push.Provider, tokens []string, title, message string) []push.ProviderSendResponse {
	responseSlice := make([]push.ProviderSendResponse, 0)

	for _, token := range tokens {
		responseSlice = append(responseSlice, p.Send(token, title, message))
	}

	return responseSlice
}
