package util

type contextKey string

const (
	CTXKeyUser                               contextKey = "user"
	CTXKeyAccessToken contextKey = "access_token"
)
