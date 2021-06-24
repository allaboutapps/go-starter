package util

import (
	"context"
	"errors"
)

type contextKey string

const (
	CTXKeyUser          contextKey = "user"
	CTXKeyAccessToken   contextKey = "access_token"
	CTXKeyRequestID     contextKey = "request_id"
	CTXKeyDisableLogger contextKey = "disable_logger"
)

func RequestIDFromContext(ctx context.Context) (string, error) {
	val := ctx.Value(CTXKeyRequestID)
	if val == nil {
		return "", errors.New("No request ID present in context")
	}

	id, ok := val.(string)
	if !ok {
		return "", errors.New("Request ID in context is not a string")
	}

	return id, nil
}

func ShouldDisableLogger(ctx context.Context) bool {
	s := ctx.Value(CTXKeyDisableLogger)
	if s == nil {
		return false
	}

	shouldDisable, ok := s.(bool)
	if !ok {
		return false
	}

	return shouldDisable
}

func DisableLogger(ctx context.Context, shouldDisable bool) context.Context {
	return context.WithValue(ctx, CTXKeyDisableLogger, shouldDisable)
}
