package util

import (
	"context"
	"errors"
	"time"
)

type contextKey string

const (
	CTXKeyUser          contextKey = "user"
	CTXKeyAccessToken   contextKey = "access_token"
	CTXKeyCacheControl  contextKey = "cache_control"
	CTXKeyRequestID     contextKey = "request_id"
	CTXKeyDisableLogger contextKey = "disable_logger"
)

type detachedContext struct {
	parent context.Context
}

func (c detachedContext) Deadline() (time.Time, bool)       { return time.Time{}, false }
func (c detachedContext) Done() <-chan struct{}             { return nil }
func (c detachedContext) Err() error                        { return nil }
func (c detachedContext) Value(key interface{}) interface{} { return c.parent.Value(key) }

// DetachContext detaches a context by returning a wrapped struct implementing the context interface, but omitting the deadline, done and error functionality.
// Mainly used to pass context information to go routines that should not be cancelled by the context.
// ! USE THIS DETACHED CONTEXT SPARINGLY, ONLY IF ABSOLUTELY NEEDED. DO *NOT* KEEP USING A DETACHED CONTEXT FOR A PROLONGED TIME OUT OF CHAIN
func DetachContext(ctx context.Context) context.Context {
	return detachedContext{ctx}
}

// RequestIDFromContext returns the ID of the (HTTP) request, returning an error if it is not present.
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

// ShouldDisableLogger checks whether the logger instance should be disabled for the provided context.
// `util.LogFromContext` will use this function to check whether it should return a default logger if
// none has been set by our logging middleware before, or fall back to the disabled logger, suppressing
// all output. Use `ctx = util.DisableLogger(ctx, true)` to disable logging for the given context.
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

// DisableLogger toggles the indication whether `util.LogFromContext` should return a disabled logger
// for a context if none has been set by our logging middleware before. Whilst the usecase for a disabled
// logger are relatively minimal (we almost always want to have some log output, even if the context
// was not directly derived from a HTTP request), this functionality was provideds so you can switch back
// to the old zerolog behavior if so desired.
func DisableLogger(ctx context.Context, shouldDisable bool) context.Context {
	return context.WithValue(ctx, CTXKeyDisableLogger, shouldDisable)
}
